package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/fmowl10/testJwt/utils"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// TokenDurationTime is duration time of token
//
var TokenDurationTime = time.Minute * 5

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var hub = utils.NewClients()

func jwtAuthenticationFunc(ctx echo.Context) error {
	u := new(utils.User)
	if err := ctx.Bind(u); err != nil {
		return ctx.String(http.StatusForbidden, err.Error())
	}
	claims := &utils.JwtClaim{
		*u,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte("jwtTokenKey"))
	if err != nil {
		return ctx.String(http.StatusServiceUnavailable, "no")
	}

	return ctx.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

// HelloAPI write Hello world.
//
func HelloAPI(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "hello")
}

func websocketHandler(ctx echo.Context) error {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	subProtocol := ctx.Request().Header.Get("Sec-WebSocket-Protocol")
	token, _ := jwt.Parse(subProtocol, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error")
		}
		return []byte("jwtTokenKey"), nil
	})
	if !token.Valid {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"message": "invalid or expired jwt"})
	}
	upgrader.Subprotocols = []string{subProtocol}

	ws, err := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = ws.WriteMessage(1, []byte("Hi Client!"))
	return read(ws, ctx)
}

func read(conn *websocket.Conn, ctx echo.Context) error {
	var name string
	messageType, p, err := conn.ReadMessage()
	if err != nil {
		return nil
	}
	if messageType == websocket.TextMessage {
		name = string(p)
	} else {
		return nil
	}
	clientChan := make(chan string)
	if err != nil {
		return nil
	}
	hub.Sign(clientChan)

	go write(clientChan, conn, name)
	go hub.Broadcast("Client " + name + " conneted")

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			hub.Broadcast("Client " + name + " disconneted")
			hub.Unsign(clientChan)
			conn.Close()
			return nil
		}
		hub.Broadcast(name + " say: " + string(p))
	}
}

func write(clientChan chan string, conn *websocket.Conn, name string) {
	for {
		s := <-clientChan
		err := conn.WriteMessage(websocket.TextMessage, []byte(s))
		if err != nil {
			return
		}
	}
}

func main() {
	// init server
	e := echo.New()
	host := flag.String("host", ":3000", "set host")

	// run hub
	go hub.Hub()

	// echo Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// web page setting
	e.Static("/", "frontend/build")
	//e.File("/favicon.ico", "frontend/build/favicon.ico")
	e.Static("/static", "frontend/build/static")

	// jwt config
	helloRoute := e.Group("/api")

	// claims setting
	config := middleware.JWTConfig{
		Claims:     &utils.JwtClaim{},
		SigningKey: []byte("jwtTokenKey"),
	}
	helloRoute.Use(middleware.JWTWithConfig(config))
	helloRoute.GET("", HelloAPI)
	// routing
	e.POST("/jwt", jwtAuthenticationFunc)

	// api and websocket routing
	//e.GET("/api", HelloAPI)
	e.GET("/ws", websocketHandler)

	// running
	log.Fatal(e.Start(*host))
}

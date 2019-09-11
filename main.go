package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"testJwt/backend"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

// TokenDurationTime is duration time of token
//
var TokenDurationTime = time.Minute * 5

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var hub = backend.Clients{}

func jwtAuthenticationFunc(ctx echo.Context) error {
	length, _ := strconv.Atoi(ctx.Request().Header.Get("Content-Length"))
	body := make([]byte, length)
	ctx.Request().Body.Read(body)
	var Data map[string]interface{}
	json.Unmarshal(body, &Data)
	now := time.Now()
	Data["exp"] = now.Add(TokenDurationTime).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(Data))
	key := []byte("test")
	ss, _ := token.SignedString(key)
	return ctx.String(http.StatusOK, ss)
}

// HelloAPI write Hello world.
//
func HelloAPI(ctx echo.Context) error {
	rawToken := ctx.Request().Header.Get("JWT")
	if len(rawToken) == 0 {
		return ctx.String(http.StatusForbidden, "Token not included")
	}
	if backend.VaildToken(rawToken) {
		return ctx.String(http.StatusOK, "Hello World!")
	} else {
		return ctx.String(http.StatusBadRequest, "Token error")
	}
}

func websocketHandler(ctx echo.Context) error {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	subProtocol := ctx.Request().Header.Get("Sec-WebSocket-Protocol")
	if !backend.VaildToken(subProtocol) {
		return ctx.String(http.StatusForbidden, "Token error")
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
	e := echo.New()
	host := flag.String("host", ":3000", "set host")
	hub.Init()
	go hub.Hub()
	e.Static("/", "frontend/build")
	e.File("/favicon.ico", "frontend/build/favicon.ico")
	e.Static("/static", "frontend/build/static")
	e.POST("/jwt", jwtAuthenticationFunc)

	e.GET("/api", HelloAPI)
	e.GET("/ws", websocketHandler)
	log.Fatal(e.Start(*host))
}

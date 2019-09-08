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
)

// https://ops.tips/blog/nginx-http2-server-push/
// http/2 지원관련 https://blog.cloudflare.com/tools-for-debugging-testing-and-using-http-2/
// https://github.com/gaurav-gogia/simple-http2-server
// https://stackoverflow.com/questions/37321760/how-to-set-up-lets-encrypt-for-a-go-server-application
// delev 사용법 https://github.com/campoy/go-tooling-workshop/blob/master/3-dynamic-analysis/1-debugging/1-delve.md

// TokenDurationTime is duration time of token
//
var TokenDurationTime = time.Minute * 5

// wait time 설정하자
// nginx 의 설정이다
// https://stackoverflow.com/questions/28828332/gorilla-websocket-disconnects-after-a-minute

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var hub = backend.Clients{}

func jwtAuthenticationFunc(w http.ResponseWriter, r *http.Request) {
	length, _ := strconv.Atoi(r.Header.Get("Content-Length"))
	body := make([]byte, length)
	r.Body.Read(body)
	var Data map[string]interface{}
	json.Unmarshal(body, &Data)
	now := time.Now()
	Data["exp"] = now.Add(TokenDurationTime).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(Data))
	key := []byte("test")
	ss, _ := token.SignedString(key)
	w.Write([]byte(ss))
}

// HelloAPI write Hello world.
//
func HelloAPI(w http.ResponseWriter, r *http.Request) {
	rawToken := r.Header.Get("JWT")
	if len(rawToken) == 0 {
		backend.ResponseError(w, http.StatusForbidden, "Can not found token")
		return
	}
	if backend.VaildToken(rawToken) {
		w.Write([]byte("Hello World!"))
	} else {
		backend.ResponseError(w, http.StatusBadRequest, "Token error")
	}
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	subProtocol := r.Header.Get("Sec-WebSocket-Protocol")
	if !backend.VaildToken(subProtocol) {
		backend.ResponseError(w, http.StatusForbidden, "Token error")
		return
	}
	upgrader.Subprotocols = []string{subProtocol}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
	}

	err = ws.WriteMessage(1, []byte("Hi Client!"))
	go read(ws)
}

func read(conn *websocket.Conn) {
	var name string
	messageType, p, err := conn.ReadMessage()
	if err != nil {
		return
	}
	if messageType == websocket.TextMessage {
		name = string(p)
	} else {
		return
	}
	clientChan := make(chan string)
	if err != nil {
		return
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
			return
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
	host := flag.String("host", ":3000", "set host")
	hub.Init()
	go hub.Hub()
	//https://blog.kowalczyk.info/article/Jl3G/https-for-free-in-go-with-little-help-of-lets-encrypt.html
	//https://stackoverflow.com/questions/15394904/nginx-load-balance-with-upstream-ssl
	http.Handle("/", http.FileServer(http.Dir("./frontend/build")))
	http.Handle("/static", http.StripPrefix("/static", http.FileServer(http.Dir("./frontend/build/static"))))

	http.HandleFunc("/jwt", jwtAuthenticationFunc)
	http.HandleFunc("/api", HelloAPI)
	http.HandleFunc("/ws", websocketHandler)
	log.Fatal(http.ListenAndServe(*host, nil))
}

//https://www.youtube.com/watch?v=mML6GiOAM1w

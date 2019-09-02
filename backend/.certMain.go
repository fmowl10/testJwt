package main


import (
	//"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"

	//"github.com/pkg/errors"
	"golang.org/x/crypto/acme/autocert"
)

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

func responseError(w http.ResponseWriter, statusCode int, responseBody string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(responseBody))
	return
}
func vaildToken(rawToken string) bool {
	token, _ := jwt.Parse(rawToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte("test"), nil
	})
	if claim, ok := token.Claims.(jwt.MapClaims); ok {
		if err := claim.Valid(); err != nil {
			return false
		}
	}
	return true
}

// HelloAPI write Hello world.
//
func HelloAPI(w http.ResponseWriter, r *http.Request) {
	rawToken := r.Header.Get("JWT")
	if len(rawToken) == 0 {
		responseError(w, http.StatusForbidden, "Can not found token")
		return
	}
	if vaildToken(rawToken) {
		w.Write([]byte("Hello World!"))
	} else {
		responseError(w, http.StatusBadRequest, "Token error")
	}
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	subProtocol := r.Header.Get("Sec-WebSocket-Protocol")
	if !vaildToken(subProtocol) {
		responseError(w, http.StatusForbidden, "Token error")
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
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(errors.Wrap(err, "error"))
			return
		}
		fmt.Println(time.Now())
		fmt.Println(string(p))
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(errors.Wrap(err, "error"))
			return
		}
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	if pusher, ok := w.(http.Pusher); ok {
		fmt.Println(ok)
		if err := pusher.Push("../frontend/build/static/css/2.621b5bde.chunk.css", nil); err != nil {
			return
		}
		if err := pusher.Push("../frontend/build/static/js/2.13299181.chunk.js", nil); err != nil {
			return
		}
		if err := pusher.Push("../frontend/build/static/js/main.4aabe457.chunk.js", nil); err != nil {
			return
		}
		if err := pusher.Push("../frontend/build/mainfest.json", nil); err != nil {
			return
		}
	} else {
		http.FileServer(http.Dir("../frontend/build/"))
	}
}

func makeHTTP() *http.Server {
	mux := &http.ServeMux{}

	mux.HandleFunc("/", index)
	mux.Handle("/static", http.StripPrefix("/static", http.FileServer(http.Dir("../frontend/build/static"))))

	mux.HandleFunc("/jwt", jwtAuthenticationFunc)
	mux.HandleFunc("/api", HelloAPI)
	mux.HandleFunc("/ws", websocketHandler)
	return &http.Server{
		Handler: mux,
	}
}

func main() {
	host := flag.String("host", ":3000", "set host")
	//https://blog.kowalczyk.info/article/Jl3G/https-for-free-in-go-with-little-help-of-lets-encrypt.html
	//https://stackoverflow.com/questions/15394904/nginx-load-balance-with-upstream-ssl
	/*	hostPolicy := func(ctx context.Context, host string) error {
			allowedHost := "www.fmowl.com"
			if host == allowedHost {
				return nil
			}
			return fmt.Errorf("no")
		}
		m := &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: hostPolicy,
			Cache:      autocert.DirCache("."),
		}
		server := makeHttp()
		server.Addr = ":3000"
		server.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}

		log.Fatal(server.ListenAndServeTLS("", ""))
	*/

	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("test.fmowl.com"),
		Cache:      autocert.DirCache("."),
	}

	server := &http.Server{
		Addr: *host,
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}

	http.Handle("/", http.FileServer(http.Dir("../frontend/build")))
	http.Handle("/static", http.StripPrefix("/static", http.FileServer(http.Dir("../frontend/build/static"))))

	http.HandleFunc("/jwt", jwtAuthenticationFunc)
	http.HandleFunc("/api", HelloAPI)
	http.HandleFunc("/ws", websocketHandler)
	go http.ListenAndServe("2999", certManager.HTTPHandler(nil))
	log.Fatal(server.ListenAndServeTLS("", ""))
}

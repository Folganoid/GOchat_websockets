package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net/http"
	"os"
)

type (
	Msg struct {
		clientKey   string
		messageText string
	}

	NewClientEvent struct {
		clientKey string
		msgChan   chan Msg
	}
)

const MAXBACKLOG = 100

var (
	dirPath           string
	clientRequest     = make(chan *NewClientEvent, 100)
	clientDisconnects = make(chan string)
)

func IndexPage(w http.ResponseWriter, req *http.Request, filename string) {

	fp, err := os.Open(dirPath + "/" + filename)
	if err != nil {
		log.Println("Could not open file", err.Error())
		w.Write([]byte("500 internal server error"))
		return
	}

	defer fp.Close()

	_, err = io.Copy(w, fp)
	if err != nil {
		log.Println("Could not send file contents", err.Error())
		w.Write([]byte("500 internal server error"))
		return
	}

}

func router() {
	clients := make(map[string]chan Msg)

	for {
		select {
		case req := <-clientRequest:
			clients[req.clientKey] = req.msgChan
			log.Println("Websocket connected: " + req.clientKey)
		case clientKey := <-clientDisconnects:
			delete(clients, clientKey)
			log.Println("Websocket disconnected: " + clientKey)
		}
	}
}

// Echo the data received on the WebSocket.
func EchoServer(ws *websocket.Conn) {

	msgChan := make(chan Msg, 100)
	clientKey := ws.RemoteAddr().String()
	clientRequest <- &NewClientEvent{clientKey, msgChan}
	defer func() { clientDisconnects <- clientKey }()

	_, err := io.Copy(ws, ws)
	if err != nil {
		log.Println("Client error: " + err.Error())
	}
}

func main() {

	if len(os.Args) < 2 {
		log.Fatal("Usage: chatExample <dir>")
	}

	dirPath = os.Args[1]

	fmt.Println("Starting...")

	go router()

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		IndexPage(w, req, "index.html")
	})
	http.HandleFunc("/index.js", func(w http.ResponseWriter, req *http.Request) {
		IndexPage(w, req, "index.js")
	})
	http.Handle("/ws", websocket.Handler(EchoServer))
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal("Error: ", err)
	}
}

package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net/http"
	"os"
)

var dirPath string

func IndexPage(w http.ResponseWriter, req *http.Request) {

	fp, err := os.Open(dirPath + "/index.html")
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

// Echo the data received on the WebSocket.
func EchoServer(ws *websocket.Conn) {
	log.Println("Websocket connected: " + ws.RemoteAddr().String())
	defer log.Println("Websocket disconnected: " + ws.RemoteAddr().String())

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

	http.HandleFunc("/", IndexPage)
	http.Handle("/ws", websocket.Handler(EchoServer))
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal("Error: ", err)
	}
}

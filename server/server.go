package server

import (
	"github.com/gorilla/websocket"
	"net/http"
	"log"
	"fmt"
	"html/template"
	"errors"
	"bytes"
)

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan Message)           // broadcast channel
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}                                            // 将http升级为websocket

type Message struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func onConnect(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(string(p))
		fmt.Println(messageType)
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

func HomePage(res http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles("I:\\ubuntu14.04\\share\\docker\\go\\www\\src\\github.com\\kong36088\\danmu\\server\\client\\index.htmlgit ")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = WriteTemplateToHttpResponse(res, t)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func WriteTemplateToHttpResponse(res http.ResponseWriter, t *template.Template) error {
	if t == nil || res == nil {
		return errors.New("WriteTemplateToHttpResponse: t must not be nil.")
	}
	var buf bytes.Buffer
	err := t.Execute(&buf, nil)
	if err != nil {
		return err
	}
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err = res.Write(buf.Bytes())
	return err
}

func StartServer() {
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/ws", onConnect)
	err := http.ListenAndServe(":9500", nil)
	if err != nil {
		log.Fatal(err)
	}
}

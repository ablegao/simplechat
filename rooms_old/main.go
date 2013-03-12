package main

import (
	"code.google.com/p/go.net/websocket"

	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
)

var serverAddr string

func init() {
	http.Handle("/games", websocket.Handler(RoomBuildConnection))
}

//测试服务器入口
func RoomTestServerStart() {
	server := httptest.NewServer(nil)
	serverAddr = server.Listener.Addr().String()
	log.Print("Test WebSocket server listening on ", serverAddr)

}

func RoomServerStart() {
	if serverAddr := http.ListenAndServe(":1234", nil); serverAddr != nil {
		log.Fatal("ListenAndServe:", serverAddr)
	}
}

func main() {
	//RoomTestServerStart()
	RoomServerStart()
}

var Mess string

// 测试第一步骤， 尝试一个可以分room 的聊天室
func BuildRoom(ws *websocket.Conn) {
	var content string
	var err error
	for {
		token := ws.Request().URL.Query().Get("token") //用户token 
		roomName := ws.Request().URL.Query().Get("r")  //用户token 
		err = websocket.Message.Receive(ws, &content)
		// If user closes or refreshes the browser, a err will occur
		if err != nil {
			fmt.Printf(" error 1: %v ", err)

			break
		}
		Mess = "callback :  " + token + " " + roomName + " " + content
		err = websocket.Message.Send(ws, Mess)
		if err != nil {
			fmt.Printf(" error 2: %v ", err)
			break
		}
	}
}

package main

import (
	"code.google.com/p/go.net/websocket"
	//"net/http/httptest"
	"bytes"
	"fmt"
	"net"
	"sync"
	"testing"
)

type Count struct {
	S string
	N int
}

var once sync.Once

//客户连接配置
func newConfig(t *testing.T, path string) *websocket.Config {
	config, _ := websocket.NewConfig(fmt.Sprintf("ws://%s%s", serverAddr, path), "http://localhost")
	return config
}

func TestServer(t *testing.T) {
	once.Do(RoomTestServerStart)
	// websocket.Dial()
	client, err := net.Dial("tcp", serverAddr)
	if err != nil {
		t.Fatal("dialing", err)
	}

	//建立连接
	conn, err := websocket.NewClient(newConfig(t, "/games?token=able&r=room1"), client)
	if err != nil {
		t.Errorf("WebSocket handshake error: %v", err)
		return
	}
	var msg []byte
	msg = []byte("hello, world\n")
	if _, err := conn.Write(msg); err != nil {
		t.Errorf("Write: %v", err)
	}
	var actual_msg = make([]byte, 512)
	n, err := conn.Read(actual_msg)
	if err != nil {
		t.Errorf("Read: %v", err)
	}
	actual_msg = actual_msg[0:n]
	if !bytes.Equal(msg, actual_msg) {
		t.Errorf("Echo: expected %q got %q", msg, actual_msg)
	}

	msg = []byte("测试2\n")
	if _, err := conn.Write(msg); err != nil {
		t.Errorf("Write: %v\n", err)
	}
	actual_msg = make([]byte, 512)
	n, err = conn.Read(actual_msg)
	if err != nil {
		t.Errorf("Read: %v\n", err)
	}
	actual_msg = actual_msg[0:n]
	if !bytes.Equal(msg, actual_msg) {
		t.Errorf("Echo: expected %q got %q\n", msg, actual_msg)
	}

}

package main

// 主界面

import (
	//"code.google.com/p/go.net/websocket"
	//"io"
	"net/http"
	"net/http/httptest"
	"runtime"
)

var serverAddr string

var ActiveChannel *RoomChannel = &RoomChannel{}

func init() {
	ActiveChannel = &RoomChannel{
		ServerName:   SERVER_NAME,
		ChannelClose: make(chan bool),
		Rooms:        make(map[string]*RoomBase),
		RoomId:       0,
	}
	go ActiveChannel.Run()
	//http.Handle("/games", websocket.Handler(RoomBuildConnection))
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)

	//Log.SetOutput(io.MultiWriter(logf, os.Stdout))

	Log.Println("ccccccc")

	RoomServerStart()
}

//测试服务器入口
func RoomTestServerStart() {
	server := httptest.NewServer(nil)
	serverAddr = server.Listener.Addr().String()
	Log.Print("Test WebSocket server listening on ", serverAddr)

}

//正式服务器
func RoomServerStart() {
	if serverAddr := http.ListenAndServe(":1234", nil); serverAddr != nil {
		Log.Fatal("ListenAndServe:", serverAddr)
	}
}

// room 控制服务器 
// 主要功能是针对命令创建房间。 

type RoomChannel struct {
	ServerName   string
	ChannelClose chan bool //服务停止信号
	Rooms        map[string]*RoomBase
	RoomId       int // room id 自增长数据
}

//房间销毁或者后台服务器终止操作 
func (rc *RoomChannel) Run() {

	for {
		select {
		case out := <-rc.ChannelClose:
			if out == true {
				close(rc.ChannelClose)
				return
			}
		}

	}
}

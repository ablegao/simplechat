package main

//rooms 算是一个点，作为一台服务器的一个点, 一台服务器上可以简历多个room， 每次简历后，将数据提交给channel ， 做统一播放。 

//

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
)

var runingRooms map[string]Room

func init() {
	//runingRooms = make(map[string]Room)
	runingRooms["room1"] = Room{
		Name:        "room1",
		MapInfo:     "map test",
		Broadcast:   make(chan Message),
		CloseSign:   make(chan bool),
		IsStart:     make(chan bool),
		OnlineUsers: make(map[string]*OnlineUser),
	}

}

type Room struct {
	Name        string                 //房间名称
	MapInfo     string                 //地图信息
	Broadcast   chan Message           //消息传递
	CloseSign   chan bool              //关闭命令
	IsStart     chan bool              //游戏状态， 是否已经开始运行。
	OnlineUsers map[string]*OnlineUser //房间内的用户  
}

//核心消息处理线程， 这会是一条独立的线程， 任何被放入公共消息池子的内容， 将会被循环写入到onlineUser.Send 当中。
func (this *Room) Run() {
	for {
		select {
		case b := <-this.Broadcast:
			for _, online := range this.OnlineUsers {
				online.Send <- b
			}

		//这个消息应该会终止后台任务.
		case c := <-this.CloseSign:
			if c == true {
				close(this.Broadcast)
				close(this.CloseSign)
				return
			}
		}
	}
}

// 用户登陆处理， 建立socket 连接
func RoomBuildConnection(ws *websocket.Conn) {
	token := ws.Request().URL.Query().Get("token") //用户token 

	roomId := ws.Request().URL.Query().Get("r") //room 

	fmt.Printf("%s RoomBuildConnection connect web socket!\n", humanCreatedAt())
	if token == "" || roomId == "" {
		fmt.Printf("%s %s \n", roomId, token)
		return
	}
	fmt.Printf("%s ", roomId)

	/*
		// 创建一个user 对象
		onlineUser := &OnlineUser{
			InRoom:     &runingRooms[roomId],
			Connection: ws,
			Send:       make(chan Message, 256), //存放256条消息
			UserInfo: &User{
				Id:     1,
				Name:   token,
				Face:   "xxxxx",
				Token:  token,
				RoomId: roomId,
			},
		}

		runingRooms[roomId].OnlineUsers[token] = onlineUser

		//写入一条消息队列
		m := Message{
			MType:      STATUS_MTYPE,
			UserStatus: UserStatus{},
		}
		runingRooms[roomId].Broadcast <- m

		//推送消息到客户端， 这在一个新的线程中操作， 避免主线的消息获取受阻拦。  
		go onlineUser.PushToClient()

		// 获取客户端消息 , PullFromClient 是一个循环锁，一只等待用户自己的信息反馈， 当用户关闭url 的时候， 这个锁将被退出 ，进入killUserResource 
		onlineUser.PullFromClient()
		// 关闭连接
		onlineUser.killUserResource()
	*/
}

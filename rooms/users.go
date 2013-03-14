package main

import (
	"code.google.com/p/go.net/websocket"
)

//用户类型
type OnlineUser struct {
	Send       chan interface{}
	InRoom     string
	Connection *websocket.Conn
	UserInfo   User
}

type User struct {
	Name string
}

//建立socket 连接
func (this *OnlineUser) PullFromClient() {
	for {
		var content string
		err := websocket.Message.Receive(this.Connection, &content)
		// If user closes or refreshes the browser, a err will occur
		if err != nil {
			return
		}

		m := Message{
			MType:       TEXT_MTYPE,
			UserInfo:    this.UserInfo,
			Time:        CreatedAt(),
			Content:     content,
			OnlineUsers: make(map[string]string),
		}

		//客户端发送一条信息， 格式化后， 写入到Broadcast  , 这个ActionRoom 里面的一条公共消息池子 
		ActiveChannel.Rooms[this.InRoom].Broadcast <- m
		//this.InRoom.Broadcast <- m
	}
}

// 向客户端发送信息
func (this *OnlineUser) PushToClient() {
	for b := range this.Send {
		err := websocket.JSON.Send(this.Connection, b)
		if err != nil {
			break
		}
	}
}

//关闭用户连接
func (this *OnlineUser) killUserResource() {
	Log.Printf(" %s is gone !\n", this.UserInfo.Name)
	this.Connection.Close()

	//离线消息
	m := Message{
		MType:       STATUS_MTYPE,
		UserInfo:    this.UserInfo,
		Time:        CreatedAt(),
		Content:     "用户【" + this.UserInfo.Name + "】离开了房间!",
		OnlineUsers: make(map[string]string),
	}
	close(this.Send)
	delete(ActiveChannel.Rooms[this.InRoom].OnlineUsers, this.UserInfo.Name)

	m.OnlineUsers = ActiveChannel.Rooms[this.InRoom].getOnlineUsers()
	ActiveChannel.Rooms[this.InRoom].Broadcast <- m

}

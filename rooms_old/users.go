package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
)

//在线用户类型
type OnlineUser struct {
	InRoom     *Room
	Connection *websocket.Conn
	UserInfo   *User
	Send       chan Message
}

//建立socket 连接
func (this *OnlineUser) PullFromClient() {
	fmt.Printf("%s PullFromClient !\n", humanCreatedAt())
	for {
		var content string
		err := websocket.Message.Receive(this.Connection, &content)
		// If user closes or refreshes the browser, a err will occur
		if err != nil {
			return
		}

		m := Message{
			MType: TEXT_MTYPE,
			TextMessage: TextMessage{
				UserInfo: this.UserInfo,
				Time:     humanCreatedAt(),
				Content:  content,
			},
		}
		//客户端发送一条信息， 格式化后， 写入到Broadcast  , 这个ActionRoom 里面的一条公共消息池子 
		this.InRoom.Broadcast <- m
	}
}

// 向客户端发送信息
func (this *OnlineUser) PushToClient() {
	fmt.Printf("%s PushToClient\n ", humanCreatedAt())
	for b := range this.Send {
		err := websocket.JSON.Send(this.Connection, b)
		if err != nil {
			break
		}
	}
}

//关闭用户连接
func (this *OnlineUser) killUserResource() {
	fmt.Printf("%s  close !\n", humanCreatedAt())
	this.Connection.Close()
	roomId := this.UserInfo.Token
	delete(this.InRoom.OnlineUsers, this.UserInfo.Token)
	close(this.Send)

	m := Message{
		MType:      STATUS_MTYPE,
		UserStatus: UserStatus{
		//Users: run.GetOnlineUsers(),
		},
	}
	runingRooms[roomId].Broadcast <- m

}

// 用户属性
type User struct {
	Id     int64  //数据库绑定id 
	Name   string //用户名称
	Face   string //头像地址
	Token  string // 交互标识
	RoomId string // 房间号

}

//每个用户的登陆状态
type UserStatus struct {
	Users []*User
}

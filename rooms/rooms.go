package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/room/add", AddRoom)
	http.Handle("/room/i", websocket.Handler(BuildRoomSocket))
}

// 动态添加room 添加room 的办法
func AddRoom(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	roomName := r.Form.Get("rn")          // 房间名称
	mapId := r.Form.Get("mapid")          // 地图 
	userIp := r.Header.Get("REMOTE_ADDR") //ip 地址

	Name := []byte(userIp + roomName)
	str := string(md5Encode(Name))

	//创建房间
	Room := RoomBase{
		Name:        roomName,
		MapId:       mapId,
		Broadcast:   make(chan Message),
		RoomClose:   make(chan bool),
		OnlineUsers: make(map[string]*OnlineUser),
		Token:       str,
	}
	go Room.Run()

	if _, ok := ActiveChannel.Rooms[str]; ok == false {
		ActiveChannel.Rooms[str] = &Room
		fmt.Fprintf(w, `{"token":"%s" , "error":0}`, str)
	} else {
		fmt.Fprintf(w, `{"token":"token exists" , "error":1}`)
	}

}

//room 创建消息 结构
type Message struct {
	MType       string
	Content     string
	UserInfo    User
	Time        string
	OnlineUsers map[string]string //在线用户id 
}

// roombase 信息。 
type RoomBase struct {
	Name        string                 //room 名称
	MapId       string                 // 使用地图
	Broadcast   chan Message           // 消息队列
	RoomClose   chan bool              // 关闭信号
	Token       string                 // 房间id 
	OnlineUsers map[string]*OnlineUser //在线用户
}

func (this *RoomBase) Run() {
	for {
		select {
		case b := <-this.Broadcast:
			// 取出消息， 告诉room 内所有人。 
			for _, v := range this.OnlineUsers {
				v.Send <- b
			}
		case rc := <-this.RoomClose:
			if rc == true {
				close(this.RoomClose)
				close(this.Broadcast)
				delete(ActiveChannel.Rooms, this.Token)
				return
			}
		}
	}
}

func (this *RoomBase) getOnlineUsers() map[string]string {
	users := make(map[string]string)
	for k, v := range this.OnlineUsers {
		users[k] = v.UserInfo.Name
	}
	return users
}

//用户类型
type OnlineUser struct {
	Send       chan Message
	InRoom     string
	Connection *websocket.Conn
	UserInfo   User
}

type User struct {
	Name string
}

// 建立客户端连接
func BuildRoomSocket(ws *websocket.Conn) {
	room := ws.Request().URL.Query().Get("room")
	user := ws.Request().URL.Query().Get("user")

	if room == "" {
		return
	}
	if _, ok := ActiveChannel.Rooms[room]; ok {

		onlinUser := OnlineUser{
			Send:       make(chan Message),
			InRoom:     room,
			Connection: ws,
			UserInfo:   User{Name: user},
		}
		ActiveChannel.Rooms[room].OnlineUsers[user] = &onlinUser

		m := Message{
			MType:       STATUS_MTYPE,
			UserInfo:    onlinUser.UserInfo,
			Time:        CreatedAt(),
			Content:     "用户[" + user + "]进入了房间!",
			OnlineUsers: make(map[string]string),
		}
		m.OnlineUsers = ActiveChannel.Rooms[room].getOnlineUsers()

		// 消息提示
		ActiveChannel.Rooms[room].Broadcast <- m

		go onlinUser.PushToClient()
		onlinUser.PullFromClient()
		//当用户关闭socket 连接时， onlinUser.PullFromClient 将会因用户断开向下执行。 
		onlinUser.killUserResource()
	} else {
		fmt.Println("On BuildRoomSocket: room[" + room + "] not exists!")
	}
}

//建立socket 连接
func (this *OnlineUser) PullFromClient() {
	fmt.Printf("%s PullFromClient !\n", CreatedAt())
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
	fmt.Printf("%s PushToClient\n ", CreatedAt())
	for b := range this.Send {
		err := websocket.JSON.Send(this.Connection, b)
		if err != nil {
			break
		}
	}
}

//关闭用户连接
func (this *OnlineUser) killUserResource() {
	fmt.Printf("%s  close !\n", CreatedAt())
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

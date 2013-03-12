package main

/*
import (
	//"code.google.com/p/go.net/websocket"
	//"fmt"
	"html/template"
	"net/http"
	//"strings"
	//"time"
)


const (
	TEXT_MTYPE   = "text_mtype"
	STATUS_MTYPE = "status_mtype"
	TIME_FORMAT  = "01-02 15:04:05"
)

// 消息环境
var runningActiveRoom *ActiveRoom = &ActiveRoom{}

//聊天室环境
type ActiveRoom struct {
	CloseSign   chan bool
	Broadcast   chan Message
	OnlineUsers map[string]*OnlineUser
}

//消息类型
type Message struct {
	MType       string
	TextMessage TextMessage
	UserStatus  UserStatus
}

//在线用户类型
type OnlineUser struct {
	InRoom     *ActiveRoom
	Connection *websocket.Conn
	UserInfo   *User
	Send       chan Message
}

//用户类型
type User struct {
	Name     string
	Email    string
	Gravatar string
}
type UserStatus struct {
	Users []*User
}

//文本消息类型
type TextMessage struct {
	Content  string
	UserInfo *User
	Time     string
}

//核心消息处理线程， 这会是一条独立的线程， 任何被放入公共消息池子的内容， 将会被循环写入到onlineUser.Send 当中。
func (this *ActiveRoom) Run() {
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

// 聊天室监控处理 , 用户奋发各项数据 
func InitChatRoom() {
	runningActiveRoom = &ActiveRoom{
		OnlineUsers: make(map[string]*OnlineUser),
		Broadcast:   make(chan Message),
		CloseSign:   make(chan bool),
	}
	go runningActiveRoom.Run()
}

// 用户登陆处理， 建立socket 连接
func BuildConnection(ws *websocket.Conn) {
	email := ws.Request().URL.Query().Get("email")
	fmt.Printf("%s connect web socket!\n", humanCreatedAt())
	if email == "" {
		return
	}

	// 创建一个user 对象
	onlineUser := &OnlineUser{
		InRoom:     runningActiveRoom,
		Connection: ws,
		Send:       make(chan Message, 256), //存放256条消息
		UserInfo: &User{
			Email:    email,
			Name:     strings.Split(email, "@")[0],
			Gravatar: email,
		},
	}

	//向actionRoom 中， 追加在线用户信息
	runningActiveRoom.OnlineUsers[email] = onlineUser

	//写入一条消息队列
	m := Message{
		MType: STATUS_MTYPE,
		UserStatus: UserStatus{
			Users: runningActiveRoom.GetOnlineUsers(),
		},
	}
	runningActiveRoom.Broadcast <- m

	//推送消息到客户端， 这在一个新的线程中操作， 避免主线的消息获取受阻拦。  
	go onlineUser.PushToClient()

	// 获取客户端消息 , PullFromClient 是一个循环锁，一只等待用户自己的信息反馈， 当用户关闭url 的时候， 这个锁将被退出 ，进入killUserResource 
	onlineUser.PullFromClient()
	// 关闭连接
	onlineUser.killUserResource()
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

	delete(this.InRoom.OnlineUsers, this.UserInfo.Email)
	close(this.Send)

	m := Message{
		MType: STATUS_MTYPE,
		UserStatus: UserStatus{
			Users: runningActiveRoom.GetOnlineUsers(),
		},
	}
	runningActiveRoom.Broadcast <- m
}

//获取在线用户列表
func (this *ActiveRoom) GetOnlineUsers() (users []*User) {
	for _, online := range this.OnlineUsers {
		users = append(users, online.UserInfo)
	}
	return users
}

func humanCreatedAt() string {
	return time.Now().Format(TIME_FORMAT)
}

func ChatTest(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	t := template.Must(template.ParseFiles("tpl/boojie_web/footer.html", "tpl/boojie_web/header.html", "tpl/boojie_web/test.html"))

	type TplInfo struct {
		Email string
	}

	info := &TplInfo{email}
	t.ExecuteTemplate(w, "content", info)
}

func init() {
	//InitChatRoom()
	//http.Handle("/chat/serv", websocket.Handler(BuildConnection))
	http.HandleFunc("/chat", ChatTest)
}
*/

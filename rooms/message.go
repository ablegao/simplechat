package main

//创建消息 结构
type Message struct {
	MType       string
	Content     string
	UserInfo    User
	Time        string
	OnlineUsers map[string]string //在线用户id 
}

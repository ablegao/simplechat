package main

import (
	"time"
)

// 定义消息类型，和时间格式
const (
	TEXT_MTYPE   = "text"   //标准文字消息
	STATUS_MTYPE = "status" // 状态消息
	GAME_MTTYPE  = "game"   // 游戏指令消息
	TIME_FORMAT  = "01-02 15:04:05"
)

//一条消息内容的组成
type Message struct {
	MType       string
	TextMessage TextMessage
	UserStatus  UserStatus
}

//文本消息结构
type TextMessage struct {
	Content  string
	UserInfo *User
	Time     string
}

func humanCreatedAt() string {
	return time.Now().Format(TIME_FORMAT)
}

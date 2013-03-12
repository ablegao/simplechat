package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type HTMLData struct {
	DataRows map[string]string
}

func init() {
	fmt.Println("insert test website")

	http.HandleFunc("/test/room/i", WEBRoomHome)
	http.HandleFunc("/test/home", WEBCreateRoom)
}

// 这个页面用来存放一些测试用html5界面
func WEBCreateRoom(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("tmpl/header.html", "tmpl/createtoom.html", "tmpl/footer.html"))

	Datas := HTMLData{DataRows: make(map[string]string)}
	for k, v := range ActiveChannel.Rooms {
		Datas.DataRows[k] = v.Name
	}

	tpl.ExecuteTemplate(w, "content", Datas)
}

//room 主页 ， 游戏为开始时
func WEBRoomHome(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("tmpl/header.html", "tmpl/room.html", "tmpl/footer.html"))

	type Data struct {
		Room string
		User string
	}
	tmplData := Data{
		Room: r.URL.Query().Get("room"),
		User: r.URL.Query().Get("user"),
	}
	tpl.ExecuteTemplate(w, "content", tmplData)
}

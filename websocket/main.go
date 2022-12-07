package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"time"
)

func ping(c *gin.Context) {
	// 升级get请求为webSocket请求
	upgrader := websocket.Upgrader{
		HandshakeTimeout: 5 * time.Second,
	}
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			return
		}
	}(ws)
	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		if string(message) == "ping" {
			message = []byte("pong")
		}
		// 写入ws数据
		err = ws.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}
}

func main() {
	//bindAddress := "localhost:2303"
	//r := gin.Default()
	//r.LoadHTMLGlob("/Users/liaobowen/src/go-stu/websocket/*")
	////监听 get请求  /test路径
	//r.GET("/test", func(c *gin.Context) {
	//	c.HTML(http.StatusOK, "test.html", nil)
	//})
	//r.GET("/ping", ping)
	//r.Run(bindAddress)
	var s fmt.Stringer
	s = "s"
	fmt.Println(s)
}

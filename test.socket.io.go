package main

import (
	"fmt"
	"os"
	socket_io "socket.io.client/websocket"
)

func main() {
	c := make(chan os.Signal)

	//连接配置
	options := socket_io.SocketIOOptions{
		Path:      "/socket.io",
		Transport: socket_io.WebSocket,
		EIO:       socket_io.V4,
	}

	//创建 socket.io
	var socketIoClient = socket_io.NewSocketIOClient("http://localhost:8000", reconnect, callback, options)

	//打开连接
	err := socketIoClient.Dail()

	if err != nil {
		println(err.Error())
	}

	<-c
}

func reconnect(e error) {
	fmt.Println(e)
}

func callback(eventName string, data string) {
	fmt.Println(eventName, data)
}

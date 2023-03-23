# Socket.IO Client Golang 

[Socket.IO 官方文档](https://pkg.go.dev/github.com/gorilla/websocket?tab=doc)

---

 **Socket.IO Client 在官方 Golang 版本中并没有被支持，对于项目逼不得已需要用到 Sokcet.IO Client 而且必须需要使用Golang的朋友可以使用改项目来连接 Socket.IO Server**

 ---

**因为这个项目只是简单的实现的连接和消息通讯，但实际上没有官方其他开发语言中提供的库多全，如果你对 Socket.io Client 有更多的需求你可以在改项目进行改进**

 ---

### 使用方法
```go
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
```
 

## 项目简介

**开发者：JD.PENG**

![touxiang](https://avatars.githubusercontent.com/u/41544772?s=400&u=a27b4260f60bf261afd39c0c62f0dd0876c22ef3&v=4)


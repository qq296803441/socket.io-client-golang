package socket_io

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
	"regexp"
	"socket.io.client/websocket/openedMessage"
	"strconv"
	"time"
)

const (
	Polling   = "polling"
	WebSocket = "websocket"
)

const (
	V3 = 3
	V4 = 4
)

type SocketIOOptions struct {
	Path      string
	Transport string
	Header    map[string]string
	EIO       int
}

type SocketIOClient struct {
	Conn        *websocket.Conn
	Addr        string
	SendMsgChan chan string
	recvMsgChan chan string
	Options     SocketIOOptions
	IsAlive     bool
	reconnect   func(e error)
	callback    func(eventName string, data string)
}

//NewSocketIOClient 创建 SocketIO
func NewSocketIOClient(
	addr string,
	reconnect func(e error),
	callback func(eventName string, data string),
	options SocketIOOptions,
) *SocketIOClient {
	var conn *websocket.Conn
	return &SocketIOClient{
		Addr:        addr,
		Conn:        conn,
		SendMsgChan: make(chan string, 10),
		recvMsgChan: make(chan string, 10),
		IsAlive:     false,
		Options:     options,
		reconnect:   reconnect,
		callback:    callback,
	}
}

//SendMessage 发送消息
func (s *SocketIOClient) SendMessage(event string, message interface{}) error {
	sendMessage := make([]interface{}, 2)
	sendMessage[0] = event
	sendMessage[1] = message

	sJson, err := json.Marshal(sendMessage)

	if err != nil {
		return err
	}
	socketIoMessage := fmt.Sprintf("%v%v", openedMessage.EventMessage, string(sJson))
	s.SendMsgChan <- socketIoMessage

	return nil
}

//sendPing 发送ping 保持在线
func (s *SocketIOClient) sendPing() {

	go func() {
		for {
			if s.Conn == nil {
				return
			}
			socketIoMessage := fmt.Sprintf("%v", openedMessage.Ping)
			s.SendMsgChan <- socketIoMessage

			time.Sleep(1 * time.Second)
		}
	}()

}

//Dail 连接
func (s *SocketIOClient) Dail() error {
	var err error
	_he := http.Header{}

	for key, item := range s.Options.Header {
		_he.Add(key, item)
	}
	//构建 websocket 连接
	wsUrl, err := s.getServerUri()

	if err != nil {
		return err
	}

	s.Conn, _, err = websocket.DefaultDialer.Dial(*wsUrl, _he)
	if err != nil {
		return err
	}
	s.IsAlive = true

	s.dataCallBack()
	s.readMsgThread()
	s.sendMsgThread()
	//s.sendPing()
	return nil
}

//getServerUri GetServerUri 构建 websocket 连接
func (s *SocketIOClient) getServerUri() (*string, error) {

	url, err := url.Parse(s.Addr)
	if err != nil {
		return nil, err
	}

	var wsUrl = ""

	var wsBool = s.Options.Transport == WebSocket

	if url.Scheme == "https" || url.Scheme == "wss" {
		if wsBool {
			wsUrl += "wss://"
		} else {
			wsUrl += "https://"
		}
	} else if url.Scheme == "http" || url.Scheme == "ws" {
		if wsBool {
			wsUrl += "ws://"
		} else {
			wsUrl += "http://"
		}
	} else {
		return nil, errors.New("")
	}

	wsUrl += url.Host

	if s.Options.Path == "" {
		wsUrl += "/socket.io"
	} else {
		wsUrl += s.Options.Path
	}

	if url.RawQuery != "" {
		wsUrl += fmt.Sprintf("/?%v&", url.RawQuery)
	} else {
		wsUrl += "/?"
	}

	wsUrl += fmt.Sprintf("EIO=%v&transport=%v", s.Options.EIO, s.Options.Transport)

	return &wsUrl, nil
}

//sendMsgThread 发送消息
func (s *SocketIOClient) sendMsgThread() {
	go func() {
		for {
			msg := <-s.SendMsgChan
			_logger1 := "\t\n================ ws send ==================="
			_logger2 := "%v\t\n ws : %v' \t\n data : %v %v"
			_logger3 := "\t\n============================================"
			fmt.Println(fmt.Sprintf(_logger2, _logger1, s.Addr, msg, _logger3))
			err := s.Conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				fmt.Println(fmt.Sprintf("%v write error : %v", s.Addr, err.Error()))
				continue
			}
		}
	}()
}

//readMsgThread 读取消息
func (s *SocketIOClient) readMsgThread() {
	go func() {
		for {
			if s.Conn != nil {
				_, message, err := s.Conn.ReadMessage()
				if err != nil {
					fmt.Println(fmt.Sprintf("%v 连接异常断开 : %v", s.Addr, err.Error()))
					s.IsAlive = false
					if s.reconnect != nil {
						s.reconnect(err)
					}
					return
				}

				// 需要读取数据，不然会阻塞
				s.recvMsgChan <- string(message)
				_logger1 := "\t\n============== ws callback ================="
				_logger2 := "%v\t\n ws : %v' \t\n data : %v %v"
				_logger3 := "\t\n============================================"
				fmt.Println(fmt.Sprintf(_logger2, _logger1, s.Addr, string(message), _logger3))
			}
		}
	}()
}

//dataCallBack 回调
func (s *SocketIOClient) dataCallBack() {
	go func() {
		for {
			msg := <-s.recvMsgChan
			if s.callback != nil {
				go func() {
					s.analysisMessage(msg)
				}()
			}
		}
	}()
}

//getMessageTypeValue 获取消息类型
func (s *SocketIOClient) getMessageTypeValue(mes string) string {
	return regexp.MustCompile(`\d+`).FindString(mes)
}

//analysisMessage 解析数据类型
func (s *SocketIOClient) analysisMessage(message string) {

	messageType := s.getMessageTypeValue(message)

	messageTypeValue, err := strconv.Atoi(messageType)
	if err != nil {
		fmt.Println(fmt.Sprintf("socketIO 一条数据没有被正常解析，请注意, strconv.Atoi(messageType) message : %v", message))
	}

	switch messageTypeValue {
	case openedMessage.Opened:
		//0
		s.callback("opened", message[1:])
	case openedMessage.Ping:
		break
	case openedMessage.Pong:
		//目前没有需求需要pong结果,可以过滤
		break
	case openedMessage.Connected:

		break
	case openedMessage.Disconnected:

		break
	case openedMessage.EventMessage:
		//42
		var newMessage = message[2:]
		var eventMessage = make([]interface{}, 0)
		err := json.Unmarshal([]byte(newMessage), &eventMessage)
		if err != nil {
			fmt.Println(fmt.Sprintf("analysisMessage error : %v  message : %v", err.Error(), message))
		}

		sJson, err := json.Marshal(eventMessage[1])
		if err != nil {
			fmt.Println(fmt.Sprintf("analysisMessage error : %v  message : %v", err.Error(), message))
		}
		s.callback(eventMessage[0].(string), string(sJson))
		break
	case openedMessage.AckMessage:

		break
	case openedMessage.ErrorMessage:

		break
	case openedMessage.BinaryMessage:

		break
	case openedMessage.BinaryAckMessage:

		break
	default:
		fmt.Println(fmt.Sprintf("socketIO 一条数据没有被正常解析，请注意, message : %v", message))
		break
	}
}

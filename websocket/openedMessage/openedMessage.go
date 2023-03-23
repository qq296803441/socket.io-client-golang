package openedMessage

type OpenedMessage struct {
	sid          string
	upgrades     []string
	pingInterval int
	pingTimeout  int
	Type         int
}

// 消息类型
const (
	Opened           = 0
	Ping             = 2
	Pong             = 3
	Connected        = 40
	Disconnected     = 41
	EventMessage     = 42
	AckMessage       = 43
	ErrorMessage     = 44
	BinaryMessage    = 45
	BinaryAckMessage = 46
)

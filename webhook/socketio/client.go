package socketio

import socketio "github.com/googollee/go-socket.io"

type SocketIOClient struct {
	Id   string
	Conn socketio.Conn
}

package socketio

import (
	"fmt"
	"log"

	s "github.com/googollee/go-socket.io"
	socketio "github.com/googollee/go-socket.io"
)

func NewSocketIoServer(manager *SocketIOConnectionManager) *s.Server {
	server := s.NewServer(nil)

	server.OnConnect("/", func(s s.Conn) error {
		id := s.RemoteHeader().Get("Sec-WebSocket-Protocol")
		// ctx := context.WithValue(context.Background(), "id", id)
		// s.SetContext(ctx)
		var client = SocketIOClient{
			Id:   id,
			Conn: s,
		}
		manager.Register <- &client
		return nil
	})

	server.OnError("/", func(s s.Conn, e error) {
		id := s.RemoteHeader().Get("Sec-WebSocket-Protocol")
		manager.Unregister <- manager.Clients[id]
		log.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s s.Conn, reason string) {
		id := s.RemoteHeader().Get("Sec-WebSocket-Protocol")
		manager.Unregister <- manager.Clients[id]
		log.Println("closed", reason)
	})

	server.OnEvent("/", "msg", func(s socketio.Conn, msg string) string {
		fmt.Println("socket io receive a new message ", msg)
		return "recv " + msg
	})

	return server
}

// server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
// 	log.Println("notice:", msg)
// 	s.Emit("reply", "have "+msg)
// })

// server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
// 	s.SetContext(msg)
// 	return "recv " + msg
// })

// server.OnEvent("/", "bye", func(s socketio.Conn) string {
// 	last := s.Context().(string)
// 	s.Emit("bye", last)
// 	s.Close()
// 	return last
// })

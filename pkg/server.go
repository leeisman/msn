package pkg

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	Message chan string
}

func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
}

func (s *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
	}
	fmt.Println(">>>>>>聊天室服務器啟動成功...")
	defer listener.Close()

	// start listen message
	go s.ListenMessage()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
		}

		// do handler
		go s.Handler(conn)
	}

}

func (s *Server) ListenMessage() {
	for {
		msg := <-s.Message
		s.mapLock.Lock()
		for _, cli := range s.OnlineMap {
			cli.c <- msg
		}
		s.mapLock.Unlock()
	}
}

func (s *Server) Broadcast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	s.Message <- sendMsg
}

func (s *Server) Handler(conn net.Conn) {
	// join onlineMap
	user := NewUSer(conn, s)
	user.Online()

	isLive := make(chan struct{})

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil {
				fmt.Println("Conn Read err:", err)
				return
			}

			msg := string(buf[:n-1])
			user.DoMessage(msg)

			isLive <- struct{}{}
		}
	}()

	for {
		select {
		case <-isLive:
			// reset timer
		case <-time.After(time.Second * 10000):
			// timeout
			user.SendMsg("你被踢了")
			close(user.c)
			conn.Close()
			return
		}
	}
}

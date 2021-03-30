package pkg

import (
	"fmt"
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	c    chan string
	conn net.Conn

	server *Server
}

func NewUSer(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		c:      make(chan string),
		conn:   conn,
		server: server,
	}
	go user.ListenMessage()
	return user
}

// To listen
func (u *User) ListenMessage() {
	for {
		msg, ok := <-u.c
		// will be close
		if !ok {
			fmt.Println("safety closed")
			return
		}
		_, err := u.conn.Write([]byte(msg + "\n"))
		if err != nil {
			fmt.Println("user write err", err)
		}
	}
}

func (u *User) Online() {
	u.server.mapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.mapLock.Unlock()

	u.server.Broadcast(u, "已上線\n")
}

func (u *User) Offline() {
	u.server.mapLock.Lock()
	delete(u.server.OnlineMap, u.Name)
	u.server.mapLock.Unlock()

	u.server.Broadcast(u, "下線\n")
}

func (u *User) SendMsg(msg string) {
	u.conn.Write([]byte(msg))
}

func (u *User) DoMessage(msg string) {
	if msg == "@who" {
		u.server.mapLock.Lock()
		for _, user := range u.server.OnlineMap {
			onlineMsg := "[" + user.Addr + ":" + user.Name + "在線...\n"
			u.SendMsg(onlineMsg)
		}
		u.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newName := strings.Split(msg, "|")[1]

		if _, ok := u.server.OnlineMap[newName]; ok {
			u.SendMsg("當前用戶名已被使用" + "\n")
		} else {
			u.server.mapLock.Lock()
			delete(u.server.OnlineMap, u.Name)
			u.server.OnlineMap[newName] = u
			u.server.mapLock.Unlock()

			u.Name = newName
			u.SendMsg("已經更新用戶名:" + u.Name + "\n")
		}

	} else if len(msg) > 4 && msg[:3] == "to|" {
		// msg format: to|name|message
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			u.SendMsg("消息格式不正確，請使用\"to|name|msg|\"\n")
			return
		}
		// get user
		remoteUser, ok := u.server.OnlineMap[remoteName]
		if !ok {
			u.SendMsg("無此用戶，請重發\n")
			return
		}

		content := strings.Split(msg, "|")[2]
		if content == "" {
			u.SendMsg("無消息內容,請重發\n")
			return
		}
		remoteUser.SendMsg(u.Name + "對您說:" + content + "\n")
	} else {
		u.server.Broadcast(u, msg)
	}
}

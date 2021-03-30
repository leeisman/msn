package pkg

import (
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	serverIp   string
	serverPort int
	conn       net.Conn
	flag       int
	Name       string
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		serverIp:   serverIp,
		serverPort: serverPort,
		flag:       999,
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial err:", err)
	}
	client.conn = conn
	return client
}

func (c *Client) menu() bool {
	flag := 999
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用路名")
	fmt.Println("0.退出")
	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 3 {
		c.flag = flag
		return true
	}
	fmt.Println("不合法輸入")
	return false
}

func (c *Client) Run() {
	for c.flag != 0 {
		for c.menu() == true {
			switch c.flag {
			case 1:
				// 公聊
				fmt.Println("公聊模式選擇...")
				c.PublicChat()
				break
			case 2:
				// 私聊
				fmt.Println("私聊模式選擇...")
				c.PrivateChat()
				break
			case 3:
				// 更新用戶名
				fmt.Println("更新用戶名選擇...")
				c.UpdateName()
			case 0:
				fmt.Println("離開...")
				return
			}
		}
	}
}

func (c *Client) UpdateName() bool {
	fmt.Println(">>>>>>>請輸入用戶名:")
	fmt.Scanln(&c.Name)
	sendMsg := "rename|" + c.Name + "\n"
	_, err := c.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}
	return true
}

func (c *Client) DoResponse() {
	// 永久阻塞打印
	io.Copy(os.Stdout, c.conn)
}

func (c *Client) PublicChat() {
	fmt.Println(">>>>>請輸入聊天內容, exit退出:")
	var chatMsg string
	fmt.Scanln(&chatMsg)
	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := c.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn write err:", err)
				break
			}
		}
		chatMsg = ""
		fmt.Println(">>>>>請輸入聊天內容, exit退出:")
		fmt.Scanln(&chatMsg)
	}
}

func (c *Client) SelectUsers() {
	sendMsg := "@who\n"
	_, err := c.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn write err:", err)
		return
	}
}

func (c *Client) PrivateChat() {
	var remoteUser string
	var charMsg string
	c.SelectUsers()
	fmt.Println(">>>>請輸入聊天對用(用戶名),exit退出:")
	fmt.Scanln(&remoteUser)

	for remoteUser != "exit" {
		fmt.Println(">>>>請輸入私聊內容,exit退出:")
		fmt.Scanln(&charMsg)
		for charMsg != "exit" {
			if len(charMsg) != 0 {
				sendMsg := "to|" + remoteUser + "|" + charMsg + "\n"
				_, err := c.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn write err:", err)
					return
				}
				charMsg = ""
				fmt.Println(">>>>請輸入私聊內容,exit退出:")
				fmt.Scanln(&charMsg)
			}
		}
		remoteUser = ""
		fmt.Println(">>>>請輸入聊天對用(用戶名),exit退出:")
		fmt.Scanln(&remoteUser)
	}
}

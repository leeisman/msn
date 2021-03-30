package main

import (
	"flag"
	"fmt"
	"msn/pkg"
)

var serverIp string
var serverPort int
var serverOrClient string

func main() {
	flag.Parse()
	if serverOrClient == "server" {
		server := pkg.NewServer(serverIp, serverPort)
		server.Start()
	}
	if serverOrClient == "client" {
		client := pkg.NewClient(serverIp, serverPort)
		if client == nil {
			fmt.Println(">>>>>>連線服務器失敗...")
			return
		}
		fmt.Println(">>>>>>連線服務器成功...")
		go client.DoResponse()
		client.Run()
	}
}

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "設址服務器IP地址(默認127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "設址服務器端口(默認8888)")
	flag.StringVar(&serverOrClient, "s", "server", "預設server")
}

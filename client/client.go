package client

import (
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Conn       net.Conn
}

func NewClient(ip string, port int) *Client {
	client := &Client{ServerIp: ip, ServerPort: port}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", client.ServerIp, client.ServerPort))
	if err != nil {
		fmt.Println("net Dial err:", err)
		return nil
	}
	client.Conn = conn
	return client
}

func (this *Client) Start() {
	go this.DealResponse()
	fmt.Println("请输入需要发布的内容，输入exit退出")
	for {
		msg := ""
		fmt.Scanln(&msg)
		if msg == "exit" {
			return
		}
		if len(msg) != 0 {
			msg += "\n"
			_, err := this.Conn.Write([]byte(msg))
			if err != nil {
				fmt.Println("输入错误")
				return
			}
		}
	}
}
func (this Client) DealResponse() {
	io.Copy(os.Stdout, this.Conn)
}

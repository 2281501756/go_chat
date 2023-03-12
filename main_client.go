package main

import (
	"fmt"
	"go_chat/client"
)

func main() {
	client := client.NewClient("127.0.0.1", 8000)
	if client == nil {
		fmt.Println("链接失败")
		return
	}
	fmt.Println("链接成功")
	client.Start()
}

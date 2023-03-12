package server

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip        string
	Port      int
	OnlineMap map[string]*User
	mapLock   sync.Mutex
	message   chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		message:   make(chan string),
	}
	return server
}

func (this *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.listen err:", err)
		return
	}
	fmt.Println("服务器运行在", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	defer listener.Close()

	go this.ListenMessage()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err", err)
			return
		}
		go this.Handle(conn)
	}
}
func (this *Server) Handle(conn net.Conn) {
	user := NewUser(conn)
	fmt.Println(user.name, "上线")
	this.mapLock.Lock()
	this.OnlineMap[user.name] = user
	this.mapLock.Unlock()
	this.BroadCast(user, "上线")
	go this.ListenUserInputData(user)

}

func (this *Server) BroadCast(user *User, message string) {
	msg := "[" + user.name + "]:" + message
	this.message <- msg
}
func (this *Server) ListenMessage() {
	for {
		msg := <-this.message
		this.mapLock.Lock()
		for _, user := range this.OnlineMap {
			user.C <- msg
		}
		this.mapLock.Unlock()
	}
}
func (this *Server) ListenUserInputData(user *User) {
	but := make([]byte, 4096)
	for {
		n, err := user.conn.Read(but)
		if n == 0 {
			this.BroadCast(user, "下线了")
			return
		}
		if err != nil && err != io.EOF {
			fmt.Println("conn read err:", err)
			return
		}
		msg := string(but[:n-1])
		this.BroadCast(user, msg)

	}
}

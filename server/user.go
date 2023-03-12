package server

import (
	"fmt"
	"io"
	"net"
)

type User struct {
	add    string
	name   string
	C      chan string
	conn   net.Conn
	server *Server
}

func NewUser(conn net.Conn, server *Server) (res *User) {
	res = &User{
		add:    conn.RemoteAddr().String(),
		name:   conn.RemoteAddr().String(),
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	go res.ListenMessage()
	return res
}

func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}

func (this *User) Online() {
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.name] = this
	this.server.mapLock.Unlock()
	this.server.BroadCast(this, "上线")
	go this.ListenUserInputData()
}
func (this *User) Offline() {
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.name)
	this.server.mapLock.Unlock()
	this.server.BroadCast(this, "下线了")
}

func (this *User) ListenUserInputData() {
	but := make([]byte, 4096)
	for {
		n, err := this.conn.Read(but)
		if n == 0 {
			this.Offline()
			return
		}
		if err != nil && err != io.EOF {
			fmt.Println("conn read err:", err)
			return
		}
		msg := string(but[:n-1])
		this.server.BroadCast(this, msg)
	}
}

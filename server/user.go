package server

import "net"

type User struct {
	add  string
	name string
	C    chan string
	conn net.Conn
}

func NewUser(conn net.Conn) (res *User) {
	res = &User{
		add:  conn.RemoteAddr().String(),
		name: conn.RemoteAddr().String(),
		C:    make(chan string),
		conn: conn,
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

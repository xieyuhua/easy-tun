package main

import (
	"github.com/fatih/color"
	"io"
	"net"
	"flag"
)

var clients = make([]net.Conn, 0)

var inSer = flag.String("ip", ":8006", "server address")

func main() {
    flag.Parse()
    
	listener, err := net.Listen("tcp", *inSer)
	if err != nil {
		color.Red(err.Error())
		return
	}
	color.Cyan("server address" + *inSer)
	for {
		// 建立连接
		conn, err := listener.Accept()
		if err != nil {
			color.Red(err.Error())
			break
		}
		// 将客户端进行存储
		clients = append(clients, conn)
		color.Cyan("Accept connections from clients %s", conn.RemoteAddr().String())
		// 将客户端的请求，除了自生排外转发至所有的客户端，如果客户端是目的地IP将会进行响应
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	buff := make([]byte, 65542)
	for {
		nr, err := conn.Read(buff)
		if err != nil {
			if err != io.EOF {
				color.Red(err.Error())
			}
			break
		}
		color.Red("server read data")
        // 	sync.Lock()	 var sync sync.Mutex 	sync.Unlock()
		// 广播 TCPConnections =  make(map[string]net.Conn)   var  TCPConnections map[string]net.Conn 
		for _, c := range clients {
		    //将客户端的请求，除了自生排外转发至所有的客户端，如果客户端是目的地IP将会进行响应
			if c.RemoteAddr().String() != conn.RemoteAddr().String() {
				color.Red("server execute broadcast [server->%s]", c.RemoteAddr().String())
				c.Write(buff[:nr])
			}
		}
	}
}

package ipc

import (
	"fmt"
	"io"
	"log"
	"net"
)

type IPC struct {
}

func New() *IPC {
	return &IPC{}
}

func (i IPC) Load() error {
	fmt.Println("IPC module has been loaded")
	return nil
}

func (i IPC) Run() error {
	go i.tcpListen()

	fmt.Println("IPC module has been run")
	return nil
}

func (i IPC) tcpListen() {
	// Listen on TCP port 2000 on all available unicast and
	// anycast IP addresses of the local system.
	l, err := net.Listen("tcp", ":2000")
	if err != nil {
		log.Printf("Couldn't listen to tcp port: ", err)
		return
	}
	log.Printf("Listening on 2000...")
	defer l.Close()

	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			log.Printf("Error in accepting connection: ", err)
			return
		}
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go func(c net.Conn) {
			// Echo all incoming data.
			io.Copy(c, c)
			// Shut down the connection.
			c.Close()
		}(conn)
	}
}

package main

import (
	"fmt"
	"log"
	"net"
)

type Message struct {
	from    string
	payload []byte
}

type Server struct {
	listernAddr string
	ln          net.Listener
	quitch      chan struct{}
	msgch       chan Message
}

func NewServer(listernAddr string) *Server {
	return &Server{
		listernAddr: listernAddr,
		quitch:      make(chan struct{}),
		msgch:       make(chan Message, 10),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listernAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	s.ln = ln

	go s.acceptLoop()

	<-s.quitch
	close(s.msgch)

	return nil
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("accept error", err)
			continue
		}

		fmt.Println("new connection to server :", conn.RemoteAddr())

		go s.readLoop(conn)
	}
}

func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read error", err)
			continue
		}

		//msg := buf[:n]
		s.msgch <- Message{
			from:    conn.RemoteAddr().String(),
			payload: buf[:n],
		}
		conn.Write([]byte("thank you for your message"))
		//fmt.Println(string(msg))
	}
}

func main() {
	server := NewServer(":3000")

	go func() {
		for msg := range server.msgch {
			fmt.Printf("message recevied from connection (%s):%s \n", msg.from, string(msg.payload))
		}
	}()

	log.Fatal(server.Start())
}

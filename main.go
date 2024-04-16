package main

import (
	"fmt"
	"log"

	"github.com/sounishnath003/dooms-distributed-filestore/transport"
)

func OnPeerFunc(p transport.Peer) error {
	p.Close()
	log.Printf("doing some processing logic with the peer outside TCPTransport")
	// return fmt.Errorf("main.go: failed to call OnPeer func")
	return nil
}

func main() {
	fmt.Println("dooms-distributed-filestore initiated...")

	tcpOpts := transport.TCPTransportOpts{
		ListenAddr:    ":3000",
		Decoder:       transport.DefaultDecoder{},
		HandShakeFunc: transport.NOPHandshakeFunc,
		OnPeer:        OnPeerFunc,
	}
	fmt.Print("TCPTransportOpts configuration=", tcpOpts)

	tr := transport.NewTCPTransport(tcpOpts)

	go func() {
		for {
			msg := <-tr.Consume()
			log.Printf("main.go: %+v\n", msg)
		}
	}()

	if err := tr.ListenAndAcceptConn(); err != nil {
		log.Fatal(err)
	}

	select {}
}

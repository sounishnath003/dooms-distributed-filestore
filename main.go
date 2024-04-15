package main

import (
	"fmt"
	"log"

	"github.com/sounishnath003/dooms-distributed-filestore/transport"
)

func main() {
	fmt.Println("dooms-distributed-filestore initiated...")

	tcpOpts := transport.TCPTransportOpts{
		ListenAddr:    ":3000",
		Decoder:       transport.DefaultDecoder{},
		HandShakeFunc: transport.NOPHandshakeFunc,
	}
	fmt.Print("TCPTransportOpts configuration=", tcpOpts)

	tr := transport.NewTCPTransport(tcpOpts)
	if err := tr.ListenAndAcceptConn(); err != nil {
		log.Fatal(err)
	}

	select {}
}

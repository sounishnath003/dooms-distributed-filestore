package main

import (
	"fmt"
	"log"

	"github.com/sounishnath003/dooms-distributed-filestore/transport"
)

func main() {
	fmt.Println("i am running dooms-distributed-filestore")
	tr := transport.NewTCPTransport(":3000")
	if err := tr.ListenAndAcceptConn(); err != nil {
		log.Fatal(err)
	}

	select {}
}

package transport

import (
	"encoding/gob"
	"io"
	"log"
)

type Encoder struct{}

type Decoder interface {
	Decode(io.Reader, *RPC) error
}

type GOBDecoder struct{}

func (dec GOBDecoder) Decode(r io.Reader, rpc *RPC) error {
	return gob.NewDecoder(r).Decode(rpc)
}

type DefaultDecoder struct{}

func (nopDec DefaultDecoder) Decode(r io.Reader, rpc *RPC) error {
	buf := make([]byte, 4096) // allocating 4KB
	n, err := r.Read(buf)
	if err != nil {
		return err
	}

	rpc.Payload = buf[:n]
	log.Printf("encoding.go: rpc call made: %+v, new message received: %s\n", rpc, string(buf[:n]))

	return nil

}

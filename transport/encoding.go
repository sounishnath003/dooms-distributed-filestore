package transport

import (
	"encoding/gob"
	"io"
	"log/slog"
)

type Encoder struct{}

type Decoder interface {
	Decode(io.Reader, *Message) error
}

type GOBDecoder struct{}

func (dec GOBDecoder) Decode(r io.Reader, msg *Message) error {
	return gob.NewDecoder(r).Decode(msg)
}

type DefaultDecoder struct{}

func (nopDec DefaultDecoder) Decode(r io.Reader, msg *Message) error {
	buf := make([]byte, 4096) // allocating 4KB
	n, err := r.Read(buf)
	if err != nil {
		return err
	}

	msg.Payload = buf[:n]
	slog.Info("new message received", slog.String("message", string(buf[:n])))

	return nil

}

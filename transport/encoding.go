package transport

import "io"

type Encoder struct{}

type Decoder interface {
	Decode(io.Reader, any) error
}

package transport

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	listerAddr := ":4000"
	tr := NewTCPTransport(TCPTransportOpts{
		ListenAddr:    ":4000",
		HandShakeFunc: NOPHandshakeFunc,
		Decoder:       GOBDecoder{},
	})

	assert.Equal(t, tr.ListenAddr, listerAddr)

	// Server
	assert.Nil(t, tr.ListenAndAcceptConn())
}

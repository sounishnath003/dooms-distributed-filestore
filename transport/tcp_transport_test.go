package transport

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	listerAddr := ":4000"
	tr := NewTCPTransport(listerAddr)

	assert.Equal(t, tr.ListenAddr, listerAddr)

	// Server
	assert.Nil(t, tr.ListenAndAcceptConn())
}

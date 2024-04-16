package transport

import "net"

// RPC represents any arbitary message data
// that is being sent over each transport
type RPC struct {
	From    net.Addr
	Payload []byte
}

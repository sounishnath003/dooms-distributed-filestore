package transport

import "errors"

// returns if Handshake b/w local or remote node could not be established
var ErrInvalidHandshakeMessage = errors.New("invalid handshake detected. local or remote nodes connection could not be established")

// ... ?
type HandShakeFunc func(Peer) error

// NOP handshake implementation
func NOPHandshakeFunc(Peer) error { return nil }

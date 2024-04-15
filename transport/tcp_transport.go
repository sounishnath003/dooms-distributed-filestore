package transport

import (
	"log/slog"
	"net"
	"sync"
)

// represents the remote node connection established over TCP transport
type TCPPeer struct {
	// conn is the underlying connection of the peer node
	conn net.Conn

	// TRUE: when DIAL and retrieve a conn
	// FALSE: when ACCEPT andretrieve a conn
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

type TCPTransportOpts struct {
	ListenAddr    string
	HandShakeFunc HandShakeFunc
	Decoder       Decoder
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
	}
}

func (t *TCPTransport) ListenAndAcceptConn() error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}
	slog.Info("listener ln=%+v ...\n", t.listener)
	// invoking the accep loop - infinite acceptance looping
	go t.startAcceptLoop()
	return nil
}

func (t *TCPTransport) startAcceptLoop() error {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			slog.Error("TCP: accept error has occured %s\n", err)
			return err
		}
		// handle the conn
		go t.handleConn(conn)
	}
}

type Temp struct{}

func (t *TCPTransport) handleConn(conn net.Conn) {
	slog.Info("new incoming connection - %v\n", conn)

	peer := NewTCPPeer(conn, true)
	if err := t.HandShakeFunc(peer); err != nil {
		// drop connection actually!!
		peer.conn.Close()
		slog.Error("TCP error has been occured during the handshaking")
		return
	}

	// Read loop
	msg := &Message{}
	for {
		if err := t.Decoder.Decode(conn, msg); err != nil {
			slog.Error("TCP error occured: %s\n", err)
			continue
		}
		slog.Info("new message has been received: %+v\n", slog.String("message", string(msg.Payload)))
	}

}

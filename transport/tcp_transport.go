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

type TCPTransport struct {
	ListenAddr string
	listener   net.Listener
	shakeHands HandShakeFunc // ACK utility handshakes
	decoder    Decoder

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPTransport(listenAddr string) *TCPTransport {
	return &TCPTransport{
		ListenAddr: listenAddr,
		shakeHands: NOPHandshakeFunc,
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

func (t *TCPTransport) handleConn(conn net.Conn) error {
	slog.Info("new incoming connection - %v\n", conn)

	peer := NewTCPPeer(conn, true)
	if err := t.shakeHands(peer); err != nil {
		// drop connection actually!!
		peer.conn.Close()
		slog.Error("TCP error has been occured during the handshaking")
		return ErrInvalidHandshakeMessage
	}

	buf := make([]byte, 4096) // Allocate a buffer of 4KB
	// Read loop
	for {
		if err := t.decoder.Decode(conn, buf); err != nil {
			slog.Error("TCP error - %s\n", err)
			continue
		}
	}

	return nil
}

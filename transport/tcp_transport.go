package transport

import (
	"log"
	"net"
	"sync"
)

type TCPTransportOpts struct {
	ListenAddr    string
	HandShakeFunc HandShakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
}

// represents the remote node connection established over TCP transport
type TCPPeer struct {
	// conn is the underlying connection of the peer node
	conn net.Conn

	// TRUE: when DIAL and retrieve a conn
	// FALSE: when ACCEPT andretrieve a conn
	outbound bool
}

// Close implements the Peer Interface
func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener
	rpcch    chan RPC

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcch:            make(chan RPC, 3),
	}
}

// Consume implements the Transform interface
// which will only return Read-only channel , for reading incoming msgs received from
// another peer or network nodes
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
}

func (t *TCPTransport) ListenAndAcceptConn() error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}
	log.Printf("listener ln=%+v ...\n", t.listener)
	// invoking the accep loop - infinite acceptance looping
	go t.startAcceptLoop()
	return nil
}

func (t *TCPTransport) startAcceptLoop() error {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			log.Println("TCP: accept error has occured", err)
			return err
		}
		// handle the conn
		go t.handleConn(conn)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn) {
	var err error
	defer func() {
		log.Println("dropping peer connection:", err)
		conn.Close()
	}()

	log.Printf("new incoming connection - %v\n", conn)

	peer := NewTCPPeer(conn, true)
	if err = t.HandShakeFunc(peer); err != nil {
		// drop connection actually!!
		peer.conn.Close()
		return
	}

	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			log.Println("TCP error has been occured during the handshaking")
			return
		}
	}

	// Read loop
	rpc := RPC{From: conn.RemoteAddr()}
	for {
		err := t.Decoder.Decode(conn, &rpc)
		if err != nil {
			log.Println("TCP read occured", err)
			return
		}

		log.Printf("message = %+v\n", rpc)
		t.rpcch <- rpc
	}

}

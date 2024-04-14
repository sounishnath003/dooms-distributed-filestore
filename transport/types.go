package transport

// peer represents anything that remote nodes connection only implementation
type Peer interface {
}

// transport handles any communication happening in the distributed network
// can be (TCP,UDP,websockets,GRPC,...)
type Transport interface {
	ListenAndAcceptConn() error
}

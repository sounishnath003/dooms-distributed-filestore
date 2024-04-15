package transport

// represents any arbitary message data that is being sent over each transport
type Message struct {
	Payload []byte
}

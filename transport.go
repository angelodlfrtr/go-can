package can

// Transport interface can be socketcan, an serial adapter, custom implementation, etc
type Transport interface {
	// Open a connection
	Open() error

	// Close a connection
	Close() error

	// Write a frame to connection
	Write(*Frame) error

	// Read a frame from connection
	Read(*Frame) (bool, error)
}

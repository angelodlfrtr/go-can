package can

import (
	"go-can/frame"
)

// Transport interface can be socketcan, an serial adapter, custom implementation, etc
type Transport interface {
	// Open a connection
	Open() error

	// Close a connection
	Close() error

	// Write a frame to connection
	Write(*frame.Frame) error

	// Read a frame from connection
	Read(*frame.Frame) (bool, error)
}

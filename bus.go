// Package can provide utils to connect to a can bus using hardware adapters system
package can

// Bus is the main interface to interact with the Transport
type Bus struct {
	// Transport represent the "logical" communication layer
	// which can be socketcan on linux, a serial adapater, or your custom transport
	Transport Transport
}

// NewBus create a new Bus with given transport
func NewBus(transport Transport) *Bus {
	return &Bus{Transport: transport}
}

// Open call Transport#Open
func (bus *Bus) Open() error {
	return bus.Transport.Open()
}

// Close call Transport#Close
func (bus *Bus) Close() error {
	return bus.Transport.Close()
}

// Write call Transport#Write
func (bus *Bus) Write(frm *Frame) error {
	return bus.Transport.Write(frm)
}

// ReadChan call Transport#ReadChan
func (bus *Bus) ReadChan() chan *Frame {
	return bus.Transport.ReadChan()
}

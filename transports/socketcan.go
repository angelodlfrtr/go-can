package transports

import (
	"github.com/angelodlfrtr/go-can"
	brutCan "github.com/brutella/can"
)

// SocketCan define a socketcan connection to canbus
type SocketCan struct {
	// Interface is the socket can interface to connect toa. eg : can0, vcan1, etc
	Interface string

	// bus is the can.Bus socket can connection
	bus *brutCan.Bus

	// busHandler handle can frame received
	busHandler brutCan.Handler

	// frames contain unread can frames from connectionJkA
	frames []brutCan.Frame

	// readChan
	readChan chan *can.Frame
}

// Open a socketcan connection
func (t *SocketCan) Open() error {
	// Open socketcan connection
	bus, err := brutCan.NewBusForInterfaceWithName(t.Interface)

	if err != nil {
		return err
	}

	go bus.ConnectAndPublish()

	t.readChan = make(chan *can.Frame)
	t.bus = bus

	// Create handler
	t.busHandler = brutCan.NewHandler(t.handleFrame)

	// Subcribe to incoming frames
	t.bus.Subscribe(t.busHandler)

	return nil
}

// Close a socketcan connection
func (t *SocketCan) Close() error {
	// Unsubscribe for frames
	t.bus.Unsubscribe(t.busHandler)

	// Close read chan
	close(t.readChan)

	// Close connectino
	return t.bus.Disconnect()
}

// Write data to socketcan interface
func (t *SocketCan) Write(frm *can.Frame) error {
	brutCanFrm := brutCan.Frame{
		ID:     frm.ArbitrationID,
		Length: frm.DLC,
		Flags:  0,
		Res0:   0,
		Res1:   0,
		Data:   frm.Data,
	}

	return t.bus.Publish(brutCanFrm)
}

// ReadChan
func (t *SocketCan) ReadChan() chan *can.Frame {
	return t.readChan
}

// handleFrame handle incoming frames from sockercan interface
// and add them to frames buffer
func (t *SocketCan) handleFrame(brutFrm brutCan.Frame) {
	frm := &can.Frame{}

	frm.ArbitrationID = brutFrm.ID
	frm.DLC = brutFrm.Length
	frm.Data = brutFrm.Data

	select {
	case t.readChan <- frm:
	default:
	}
}

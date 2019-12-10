package transports

import (
	brutCan "github.com/brutella/can"
	"go-can/frame"
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
}

// Open a socketcan connection
func (t *SocketCan) Open() error {
	// Open socketcan connection
	bus, err := brutCan.NewBusForInterfaceWithName(t.Interface)

	if err != nil {
		return err
	}

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

	// Close connectino
	return t.bus.Disconnect()
}

// Write data to socketcan interface
func (t *SocketCan) Write(frm *frame.Frame) error {
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

// Read data from socketcan interface
func (t *SocketCan) Read(frm *frame.Frame) (bool, error) {
	if len(t.frames) == 0 {
		return false, nil
	}

	frm.ArbitrationID = t.frames[0].ID
	frm.DLC = t.frames[0].Length
	frm.Data = t.frames[0].Data

	// Remove frame
	t.frames = t.frames[1:]

	return true, nil
}

// handleFrame handle incoming frames from sockercan interface
// and add them to frames buffer
func (t *SocketCan) handleFrame(frm brutCan.Frame) {
	t.frames = append(t.frames, frm)
}

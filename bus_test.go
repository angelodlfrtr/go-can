package can

import (
	"testing"
	"time"
)

// FakeTransport
type FakeTransport struct {
	readChan chan *Frame
}

func (t *FakeTransport) Open() error {
	return nil
}

// Close a serial connection
func (t *FakeTransport) Close() error {
	return nil
}

// Write a frame to serial connection
func (t *FakeTransport) Write(frm *Frame) error {
	return nil
}

// ReadChan returns the read chan
func (t *FakeTransport) ReadChan() chan *Frame {
	return t.readChan
}

func TestNewBus(t *testing.T) {
	tr := &FakeTransport{}
	bus := NewBus(tr)

	t.Log(*bus)
}

func TestOpen(t *testing.T) {
	tr := &FakeTransport{}
	bus := NewBus(tr)

	if err := bus.Open(); err != nil {
		t.Fatal(err)
	}
}

func TestWrite(t *testing.T) {
	tr := &FakeTransport{}
	bus := NewBus(tr)

	if err := bus.Open(); err != nil {
		t.Fatal(err)
	}

	frm := &Frame{
		ArbitrationID: uint32(0x45),
		DLC:           6,
		Data:          [8]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06},
	}

	if err := bus.Write(frm); err != nil {
		t.Fatal(err)
	}
}

func TestRead(t *testing.T) {
	tr := &FakeTransport{}
	bus := NewBus(tr)

	if err := bus.Open(); err != nil {
		t.Fatal(err)
	}

	start := time.Now()
	timeout := 1 * time.Second
	ticker := time.NewTicker(timeout)

	for {
		select {
		case frm := <-bus.ReadChan():
			t.Log(time.Since(start))
			t.Log(frm)
			t.Log("")
		case <-ticker.C:
			t.Log("Timeout")
			return
		}
	}
}

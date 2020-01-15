package can

import (
	"testing"
	"time"

	"github.com/angelodlfrtr/go-can/frame"
	"github.com/angelodlfrtr/go-can/transports"
)

const TestPort string = "/dev/tty.usbserial-1410"

func TestNewBus(t *testing.T) {
	tr := &transports.USBCanAnalyzer{
		Port:     TestPort,
		BaudRate: 2000000,
	}

	bus := NewBus(tr)

	t.Log(*bus)
}

func TestOpen(t *testing.T) {
	tr := &transports.USBCanAnalyzer{
		Port:     TestPort,
		BaudRate: 2000000,
	}

	bus := NewBus(tr)

	if err := bus.Open(); err != nil {
		t.Fatal(err)
	}
}

func TestWrite(t *testing.T) {
	tr := &transports.USBCanAnalyzer{
		Port:     TestPort,
		BaudRate: 2000000,
	}

	bus := NewBus(tr)

	if err := bus.Open(); err != nil {
		t.Fatal(err)
	}

	frm := &frame.Frame{
		ArbitrationID: uint32(0x45),
		DLC:           6,
		Data:          [8]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06},
	}

	if err := bus.Write(frm); err != nil {
		t.Fatal(err)
	}
}

func TestRead(t *testing.T) {
	tr := &transports.USBCanAnalyzer{
		Port:     TestPort,
		BaudRate: 2000000,
	}

	bus := NewBus(tr)

	if err := bus.Open(); err != nil {
		t.Fatal(err)
	}

	start := time.Now()
	timeout := 5 * time.Second

	for {
		if time.Since(start) > timeout {
			break
		}

		frm := &frame.Frame{}

		if ok, _ := bus.Read(frm); ok {
			t.Log(time.Since(start))
			t.Log(frm)
			t.Log("")
		}
	}
}

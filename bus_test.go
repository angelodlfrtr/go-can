package can

import (
	"github.com/angelodlfrtr/go-can/transports"
	"testing"
)

func TestNewBus(t *testing.T) {
	tr := &transports.USBCanAnalyzer{
		Port:     "/dev/tty.usbserial-14220",
		BaudRate: 2000000,
	}

	bus := NewBus(tr)

	t.Log(*bus)
}

func TestOpen(t *testing.T) {
	tr := &transports.USBCanAnalyzer{
		Port:     "/dev/tty.usbserial-14220",
		BaudRate: 2000000,
	}

	bus := NewBus(tr)

	if err := bus.Open(); err != nil {
		t.Fatal(err)
	}
}

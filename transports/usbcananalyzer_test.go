package transports

import (
	"testing"
	"time"

	"github.com/angelodlfrtr/go-can/frame"
)

func TestOpen(t *testing.T) {
	// Configure connection
	tr := &USBCanAnalyzer{
		Port:     "/dev/tty.usbserial-14220",
		BaudRate: 2000000,
	}

	// Try to open connection
	if err := tr.Open(); err != nil {
		t.Fatal(err)
	}

	t.Log("Connection opened")
}

func TestClose(t *testing.T) {
	// Configure connection
	tr := &USBCanAnalyzer{
		Port:     "/dev/tty.usbserial-14220",
		BaudRate: 2000000,
	}

	// Try to open connection
	if err := tr.Open(); err != nil {
		t.Fatal(err)
	}

	if err := tr.Close(); err != nil {
		t.Fatal(err)
	}

	t.Log("Connection closed")
}

func TestRead(t *testing.T) {
	// Configure connection
	tr := &USBCanAnalyzer{
		Port:     "/dev/tty.usbserial-14220",
		BaudRate: 2000000,
	}

	// Try to open connection
	if err := tr.Open(); err != nil {
		t.Fatal(err)
	}

	maxTimeout := 10 * time.Second

	start := time.Now()
	maxFrames := 4
	nbFrames := 0

	for {
		frm := &frame.Frame{}
		ok, err := tr.Read(frm)

		if err != nil {
			t.Fatal(err)
		}

		if ok {
			nbFrames++
			t.Log("Frame readed")

			t.Log("DLC : ", frm.DLC)
			t.Log("ArbID : ", frm.ArbitrationID)
			t.Log("Data : ", frm.Data)

			if nbFrames >= maxFrames {
				break
			}
		}

		if time.Since(start) > maxTimeout {
			t.Fatal("Test max read timeout exceeded")
		}
	}
}

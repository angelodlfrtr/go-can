package transports

import (
	"testing"
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

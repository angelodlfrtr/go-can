package transports

import (
	"fmt"
	"testing"
	"time"
)

func TestTCPCanOpen(t *testing.T) {
	// Configure connection
	tr := &TCPCan{
		Port: 7777,
		Host: "192.168.0.100",
	}

	fmt.Println("coucou")

	// Try to open connection
	if err := tr.Open(); err != nil {
		t.Fatal(err)
	}

	t.Log("Connection opened")

	go func() {
		for frm := range tr.ReadChan() {
			t.Log(frm)
		}
	}()

	time.Sleep(5 * time.Second)
}

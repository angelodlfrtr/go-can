package main

import (
	"log"
	"time"

	"github.com/angelodlfrtr/go-can"
	"github.com/angelodlfrtr/go-can/transports"
)

// TestPort contain serial path to test port
const TestPort string = "/dev/tty.usbserial-14140"

func main() {
	// Configure transport
	tr := &transports.USBCanAnalyzer{
		Port:     TestPort,
		BaudRate: 2000000,
	}

	// Open bus
	bus := can.NewBus(tr)

	if err := bus.Open(); err != nil {
		log.Fatal(err)
	}

	// Write some frames

	log.Println("Write 10 frames")

	for i := 0; i < 9; i++ {
		frm := &can.Frame{
			ArbitrationID: uint32(i),
			Data:          [8]byte{0x00, 0x01, uint8(i)},
		}

		if err := bus.Write(frm); err != nil {
			log.Fatal(err)
		}

		log.Printf("Frame %v writed", frm)
	}

	// Read frames during

	log.Println("Wait a frame (10s timeout)")
	timer := time.NewTimer(10 * time.Second)

	select {
	case frm := <-bus.ReadChan():
		log.Println(frm)
	case <-timer.C:
		log.Println("Timeout")
	}

	if err := bus.Close(); err != nil {
		log.Fatal(err)
	}

	log.Println("done")
}

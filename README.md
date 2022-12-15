# go-can

[github.com/angelodlfrtr/go-can](https://github.com/angelodlfrtr/go-can) is a canbus golang library supporting multiple transports (serial adapater, socketcan, etc).

**Does not support extended frames, feel free to create a PR**

## Installation

```bash
go get github.com/angelodlfrtr/go-can
```

## Basic usage

```go
package main

import (
	"log"
	"time"

	"github.com/angelodlfrtr/go-can"
	"github.com/angelodlfrtr/go-can/transports"
)

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
			Data:          [8]byte{0x00, 0X01, uint8(i)},
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
```

## License

Copyright (c) 2019 The contributors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

# go-can

[github.com/angelodlfrtr/go-can](https://github.com/angelodlfrtr/go-can) is a canbus golang library supporting multiple transports (serial adapater, socketcan, etc).


## Installation

```bash
go get github.com/angelodlfrtr/go-can
```

## Basic usage

```go
package main

import (
	"github.com/angelodlfrtr/go-can"
	"github.com/angelodlfrtr/go-can/transports"
	"log"
)

func main() {
	// Set up a transport, here via USBCanAnalyzer usb to serial adapater
	transport := transports.USBCanAnalyzer{
		Port:     "/dev/someusbtty",
		BaudRate: 2000000,
	}
	
	// Set up bus
	bus := can.NewBus(Transport: transport)
	
	// Try to open bus
	if err := can.Open(); err != nil {
		log.Fatal(err)
	}
	
	log.Println("Bus opened")
	
	// Write a frame
	frame := &can.Frame{
		ArbitrationID: uint32(0x45),
		DLC:           4,
		Data:          [8]byte{0x01, 0x02, 0x03, 0x04}
	}
	
	if err := bus.Write(frame); err != nil {
		Log.Fatal(err)
	}
	
	log.Println("Frame writed")
	
	// Read frames
	maxFrames := 5 // Read 5 frames
	nbFrames = 0
	
	log.Println("Listen for frames")
	
	for {
		frame := &can.Frame{}
		ok, err := bus.Read(frame)
		
		if err != nil {
			log.Fatal(err)
		}
		
		if ok {
			log.Println(frame)
			nbFrames++
			
			if nbFrames >= maxFrames {
				break
			}
		}
	}
	
	// Close the bus
	if err := bus.Close(); err != nil {
		Log.Fatal(err)
	}
	
	log.Println("Bus closed")
}
```

## License

Copyright (c) 2019 angelodlfrtr

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
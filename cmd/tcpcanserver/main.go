// Package main provide a small tcp server connection to a socket can interface and transfering
// data between server and a 'tcpcan' client. It allow client to do can over tcp
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/angelodlfrtr/go-can"
	"github.com/angelodlfrtr/go-can/transports"
)

// Flags
var socketCanInterface string
var port int

var clientConn net.Conn

func main() {
	// Parse flags
	flag.IntVar(&port, "port", 7777, "tcp listener port")
	flag.StringVar(&socketCanInterface, "ifr", "can0", "socket can interface to use")
	flag.Parse()

	// Open socket can bus
	log.Printf(">> Opening socketcan connection to %s\n", socketCanInterface)
	tr := &transports.SocketCan{Interface: socketCanInterface}

	bus := can.NewBus(tr)

	if err := bus.Open(); err != nil {
		log.Fatal(err)
	}
	log.Printf(">> Socketcan connection to %s success\n", socketCanInterface)

	// Start a tcp listener
	log.Printf(">> Staring a tcp listener on :%d\n", port)
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}

	// Accept connections
	go func() {
		for {
			conn, err := ln.Accept()

			if err != nil {
				log.Fatal(err)
			}

			log.Printf(">> Accepted client conn : %s\n", conn.RemoteAddr())

			if clientConn == nil {
				clientConn = conn
			}

			time.Sleep(10 * time.Millisecond)
		}
	}()

	// Write data from socketcan to tcp client
	go func() {
		readCh := bus.Transport.ReadChan()

		for frm := range readCh {
			if clientConn != nil {
				frmBytes, err := json.Marshal(frm)
				if err != nil {
					log.Println("ERROR", err)
					continue
				}

				frmBytes = append(frmBytes, []byte("\r\n")...)

				if _, err := clientConn.Write(frmBytes); err != nil {
					clientConn.Close()
					clientConn = nil
				}
			}
		}
	}()

	// Write data from tcp client to socketcan bus
	go func() {
		for {
			if clientConn == nil {
				time.Sleep(10 * time.Millisecond)
				continue
			}

			data := make([]byte, 512)
			ll, err := clientConn.Read(data)

			if err != nil {
				clientConn.Close()
				clientConn = nil
			}

			dataBytes := make([]byte, ll)
			copy(dataBytes, data)

			// Split data with delimiter
			frmsBytesSlice := bytes.Split(dataBytes, []byte("\r\n"))
			if len(frmsBytesSlice) == 0 {
				continue
			}

			for _, frmBytes := range frmsBytesSlice {
				// Try to convert bytes to a frame
				frm := &can.Frame{}
				if err := json.Unmarshal(frmBytes, &frm); err != nil {
					log.Println("ERROR", err)
					continue
				}

				bus.Write(frm)
			}
		}
	}()

	wait := make(chan bool)
	<-wait
}

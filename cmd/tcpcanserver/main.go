// Package main provide a small tcp server connection to a socket can interface and transferring
// data between server and a 'tcpcan' client. It allow client to do can over tcp
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"github.com/angelodlfrtr/go-can"
	"github.com/angelodlfrtr/go-can/transports"
)

// Flags
var (
	socketCanInterface string
	port               int
)

var (
	clientConn net.Conn
	clientBuf  []byte
	mutx       sync.Mutex
)

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

			mutx.Lock()
			if clientConn == nil {
				clientConn = conn
			}
			mutx.Unlock()

			go func() {
				jsonDec := json.NewDecoder(clientConn)

				for {
					frm := &can.Frame{}

					if err := jsonDec.Decode(frm); err != nil {
						if err == net.ErrClosed || err == io.EOF {
							break
						}

						panic(err)
					}

					bus.Write(frm)
				}
			}()
		}
	}()

	// Write data from socketcan to tcp client
	go func() {
		readCh := bus.Transport.ReadChan()

		for frm := range readCh {
			if clientConn != nil {
				frmBytes, err := json.Marshal(frm)
				if err != nil {
					panic(err)
				}

				frmBytes = append(frmBytes, []byte("\r\n")...)

				if _, err := clientConn.Write(frmBytes); err != nil {
					log.Println("error while write to client conn:", err.Error())
					clientConn.Close()
					mutx.Lock()
					clientConn = nil
					mutx.Unlock()
				}
			}
		}
	}()

	wait := make(chan bool)
	<-wait
}

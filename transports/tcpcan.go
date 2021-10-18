// @NOTE: use it in dev mode only. See cmd/tcpcanserver for server binary
package transports

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/angelodlfrtr/go-can"
)

// TCPCan define a custom TCP connection to an supported host
type TCPCan struct {
	// Port is the remote port eg : 7777
	Port int

	// Host is the remote host eg : 192.168.1.56
	Host string

	// conn contain current TCP connection to remote
	conn net.Conn

	// readChan
	readChan chan *can.Frame

	// stopChan
	stopChan chan bool
}

// Open tcpCan transport
func (tcpCan *TCPCan) Open() error {
	tcpCan.readChan = make(chan *can.Frame)
	tcpCan.stopChan = make(chan bool)

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", tcpCan.Host, tcpCan.Port), 2*time.Second)
	if err != nil {
		return err
	}

	tcpCan.conn = conn

	go tcpCan.listen()

	return nil
}

// Close tcpCan transport
func (tcpCan *TCPCan) Close() error {
	close(tcpCan.stopChan)
	return tcpCan.conn.Close()
}

// Write a frame to tcpCan transport
func (tcpCan *TCPCan) Write(frm *can.Frame) error {
	// Convert frame to json
	jsonBytes, err := json.Marshal(frm)
	if err != nil {
		return err
	}

	// Append delimiter
	jsonBytes = append(jsonBytes, []byte("\r\n")...)

	_, err = tcpCan.conn.Write(jsonBytes)
	return err
}

// ReadChan returns channel for reading frames
func (tcpCan *TCPCan) ReadChan() chan *can.Frame {
	return tcpCan.readChan
}

func (tcpCan *TCPCan) listen() {
	jsonDec := json.NewDecoder(tcpCan.conn)

	for {
		frm := &can.Frame{}
		if err := jsonDec.Decode(frm); err != nil {
			if err == net.ErrClosed || err == io.EOF {
				break
			}

			panic(err)
		}

		tcpCan.readChan <- frm
	}
}

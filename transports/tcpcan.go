package transports

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"

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

func (tcpCan *TCPCan) Open() error {
	tcpCan.readChan = make(chan *can.Frame)
	tcpCan.stopChan = make(chan bool)

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", tcpCan.Host, tcpCan.Port))
	if err != nil {
		return err
	}

	tcpCan.conn = conn

	go tcpCan.listen()

	return nil
}

func (tcpCan *TCPCan) Close() error {
	select {
	case tcpCan.stopChan <- true:
	default:
	}

	return tcpCan.conn.Close()
}

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

func (tcpCan *TCPCan) ReadChan() chan *can.Frame {
	return tcpCan.readChan
}

func (tcpCan *TCPCan) listen() {
	for {
		select {
		case <-tcpCan.stopChan:
			return
		default:
		}

		data := make([]byte, 512)
		ll, err := tcpCan.conn.Read(data)

		if err != nil {
			continue
		}

		dataBytes := make([]byte, ll)
		copy(dataBytes, data)
		frmsBytesSlice := bytes.Split(dataBytes, []byte("\r\n"))

		if len(frmsBytesSlice) == 0 {
			continue
		}

		for _, frmBytes := range frmsBytesSlice {
			if len(frmBytes) > 0 {
				// Try to convert bytes to a frame
				frm := &can.Frame{}
				if err := json.Unmarshal(frmBytes, &frm); err != nil {
					continue
				}

				select {
				case tcpCan.readChan <- frm:
				default:
				}
			}
		}
	}
}

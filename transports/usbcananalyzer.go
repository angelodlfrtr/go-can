package transports

import (
	"encoding/binary"
	"fmt"
	"github.com/angelodlfrtr/go-can/frame"
	"github.com/angelodlfrtr/serial"
	"io"
	"time"
)

// USBCanAnalyzer define a USBCanAnalyzer connection to canbus via serial connection on USB
type USBCanAnalyzer struct {
	// Port is the serial port eg : COM0 on windows, /dev/ttytest on posix, etc
	Port string

	// BaudRate is the serial connection baud rate
	BaudRate int

	// client is the serial.Port instance
	client *serial.Port

	// dataBuf contain data received by serial connection
	dataBuf []byte
}

// Open a serial connection
// Show https://github.com/kobolt/usb-can/blob/master/canusb.c for protocol definition
func (t *USBCanAnalyzer) Open() error {
	serialConfig := &serial.Config{
		// Name of the serial port
		Name: t.Port,

		// Baud rate should normaly be 2 000 000
		Baud: t.BaudRate,

		// ReadTimeout for the connection. If zero, the Read() operation is blocking
		ReadTimeout: 100 * time.Millisecond,

		// Size is 8 databytes for USBCanAnalyzer
		Size: 8,

		// StopBits is 1 for usbCanAnalyzer
		StopBits: 1,

		// Parity none for usbCanAnalyzer
		Parity: serial.ParityNone,
	}

	port, err := serial.OpenPort(serialConfig)

	if err != nil {
		return err
	}

	t.client = port

	// Send initalization sequence (configure adapter)
	seq := []byte{
		0xAA,
		0x55,
		0x12,
		0x07,
		0x01,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x01,
		0x00,
		0x00,
		0x00,
		0x00,
		0x1B,
	}

	if _, err := t.client.Write(seq); err != nil {
		return err
	}

	// Wait 500ms (else adapater has bugs)
	time.Sleep(500 * time.Millisecond)

	return nil
}

// Close a serial connection
func (t *USBCanAnalyzer) Close() error {
	if t.client == nil {
		return nil
	}

	return t.client.Close()
}

// Write a frame to serial connection
func (t *USBCanAnalyzer) Write(frm *frame.Frame) error {
	// 0xAA : adapter start of frame
	data := []byte{0xAA}

	// DLC
	data[1] = 0xC0 | frm.DLC

	// Write arbitration id
	binary.LittleEndian.PutUint32(data[2:], frm.ArbitrationID)

	// Append data
	for i := 0; i < 8; i++ {
		data[i+6] = frm.Data[i]
	}

	// Adapater end of frame
	data[13] = 0x55

	_, err := t.client.Write(data)
	return err
}

// Read a frame from serial connection
func (t *USBCanAnalyzer) Read(frm *frame.Frame) (bool, error) {
	data := make([]byte, 64)

	// Read data
	n, err := t.client.Read(data)

	if err == io.EOF {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	if n == 0 {
		return false, nil
	}

	// Append to global data buf
	for i := 0; i < n; i++ {
		t.dataBuf = append(t.dataBuf, data[i])
	}

	// Find adapter start of frame
	for {
		if len(t.dataBuf) == 0 {
			break
		}

		if t.dataBuf[0] == 0xAA {
			break
		}

		t.dataBuf = t.dataBuf[1:]
	}

	// If no data, return
	if len(t.dataBuf) == 0 {
		return false, nil
	}

	// Check if data can contain an entire frame (min frame size is 5 in case of 0 data)
	if len(t.dataBuf) < 5 {
		return false, nil
	}

	// Read frame

	// DLC
	frm.DLC = t.dataBuf[1] - 0xC0

	// Check buffer len can contain a frame
	if len(t.dataBuf) < 5+int(frm.DLC) {
		return false, nil
	}

	// Arbitration ID
	frm.ArbitrationID = binary.LittleEndian.Uint32(t.dataBuf[2:])

	// Data
	for i := 0; i < int(frm.DLC); i++ {
		frm.Data[i] = t.dataBuf[i+4]
	}

	// Resize t.dataBuf
	lastMsgLen := 1 + 1 + 2 + frm.DLC + 1 // 0xAA (SOF) + DLC + arbId + data + 0x55 (EOF)
	t.dataBuf = t.dataBuf[lastMsgLen:]

	fmt.Println(t.dataBuf)

	return true, nil
}

package can

// Frame represent a can frame
type Frame struct {
	// ArbitrationID is the frame identifier
	ArbitrationID uint32

	// DLC represent the size of the data field
	DLC uint8

	// Data is the data to transmit in the frame
	Data [8]byte
}

// GetData read frame.DLC data from frame.Data
func (frame *Frame) GetData() []byte {
	return frame.Data[0:frame.DLC]
}

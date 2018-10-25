package battleground

import (
	"fmt"
	"io"
)

type ReverseBuffer struct {
	buffer         []byte
	remainingBytes int
}

type ErrBufferOverrun struct {
	capacity int
	overrunBy int
}

func NewReverseBuffer(buffer []byte) *ReverseBuffer {
	r := new(ReverseBuffer)
	r.buffer = buffer
	r.remainingBytes = len(buffer)
	return r
}

func (rb *ReverseBuffer) GetFilledSlice() []byte {
	return rb.buffer[rb.remainingBytes:len(rb.buffer)]
}

func (rb *ReverseBuffer) Read(p []byte) (n int, err error) {
	length := len(p)

	if length > rb.remainingBytes {
		length = rb.remainingBytes
	}

	if length == 0 {
		return 0, io.EOF
	}

	low := rb.remainingBytes - length
	rb.checkOverrun(low)
	copy(p[0:length], rb.buffer[low:rb.remainingBytes])
	rb.remainingBytes -= length

	return length, nil
}

func (rb *ReverseBuffer) Write(p []byte) (n int, err error) {
	length := len(p)

	low := rb.remainingBytes - length
	resized, err := rb.resizeIfNeeded(-low)
	if err != nil {
		return 0, err
	}

	if resized {
		low = rb.remainingBytes - length
	}

	copy(rb.buffer[low:rb.remainingBytes], p[:])
	rb.remainingBytes -= length
	return length, nil
}

func (rb *ReverseBuffer) Seek(offset int64, whence int) (int64, error) {
	newRemainingBytes := rb.remainingBytes
	switch whence {
	case io.SeekStart:
		newRemainingBytes = len(rb.buffer) - int(offset)
	case io.SeekEnd:
		newRemainingBytes = int(offset)
	case io.SeekCurrent:
		newRemainingBytes -= int(offset)
	}

	if newRemainingBytes < 0 {
		return 0, &ErrBufferOverrun{
			len(rb.buffer),
			-newRemainingBytes,
		}
	}

	rb.remainingBytes = newRemainingBytes

	return int64(len(rb.buffer) - rb.remainingBytes), nil
}

func (rb *ReverseBuffer) resizeIfNeeded(minimumIncreaseBy int) (resized bool, err error) {
	if minimumIncreaseBy <= 0 {
		return false, nil
	}

	oldLength := len(rb.buffer)
	minNewLength := oldLength + minimumIncreaseBy

	var newLength int
	if oldLength == 0 {
		newLength = minNewLength
	} else {
		newLength = oldLength
	}

	for newLength < minNewLength {
		newLength *= 2
	}

	newBuffer := make([]byte, newLength)
	copy(newBuffer[oldLength:newLength], rb.buffer[:])
	rb.remainingBytes += newLength - oldLength
	rb.buffer = newBuffer

	return true, nil
}

func (e ErrBufferOverrun) Error() string {
	return fmt.Sprintf("Buffer size exceeded: capacity %d, overrun by %d", e.capacity, e.overrunBy)
}

func (rb *ReverseBuffer) checkOverrun(remainingBytes int) (err error) {
	if remainingBytes < 0 {
		return &ErrBufferOverrun{
			len(rb.buffer),
			-remainingBytes,
		}
	}

	return nil
}
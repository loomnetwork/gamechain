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

func (r *ReverseBuffer) GetFilledSlice() []byte {
	return r.buffer[r.remainingBytes:len(r.buffer)]
}

func (r *ReverseBuffer) Read(p []byte) (n int, err error) {
	length := len(p)
	fmt.Println(length)

	if length > r.remainingBytes {
		length = r.remainingBytes
	}

	if length == 0 {
		return 0, io.EOF
	}

	low := r.remainingBytes - length
	r.checkOverrun(low)
	copy(p[0:length], r.buffer[low:r.remainingBytes])
	r.remainingBytes -= length

	return length, nil
}

func (r *ReverseBuffer) Write(p []byte) (n int, err error) {
	length := len(p)

	low := r.remainingBytes - length
	err = r.checkOverrun(low)
	if err != nil {
		return 0, err
	}

	copy(r.buffer[low:r.remainingBytes], p[:])
	r.remainingBytes -= length
	return length, nil
}

func (r *ReverseBuffer) Seek(offset int64, whence int) (int64, error) {
	newRemainingBytes := r.remainingBytes
	switch whence {
	case io.SeekStart:
		newRemainingBytes = len(r.buffer) - int(offset)
	case io.SeekEnd:
		newRemainingBytes = int(offset)
	case io.SeekCurrent:
		newRemainingBytes -= int(offset)
	}

	if newRemainingBytes < 0 {
		return 0, &ErrBufferOverrun{
			len(r.buffer),
			-newRemainingBytes,
		}
	}

	r.remainingBytes = newRemainingBytes

	return int64(len(r.buffer) - r.remainingBytes), nil
}

func (e ErrBufferOverrun) Error() string {
	return fmt.Sprintf("Buffer size exceeded: capacity %d, overrun by %d", e.capacity, e.overrunBy)
}

func (r *ReverseBuffer) checkOverrun(remainingBytes int) (err error) {
	if remainingBytes < 0 {
		return &ErrBufferOverrun{
			len(r.buffer),
			-remainingBytes,
		}
	}

	return nil
}
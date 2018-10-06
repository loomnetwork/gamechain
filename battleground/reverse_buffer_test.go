package battleground

import (
	"encoding/binary"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestReverseBufferWrite(t *testing.T) {
	rb := NewReverseBuffer(make([]byte, 64))

	_ = binary.Write(rb, binary.LittleEndian, int16(1))
	_ = binary.Write(rb, binary.LittleEndian, int16(2))
	_ = binary.Write(rb, binary.LittleEndian, int8(3))
	_ = binary.Write(rb, binary.LittleEndian, int8(4))

	slice := rb.GetFilledSlice()
	assert.Equal(t, []byte { 4, 3, 2, 0, 1, 0 }, slice)
}

func TestReverseBufferRead(t *testing.T) {
	rb := NewReverseBuffer([]byte {4, 3, 2, 0, 1, 0 })
	var num16 int16
	var num8 int8

	_ = binary.Read(rb, binary.LittleEndian, &num16)
	assert.Equal(t, int16(1), num16)

	_ = binary.Read(rb, binary.LittleEndian, &num16)
	assert.Equal(t, int16(2), num16)

	_ = binary.Read(rb, binary.LittleEndian, &num8)
	assert.Equal(t, int8(3), num8)

	_ = binary.Read(rb, binary.LittleEndian, &num8)
	assert.Equal(t, int8(4), num8)
}

func TestReverseBufferWriteString(t *testing.T) {
	rb := NewReverseBuffer(make([]byte, 64))

	str := "Test 123"
	_ = binary.Write(rb, binary.LittleEndian, int32(len(str)))
	_ = binary.Write(rb, binary.LittleEndian, []byte(str))

	slice := rb.GetFilledSlice()
	assert.Equal(t, []byte { 'T', 'e', 's', 't', ' ', '1', '2', '3', 8, 0, 0, 0 }, slice)
}

func TestReverseBufferWriteOverrun(t *testing.T) {
	rb := NewReverseBuffer(make([]byte, 3))

	var err error
	err = binary.Write(rb, binary.LittleEndian, int16(1))
	assert.Equal(t, nil, err)
	err = binary.Write(rb, binary.LittleEndian, int16(2))
	assert.NotEqual(t, nil, err)
}

func TestReverseBufferReadOverrun(t *testing.T) {
	rb := NewReverseBuffer(make([]byte, 3))
	var num32 int32
	var err error
	err = binary.Read(rb, binary.LittleEndian, &num32)
	assert.NotEqual(t, nil, err)
}

func TestReverseBufferReadEOF(t *testing.T) {
	rb := NewReverseBuffer(make([]byte, 2))
	var num16 int16
	var err error
	err = binary.Read(rb, binary.LittleEndian, &num16)
	assert.Equal(t, nil, err)
	err = binary.Read(rb, binary.LittleEndian, &num16)
	assert.Equal(t, io.EOF, err)
}
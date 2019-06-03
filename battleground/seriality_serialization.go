package battleground

import (
	"bytes"
	"encoding/binary"
	"github.com/loomnetwork/gamechain/types/zb/zb_data"
	"io"
	"math/big"
)

func deserializeRect(r io.Reader) (rect zb_data.Rect, err error) {
	position, err := deserializeVector2Int(r)
	if err != nil {
		return
	}
	rect.Position = &position

	size, err := deserializeVector2Int(r)
	if err != nil {
		return
	}
	rect.Size_ = &size

	return rect, nil
}

func deserializeVector2Int(r io.Reader) (v zb_data.Vector2Int, err error) {
	if err = binary.Read(r, binary.BigEndian, &v.X); err != nil {
		return
	}

	if err = binary.Read(r, binary.BigEndian, &v.Y); err != nil {
		return
	}

	return
}

func serializeString(w io.Writer, str string) (err error) {
	var length = uint32(len(str))

	chunkCount := length / 32
	if length % 32 > 0 {
		chunkCount++
	}

	// Write 4 bytes uint32
	if err = binary.Write(w, binary.BigEndian, &length); err != nil {
		return
	}

	// Write 28 bytes of zeros
	if err = binary.Write(w, binary.BigEndian, make([]byte, 28)); err != nil {
		return
	}

	chunkedLength := chunkCount * 32
	chunkBuffer := make([]byte, 32)
	for i := uint32(0) ; i < chunkCount; i++ {
		chunkSize := 32
		if i == chunkCount - 1 {
			chunkSize = 32 - int(chunkedLength - length)
		}

		for i := range chunkBuffer {
			chunkBuffer[i] = 0
		}

		copy(chunkBuffer[:], []byte(str)[i * 32: i * 32 + uint32(chunkSize)])

		if _, err = w.Write(chunkBuffer); err != nil {
			return
		}
	}

	return nil
}

func deserializeString(r io.Reader) (str string, err error) {
	chunkBuffer := make([]byte, 32)
	if _, err = r.Read(chunkBuffer); err != nil {
		return
	}

	length := new(big.Int).SetBytes(chunkBuffer).Uint64()

	chunkCount := length / 32
	if length % 32 > 0 {
		chunkCount++
	}

	stringBuffer := bytes.NewBuffer(make([]byte, 0, length))

	chunkedLength := chunkCount * 32

	for i := uint64(0) ; i < chunkCount; i++ {
		chunkSize := 32
		if i == chunkCount - 1 {
			chunkSize = 32 - int(chunkedLength - length)
		}

		if _, err = r.Read(chunkBuffer); err != nil {
			return
		}

		if _, err = stringBuffer.Write(chunkBuffer[:chunkSize]); err != nil {
			return
		}
	}

	str = string(stringBuffer.Bytes())

	return
}
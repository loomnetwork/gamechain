package battleground

import (
	"bytes"
	"encoding/binary"
	"github.com/loomnetwork/gamechain/types/zb"
	"io"
	"math/big"
)

func deserializeRect(r io.Reader) (rect zb.Rect, err error) {
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

func deserializeVector2Int(r io.Reader) (v zb.Vector2Int, err error) {
	if err = binary.Read(r, binary.BigEndian, &v.X); err != nil {
		return
	}

	if err = binary.Read(r, binary.BigEndian, &v.Y); err != nil {
		return
	}

	return
}

func deserializeString(r io.Reader) (str string, err error) {
	lengthBuffer := make([]byte, 32)
	if _, err = r.Read(lengthBuffer); err != nil {
		return
	}

	length := new(big.Int).SetBytes(lengthBuffer).Uint64()

	chunkCount := length / 32
	if length % 32 > 0 {
		chunkCount++
	}

	stringBuffer := bytes.NewBuffer(make([]byte, 0, length))

	chunkBuffer := make([]byte, 32)
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
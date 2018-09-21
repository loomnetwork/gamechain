package battleground

import (
	"fmt"

	"github.com/loomnetwork/go-loom/util"
)

type Versioned interface {
	GetVersion() string
	MakeKey([]byte) []byte
}

type V1 struct {
	KeyPrefix string
}

func (v *V1) GetVersion() string {
	return "v1"
}

func (v *V1) MakeKey(key []byte) []byte {
	return util.PrefixKey([]byte(v.GetVersion()), key)
}

func getVersionedObject(version string) (Versioned, error) {
	if version == "v1" {
		return &V1{}, nil
	}
	return nil, fmt.Errorf("version not found")
}

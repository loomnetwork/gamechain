package battleground

import (
	"crypto/rand"
	"encoding/json"

	"github.com/google/uuid"

	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
)

const UUIDBytes = 16

func isOwner(ctx contract.Context, username string) bool {
	ok, _ := ctx.HasPermission([]byte(username), []string{"owner"})
	return ok
}

func prepareEmitMsgJSON(address []byte, owner, method string) ([]byte, error) {
	emitMsg := struct {
		Owner  string
		Method string
		Addr   []byte
	}{owner, method, address}

	return json.Marshal(emitMsg)
}

// Generates crypto random uuids
func generateUUID() (string, error) {
	buffer := make([]byte, UUIDBytes)

	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}

	if genUuid, err := uuid.FromBytes(buffer); err != nil {
		return "", err
	} else {
		return genUuid.String(), nil
	}

}

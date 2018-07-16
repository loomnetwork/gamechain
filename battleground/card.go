package battleground

import (
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	"github.com/loomnetwork/zombie_battleground/types/zb"
)

// TODO: need to merge with the main contract file.
// All of the card functionality should be moved to call the main net.

func (z *ZombieBattleground) ListCardLibrary(ctx contract.Context, req *zb.ListCardLibraryRequest) (zb.ListCardLibraryResponse, error) {
	return zb.ListCardLibraryResponse{}, nil
}

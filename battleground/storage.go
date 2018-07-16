package battleground

import (
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	"github.com/loomnetwork/go-loom/util"
	"github.com/loomnetwork/zombie_battleground/types/zb"
)

var (
	cardListKey = []byte("cardlist")
)

func cardKey(id string) []byte {
	return util.PrefixKey([]byte("card"), []byte(id))
}

func saveCardList(ctx contract.Context, cardList *zb.CardList) error {
	for _, card := range cardList.Cards {
		if err := ctx.Set(cardKey(card.Id), card); err != nil {
			return err
		}
	}
	return nil
}

func loadCardList(ctx contract.Context) (*zb.CardList, error) {
	var cl zb.CardList
	err := ctx.Get(cardListKey, &cl)
	if err != nil {
		return nil, err
	}
	return &cl, nil
}

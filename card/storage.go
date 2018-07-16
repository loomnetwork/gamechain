package card

import (
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	"github.com/loomnetwork/zombie_battleground/types/zbcard"
)

var (
	cardListKey = []byte("cardlist")
)

func saveCardList(ctx contract.Context, cardList *zbcard.CardList) error {
	return ctx.Set(cardListKey, cardList)
}

func loadCardList(ctx contract.Context) (*zbcard.CardList, error) {
	var cl zbcard.CardList
	err := ctx.Get(cardListKey, &cl)
	if err != nil {
		return nil, err
	}
	return &cl, nil
}

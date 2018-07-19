package battleground

import (
	"strconv"

	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	"github.com/loomnetwork/go-loom/util"
	"github.com/loomnetwork/zombie_battleground/types/zb"
)

var (
	cardListKey          = []byte("cardlist")
	heroListKey          = []byte("herolist")
	defaultDeckKey       = []byte("default-deck")
	defaultCollectionKey = []byte("default-collection")
)

func cardKey(id int64) []byte {
	return util.PrefixKey([]byte("card"), []byte(strconv.FormatInt(id, 10)))
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

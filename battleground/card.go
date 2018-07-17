package battleground

import (
	"sort"

	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	"github.com/loomnetwork/zombie_battleground/types/zb"
)

// TODO: need to merge with the main contract file.
// This is just temporary for client to get all the card info.
// All of the card functionality should be moved to call the mainnet.

func (z *ZombieBattleground) ListCardLibrary(ctx contract.StaticContext, req *zb.ListCardLibraryRequest) (*zb.ListCardLibraryResponse, error) {
	var cardList zb.CardList
	if err := ctx.Get(cardListKey, &cardList); err != nil {
		return nil, err
	}
	// convert to card list to card library view grouped by element
	var category = make(map[string][]*zb.Card)
	for _, card := range cardList.Cards {
		if _, ok := category[card.Element]; !ok {
			category[card.Element] = make([]*zb.Card, 0)
		}
		category[card.Element] = append(category[card.Element], card)
	}
	// order the element by name
	var elements []string
	for k := range category {
		elements = append(elements, k)
	}
	sort.Strings(elements)

	var sets []*zb.CardSet
	for _, elem := range elements {
		cards, ok := category[elem]
		if !ok {
			continue
		}
		set := &zb.CardSet{
			Name:  elem,
			Cards: cards,
		}
		sets = append(sets, set)
	}

	return &zb.ListCardLibraryResponse{Sets: sets}, nil
}

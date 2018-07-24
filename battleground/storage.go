package battleground

import (
	"encoding/json"
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

func prepareEmitMsgJSON(address []byte, owner, method string) ([]byte, error) {
	emitMsg := struct {
		Owner  string
		Method string
		Addr   []byte
	}{owner, method, address}

	return json.Marshal(emitMsg)
}

func isUser(ctx contract.Context, userID string) bool {
	ok, _ := ctx.HasPermission([]byte(userID), []string{"user"})
	return ok
}

func deleteDeckByName(decklist []*zb.Deck, name string) ([]*zb.Deck, bool) {
	newlist := make([]*zb.Deck, 0)
	for _, deck := range decklist {
		if deck.Name != name {
			newlist = append(newlist, deck)
		}
	}
	return newlist, len(newlist) != len(decklist)
}

func getDeckByName(decklist []*zb.Deck, name string) *zb.Deck {
	for _, deck := range decklist {
		if deck.Name == name {
			return deck
		}
	}
	return nil
}

func copyAccountInfo(account *zb.Account, req *zb.UpsertAccountRequest) {
	account.PhoneNumberVerified = req.PhoneNumberVerified
	account.RewardRedeemed = req.RewardRedeemed
	account.IsKickstarter = req.IsKickstarter
	account.Image = req.Image
	account.EmailNotification = req.EmailNotification
	account.EloScore = req.EloScore
	account.CurrentTier = req.CurrentTier
	account.GameMembershipTier = req.GameMembershipTier
}

package battleground

import (
	"encoding/json"
	"strconv"
	"strings"

	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	"github.com/loomnetwork/go-loom/util"
	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/pkg/errors"
)

const (
	OwnerRole = "owner"
)

var (
	cardPrefix       = []byte("card")
	userPreifx       = []byte("user")
	heroesPrefix     = []byte("heroes")
	collectionPrefix = []byte("collection")
	decksPrefix      = []byte("decks")

	cardListKey          = []byte("cardlist")
	heroListKey          = []byte("herolist")
	defaultDeckKey       = []byte("default-deck")
	defaultCollectionKey = []byte("default-collection")
	defaultHeroesKey     = []byte("default-heroes")
)

var (
	ErrNotfound        = errors.New("not found")
	ErrUserNotVerified = errors.New("user is not verified")
)

func userAccountKey(id string) []byte {
	return util.PrefixKey(userPreifx, []byte(id))
}

func userDecksKey(id string) []byte {
	return util.PrefixKey(userPreifx, []byte(id), decksPrefix)
}

func userCardCollectionKey(id string) []byte {
	return util.PrefixKey(userPreifx, []byte(id), collectionPrefix)
}

func userHeroesKey(id string) []byte {
	return util.PrefixKey(userPreifx, []byte(id), heroesPrefix)
}

func cardKey(id int64) []byte {
	return util.PrefixKey(cardPrefix, []byte(strconv.FormatInt(id, 10)))
}

func saveCardList(ctx contract.Context, cardList *zb.CardList) error {
	for _, card := range cardList.Cards {
		if err := ctx.Set(cardKey(card.Id), card); err != nil {
			return err
		}
	}
	return nil
}

func loadCardList(ctx contract.StaticContext) (*zb.CardList, error) {
	var cl zb.CardList
	err := ctx.Get(cardListKey, &cl)
	if err != nil {
		return nil, err
	}
	return &cl, nil
}

func loadCardCollection(ctx contract.StaticContext, userID string) (*zb.CardCollectionList, error) {
	var userCollection zb.CardCollectionList
	err := ctx.Get(userCardCollectionKey(userID), &userCollection)
	if err != nil && err != contract.ErrNotFound {
		return nil, err
	}
	return &userCollection, nil
}

func saveCardCollection(ctx contract.Context, userID string, cardCollection *zb.CardCollectionList) error {
	return ctx.Set(userCardCollectionKey(userID), cardCollection)
}

func loadDecks(ctx contract.StaticContext, userID string) (*zb.DeckList, error) {
	var deckList zb.DeckList
	err := ctx.Get(userDecksKey(userID), &deckList)
	if err != nil && err != contract.ErrNotFound {
		return nil, err
	}
	return &deckList, nil
}

func saveDecks(ctx contract.Context, userID string, decks *zb.DeckList) error {
	return ctx.Set(userDecksKey(userID), decks)
}

func loadHeroes(ctx contract.StaticContext, userID string) (*zb.HeroList, error) {
	var heroes zb.HeroList
	err := ctx.Get(userHeroesKey(userID), &heroes)
	if err != nil && err != contract.ErrNotFound {
		return nil, err
	}
	return &heroes, nil
}

func saveHeroes(ctx contract.Context, userID string, heroes *zb.HeroList) error {
	return ctx.Set(userHeroesKey(userID), heroes)
}

func prepareEmitMsgJSON(address []byte, owner, method string) ([]byte, error) {
	emitMsg := struct {
		Owner  string
		Method string
		Addr   []byte
	}{owner, method, address}

	return json.Marshal(emitMsg)
}

func isOwner(ctx contract.Context, userID string) bool {
	ok, _ := ctx.HasPermission([]byte(userID), []string{OwnerRole})
	return ok
}

func deleteDeckById(deckList []*zb.Deck, id int64) ([]*zb.Deck, bool) {
	newList := make([]*zb.Deck, 0)
	for _, deck := range deckList {
		if deck.Id != id {
			newList = append(newList, deck)
		}
	}
	return newList, len(newList) != len(deckList)
}

func getDeckById(deckList []*zb.Deck, id int64) *zb.Deck {
	for _, deck := range deckList {
		if deck.Id == id {
			return deck
		}
	}
	return nil
}

func getDeckByName(deckList []*zb.Deck, name string) *zb.Deck {
	for _, deck := range deckList {
		if strings.EqualFold(deck.Name, name) {
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

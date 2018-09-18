package battleground

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	"github.com/loomnetwork/go-loom/util"
	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/pkg/errors"
)

const (
	OwnerRole = "user" // TODO: change to owner
)

var (
	cardPrefix           = []byte("card")
	userPreifx           = []byte("user")
	heroesPrefix         = []byte("heroes")
	collectionPrefix     = []byte("collection")
	decksPrefix          = []byte("decks")
	matchesPrefix        = []byte("matches")
	pendingMatchesPrefix = []byte("pending-matches")

	cardListKey                 = []byte("cardlist")
	heroListKey                 = []byte("herolist")
	defaultDeckKey              = []byte("default-deck")
	defaultCollectionKey        = []byte("default-collection")
	defaultHeroesKey            = []byte("default-heroes")
	matchCountKey               = []byte("match-count")
	playersInMatchmakingListKey = []byte("players-matchmaking")
)

var (
	ErrNotfound        = errors.New("not found")
	ErrUserNotVerified = errors.New("user is not verified")
)

// Maintain compatability with version 1.
// TODO: Remove these and the following user* prefix instead if we're about to wipe out the data
func AccountKey(userID string) []byte {
	return []byte("user:" + userID)
}

func DecksKey(userID string) []byte {
	return []byte("user:" + userID + ":decks")
}

func CardCollectionKey(userID string) []byte {
	return []byte("user:" + userID + ":collection")
}

func HeroesKey(userID string) []byte {
	return []byte("user:" + userID + ":heroes")
}

func MatchKey(matchID int64) []byte {
	return []byte(fmt.Sprintf("match:%d", matchID))
}

func GameStateKey(gameStateID int64) []byte {
	return []byte(fmt.Sprintf("gamestate:%d", gameStateID))
}

// func userAccountKey(id string) []byte {
// 	return util.PrefixKey(userPreifx, []byte(id))
// }

// func userDecksKey(id string) []byte {
// 	return util.PrefixKey(userPreifx, []byte(id), decksPrefix)
// }

// func userCardCollectionKey(id string) []byte {
// 	return util.PrefixKey(userPreifx, []byte(id), collectionPrefix)
// }

// func userHeroesKey(id string) []byte {
// 	return util.PrefixKey(userPreifx, []byte(id), heroesPrefix)
// }

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
	err := ctx.Get(CardCollectionKey(userID), &userCollection)
	if err != nil && err != contract.ErrNotFound {
		return nil, err
	}
	return &userCollection, nil
}

func saveCardCollection(ctx contract.Context, userID string, cardCollection *zb.CardCollectionList) error {
	return ctx.Set(CardCollectionKey(userID), cardCollection)
}

func loadDecks(ctx contract.StaticContext, userID string) (*zb.DeckList, error) {
	var deckList zb.DeckList
	err := ctx.Get(DecksKey(userID), &deckList)
	if err != nil && err != contract.ErrNotFound {
		return nil, err
	}
	return &deckList, nil
}

func saveDecks(ctx contract.Context, userID string, decks *zb.DeckList) error {
	return ctx.Set(DecksKey(userID), decks)
}

func loadHeroes(ctx contract.StaticContext, userID string) (*zb.HeroList, error) {
	var heroes zb.HeroList
	err := ctx.Get(HeroesKey(userID), &heroes)
	if err != nil && err != contract.ErrNotFound {
		return nil, err
	}
	return &heroes, nil
}

func saveHeroes(ctx contract.Context, userID string, heroes *zb.HeroList) error {
	return ctx.Set(HeroesKey(userID), heroes)
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

func deleteDeckByID(deckList []*zb.Deck, id int64) ([]*zb.Deck, bool) {
	newList := make([]*zb.Deck, 0)
	for _, deck := range deckList {
		if deck.Id != id {
			newList = append(newList, deck)
		}
	}
	return newList, len(newList) != len(deckList)
}

func getDeckByID(deckList []*zb.Deck, id int64) *zb.Deck {
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

func savePendingMatchList(ctx contract.Context, pendingMatchList *zb.PendingMatchList) error {
	if err := ctx.Set(pendingMatchesPrefix, pendingMatchList); err != nil {
		return err
	}
	return nil
}

func loadPendingMatchList(ctx contract.StaticContext) (*zb.PendingMatchList, error) {
	var rl zb.PendingMatchList
	err := ctx.Get(pendingMatchesPrefix, &rl)
	if err != nil && err != contract.ErrNotFound {
		return nil, err
	}
	return &rl, nil
}

func saveMatchList(ctx contract.Context, matchList *zb.MatchList) error {
	if err := ctx.Set(matchesPrefix, matchList); err != nil {
		return err
	}
	return nil
}

func loadMatchList(ctx contract.StaticContext) (*zb.MatchList, error) {
	var rl zb.MatchList
	err := ctx.Get(matchesPrefix, &rl)
	if err != nil && err != contract.ErrNotFound {
		return nil, err
	}
	return &rl, nil
}

func saveMatch(ctx contract.Context, match *zb.Match) error {
	if err := ctx.Set(MatchKey(match.Id), match); err != nil {
		return err
	}
	return nil
}

func loadMatch(ctx contract.StaticContext, matchID int64) (*zb.Match, error) {
	var m zb.Match
	err := ctx.Get(MatchKey(matchID), &m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func saveMatchCount(ctx contract.Context, ID int64) error {
	if err := ctx.Set(matchCountKey, &zb.MatchCount{CurrentId: ID}); err != nil {
		return err
	}
	return nil
}

func loadMatchCount(ctx contract.StaticContext) (int64, error) {
	var count zb.MatchCount
	err := ctx.Get(matchCountKey, &count)
	if err != nil {
		return 0, err
	}
	return count.CurrentId, nil
}

func addPlayerInMatchmakingList(ctx contract.Context, ID string) error {
	IDs, err := loadPlayersInMatchmakingList(ctx)
	if err != nil && err != contract.ErrNotFound {
		return err
	}

	IDs = append(IDs, ID)

	list := zb.PlayersInMatchmakingList{}
	list.UserIDs = IDs
	if err := ctx.Set(playersInMatchmakingListKey, &list); err != nil {
		return err
	}

	return nil
}

func loadPlayersInMatchmakingList(ctx contract.StaticContext) ([]string, error) {
	var list zb.PlayersInMatchmakingList
	err := ctx.Get(playersInMatchmakingListKey, &list)
	if err != nil {
		return nil, err
	}
	return list.UserIDs, nil
}

func saveGameState(ctx contract.Context, gs *zb.GameState) error {
	if err := ctx.Set(GameStateKey(gs.Id), gs); err != nil {
		return err
	}
	return nil
}

func loadGameState(ctx contract.StaticContext, id int64) (*zb.GameState, error) {
	var state zb.GameState
	err := ctx.Get(GameStateKey(id), &state)
	if err != nil {
		return nil, err
	}
	return &state, nil
}

package battleground

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/loomnetwork/gamechain/types/zb"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	"github.com/loomnetwork/go-loom/util"
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
	matchMakingPrefix    = []byte("matchmaking")

	cardListKey                 = []byte("cardlist")
	heroListKey                 = []byte("herolist")
	defaultDeckKey              = []byte("default-deck")
	defaultCollectionKey        = []byte("default-collection")
	defaultHeroesKey            = []byte("default-heroes")
	matchCountKey               = []byte("match-count")
	playersInMatchmakingListKey = []byte("players-matchmaking")
	gameModeListKey             = []byte("gamemode-list")

	oracleKey = []byte("oracle-key")
)

var (
	ErrNotfound        = errors.New("not found")
	ErrUserNotVerified = errors.New("user is not verified")
)

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

func UserMatchKey(userID string) []byte {
	return []byte("user:" + userID + ":match")
}

func MakeVersionedKey(version string, key []byte) []byte {
	return util.PrefixKey([]byte(version), key)
}

func saveCardList(ctx contract.Context, version string, cardList *zb.CardList) error {
	return ctx.Set(MakeVersionedKey(version, cardListKey), cardList)
}

func loadCardList(ctx contract.StaticContext, version string) (*zb.CardList, error) {
	var cl zb.CardList
	err := ctx.Get(MakeVersionedKey(version, cardListKey), &cl)
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

func saveMatchMakingInfoList(ctx contract.Context, infos *zb.MatchMakingInfoList) error {
	if err := ctx.Set(matchesPrefix, infos); err != nil {
		return err
	}
	return nil
}

func loadMatchMakingInfoList(ctx contract.Context) (*zb.MatchMakingInfoList, error) {
	var infos zb.MatchMakingInfoList
	err := ctx.Get(matchesPrefix, &infos)
	if err != nil && err != contract.ErrNotFound {
		return nil, err
	}
	return &infos, nil
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

func createMatch(ctx contract.Context, match *zb.Match) error {
	nextID, err := nextMatchID(ctx)
	if err != nil {
		return err
	}
	match.Id = nextID
	match.Topics = []string{fmt.Sprintf("match:%d", nextID)}
	return saveMatch(ctx, match)
}

func loadMatch(ctx contract.StaticContext, matchID int64) (*zb.Match, error) {
	var m zb.Match
	err := ctx.Get(MatchKey(matchID), &m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func nextMatchID(ctx contract.Context) (int64, error) {
	var count zb.MatchCount
	err := ctx.Get(matchCountKey, &count)
	if err != nil && err != contract.ErrNotFound {
		return 0, err
	}
	count.CurrentId++
	if err := ctx.Set(matchCountKey, &zb.MatchCount{CurrentId: count.CurrentId}); err != nil {
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

func saveUserMatch(ctx contract.Context, userID string, match *zb.Match) error {
	if err := ctx.Set(UserMatchKey(userID), match); err != nil {
		return err
	}
	return nil
}

func loadUserMatch(ctx contract.StaticContext, userID string) (*zb.Match, error) {
	var m zb.Match
	err := ctx.Get(UserMatchKey(userID), &m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func addGameModeToList(ctx contract.Context, gameMode *zb.GameMode) error {
	gameModeList, err := loadGameModeList(ctx)
	if gameModeList == nil {
		gameModeList = &zb.GameModeList{GameModes: []*zb.GameMode{}}
	} else if err != nil {
		return err
	}
	gameModeList.GameModes = append(gameModeList.GameModes, gameMode)

	if err = saveGameModeList(ctx, gameModeList); err != nil {
		return err
	}

	return nil
}

func saveGameModeList(ctx contract.Context, gameModeList *zb.GameModeList) error {
	if err := ctx.Set(gameModeListKey, gameModeList); err != nil {
		return err
	}

	return nil
}

func loadGameModeList(ctx contract.StaticContext) (*zb.GameModeList, error) {
	var list zb.GameModeList
	err := ctx.Get(gameModeListKey, &list)
	if err != nil && err != contract.ErrNotFound {
		return nil, err
	}

	return &list, nil
}

func getGameModeFromList(gameModeList *zb.GameModeList, ID string) *zb.GameMode {
	for _, gameMode := range gameModeList.GameModes {
		if gameMode.ID == ID {
			return gameMode
		}
	}

	return nil
}

func getGameModeFromListByName(gameModeList *zb.GameModeList, name string) *zb.GameMode {
	for _, gameMode := range gameModeList.GameModes {
		if gameMode.Name == name {
			return gameMode
		}
	}

	return nil
}

func deleteGameMode(gameModeList *zb.GameModeList, ID string) (*zb.GameModeList, bool) {
	newList := make([]*zb.GameMode, 0)
	for _, gameMode := range gameModeList.GameModes {
		if gameMode.ID != ID {
			newList = append(newList, gameMode)
		}
	}

	return &zb.GameModeList{GameModes: newList}, len(newList) != len(gameModeList.GameModes)
}

func populateDeckCards(ctx contract.Context, playerStates []*zb.PlayerState, version string) error {
	var cardList zb.CardList
	if err := ctx.Get(MakeVersionedKey(version, cardListKey), &cardList); err != nil {
		return fmt.Errorf("error getting card library: %s", err)
	}
	for _, playerState := range playerStates {
		deck := playerState.Deck
		for _, deckCard := range deck.Cards {
			cardDetails, err := getCardDetails(&cardList, deckCard)
			if err != nil {
				return fmt.Errorf("unable to get card %s from card library: %s", deckCard.CardName, err.Error())
			}

			cardInstance := &zb.CardInstance{
				//InstanceId:
				Attack:  cardDetails.Damage,
				Defense: cardDetails.Health,
				Prototype: &zb.CardPrototype{
					Name: cardDetails.Name,
				},
			}
			playerState.CardsInDeck = append(playerState.CardsInDeck, cardInstance)
		}
		for _, c := range playerState.CardsInDeck {
			ctx.Logger().Debug(fmt.Sprintf("card: name :%s, attack: %v\n", c.Prototype.Name, c.Attack))
		}
	}

	return nil
}

func getCardDetails(cardList *zb.CardList, deckCard *zb.CardCollection) (*zb.Card, error) {
	for _, card := range cardList.Cards {
		if card.Name == deckCard.CardName {
			return card, nil
		}
	}
	return nil, fmt.Errorf("card not found in card library")
}

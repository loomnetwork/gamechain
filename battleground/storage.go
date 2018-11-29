package battleground

import (
	"encoding/json"
	"fmt"
	"github.com/gogo/protobuf/proto"
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
	userPrefix           = []byte("user")
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
	playerPoolKey               = []byte("playerpool")
	taggedPlayerPoolKey         = []byte("tagged-playerpool")
	oracleKey                   = []byte("oracle-key")
	aiDecksKey                  = []byte("ai-decks")
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

func saveAIDecks(ctx contract.Context, version string, decks *zb.AIDeckList) error {
	return ctx.Set(MakeVersionedKey(version, aiDecksKey), decks)
}

func loadAIDecks(ctx contract.StaticContext, version string) (*zb.AIDeckList, error) {
	var deckList zb.AIDeckList
	err := ctx.Get(MakeVersionedKey(version, aiDecksKey), &deckList)
	if err != nil {
		return nil, err
	}
	return &deckList, nil
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

func savePlayerPool(ctx contract.Context, pool *zb.PlayerPool) error {
	return ctx.Set(playerPoolKey, pool)
}

func loadPlayerPool(ctx contract.StaticContext) (*zb.PlayerPool, error) {
	var pool zb.PlayerPool
	err := ctx.Get(playerPoolKey, &pool)
	if err != nil && err != contract.ErrNotFound {
		return nil, err
	}
	return &pool, nil
}

func saveTaggedPlayerPool(ctx contract.Context, pool *zb.PlayerPool) error {
	return ctx.Set(taggedPlayerPoolKey, pool)
}

func loadTaggedPlayerPool(ctx contract.StaticContext) (*zb.PlayerPool, error) {
	var pool zb.PlayerPool
	err := ctx.Get(taggedPlayerPoolKey, &pool)
	if err != nil && err != contract.ErrNotFound {
		return nil, err
	}
	return &pool, nil
}

func saveMatch(ctx contract.Context, match *zb.Match) error {
	if err := ctx.Set(MatchKey(match.Id), match); err != nil {
		return err
	}
	return nil
}

func createMatch(ctx contract.Context, match *zb.Match, useClientGameLogic bool) error {
	nextID, err := nextMatchID(ctx)
	if err != nil {
		return err
	}
	match.Id = nextID
	match.Topics = []string{fmt.Sprintf("match:%d", nextID)}
	match.CreatedAt = ctx.Now().Unix()
	match.UseClientGameLogic = useClientGameLogic
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

func saveUserCurrentMatch(ctx contract.Context, userID string, match *zb.Match) error {
	if err := ctx.Set(UserMatchKey(userID), match); err != nil {
		return err
	}
	return nil
}

func loadUserCurrentMatch(ctx contract.StaticContext, userID string) (*zb.Match, error) {
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

func newCardInstanceSpecificDataFromCardDetails(cardDetails *zb.Card) *zb.CardInstanceSpecificData {
	return &zb.CardInstanceSpecificData{
		Attack: cardDetails.Attack,
		Defense: cardDetails.Defense,
		Type: cardDetails.Type,
		Set: cardDetails.Set,
		GooCost: cardDetails.GooCost,
	}
}

func newCardInstanceFromCardDetails(cardDetails *zb.Card, instanceId int32, owner string) *zb.CardInstance {
	return &zb.CardInstance{
		InstanceId: instanceId,
		Owner:      owner,
		Prototype: proto.Clone(cardDetails).(*zb.Card),
		Instance: newCardInstanceSpecificDataFromCardDetails(cardDetails),
	}
}

func populateDeckCards(ctx contract.Context, cardLibrary *zb.CardList, playerStates []*zb.PlayerState) error {
	instanceId := int32(0) // unique instance IDs for cards
	for _, playerState := range playerStates {
		deck := playerState.Deck
		for _, cardAmounts := range deck.Cards {
			for i := int64(0); i < cardAmounts.Amount; i++ {
				cardDetails, err := getCardDetails(cardLibrary, cardAmounts.CardName)
				if err != nil {
					return fmt.Errorf("unable to get card %s from card library: %s", cardAmounts.CardName, err.Error())
				}

				cardInstance := newCardInstanceFromCardDetails(
					cardDetails,
					instanceId,
					playerState.Id,
				)

				playerState.CardsInDeck = append(playerState.CardsInDeck, cardInstance)
				instanceId++
			}
		}
	}

	return nil
}

func getCardLibrary(ctx contract.Context, version string) (*zb.CardList, error) {
	var cardList zb.CardList
	if err := ctx.Get(MakeVersionedKey(version, cardListKey), &cardList); err != nil {
		return nil, fmt.Errorf("error getting card library: %s", err)
	}

	return &cardList, nil
}

func getCardDetails(cardList *zb.CardList, cardName string) (*zb.Card, error) {
	for _, card := range cardList.Cards {
		if card.Name == cardName {
			return card, nil
		}
	}
	return nil, fmt.Errorf("card with name %s not found in card library", cardName)
}

package battleground

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gogo/protobuf/proto"

	"github.com/loomnetwork/gamechain/types/zb"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	"github.com/loomnetwork/go-loom/util"
	"github.com/pkg/errors"
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
	matchCountKey               = []byte("match-count")
	playersInMatchmakingListKey = []byte("players-matchmaking")
	gameModeListKey             = []byte("gamemode-list")
	playerPoolKey               = []byte("playerpool")
	taggedPlayerPoolKey         = []byte("tagged-playerpool")
	oracleKey                   = []byte("oracle-key")
	aiDecksKey                  = []byte("ai-decks")
	stateKey                    = []byte("state")
	nonceKey                    = []byte("nonce")
	currentUserIDUIntKey        = []byte("current-user-id")
	contentVersionKey           = []byte("content-version")
	pvpVersionKey               = []byte("pvp-version")
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

func InitialGameStateKey(gameStateID int64) []byte {
	return []byte(fmt.Sprintf("initial-gamestate:%d", gameStateID))
}

func UserMatchKey(userID string) []byte {
	return []byte("user:" + userID + ":match")
}

func MakeVersionedKey(version string, key []byte) []byte {
	return util.PrefixKey([]byte(version), key)
}

func saveState(ctx contract.Context, state *zb.GamechainState) error {
	if err := ctx.Set(stateKey, state); err != nil {
		return err
	}
	return nil
}

func loadState(ctx contract.StaticContext) (*zb.GamechainState, error) {
	var m zb.GamechainState
	err := ctx.Get(stateKey, &m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func RewardTutorialClaimedKey(userID string) []byte {
	return []byte("user:" + userID + ":rewardTutorialClaimed")
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

// loadCardCollectionFromAddress loads address mapping to card collection
func loadCardCollectionByAddress(ctx contract.StaticContext) (*zb.CardCollectionList, error) {
	var userCollection zb.CardCollectionList
	addr := string(ctx.Message().Sender.Local)
	err := ctx.Get(CardCollectionKey(addr), &userCollection)
	if err != nil && err != contract.ErrNotFound {
		return nil, err
	}
	return &userCollection, nil
}

// saveCardCollectionByAddress save card collection using address as a key
func saveCardCollectionByAddress(ctx contract.Context, cardCollection *zb.CardCollectionList) error {
	addr := string(ctx.Message().Sender.Local)
	return ctx.Set(CardCollectionKey(addr), cardCollection)
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

func getDeckWithRegistrationData(ctx contract.Context, registrationData *zb.PlayerProfileRegistrationData) (*zb.Deck, error) {
	if registrationData.DebugCheats.Enabled && registrationData.DebugCheats.UseCustomDeck {
		return registrationData.DebugCheats.CustomDeck, nil
	}

	// get matched player deck
	matchedDl, err := loadDecks(ctx, registrationData.UserId)
	if err != nil {
		return nil, err
	}
	matchedDeck := getDeckByID(matchedDl.Decks, registrationData.DeckId)
	if matchedDeck == nil {
		return nil, fmt.Errorf("deck id %d not found", registrationData.DeckId)
	}

	return matchedDeck, nil
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

func loadPlayerPool(ctx contract.Context) (*zb.PlayerPool, error) {
	return loadPlayerPoolInternal(ctx, playerPoolKey)
}

func saveTaggedPlayerPool(ctx contract.Context, pool *zb.PlayerPool) error {
	return ctx.Set(taggedPlayerPoolKey, pool)
}

func loadTaggedPlayerPool(ctx contract.Context) (*zb.PlayerPool, error) {
	return loadPlayerPoolInternal(ctx, taggedPlayerPoolKey)
}

func loadPlayerPoolInternal(ctx contract.Context, poolKey []byte) (*zb.PlayerPool, error) {
	var pool zb.PlayerPool
	err := ctx.Get(poolKey, &pool)
	if err != nil && err != contract.ErrNotFound {
		// Try to reset the pool
		ctx.Logger().Error("error loading pool, clearing", "key", string(taggedPlayerPoolKey), "err", err)
		pool = zb.PlayerPool{}
		if err = ctx.Set(taggedPlayerPoolKey, &pool); err != nil {
			return nil, err
		}

		return &pool, nil
	}
	return &pool, nil
}

func saveMatch(ctx contract.Context, match *zb.Match) error {
	if err := ctx.Set(MatchKey(match.Id), match); err != nil {
		return err
	}
	return nil
}

func createMatch(ctx contract.Context, match *zb.Match, useBackendGameLogic bool) error {
	nextID, err := nextMatchID(ctx)
	if err != nil {
		return err
	}
	match.Id = nextID
	match.Topics = []string{fmt.Sprintf("match:%d", nextID)}
	match.CreatedAt = ctx.Now().Unix()
	match.UseBackendGameLogic = useBackendGameLogic
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

func setRewardTutorialClaimed(ctx contract.Context, userID string, claim *zb.RewardTutorialClaimed) error {
	return ctx.Set(RewardTutorialClaimedKey(userID), claim)
}

func getRewardTutorialClaimed(ctx contract.Context, userID string) (*zb.RewardTutorialClaimed, error) {
	var rewardClaimed zb.RewardTutorialClaimed
	err := ctx.Get(RewardTutorialClaimedKey(userID), &rewardClaimed)
	if err != nil && err != contract.ErrNotFound {
		return nil, err
	}
	return &rewardClaimed, nil
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

func saveInitialGameState(ctx contract.Context, gs *zb.GameState) error {
	if err := ctx.Set(InitialGameStateKey(gs.Id), gs); err != nil {
		return err
	}
	return nil
}

func loadInitialGameState(ctx contract.StaticContext, id int64) (*zb.GameState, error) {
	var state zb.GameState
	err := ctx.Get(InitialGameStateKey(id), &state)
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
		Attack:  cardDetails.Attack,
		Defense: cardDetails.Defense,
		Type:    cardDetails.Type,
		Faction: cardDetails.Faction,
		GooCost: cardDetails.GooCost,
	}
}

func newCardInstanceFromCardDetails(cardDetails *zb.Card, instanceID *zb.InstanceId, owner string, ownerIndex int32) *zb.CardInstance {
	instance := newCardInstanceSpecificDataFromCardDetails(cardDetails)
	var abilities []*zb.CardAbilityInstance
	for _, raw := range cardDetails.Abilities {
		switch raw.Type {
		case zb.CardAbilityType_Rage:
			abilities = append(abilities, &zb.CardAbilityInstance{
				IsActive: true,
				Trigger:  raw.Trigger,
				AbilityType: &zb.CardAbilityInstance_Rage{
					Rage: &zb.CardAbilityRage{
						AddedAttack: raw.Value,
					},
				},
			})
		case zb.CardAbilityType_PriorityAttack:
			abilities = append(abilities, &zb.CardAbilityInstance{
				IsActive: true,
				Trigger:  raw.Trigger,
				AbilityType: &zb.CardAbilityInstance_PriorityAttack{
					PriorityAttack: &zb.CardAbilityPriorityAttack{},
				},
			})
		case zb.CardAbilityType_ReanimateUnit:
			abilities = append(abilities, &zb.CardAbilityInstance{
				IsActive: true,
				Trigger:  raw.Trigger,
				AbilityType: &zb.CardAbilityInstance_Reanimate{
					Reanimate: &zb.CardAbilityReanimate{
						DefaultAttack:  cardDetails.Attack,
						DefaultDefense: cardDetails.Defense,
					},
				},
			})
		case zb.CardAbilityType_ChangeStat:
			abilities = append(abilities, &zb.CardAbilityInstance{
				IsActive: true,
				Trigger:  raw.Trigger,
				AbilityType: &zb.CardAbilityInstance_ChangeStat{
					ChangeStat: &zb.CardAbilityChangeStat{
						StatAdjustment: raw.Value,
						Stat:           raw.Stat,
					},
				},
			})
		case zb.CardAbilityType_AttackOverlord:
			abilities = append(abilities, &zb.CardAbilityInstance{
				IsActive: true,
				Trigger:  raw.Trigger,
				AbilityType: &zb.CardAbilityInstance_AttackOverlord{
					AttackOverlord: &zb.CardAbilityAttackOverlord{
						Damage: raw.Value,
					},
				},
			})
		case zb.CardAbilityType_ReplaceUnitsWithTypeOnStrongerOnes:
			abilities = append(abilities, &zb.CardAbilityInstance{
				IsActive: true,
				Trigger:  raw.Trigger,
				AbilityType: &zb.CardAbilityInstance_ReplaceUnitsWithTypeOnStrongerOnes{
					ReplaceUnitsWithTypeOnStrongerOnes: &zb.CardAbilityReplaceUnitsWithTypeOnStrongerOnes{
						Faction: cardDetails.Faction,
					},
				},
			})
		case zb.CardAbilityType_DealDamageToThisAndAdjacentUnits:
			abilities = append(abilities, &zb.CardAbilityInstance{
				IsActive: true,
				Trigger:  raw.Trigger,
				AbilityType: &zb.CardAbilityInstance_DealDamageToThisAndAdjacentUnits{
					DealDamageToThisAndAdjacentUnits: &zb.CardAbilityDealDamageToThisAndAdjacentUnits{
						AdjacentDamage: cardDetails.Attack,
					},
				},
			})
		}

	}
	return &zb.CardInstance{
		InstanceId:         proto.Clone(instanceID).(*zb.InstanceId),
		Owner:              owner,
		Prototype:          proto.Clone(cardDetails).(*zb.Card),
		Instance:           instance,
		AbilitiesInstances: abilities,
		Zone:               zb.Zone_DECK, // default to deck
		OwnerIndex:         ownerIndex,
	}
}

func getInstanceIdsFromCardInstances(cards []*zb.CardInstance) []*zb.InstanceId {
	var instanceIds = make([]*zb.InstanceId, len(cards), len(cards))
	for i := range cards {
		instanceIds[i] = cards[i].InstanceId
	}

	return instanceIds
}

func populateDeckCards(cardLibrary *zb.CardList, playerStates []*zb.PlayerState, useBackendGameLogic bool) error {
	for playerIndex, playerState := range playerStates {
		deck := playerState.Deck
		if deck == nil {
			return fmt.Errorf("no card deck fro player %s", playerState.Id)
		}
		for _, cardAmounts := range deck.Cards {
			for i := int64(0); i < cardAmounts.Amount; i++ {
				cardDetails, err := getCardDetails(cardLibrary, cardAmounts.CardName)
				if err != nil {
					return fmt.Errorf("unable to get card %s from card library: %s", cardAmounts.CardName, err.Error())
				}

				cardInstance := newCardInstanceFromCardDetails(
					cardDetails,
					nil,
					playerState.Id,
					int32(playerIndex),
				)
				playerState.CardsInDeck = append(playerState.CardsInDeck, cardInstance)
			}
		}
	}

	removeUnsupportedCardFeatures(useBackendGameLogic, playerStates)

	return nil
}

func removeUnsupportedCardFeatures(useBackendGameLogic bool, playerStates []*zb.PlayerState) {
	if !useBackendGameLogic {
		return
	}

	for _, playerState := range playerStates {
		filteredCards := make([]*zb.CardInstance, 0, 0)

		for _, card := range playerState.CardsInDeck {
			filteredAbilities := make([]*zb.CardAbility, 0, 0)
			for _, ability := range card.Prototype.Abilities {
				switch ability.Type {
				case zb.CardAbilityType_Rage:
					fallthrough
				case zb.CardAbilityType_PriorityAttack:
					fallthrough
				case zb.CardAbilityType_ReanimateUnit:
					fallthrough
				case zb.CardAbilityType_ChangeStat:
					fallthrough
				case zb.CardAbilityType_AttackOverlord:
					fallthrough
				case zb.CardAbilityType_ReplaceUnitsWithTypeOnStrongerOnes:
					filteredAbilities = append(filteredAbilities, ability)
				default:
					fmt.Printf("Unsupported CardAbilityType value %s, removed (card '%s')\n", zb.CardAbilityType_Enum_name[int32(ability.Type)], card.Prototype.Name)
				}
			}

			card.Prototype.Abilities = filteredAbilities

			switch card.Prototype.Type {
			case zb.CreatureType_Feral:
				fallthrough
			case zb.CreatureType_Heavy:
				fmt.Printf("Unsupported CreatureType value %s, fallback to WALKER (card %s)\n", zb.CreatureType_Enum_name[int32(card.Prototype.Type)], card.Prototype.Name)
				card.Prototype.Type = zb.CreatureType_Walker
			}

			switch card.Instance.Type {
			case zb.CreatureType_Feral:
				fallthrough
			case zb.CreatureType_Heavy:
				fmt.Printf("Unsupported CreatureType value %s, fallback to WALKER (card %s)\n", zb.CreatureType_Enum_name[int32(card.Instance.Type)], card.Prototype.Name)
				card.Instance.Type = zb.CreatureType_Walker
			}

			switch card.Prototype.Kind {
			case zb.CardKind_Creature:
				filteredCards = append(filteredCards, card)
			default:
				fmt.Printf("Unsupported CardKind value %s, removed (card '%s')\n", zb.CardKind_Enum_name[int32(card.Prototype.Kind)], card.Prototype.Name)
			}

			switch card.Prototype.Rank {
			case zb.CreatureRank_Officer:
				fallthrough
			case zb.CreatureRank_Commander:
				fallthrough
			case zb.CreatureRank_General:
				fmt.Printf("Unsupported CreatureRank value %s, fallback to MINION (card %s)\n", zb.CreatureRank_Enum_name[int32(card.Prototype.Rank)], card.Prototype.Name)
				card.Prototype.Rank = zb.CreatureRank_Minion
			}
		}

		playerState.CardsInDeck = filteredCards
	}
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

func findCardByMouldID(cardList *zb.CardList, mouldID int64) (*zb.Card, bool) {
	for _, card := range cardList.Cards {
		if card.MouldId == mouldID {
			return card, true
		}
	}
	return nil, false
}

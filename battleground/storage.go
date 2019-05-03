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
	overlordsPrefix      = []byte("heroes")
	collectionPrefix     = []byte("collection")
	decksPrefix          = []byte("decks")
	matchesPrefix        = []byte("matches")
	pendingMatchesPrefix = []byte("pending-matches")
	matchMakingPrefix    = []byte("matchmaking")

	cardListKey                 = []byte("cardlist")
	overlordListKey             = []byte("herolist")
	defaultDecksKey             = []byte("default-deck")
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
	overlordLevelingDataKey     = []byte("overlord-experience")
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

func OverlordsKey(userID string) []byte {
	return []byte("user:" + userID + ":heroes")
}

func UserNotificationsKey(userID string) []byte {
	return []byte("user:" + userID + ":notifications")
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

func loadCardCollectionByUserId(ctx contract.Context, userID string, version string) (*zb.CardCollectionList, error) {
	var userCollection zb.CardCollectionList
	err := ctx.Get(CardCollectionKey(userID), &userCollection)
	if err != nil {
		if err == contract.ErrNotFound {
			userCollection.Cards = []*zb.CardCollectionCard{}
		} else {
			return nil, err
		}
	}

	// Update data
	err = validateAndUpdateCardCollectionList(
		ctx,
		version,
		&userCollection,
		false,
		func(collection *zb.CardCollectionList) error {
			err := saveCardCollectionByUserId(ctx, userID, collection)
			if err != nil {
				return err
			}

			return nil
		})

	if err != nil {
		return nil, err
	}

	return &userCollection, nil
}

func saveCardCollectionByUserId(ctx contract.Context, userID string, cardCollection *zb.CardCollectionList) error {
	return ctx.Set(CardCollectionKey(userID), cardCollection)
}

// loadCardCollectionFromAddress loads address mapping to card collection
func loadCardCollectionByAddress(ctx contract.Context, version string) (*zb.CardCollectionList, error) {
	var userCollection zb.CardCollectionList
	addr := string(ctx.Message().Sender.Local)
	err := ctx.Get(CardCollectionKey(addr), &userCollection)
	if err != nil {
		if err == contract.ErrNotFound {
			userCollection.Cards = []*zb.CardCollectionCard{}
		} else {
			return nil, err
		}
	}

	// Update data
	err = validateAndUpdateCardCollectionList(
		ctx,
		version,
		&userCollection,
		false,
		func(collection *zb.CardCollectionList) error {
			err :=  saveCardCollectionByAddress(ctx, collection)
			if err != nil {
				return err
			}

			return nil
		})

	if err != nil {
		return nil, err
	}

	return &userCollection, nil
}

// saveCardCollectionByAddress save card collection using address as a key
func saveCardCollectionByAddress(ctx contract.Context, cardCollection *zb.CardCollectionList) error {
	addr := string(ctx.Message().Sender.Local)
	return ctx.Set(CardCollectionKey(addr), cardCollection)
}

func validateAndUpdateGenericCardList(
	ctx contract.StaticContext,
	cardList []interface{},
	version string,
	validateMouldId bool,
	getMouldIdFunc func(card interface{}) int64,
	setMouldIdFunc func(card interface{}, mouldId int64),
	getCardNameFunc func(card interface{}) string,
	setEmptyCardNameFunc func(card interface{}),
	) (changed bool, changedList []interface{}, err error) {
	changed = false
	err = nil

	var cardLibrary *zb.CardList = nil
	if validateMouldId {
		cardLibrary, err = getCardLibrary(ctx, version)
		if err != nil {
			return false, nil, err
		}
	}

	newCardList := make([]interface{}, 0)
	for _, card := range cardList {
		if getMouldIdFunc(card) != 0 {
			if validateMouldId {
				// Check if mould ID exists in card library
				mouldIdMatchFound := false
				for _, libraryCard := range cardLibrary.Cards {
					if libraryCard.MouldId == getMouldIdFunc(card) {
						mouldIdMatchFound = true
						break
					}
				}

				if mouldIdMatchFound {
					newCardList = append(newCardList, card)
				} else {
					// If a match is not found, remove the card from list
					ctx.Logger().Warn(fmt.Sprintf("card with mould id %d not found in card library, removing", getMouldIdFunc(card)))
					changed = true
				}
			} else {
				newCardList = append(newCardList, card)
			}

			continue
		}

		if !validateMouldId && cardLibrary == nil {
			cardLibrary, err = getCardLibrary(ctx, version)
			if err != nil {
				return false, nil, err
			}
		}

		// Convert card name to mould ID
		nameMatchFound := false
		for _, libraryCard := range cardLibrary.Cards {
			if strings.EqualFold(libraryCard.Name, getCardNameFunc(card)) {
				nameMatchFound = true
				setMouldIdFunc(card, libraryCard.MouldId)
				setEmptyCardNameFunc(card)
				break
			}
		}

		if !nameMatchFound {
			ctx.Logger().Warn(fmt.Sprintf("card with name %s not found in card library, removing", getCardNameFunc(card)))
			continue
		}

		changed = true
		newCardList = append(newCardList, card)
	}

	return changed, newCardList, nil
}

func validateAndUpdateCardCollectionList(
	ctx contract.StaticContext,
	version string,
	cardCollectionList *zb.CardCollectionList,
	validateMouldId bool,
	saveChangedCollectionFunc func(collection *zb.CardCollectionList) error,
) error {
	var collectionCardsInterface = make([]interface{}, len(cardCollectionList.Cards))
	for i, d := range cardCollectionList.Cards {
		collectionCardsInterface[i] = d
	}
	changed, changedCardCollectionList, err :=
		validateAndUpdateGenericCardList(
			ctx,
			collectionCardsInterface,
			version,
			validateMouldId,
			func(card interface{}) int64 {
				return card.(*zb.CardCollectionCard).MouldId
			},
			func(card interface{}, mouldId int64) {
				card.(*zb.CardCollectionCard).MouldId = mouldId
			},
			func(card interface{}) string {
				return card.(*zb.CardCollectionCard).CardNameDeprecated
			},
			func(card interface{}) {
				card.(*zb.CardCollectionCard).CardNameDeprecated = ""
			},
		)
	if err != nil {
		return err
	}

	if changed {
		cardCollectionList.Cards = make([]*zb.CardCollectionCard, len(changedCardCollectionList))
		for i := range changedCardCollectionList {
			cardCollectionList.Cards[i] = changedCardCollectionList[i].(*zb.CardCollectionCard)
		}

		err = saveChangedCollectionFunc(cardCollectionList)
		if err != nil {
			return err
		}
	}

	return nil
}

func loadDecks(ctx contract.Context, userID string, version string) (*zb.DeckList, error) {
	var deckList zb.DeckList
	err := ctx.Get(DecksKey(userID), &deckList)
	if err != nil {
		if err == contract.ErrNotFound {
			deckList.Decks = []*zb.Deck{}
		} else {
			return nil, err
		}
	}

	deckListChanged := false
	for _, deck := range deckList.Decks {
		// Update data
		var deckCardsInterface = make([]interface{}, len(deck.Cards))
		for i, d := range deck.Cards {
			deckCardsInterface[i] = d
		}
		changed, changedDeckCards, err :=
			validateAndUpdateGenericCardList(
				ctx,
				deckCardsInterface,
				version,
				true,
				func(card interface{}) int64 {
					return card.(*zb.DeckCard).MouldId
				},
				func(card interface{}, mouldId int64) {
					card.(*zb.DeckCard).MouldId = mouldId
				},
				func(card interface{}) string {
					return card.(*zb.DeckCard).CardNameDeprecated
				},
				func(card interface{}) {
					card.(*zb.DeckCard).CardNameDeprecated = ""
				},
			)
		if err != nil {
			return nil, err
		}

		if changed {
			deckListChanged = true
			deck.Cards = make([]*zb.DeckCard, len(changedDeckCards))
			for i := range changedDeckCards {
				deck.Cards[i] = changedDeckCards[i].(*zb.DeckCard)
			}
		}
	}

	if deckListChanged {
		err = saveDecks(ctx, userID, &deckList)
		if err != nil {
			return nil, err
		}
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

func loadOverlords(ctx contract.StaticContext, userID string) (*zb.OverlordList, error) {
	var overlords zb.OverlordList
	err := ctx.Get(OverlordsKey(userID), &overlords)
	if err != nil {
		if err == contract.ErrNotFound {
			overlords.Overlords = []*zb.Overlord{}
		} else {
			return nil, err
		}
	}
	return &overlords, nil
}

func saveOverlords(ctx contract.Context, userID string, overlords *zb.OverlordList) error {
	return ctx.Set(OverlordsKey(userID), overlords)
}

func loadOverlordLevelingData(ctx contract.Context, version string) (*zb.OverlordLevelingData, error) {
	var overlordLevelingData zb.OverlordLevelingData
	if err := ctx.Get(MakeVersionedKey(version, overlordLevelingDataKey), &overlordLevelingData); err != nil {
		return nil, errors.Wrap(err, "error getting overlord leveling data")
	}

	return &overlordLevelingData, nil
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

func getDeckWithRegistrationData(ctx contract.Context, registrationData *zb.PlayerProfileRegistrationData, version string) (*zb.Deck, error) {
	if registrationData.DebugCheats.Enabled && registrationData.DebugCheats.UseCustomDeck {
		return registrationData.DebugCheats.CustomDeck, nil
	}

	// get matched player deck
	matchedDl, err := loadDecks(ctx, registrationData.UserId, version)
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
	if err != nil {
		if err == contract.ErrNotFound {
			pool.PlayerProfiles = []*zb.PlayerProfile{}
		} else {
			// Try to reset the pool
			ctx.Logger().Error("error loading pool, clearing", "key", string(poolKey), "err", err)
			pool = zb.PlayerPool{}
			if err = ctx.Set(poolKey, &pool); err != nil {
				return nil, err
			}

			return &pool, nil
		}
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

func saveUserNotifications(ctx contract.Context, userID string, notifications *zb.NotificationList) error {
	if err := ctx.Set(UserNotificationsKey(userID), notifications); err != nil {
		return err
	}
	return nil
}

func loadUserNotifications(ctx contract.StaticContext, userID string) (*zb.NotificationList, error) {
	var notificationList zb.NotificationList
	err := ctx.Get(UserNotificationsKey(userID), &notificationList)
	if err != nil {
		if err == contract.ErrNotFound {
			notificationList.Notifications = []*zb.Notification{}
		} else {
			return nil, err
		}
	}
	return &notificationList, nil
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
	if err != nil {
		if err == contract.ErrNotFound {
			list.GameModes = []*zb.GameMode{}
		} else {
			return nil, err
		}
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
	cardDetails = proto.Clone(cardDetails).(*zb.Card)
	return &zb.CardInstanceSpecificData{
		Damage:    cardDetails.Damage,
		Defense:   cardDetails.Defense,
		Type:      cardDetails.Type,
		Faction:   cardDetails.Faction,
		Cost:      cardDetails.Cost,
		Abilities: cardDetails.Abilities,
	}
}

func newCardInstanceFromCardDetails(cardDetails *zb.Card, instanceID *zb.InstanceId, owner string, ownerIndex int32) *zb.CardInstance {
	instance := newCardInstanceSpecificDataFromCardDetails(cardDetails)
	var abilities []*zb.CardAbilityInstance
	for _, raw := range cardDetails.Abilities {
		switch raw.Ability {
		case zb.AbilityType_Rage:
			abilities = append(abilities, &zb.CardAbilityInstance{
				IsActive: true,
				Trigger:  raw.Trigger,
				AbilityType: &zb.CardAbilityInstance_Rage{
					Rage: &zb.CardAbilityRage{
						AddedDamage: raw.Value,
					},
				},
			})
		case zb.AbilityType_PriorityAttack:
			abilities = append(abilities, &zb.CardAbilityInstance{
				IsActive: true,
				Trigger:  raw.Trigger,
				AbilityType: &zb.CardAbilityInstance_PriorityAttack{
					PriorityAttack: &zb.CardAbilityPriorityAttack{},
				},
			})
		case zb.AbilityType_ReanimateUnit:
			abilities = append(abilities, &zb.CardAbilityInstance{
				IsActive: true,
				Trigger:  raw.Trigger,
				AbilityType: &zb.CardAbilityInstance_Reanimate{
					Reanimate: &zb.CardAbilityReanimate{
						DefaultDamage:  cardDetails.Damage,
						DefaultDefense: cardDetails.Defense,
					},
				},
			})
		case zb.AbilityType_ChangeStat:
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
		case zb.AbilityType_AttackOverlord:
			abilities = append(abilities, &zb.CardAbilityInstance{
				IsActive: true,
				Trigger:  raw.Trigger,
				AbilityType: &zb.CardAbilityInstance_AttackOverlord{
					AttackOverlord: &zb.CardAbilityAttackOverlord{
						Damage: raw.Value,
					},
				},
			})
		case zb.AbilityType_ReplaceUnitsWithTypeOnStrongerOnes:
			abilities = append(abilities, &zb.CardAbilityInstance{
				IsActive: true,
				Trigger:  raw.Trigger,
				AbilityType: &zb.CardAbilityInstance_ReplaceUnitsWithTypeOnStrongerOnes{
					ReplaceUnitsWithTypeOnStrongerOnes: &zb.CardAbilityReplaceUnitsWithTypeOnStrongerOnes{
						Faction: cardDetails.Faction,
					},
				},
			})
		case zb.AbilityType_DealDamageToThisAndAdjacentUnits:
			abilities = append(abilities, &zb.CardAbilityInstance{
				IsActive: true,
				Trigger:  raw.Trigger,
				AbilityType: &zb.CardAbilityInstance_DealDamageToThisAndAdjacentUnits{
					DealDamageToThisAndAdjacentUnits: &zb.CardAbilityDealDamageToThisAndAdjacentUnits{
						AdjacentDamage: cardDetails.Damage,
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
				cardDetails, err := getCardByMouldId(cardLibrary, cardAmounts.MouldId)
				if err != nil {
					return fmt.Errorf("unable to get card %d from card library: %s", cardAmounts.MouldId, err.Error())
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
			filteredAbilities := make([]*zb.AbilityData, 0, 0)
			for _, ability := range card.Prototype.Abilities {
				switch ability.Ability {
				case zb.AbilityType_Rage:
					fallthrough
				case zb.AbilityType_PriorityAttack:
					fallthrough
				case zb.AbilityType_ReanimateUnit:
					fallthrough
				case zb.AbilityType_ChangeStat:
					fallthrough
				case zb.AbilityType_AttackOverlord:
					fallthrough
				case zb.AbilityType_ReplaceUnitsWithTypeOnStrongerOnes:
					filteredAbilities = append(filteredAbilities, ability)
				default:
					fmt.Printf("Unsupported AbilityType value %s, removed (card '%s')\n", zb.AbilityType_Enum_name[int32(ability.Ability)], card.Prototype.Name)
				}
			}

			card.Prototype.Abilities = filteredAbilities

			switch card.Prototype.Type {
			case zb.CardType_Feral:
				fallthrough
			case zb.CardType_Heavy:
				fmt.Printf("Unsupported CardType value %s, fallback to WALKER (card %s)\n", zb.CardType_Enum_name[int32(card.Prototype.Type)], card.Prototype.Name)
				card.Prototype.Type = zb.CardType_Walker
			}

			switch card.Instance.Type {
			case zb.CardType_Feral:
				fallthrough
			case zb.CardType_Heavy:
				fmt.Printf("Unsupported CardType value %s, fallback to WALKER (card %s)\n", zb.CardType_Enum_name[int32(card.Instance.Type)], card.Prototype.Name)
				card.Instance.Type = zb.CardType_Walker
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

func getCardLibrary(ctx contract.StaticContext, version string) (*zb.CardList, error) {
	var cardList zb.CardList
	if err := ctx.Get(MakeVersionedKey(version, cardListKey), &cardList); err != nil {
		return nil, fmt.Errorf("error getting card library: %s", err)
	}

	return &cardList, nil
}

func getCardByName(cardList *zb.CardList, cardName string) (*zb.Card, error) {
	for _, card := range cardList.Cards {
		if card.Name == cardName {
			return card, nil
		}
	}
	return nil, fmt.Errorf("card with name %s not found in card library", cardName)
}

func getCardByMouldId(cardList *zb.CardList, mouldId int64) (*zb.Card, error) {
	for _, card := range cardList.Cards {
		if card.MouldId == mouldId {
			return card, nil
		}
	}
	return nil, fmt.Errorf("card with mould id %d not found in card library", mouldId)
}
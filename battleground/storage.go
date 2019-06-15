package battleground

import (
	"encoding/json"
	"fmt"
	battleground_proto "github.com/loomnetwork/gamechain/battleground/proto"
	"github.com/loomnetwork/gamechain/types/zb/zb_calls"
	"github.com/loomnetwork/gamechain/types/zb/zb_data"
	"github.com/loomnetwork/gamechain/types/zb/zb_enums"
	loom "github.com/loomnetwork/go-loom"
	"strings"

	"github.com/gogo/protobuf/proto"

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
	overlordPrototypeListKey    = []byte("overlord-prototype-list")
	defaultDecksKey             = []byte("default-deck")
	defaultCollectionKey        = []byte("default-collection")
	matchCountKey               = []byte("match-count")
	playersInMatchmakingListKey = []byte("players-matchmaking")
	gameModeListKey             = []byte("gamemode-list")
	playerPoolKey               = []byte("playerpool")
	taggedPlayerPoolKey         = []byte("tagged-playerpool")
	oracleKey                   = []byte("oracle-key")
	aiDecksKey                  = []byte("ai-decks")
	contractStateKey            = []byte("contract-state")
	contractConfigurationKey    = []byte("contract-configuration")
	nonceKey                    = []byte("nonce")
	currentUserIDUIntKey        = []byte("current-user-id")
	overlordLevelingDataKey     = []byte("overlord-leveling")
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

func OverlordsUserDataKey(userID string) []byte {
	return []byte("user:" + userID + ":overlordsuserdata")
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

func MakeAddressToUserIdKey(address loom.Address) []byte {
	return []byte("address-to-userid:" + string(address.Local.Hex()))
}

func getUserIdByAddress(ctx contract.StaticContext, address loom.Address) (string, error) {
	var userId zb_data.UserIdContainer

	err := ctx.Get(MakeAddressToUserIdKey(address), &userId)
	if err != nil {
		return "", err
	}

	return userId.UserId, nil
}

func setUserIdAddress(ctx contract.Context, userId string, address loom.Address) error {
	userIdContainer := zb_data.UserIdContainer{UserId: userId}
	err := ctx.Set(MakeAddressToUserIdKey(address), &userIdContainer)
	if err != nil {
		return err
	}

	return nil
}

func saveContractState(ctx contract.Context, state *zb_data.ContractState) error {
	if err := ctx.Set(contractStateKey, state); err != nil {
		return err
	}
	return nil
}

func loadContractState(ctx contract.StaticContext) (*zb_data.ContractState, error) {
	var m zb_data.ContractState
	err := ctx.Get(contractStateKey, &m)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get contract state")
	}
	return &m, nil
}

func saveContractConfiguration(ctx contract.Context, state *zb_data.ContractConfiguration) error {
	if err := ctx.Set(contractConfigurationKey, state); err != nil {
		return err
	}
	return nil
}

func loadContractConfiguration(ctx contract.StaticContext) (*zb_data.ContractConfiguration, error) {
	var m zb_data.ContractConfiguration
	err := ctx.Get(contractConfigurationKey, &m)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get contract configuration")
	}

	return &m, nil
}

func RewardTutorialClaimedKey(userID string) []byte {
	return []byte("user:" + userID + ":rewardTutorialClaimed")
}

func loadCardCollectionByUserId(ctx contract.Context, userID string, version string) (*zb_data.CardCollectionList, error) {
	var userCollection zb_data.CardCollectionList
	err := ctx.Get(CardCollectionKey(userID), &userCollection)
	if err != nil {
		if err == contract.ErrNotFound {
			userCollection.Cards = []*zb_data.CardCollectionCard{}
		} else {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	return &userCollection, nil
}

func saveCardCollectionByUserId(ctx contract.Context, userID string, cardCollection *zb_data.CardCollectionList) error {
	return ctx.Set(CardCollectionKey(userID), cardCollection)
}

// loadCardCollectionFromAddress loads address mapping to card collection
func loadCardCollectionByAddress(ctx contract.Context, version string) (*zb_data.CardCollectionList, error) {
	var userCollection zb_data.CardCollectionList
	addr := string(ctx.Message().Sender.Local)
	err := ctx.Get(CardCollectionKey(addr), &userCollection)
	if err != nil {
		if err == contract.ErrNotFound {
			userCollection.Cards = []*zb_data.CardCollectionCard{}
		} else {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	return &userCollection, nil
}

// saveCardCollectionByAddress save card collection using address as a key
func saveCardCollectionByAddress(ctx contract.Context, cardCollection *zb_data.CardCollectionList) error {
	addr := string(ctx.Message().Sender.Local)
	return ctx.Set(CardCollectionKey(addr), cardCollection)
}

func loadDecks(ctx contract.Context, userID string, version string) (*zb_data.DeckList, error) {
	var deckList zb_data.DeckList
	err := ctx.Get(DecksKey(userID), &deckList)
	if err != nil {
		if err == contract.ErrNotFound {
			deckList.Decks = []*zb_data.Deck{}
		} else {
			return nil, err
		}
	}

	deckListChanged, err := fixDeckListCardEditions(ctx, &deckList, version)
	if err != nil {
		return nil, err
	}

	if deckListChanged {
		err = saveDecks(ctx, version, userID, &deckList)
		if err != nil {
			return nil, err
		}
	}

	return &deckList, nil
}

func saveDecks(ctx contract.Context, version string, userID string, decks *zb_data.DeckList) error {
	_, err := fixDeckListCardEditions(ctx, decks, version)
	if err != nil {
		return err
	}

	return ctx.Set(DecksKey(userID), decks)
}

func saveAIDecks(ctx contract.Context, version string, decks *zb_data.AIDeckList) error {
	return ctx.Set(MakeVersionedKey(version, aiDecksKey), decks)
}

func loadAIDecks(ctx contract.StaticContext, version string) (*zb_data.AIDeckList, error) {
	var deckList zb_data.AIDeckList
	err := ctx.Get(MakeVersionedKey(version, aiDecksKey), &deckList)
	if err != nil {
		return nil, err
	}
	return &deckList, nil
}

func loadOverlordPrototypes(ctx contract.StaticContext, version string) (*zb_data.OverlordPrototypeList, error) {
	var overlordPrototypes zb_data.OverlordPrototypeList
	err := ctx.Get(MakeVersionedKey(version, overlordPrototypeListKey), &overlordPrototypes)
	if err != nil {
		return nil, errors.Wrap(err, "error loading overlord prototypes")
	}
	return &overlordPrototypes, nil
}

func saveOverlordPrototypes(ctx contract.Context, version string, overlordPrototypes *zb_data.OverlordPrototypeList) error {
	return ctx.Set(MakeVersionedKey(version, overlordPrototypeListKey), overlordPrototypes)
}

func loadOverlordsUserData(ctx contract.StaticContext, userID string) (*zb_data.OverlordUserDataList, error) {
	var overlordsUserData zb_data.OverlordUserDataList
	err := ctx.Get(OverlordsUserDataKey(userID), &overlordsUserData)
	if err != nil {
		if err == contract.ErrNotFound {
			overlordsUserData.OverlordsUserData = []*zb_data.OverlordUserData{}
		} else {
			return nil, errors.Wrap(err, "error loading overlords user data")
		}
	}
	return &overlordsUserData, nil
}

func saveOverlordsUserData(ctx contract.Context, userID string, overlordsUserData *zb_data.OverlordUserDataList) error {
	return ctx.Set(OverlordsUserDataKey(userID), overlordsUserData)
}

func loadOverlordUserInstances(ctx contract.StaticContext, version string, userID string) ([]*zb_data.OverlordUserInstance, error) {
	prototypeList, err := loadOverlordPrototypes(ctx, version)
	if err != nil {
		return nil, errors.Wrap(err, "error loading overlordUserData user instances")
	}

	userDataList, err := loadOverlordsUserData(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "error loading overlord user instances")
	}

	idToUserData := make(map[int64]*zb_data.OverlordUserData)
	for _, overlordUserData := range userDataList.OverlordsUserData {
		idToUserData[overlordUserData.PrototypeId] = overlordUserData
	}

	var userInstances []*zb_data.OverlordUserInstance
	for _, overlord := range prototypeList.Overlords {
		overlordUserData, exists := idToUserData[overlord.Id]
		if !exists {
			overlordUserData = &zb_data.OverlordUserData{
				PrototypeId: overlord.Id,
				Level: 1,
			}
		}

		userInstances = append(userInstances, &zb_data.OverlordUserInstance{
			Prototype: overlord,
			UserData: overlordUserData,
		})
	}
	return userInstances, nil
}

func getOverlordsUserDataFromOverlordUserInstances(overlordUserInstances []*zb_data.OverlordUserInstance) []*zb_data.OverlordUserData {
	overlordsUserData := make([]*zb_data.OverlordUserData, len(overlordUserInstances), len(overlordUserInstances))
	for index := range overlordUserInstances {
		overlordsUserData[index] = overlordUserInstances[index].UserData
	}

	return overlordsUserData
}

/*
func saveContractState(ctx contract.Context, state *zb.ContractState) error {
	if err := ctx.Set(contractStateKey, state); err != nil {
		return err
	}
	return nil
}

func loadContractState(ctx contract.StaticContext) (*zb.ContractState, error) {
	var m zb.ContractState
	err := ctx.Get(contractStateKey, &m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
 */

func loadOverlordLevelingData(ctx contract.StaticContext, version string) (*zb_data.OverlordLevelingData, error) {
	var overlordLevelingData zb_data.OverlordLevelingData
	if err := ctx.Get(MakeVersionedKey(version, overlordLevelingDataKey), &overlordLevelingData); err != nil {
		return nil, errors.Wrap(err, "error getting overlord leveling data")
	}

	return &overlordLevelingData, nil
}

func saveOverlordLevelingData(ctx contract.Context, version string, overlordLevelingData *zb_data.OverlordLevelingData) error {
	if err := ctx.Set(MakeVersionedKey(version, overlordLevelingDataKey), overlordLevelingData); err != nil {
		return errors.Wrap(err, "error setting overlord leveling data")
	}

	return nil
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

func deleteDeckByID(deckList []*zb_data.Deck, id int64) ([]*zb_data.Deck, bool) {
	newList := make([]*zb_data.Deck, 0)
	for _, deck := range deckList {
		if deck.Id != id {
			newList = append(newList, deck)
		}
	}
	return newList, len(newList) != len(deckList)
}

func getDeckWithRegistrationData(ctx contract.Context, registrationData *zb_data.PlayerProfileRegistrationData, version string) (*zb_data.Deck, error) {
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

func getDeckByID(deckList []*zb_data.Deck, id int64) *zb_data.Deck {
	for _, deck := range deckList {
		if deck.Id == id {
			return deck
		}
	}
	return nil
}

func getDeckByName(deckList []*zb_data.Deck, name string) *zb_data.Deck {
	for _, deck := range deckList {
		if strings.EqualFold(deck.Name, name) {
			return deck
		}
	}
	return nil
}

func copyAccountInfo(account *zb_data.Account, req *zb_calls.UpsertAccountRequest) {
	account.PhoneNumberVerified = req.PhoneNumberVerified
	account.RewardRedeemed = req.RewardRedeemed
	account.IsKickstarter = req.IsKickstarter
	account.Image = req.Image
	account.EmailNotification = req.EmailNotification
	account.EloScore = req.EloScore
	account.CurrentTier = req.CurrentTier
	account.GameMembershipTier = req.GameMembershipTier
}

func savePlayerPool(ctx contract.Context, pool *zb_data.PlayerPool) error {
	return ctx.Set(playerPoolKey, pool)
}

func loadPlayerPool(ctx contract.Context) (*zb_data.PlayerPool, error) {
	return loadPlayerPoolInternal(ctx, playerPoolKey)
}

func saveTaggedPlayerPool(ctx contract.Context, pool *zb_data.PlayerPool) error {
	return ctx.Set(taggedPlayerPoolKey, pool)
}

func loadTaggedPlayerPool(ctx contract.Context) (*zb_data.PlayerPool, error) {
	return loadPlayerPoolInternal(ctx, taggedPlayerPoolKey)
}

func loadPlayerPoolInternal(ctx contract.Context, poolKey []byte) (*zb_data.PlayerPool, error) {
	var pool zb_data.PlayerPool
	err := ctx.Get(poolKey, &pool)
	if err != nil {
		if err == contract.ErrNotFound {
			pool.PlayerProfiles = []*zb_data.PlayerProfile{}
		} else {
			// Try to reset the pool
			ctx.Logger().Error("error loading pool, clearing", "key", string(poolKey), "err", err)
			pool = zb_data.PlayerPool{}
			if err = ctx.Set(poolKey, &pool); err != nil {
				return nil, err
			}

			return &pool, nil
		}
	}
	return &pool, nil
}

func saveMatch(ctx contract.Context, match *zb_data.Match) error {
	if err := ctx.Set(MatchKey(match.Id), match); err != nil {
		return err
	}
	return nil
}

func createMatch(ctx contract.Context, match *zb_data.Match, useBackendGameLogic bool) error {
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

func loadMatch(ctx contract.StaticContext, matchID int64) (*zb_data.Match, error) {
	var m zb_data.Match
	err := ctx.Get(MatchKey(matchID), &m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func nextMatchID(ctx contract.Context) (int64, error) {
	var count zb_data.MatchCount
	err := ctx.Get(matchCountKey, &count)
	if err != nil && err != contract.ErrNotFound {
		return 0, err
	}
	count.CurrentId++
	if err := ctx.Set(matchCountKey, &zb_data.MatchCount{CurrentId: count.CurrentId}); err != nil {
		return 0, err
	}
	return count.CurrentId, nil
}

func saveGameState(ctx contract.Context, gs *zb_data.GameState) error {
	if err := ctx.Set(GameStateKey(gs.Id), gs); err != nil {
		return err
	}
	return nil
}

func loadGameState(ctx contract.StaticContext, id int64) (*zb_data.GameState, error) {
	var state zb_data.GameState
	err := ctx.Get(GameStateKey(id), &state)
	if err != nil {
		return nil, err
	}
	return &state, nil
}

func saveInitialGameState(ctx contract.Context, gs *zb_data.GameState) error {
	if err := ctx.Set(InitialGameStateKey(gs.Id), gs); err != nil {
		return err
	}
	return nil
}

func loadInitialGameState(ctx contract.StaticContext, id int64) (*zb_data.GameState, error) {
	var state zb_data.GameState
	err := ctx.Get(InitialGameStateKey(id), &state)
	if err != nil {
		return nil, err
	}
	return &state, nil
}

func saveUserCurrentMatch(ctx contract.Context, userID string, match *zb_data.Match) error {
	if err := ctx.Set(UserMatchKey(userID), match); err != nil {
		return err
	}
	return nil
}

func loadUserCurrentMatch(ctx contract.StaticContext, userID string) (*zb_data.Match, error) {
	var m zb_data.Match
	err := ctx.Get(UserMatchKey(userID), &m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func saveUserNotifications(ctx contract.Context, userID string, notifications *zb_data.NotificationList) error {
	if err := ctx.Set(UserNotificationsKey(userID), notifications); err != nil {
		return err
	}
	return nil
}

func loadUserNotifications(ctx contract.StaticContext, userID string) (*zb_data.NotificationList, error) {
	var notificationList zb_data.NotificationList
	err := ctx.Get(UserNotificationsKey(userID), &notificationList)
	if err != nil {
		if err == contract.ErrNotFound {
			notificationList.Notifications = []*zb_data.Notification{}
		} else {
			return nil, err
		}
	}
	return &notificationList, nil
}

func addGameModeToList(ctx contract.Context, gameMode *zb_data.GameMode) error {
	gameModeList, err := loadGameModeList(ctx)
	if gameModeList == nil {
		gameModeList = &zb_data.GameModeList{GameModes: []*zb_data.GameMode{}}
	} else if err != nil {
		return err
	}
	gameModeList.GameModes = append(gameModeList.GameModes, gameMode)

	if err = saveGameModeList(ctx, gameModeList); err != nil {
		return err
	}

	return nil
}

func saveGameModeList(ctx contract.Context, gameModeList *zb_data.GameModeList) error {
	if err := ctx.Set(gameModeListKey, gameModeList); err != nil {
		return err
	}

	return nil
}

func loadGameModeList(ctx contract.StaticContext) (*zb_data.GameModeList, error) {
	var list zb_data.GameModeList
	err := ctx.Get(gameModeListKey, &list)
	if err != nil {
		if err == contract.ErrNotFound {
			list.GameModes = []*zb_data.GameMode{}
		} else {
			return nil, err
		}
	}

	return &list, nil
}

func getGameModeFromList(gameModeList *zb_data.GameModeList, ID string) *zb_data.GameMode {
	for _, gameMode := range gameModeList.GameModes {
		if gameMode.ID == ID {
			return gameMode
		}
	}

	return nil
}

func getGameModeFromListByName(gameModeList *zb_data.GameModeList, name string) *zb_data.GameMode {
	for _, gameMode := range gameModeList.GameModes {
		if gameMode.Name == name {
			return gameMode
		}
	}

	return nil
}

func deleteGameMode(gameModeList *zb_data.GameModeList, ID string) (*zb_data.GameModeList, bool) {
	newList := make([]*zb_data.GameMode, 0)
	for _, gameMode := range gameModeList.GameModes {
		if gameMode.ID != ID {
			newList = append(newList, gameMode)
		}
	}

	return &zb_data.GameModeList{GameModes: newList}, len(newList) != len(gameModeList.GameModes)
}

func cardKeyFromCardTokenId(cardTokenId int64) battleground_proto.CardKey {
	return battleground_proto.CardKey{
		MouldId: cardTokenId / 10,
		Variant: zb_enums.CardVariant_Enum(cardTokenId % 10),
	}
}

func newCardInstanceSpecificDataFromCardDetails(cardDetails *zb_data.Card) *zb_data.CardInstanceSpecificData {
	cardDetails = proto.Clone(cardDetails).(*zb_data.Card)
	return &zb_data.CardInstanceSpecificData{
		Damage:    cardDetails.Damage,
		Defense:   cardDetails.Defense,
		Type:      cardDetails.Type,
		Faction:   cardDetails.Faction,
		Cost:      cardDetails.Cost,
		Abilities: cardDetails.Abilities,
	}
}

func newCardInstanceFromCardDetails(cardDetails *zb_data.Card, instanceID *zb_data.InstanceId, owner string, ownerIndex int32) *zb_data.CardInstance {
	instance := newCardInstanceSpecificDataFromCardDetails(cardDetails)
	var abilities []*zb_data.CardAbilityInstance
	for _, raw := range cardDetails.Abilities {
		switch raw.Ability {
		case zb_enums.AbilityType_Rage:
			abilities = append(abilities, &zb_data.CardAbilityInstance{
				IsActive: true,
				Trigger:  raw.Trigger,
				AbilityType: &zb_data.CardAbilityInstance_Rage{
					Rage: &zb_data.CardAbilityRage{
						AddedDamage: raw.Value,
					},
				},
			})
		case zb_enums.AbilityType_PriorityAttack:
			abilities = append(abilities, &zb_data.CardAbilityInstance{
				IsActive: true,
				Trigger:  raw.Trigger,
				AbilityType: &zb_data.CardAbilityInstance_PriorityAttack{
					PriorityAttack: &zb_data.CardAbilityPriorityAttack{},
				},
			})
		case zb_enums.AbilityType_ReanimateUnit:
			abilities = append(abilities, &zb_data.CardAbilityInstance{
				IsActive: true,
				Trigger:  raw.Trigger,
				AbilityType: &zb_data.CardAbilityInstance_Reanimate{
					Reanimate: &zb_data.CardAbilityReanimate{
						DefaultDamage:  cardDetails.Damage,
						DefaultDefense: cardDetails.Defense,
					},
				},
			})
		case zb_enums.AbilityType_ChangeStat:
			abilities = append(abilities, &zb_data.CardAbilityInstance{
				IsActive: true,
				Trigger:  raw.Trigger,
				AbilityType: &zb_data.CardAbilityInstance_ChangeStat{
					ChangeStat: &zb_data.CardAbilityChangeStat{
						StatAdjustment: raw.Value,
						Stat:           raw.Stat,
					},
				},
			})
		case zb_enums.AbilityType_AttackOverlord:
			abilities = append(abilities, &zb_data.CardAbilityInstance{
				IsActive: true,
				Trigger:  raw.Trigger,
				AbilityType: &zb_data.CardAbilityInstance_AttackOverlord{
					AttackOverlord: &zb_data.CardAbilityAttackOverlord{
						Damage: raw.Value,
					},
				},
			})
		case zb_enums.AbilityType_ReplaceUnitsWithTypeOnStrongerOnes:
			abilities = append(abilities, &zb_data.CardAbilityInstance{
				IsActive: true,
				Trigger:  raw.Trigger,
				AbilityType: &zb_data.CardAbilityInstance_ReplaceUnitsWithTypeOnStrongerOnes{
					ReplaceUnitsWithTypeOnStrongerOnes: &zb_data.CardAbilityReplaceUnitsWithTypeOnStrongerOnes{
						Faction: cardDetails.Faction,
					},
				},
			})
		case zb_enums.AbilityType_DealDamageToThisAndAdjacentUnits:
			abilities = append(abilities, &zb_data.CardAbilityInstance{
				IsActive: true,
				Trigger:  raw.Trigger,
				AbilityType: &zb_data.CardAbilityInstance_DealDamageToThisAndAdjacentUnits{
					DealDamageToThisAndAdjacentUnits: &zb_data.CardAbilityDealDamageToThisAndAdjacentUnits{
						AdjacentDamage: cardDetails.Damage,
					},
				},
			})
		}

	}
	return &zb_data.CardInstance{
		InstanceId:         proto.Clone(instanceID).(*zb_data.InstanceId),
		Owner:              owner,
		Prototype:          proto.Clone(cardDetails).(*zb_data.Card),
		Instance:           instance,
		AbilitiesInstances: abilities,
		Zone:               zb_enums.Zone_DECK, // default to deck
		OwnerIndex:         ownerIndex,
	}
}

func getInstanceIdsFromCardInstances(cards []*zb_data.CardInstance) []*zb_data.InstanceId {
	var instanceIds = make([]*zb_data.InstanceId, len(cards), len(cards))
	for i := range cards {
		instanceIds[i] = cards[i].InstanceId
	}

	return instanceIds
}

func populateDeckCards(cardLibrary *zb_data.CardList, playerStates []*zb_data.PlayerState, useBackendGameLogic bool) error {
	for playerIndex, playerState := range playerStates {
		deck := playerState.Deck
		if deck == nil {
			return fmt.Errorf("no card deck fro player %s", playerState.Id)
		}
		for _, cardAmounts := range deck.Cards {
			for i := int64(0); i < cardAmounts.Amount; i++ {
				cardDetails, err := getCardByCardKey(cardLibrary, cardAmounts.CardKey)
				if err != nil {
					return fmt.Errorf("unable to get card [%s] from card library: %s", cardAmounts.CardKey.String(), err.Error())
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

func removeUnsupportedCardFeatures(useBackendGameLogic bool, playerStates []*zb_data.PlayerState) {
	if !useBackendGameLogic {
		return
	}

	for _, playerState := range playerStates {
		filteredCards := make([]*zb_data.CardInstance, 0, 0)

		for _, card := range playerState.CardsInDeck {
			filteredAbilities := make([]*zb_data.AbilityData, 0, 0)
			for _, ability := range card.Prototype.Abilities {
				switch ability.Ability {
				case zb_enums.AbilityType_Rage:
					fallthrough
				case zb_enums.AbilityType_PriorityAttack:
					fallthrough
				case zb_enums.AbilityType_ReanimateUnit:
					fallthrough
				case zb_enums.AbilityType_ChangeStat:
					fallthrough
				case zb_enums.AbilityType_AttackOverlord:
					fallthrough
				case zb_enums.AbilityType_ReplaceUnitsWithTypeOnStrongerOnes:
					filteredAbilities = append(filteredAbilities, ability)
				default:
					fmt.Printf("Unsupported AbilityType value %s, removed (card '%s')\n", zb_enums.AbilityType_Enum_name[int32(ability.Ability)], card.Prototype.Name)
				}
			}

			card.Prototype.Abilities = filteredAbilities

			switch card.Prototype.Type {
			case zb_enums.CardType_Feral:
				fallthrough
			case zb_enums.CardType_Heavy:
				fmt.Printf("Unsupported CardType value %s, fallback to WALKER (card %s)\n", zb_enums.CardType_Enum_name[int32(card.Prototype.Type)], card.Prototype.Name)
				card.Prototype.Type = zb_enums.CardType_Walker
			}

			switch card.Instance.Type {
			case zb_enums.CardType_Feral:
				fallthrough
			case zb_enums.CardType_Heavy:
				fmt.Printf("Unsupported CardType value %s, fallback to WALKER (card %s)\n", zb_enums.CardType_Enum_name[int32(card.Instance.Type)], card.Prototype.Name)
				card.Instance.Type = zb_enums.CardType_Walker
			}

			switch card.Prototype.Kind {
			case zb_enums.CardKind_Creature:
				filteredCards = append(filteredCards, card)
			default:
				fmt.Printf("Unsupported CardKind value %s, removed (card '%s')\n", zb_enums.CardKind_Enum_name[int32(card.Prototype.Kind)], card.Prototype.Name)
			}

			switch card.Prototype.Rank {
			case zb_enums.CreatureRank_Officer:
				fallthrough
			case zb_enums.CreatureRank_Commander:
				fallthrough
			case zb_enums.CreatureRank_General:
				fmt.Printf("Unsupported CreatureRank value %s, fallback to MINION (card %s)\n", zb_enums.CreatureRank_Enum_name[int32(card.Prototype.Rank)], card.Prototype.Name)
				card.Prototype.Rank = zb_enums.CreatureRank_Minion
			}
		}

		playerState.CardsInDeck = filteredCards
	}
}

func getCardLibrary(ctx contract.StaticContext, version string) (*zb_data.CardList, error) {
	var cardList zb_data.CardList
	if err := ctx.Get(MakeVersionedKey(version, cardListKey), &cardList); err != nil {
		return nil, fmt.Errorf("error getting card library: %s", err)
	}

	mouldIdToCard, err := getCardKeyToCardMap(cardList.Cards)
	if err != nil {
		return nil, err
	}

	for _, card := range cardList.Cards {
		err = applySourceMouldIdAndOverrides(card, mouldIdToCard)
		if err != nil {
			return nil, err
		}
	}

	return &cardList, nil
}

func getCardByCardKey(cardList *zb_data.CardList, cardKey battleground_proto.CardKey) (*zb_data.Card, error) {
	for _, card := range cardList.Cards {
		if card.CardKey == cardKey {
			return card, nil
		}
	}
	return nil, fmt.Errorf("card with card key [%s] not found in card library", cardKey.String())
}

func fixDeckListCardEditions(ctx contract.Context, deckList *zb_data.DeckList, version string) (changed bool, err error) {
	cardLibrary, err := getCardLibrary(ctx, version)
	if err != nil {
		return false, err
	}

	cardKeyToCardMap, err := getCardKeyToCardMap(cardLibrary.Cards)
	if err != nil {
		return false, err
	}

	changed = false
	for _, deck := range deckList.Decks {
		if fixDeckCardEditions(deck, cardKeyToCardMap) {
			changed = true
		}
	}

	return changed, nil
}

func fixDeckCardEditions(deck *zb_data.Deck, cardKeyToCardMap map[battleground_proto.CardKey]*zb_data.Card) (changed bool) {
	var newDeckCards = make([]*zb_data.DeckCard, 0)
	var cardKeyToDeckCard = make(map[battleground_proto.CardKey]*zb_data.DeckCard)
	for _, deckCard := range deck.Cards {
		cardKeyToDeckCard[deckCard.CardKey] = deckCard
	}

	for _, deckCard := range deck.Cards {
		// Check if this specific edition of a card exists in card library
		_, editionExists := cardKeyToCardMap[deckCard.CardKey]
		if !editionExists {
			// If this edition is not in card library, try to fallback to Normal edition
			normalEditionCardKey := battleground_proto.CardKey{
				MouldId: deckCard.CardKey.MouldId,
				Variant: zb_enums.CardVariant_Standard,
			}

			_, normalEditionExists := cardKeyToCardMap[normalEditionCardKey]

			// If normal edition doesn't exist in card library too, just remove card from the deck completely
			if !normalEditionExists {
				changed = true
			} else {
				normalEditionDeckCard, normalEditionDeckCardExists := cardKeyToDeckCard[normalEditionCardKey]
				// If normal edition exists in card library AND in the deck,
				// add the amount of the special edition to normal edition
				if normalEditionDeckCardExists {
					normalEditionDeckCard.Amount += deckCard.Amount
					changed = true
				} else {
					// If normal edition exists in card library, but NOT in the deck,
					// create a normal edition deck card and add special edition amount to it
					normalEditionDeckCard = &zb_data.DeckCard{
						CardKey: normalEditionCardKey,
						Amount: deckCard.Amount,
					}

					newDeckCards = append(newDeckCards, normalEditionDeckCard)
					cardKeyToDeckCard[normalEditionCardKey] = normalEditionDeckCard
					changed = true
				}
			}
		} else {
			newDeckCards = append(newDeckCards, deckCard)
		}
	}

	deck.Cards = newDeckCards
	return changed
}
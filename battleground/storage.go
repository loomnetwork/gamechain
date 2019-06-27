package battleground

import (
	"encoding/json"
	"fmt"
	battleground_proto "github.com/loomnetwork/gamechain/battleground/proto"
	"github.com/loomnetwork/gamechain/types/oracle"
	"github.com/loomnetwork/gamechain/types/zb/zb_calls"
	"github.com/loomnetwork/gamechain/types/zb/zb_data"
	"github.com/loomnetwork/gamechain/types/zb/zb_enums"
	"github.com/loomnetwork/go-loom"
	"math/big"
	"strconv"
	"strings"

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

	cardLibraryKey              = []byte("cardlist")
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
	oracleCommandRequestListKey = []byte("oracle-command-request-list")
)

var (
	ErrNotFound                    = errors.New("not found")
	ErrUserNotVerified             = errors.New("user is not verified")
	ErrNotOwnerOrOracleNotVerified = errors.New("sender is not user owner or oracle")
)

type cardKeyToAmountChangeMap map[battleground_proto.CardKey]int64
type addressToCardKeyToAmountChangesMap map[string]map[battleground_proto.CardKey]int64

const defaultUserIdPrefix = "ZombieSlayer_"

func AccountKey(userID string) []byte {
	return []byte("user:" + userID)
}

func UserPersistentDataKey(userID string) []byte {
	return []byte("user:" + userID + ":persistent-data")
}

func DecksKey(userID string) []byte {
	return []byte("user:" + userID + ":decks")
}

func PendingMintingTransactionReceiptCollectionKey(userID string) []byte {
	return []byte("user:" + userID + ":pending-minting-transaction-receipts")
}

func AddressToPendingCardAmountChangeItemsKey(address loom.Address) []byte {
	return []byte("address" + string(address.Local) + ":address-to-pending-card-amount-change-items")
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

func MakeAddressToUserIdLinkKey(address loom.Address) []byte {
	return []byte("address-to-user-id:" + string(address.Local))
}

func loadPendingCardAmountChangesContainerByAddress(ctx contract.StaticContext, address loom.Address) (*zb_data.CardAmountChangeItemsContainer, error) {
	var container zb_data.CardAmountChangeItemsContainer

	err := ctx.Get(AddressToPendingCardAmountChangeItemsKey(address), &container)
	if err != nil {
		if err == contract.ErrNotFound {
			container = zb_data.CardAmountChangeItemsContainer{
				CardAmountChange: []*zb_data.CardAmountChangeItem{},
			}
		} else {
			return nil, err
		}
	}

	return &container, nil
}

func savePendingCardAmountChangeItemsContainerByAddress(ctx contract.Context, address loom.Address, container *zb_data.CardAmountChangeItemsContainer) error {
	if err := ctx.Set(AddressToPendingCardAmountChangeItemsKey(address), container); err != nil {
		return errors.Wrap(err, "error saving PendingCardAmountChangeItemsContainerByAddress")
	}

	return nil
}

func loadUserIdByAddress(ctx contract.StaticContext, address loom.Address) (string, error) {
	var userId zb_data.UserIdContainer
	if err := ctx.Get(MakeAddressToUserIdLinkKey(address), &userId); err != nil {
		return "", errors.Wrap(err, "error loading user id by address")
	}

	return userId.UserId, nil
}

func saveAddressToUserIdLink(ctx contract.Context, userId string, address loom.Address) error {
	userIdContainer := zb_data.UserIdContainer{UserId: userId}
	if err := ctx.Set(MakeAddressToUserIdLinkKey(address), &userIdContainer); err != nil {
		return errors.Wrap(err, "error saving user id for address")
	}

	return nil
}

func saveOracleCommandRequestList(ctx contract.Context, commandRequestList *oracle.OracleCommandRequestList) error {
	if err := ctx.Set(oracleCommandRequestListKey, commandRequestList); err != nil {
		return errors.Wrap(err, "error saving oracle command request list")
	}
	return nil
}

func loadOracleCommandRequestList(ctx contract.StaticContext) (*oracle.OracleCommandRequestList, error) {
	var commandRequestList oracle.OracleCommandRequestList
	if err := ctx.Get(oracleCommandRequestListKey, &commandRequestList); err != nil {
		if err == contract.ErrNotFound {
			commandRequestList.Commands = []*oracle.OracleCommandRequest{}
		} else {
			return nil, errors.Wrap(err, "error loading oracle command request list")
		}
	}
	return &commandRequestList, nil
}

func saveContractState(ctx contract.Context, state *zb_data.ContractState) error {
	if err := ctx.Set(contractStateKey, state); err != nil {
		return errors.Wrap(err, "error saving contract state")
	}
	return nil
}

func loadContractState(ctx contract.StaticContext) (*zb_data.ContractState, error) {
	var state zb_data.ContractState
	if err := ctx.Get(contractStateKey, &state); err != nil {
		return nil, errors.Wrap(err, "error loading contract state")

	}
	return &state, nil
}

func loadContractConfiguration(ctx contract.StaticContext) (*zb_data.ContractConfiguration, error) {
	var configuration zb_data.ContractConfiguration
	if err := ctx.Get(contractConfigurationKey, &configuration); err != nil {
		return nil, errors.Wrap(err, "error loading contract configuration")
	}

	return &configuration, nil
}

func saveContractConfiguration(ctx contract.Context, state *zb_data.ContractConfiguration) error {
	if err := ctx.Set(contractConfigurationKey, state); err != nil {
		return errors.Wrap(err, "error saving contract configuration")
	}
	return nil
}

func loadCardCollectionRaw(ctx contract.Context, userID string) (*zb_data.CardCollectionList, error) {
	var userCollection zb_data.CardCollectionList
	if err := ctx.Get(CardCollectionKey(userID), &userCollection); err != nil {
		if err == contract.ErrNotFound {
			userCollection.Cards = []*zb_data.CardCollectionCard{}
		} else {
			return nil, errors.Wrap(err, "error loading card collection")
		}
	}

	return &userCollection, nil
}

func saveCardCollection(ctx contract.Context, userID string, cardCollection *zb_data.CardCollectionList) error {
	if err := ctx.Set(CardCollectionKey(userID), cardCollection); err != nil {
		return errors.Wrap(err, "error saving card collection")
	}
	return nil
}

func loadPendingMintingTransactionReceipts(ctx contract.StaticContext, userID string) (*zb_data.MintingTransactionReceiptCollection, error) {
	var receiptCollection zb_data.MintingTransactionReceiptCollection
	if err := ctx.Get(PendingMintingTransactionReceiptCollectionKey(userID), &receiptCollection); err != nil {
		if err == contract.ErrNotFound {
			receiptCollection.Receipts = []*zb_data.MintingTransactionReceipt{}
		} else {
			return nil, errors.Wrap(err, "failed to load minting transaction receipts")
		}
	}

	return &receiptCollection, nil
}

func loadUserPersistentData(ctx contract.StaticContext, userID string) (*zb_data.UserPersistentData, error) {
	var userPersistentData zb_data.UserPersistentData
	if err := ctx.Get(UserPersistentDataKey(userID), &userPersistentData); err != nil {
		if err == contract.ErrNotFound {
			userPersistentData.ExecutedDataWipesVersions = []string{}
		} else {
			return nil, errors.Wrap(err, "error loading user's persistent data")
		}
	}
	return &userPersistentData, nil
}

func saveUserPersistentData(ctx contract.Context, userID string, userPersistentData *zb_data.UserPersistentData) error {
	if err := ctx.Set(UserPersistentDataKey(userID), userPersistentData); err != nil {
		return errors.Wrap(err, "error saving user's persistent data")
	}
	return nil
}

func savePendingMintingTransactionReceipts(ctx contract.Context, userID string, receiptCollection *zb_data.MintingTransactionReceiptCollection) error {
	if err := ctx.Set(PendingMintingTransactionReceiptCollectionKey(userID), receiptCollection); err != nil {
		return errors.Wrap(err, "failed to save minting transaction receipts")
	}

	return nil
}

func saveDecks(ctx contract.Context, version string, userID string, decks *zb_data.DeckList) error {
	_, err := fixDeckListCardVariants(ctx, decks, version)
	if err != nil {
		return errors.Wrap(err, "error saving decks")
	}

	if err = ctx.Set(DecksKey(userID), decks); err != nil {
		return errors.Wrap(err, "error saving decks")
	}
	return nil
}

func loadDecks(ctx contract.Context, userID string, version string) (*zb_data.DeckList, error) {
	var deckList zb_data.DeckList
	if err := ctx.Get(DecksKey(userID), &deckList); err != nil {
		if err == contract.ErrNotFound {
			deckList.Decks = []*zb_data.Deck{}
		} else {
			return nil, errors.Wrap(err, "error loading decks")
		}
	}

	deckListChanged, err := fixDeckListCardVariants(ctx, &deckList, version)
	if err != nil {
		return nil, errors.Wrap(err, "error loading decks")
	}

	if deckListChanged {
		err = saveDecks(ctx, version, userID, &deckList)
		if err != nil {
			return nil, errors.Wrap(err, "error loading decks")
		}
	}

	return &deckList, nil
}

func saveAIDecks(ctx contract.Context, version string, decks *zb_data.AIDeckList) error {
	if err := ctx.Set(MakeVersionedKey(version, aiDecksKey), decks); err != nil {
		return errors.Wrap(err, "error saving AI decks")
	}
	return nil
}

func loadAIDecks(ctx contract.StaticContext, version string) (*zb_data.AIDeckList, error) {
	var deckList zb_data.AIDeckList
	if err := ctx.Get(MakeVersionedKey(version, aiDecksKey), &deckList); err != nil {
		return nil, err
	}
	return &deckList, nil
}

func loadOverlordPrototypes(ctx contract.StaticContext, version string) (*zb_data.OverlordPrototypeList, error) {
	var overlordPrototypes zb_data.OverlordPrototypeList
	if err := ctx.Get(MakeVersionedKey(version, overlordPrototypeListKey), &overlordPrototypes); err != nil {
		return nil, errors.Wrap(err, "error loading overlord prototypes")
	}
	return &overlordPrototypes, nil
}

func saveOverlordPrototypes(ctx contract.Context, version string, overlordPrototypes *zb_data.OverlordPrototypeList) error {
	if err := ctx.Set(MakeVersionedKey(version, overlordPrototypeListKey), overlordPrototypes); err != nil {
		return errors.Wrap(err, "error saving overlord prototypes")
	}
	return nil
}

func loadOverlordUserDataList(ctx contract.StaticContext, userID string) (*zb_data.OverlordUserDataList, error) {
	var overlordsUserData zb_data.OverlordUserDataList
	if err := ctx.Get(OverlordsUserDataKey(userID), &overlordsUserData); err != nil {
		if err == contract.ErrNotFound {
			overlordsUserData.OverlordsUserData = []*zb_data.OverlordUserData{}
		} else {
			return nil, errors.Wrap(err, "error loading overlords user data list")
		}
	}
	return &overlordsUserData, nil
}

func saveOverlordUserDataList(ctx contract.Context, userID string, overlordsUserData *zb_data.OverlordUserDataList) error {
	if err := ctx.Set(OverlordsUserDataKey(userID), overlordsUserData); err != nil {
		return errors.Wrap(err, "error saving overlord user data lis")
	}
	return nil
}

func loadOverlordUserInstances(ctx contract.StaticContext, version string, userID string) ([]*zb_data.OverlordUserInstance, error) {
	prototypeList, err := loadOverlordPrototypes(ctx, version)
	if err != nil {
		return nil, errors.Wrap(err, "error loading overlord user instances")
	}

	userDataList, err := loadOverlordUserDataList(ctx, userID)
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

func isOwner(ctx contract.StaticContext, userID string) bool {
	ok, _ := ctx.HasPermissionFor(ctx.Message().Sender, []byte(userID), []string{OwnerRole})
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

func loadCardLibraryRaw(ctx contract.StaticContext, version string) (*zb_data.CardList, error) {
	var cardList zb_data.CardList
	if err := ctx.Get(MakeVersionedKey(version, cardLibraryKey), &cardList); err != nil {
		return nil, errors.Wrap(err, "error loading card library")
	}

	return &cardList, nil
}

func loadCardLibrary(ctx contract.StaticContext, version string) (*zb_data.CardList, error) {
	cardList, err := loadCardLibraryRaw(ctx, version)
	if err != nil {
		return nil, err
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

	return cardList, nil
}

func saveCardLibrary(ctx contract.Context, version string, cardList *zb_data.CardList) error {
	if err := ctx.Set(MakeVersionedKey(version, cardLibraryKey), cardList); err != nil {
		return errors.Wrap(err, "error saving card library")
	}

	return nil
}

func loadDefaultCardCollection(ctx contract.StaticContext, version string) (*zb_data.CardCollectionList, error) {
	var cardCollectionList zb_data.CardCollectionList
	if err := ctx.Get(MakeVersionedKey(version, defaultCollectionKey), &cardCollectionList); err != nil {
		return nil, errors.Wrap(err, "error loading default card collection")
	}

	return &cardCollectionList, nil
}

func saveDefaultCardCollection(ctx contract.Context, version string, cardCollectionList *zb_data.CardCollectionList) error {
	if err := ctx.Set(MakeVersionedKey(version, defaultCollectionKey), cardCollectionList); err != nil {
		return errors.Wrap(err, "error saving default card collection")
	}

	return nil
}

func loadDefaultDecks(ctx contract.StaticContext, version string) (*zb_data.DeckList, error) {
	var deckList zb_data.DeckList
	if err := ctx.Get(MakeVersionedKey(version, defaultDecksKey), &deckList); err != nil {
		return nil, errors.Wrap(err, "error loading default decks")
	}

	return &deckList, nil
}

func saveDefaultDecks(ctx contract.Context, version string, deckList *zb_data.DeckList) error {
	if err := ctx.Set(MakeVersionedKey(version, defaultDecksKey), deckList); err != nil {
		return errors.Wrap(err, "error saving default decks")
	}

	return nil
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

func getCardByCardKey(cardList *zb_data.CardList, cardKey battleground_proto.CardKey) (*zb_data.Card, error) {
	for _, card := range cardList.Cards {
		if card.CardKey == cardKey {
			return card, nil
		}
	}
	return nil, fmt.Errorf("card with card key [%s] not found in card library", cardKey.String())
}

func fixDeckListCardVariants(ctx contract.Context, deckList *zb_data.DeckList, version string) (changed bool, err error) {
	cardLibrary, err := loadCardLibrary(ctx, version)
	if err != nil {
		return false, errors.Wrap(err, "error fixing card variants")
	}

	cardKeyToCardMap, err := getCardKeyToCardMap(cardLibrary.Cards)
	if err != nil {
		return false, errors.Wrap(err, "error fixing card variants")
	}

	changed = false
	for _, deck := range deckList.Decks {
		if fixDeckCardVariants(deck, cardKeyToCardMap) {
			changed = true
		}
	}

	return changed, nil
}

func fixDeckCardVariants(deck *zb_data.Deck, cardKeyToCardMap map[battleground_proto.CardKey]*zb_data.Card) (changed bool) {
	var newDeckCards = make([]*zb_data.DeckCard, 0)
	var cardKeyToDeckCard = make(map[battleground_proto.CardKey]*zb_data.DeckCard)
	for _, deckCard := range deck.Cards {
		cardKeyToDeckCard[deckCard.CardKey] = deckCard
	}

	for _, deckCard := range deck.Cards {
		// Check if this specific variant of a card exists in card library
		_, variantExists := cardKeyToCardMap[deckCard.CardKey]
		if !variantExists {
			// If this variant is not in card library, try to fallback to Normal variant
			normalVariantCardKey := battleground_proto.CardKey{
				MouldId: deckCard.CardKey.MouldId,
				Variant: zb_enums.CardVariant_Standard,
			}

			_, normalVariantExists := cardKeyToCardMap[normalVariantCardKey]

			// If normal variant doesn't exist in card library too, just remove card from the deck completely
			if !normalVariantExists {
				changed = true
			} else {
				normalVariantDeckCard, normalVariantDeckCardExists := cardKeyToDeckCard[normalVariantCardKey]
				// If normal variant exists in card library AND in the deck,
				// add the amount of the special variant to normal variant
				if normalVariantDeckCardExists {
					normalVariantDeckCard.Amount += deckCard.Amount
					changed = true
				} else {
					// If normal variant exists in card library, but NOT in the deck,
					// create a normal variant deck card and add special variant amount to it
					normalVariantDeckCard = &zb_data.DeckCard{
						CardKey: normalVariantCardKey,
						Amount: deckCard.Amount,
					}

					newDeckCards = append(newDeckCards, normalVariantDeckCard)
					cardKeyToDeckCard[normalVariantCardKey] = normalVariantDeckCard
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

func parseUserIdToNumber(userIdString string) *big.Int {
	userIdNumberString := strings.Replace(userIdString, defaultUserIdPrefix, "", -1)
	userIdNumber, err := strconv.ParseInt(userIdNumberString, 10, 64)
	if err != nil {
		userIdNumber = 0
	}

	return big.NewInt(userIdNumber)
}
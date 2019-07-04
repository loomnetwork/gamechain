package battleground

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/loomnetwork/gamechain/tools/battleground_utility"
	"github.com/loomnetwork/gamechain/types/zb/zb_calls"
	"github.com/loomnetwork/gamechain/types/zb/zb_data"
	"github.com/loomnetwork/gamechain/types/zb/zb_enums"
	"math/big"
	"os"
	"sort"
	"time"
	"unicode/utf8"

	"github.com/gogo/protobuf/proto"
	orctype "github.com/loomnetwork/gamechain/types/oracle"
	"github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/plugin"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	"github.com/loomnetwork/go-loom/types"
	"github.com/pkg/errors"
)

const (
	MaxGameModeNameChar        = 48
	MaxGameModeDescriptionChar = 255
	MaxGameModeVersionChar     = 16
	TurnTimeout                = 120 * time.Second
	KeepAliveTimeout           = 60 * time.Second // client keeps sending keepalive every 30 second. have to make sure we have some buffer for network delays
)

const (
	TopicCreateAccountEvent      = "createaccount"
	TopicUpdateAccountEvent      = "updateaccount"
	TopicCreateDeckEvent         = "createdeck"
	TopicEditDeckEvent           = "editdeck"
	TopicDeleteDeckEvent         = "deletedeck"
	TopicAddOverlordExpEvent     = "addheroexperience"
	TopicRegisterPlayerPoolEvent = "registerplayerpool"
	TopicFindMatchEvent          = "findmatch"
	TopicAcceptMatchEvent        = "acceptmatch"
	// match pattern match:id e.g. match:1, match:2, ...
	TopicMatchEventPrefix      = "match:"
	TopicUserEventPrefix       = "user:"
)

const (
	OracleRole = "oracle"
	OwnerRole  = "user"
)

var (
	// secret
	secret                             string
	_, debugEnabled                    = os.LookupEnv("RL_DEBUG")
	purchaseGatewayPrivateKeyHexString = os.Getenv("RL_PURCHASE_GATEWAY_PRIVATE_KEY")
	// Error list
	ErrOracleNotSpecified = errors.New("oracle not specified")
	ErrOracleNotVerified  = errors.New("oracle not verified")
	ErrInvalidEventBatch  = errors.New("invalid event batch")
	ErrVersionNotSet      = errors.New("data version not set")
	ErrDebugNotEnabled    = errors.New("debug mode not enabled")
)

type ZombieBattleground struct {
}

func (z *ZombieBattleground) Meta() (plugin.Meta, error) {
	return plugin.Meta{
		Name:    "ZombieBattleground",
		Version: "1.0.0",
	}, nil
}

func (z *ZombieBattleground) Init(ctx contract.Context, req *zb_calls.InitRequest) error {
	secret = os.Getenv("SECRET_KEY")
	if secret == "" {
		secret = "justsowecantestwithoutenvvar"
	}

	if req.Oracle != nil {
		ctx.GrantPermissionTo(loom.UnmarshalAddressPB(req.Oracle), []byte(req.Oracle.String()), OracleRole)
		if err := ctx.Set(oracleKey, req.Oracle); err != nil {
			return errors.Wrap(err, "error setting oracle")
		}
	}

	// initialize card library
	cardLibrary := zb_data.CardList{
		Cards: req.Cards,
	}

	if err := saveCardLibrary(ctx, req.Version, &cardLibrary); err != nil {
		return err
	}

	// initialize overlords
	overlordPrototypeList := zb_data.OverlordPrototypeList{
		Overlords: req.Overlords,
	}
	if err := saveOverlordPrototypes(ctx, req.Version, &overlordPrototypeList); err != nil {
		return err
	}

	// initialize card collection
	defaultCardCollection := zb_data.CardCollectionList{
		Cards: req.DefaultCollection,
	}
	if err := saveDefaultCardCollection(ctx, req.Version, &defaultCardCollection); err != nil {
		return err
	}

	// initialize default deck
	defaultDecks := zb_data.DeckList{
		Decks: req.DefaultDecks,
	}
	if err := saveDefaultDecks(ctx, req.Version, &defaultDecks); err != nil {
		return err
	}

	// initialize AI decks
	aiDecks := zb_data.AIDeckList{
		Decks: req.AiDecks,
	}

	if err := saveAIDecks(ctx, req.Version, &aiDecks); err != nil {
		return err
	}

	// initialize overlord leveling
	overlordLevelingData := req.OverlordLeveling
	if overlordLevelingData == nil {
		overlordLevelingData = &zb_data.OverlordLevelingData{}
	}
	if err := saveOverlordLevelingData(ctx, req.Version, overlordLevelingData); err != nil {
		return err
	}

	return nil
}

func (z *ZombieBattleground) UpdateInit(ctx contract.Context, req *zb_calls.UpdateInitRequest) error {
	initData := req.InitData

	var overlordPrototypeList zb_data.OverlordPrototypeList
	var cardList zb_data.CardList
	var defaultCardCollectionList zb_data.CardCollectionList
	var defaultDecks zb_data.DeckList
	var aiDeckList zb_data.AIDeckList
	var overlordLevelingData *zb_data.OverlordLevelingData

	// load data
	// card library
	cardList.Cards = initData.Cards
	if cardList.Cards == nil {
		return fmt.Errorf("'cards' key missing")
	}

	// overlords
	overlordPrototypeList.Overlords = initData.Overlords
	if overlordPrototypeList.Overlords == nil {
		return fmt.Errorf("'overlords' key missing")
	}

	// default collection
	defaultCardCollectionList.Cards = initData.DefaultCollection
	if defaultCardCollectionList.Cards == nil {
		// HACK: for some reason, empty message are converted to nil
		// Allow empty card collection for now, since it is not used anyway
		defaultCardCollectionList.Cards = []*zb_data.CardCollectionCard{}
	}
	if defaultCardCollectionList.Cards == nil {
		return fmt.Errorf("'defaultCollection' key missing")
	}

	// default decks
	defaultDecks.Decks = initData.DefaultDecks
	if defaultDecks.Decks == nil {
		return fmt.Errorf("'defaultDecks' key missing")
	}

	// AI decks
	aiDeckList.Decks = initData.AiDecks
	if aiDeckList.Decks == nil {
		return fmt.Errorf("'aiDecks' key missing")
	}

	// overlord experience info
	overlordLevelingData = initData.OverlordLeveling
	if overlordLevelingData == nil {
		return fmt.Errorf("'overlordLeveling' key missing")
	}

	// validate data
	// card library
	err := validateCardLibraryCards(cardList.Cards)
	if err != nil {
		return errors.Wrap(err, "error while validating card library")
	}

	// default decks
	for _, deck := range defaultDecks.Decks {
		if err := validateDeckAgainstCardLibrary(cardList.Cards, deck.Cards); err != nil {
			return errors.Wrap(err, fmt.Sprintf("error validating default deck with id %d", deck.Id))
		}
	}

	// ai decks
	for index, deck := range aiDeckList.Decks {
		if err := validateDeckAgainstCardLibrary(cardList.Cards, deck.Deck.Cards); err != nil {
			return errors.Wrap(err, fmt.Sprintf("error validating AI deck %d", index))
		}
	}

	// initialize card library
	if err := saveCardLibrary(ctx, initData.Version, &cardList); err != nil {
		return err
	}

	// initialize overlords
	if err := saveOverlordPrototypes(ctx, initData.Version, &overlordPrototypeList); err != nil {
		return errors.Wrap(err, "error updating overlord list")
	}

	// initialize default collection
	if err := saveDefaultCardCollection(ctx, initData.Version, &defaultCardCollectionList); err != nil {
		return errors.Wrap(err, "error updating default collection")
	}

	// initialize default deck
	if err := saveDefaultDecks(ctx, initData.Version, &defaultDecks); err != nil {
		return errors.Wrap(err, "error updating default decks")
	}

	// initialize AI decks
	if err := saveAIDecks(ctx, initData.Version, &aiDeckList); err != nil {
		return errors.Wrap(err, "error updating ai decks")
	}

	// initialize overlord experience
	if err := saveOverlordLevelingData(ctx, initData.Version, overlordLevelingData); err != nil {
		return err
	}

	return nil
}

func (z *ZombieBattleground) GetInit(ctx contract.StaticContext, req *zb_calls.GetInitRequest) (*zb_calls.GetInitResponse, error) {
	var cardLibrary *zb_data.CardList
	var overlordPrototypeList *zb_data.OverlordPrototypeList
	var defaultCardCollection *zb_data.CardCollectionList
	var defaultDecks *zb_data.DeckList
	var aiDeckList *zb_data.AIDeckList
	var overlordLevelingData *zb_data.OverlordLevelingData

	cardLibrary, err := loadCardLibraryRaw(ctx, req.Version)
	if err != nil {
		return nil, err
	}

	overlordPrototypeList, err = loadOverlordPrototypes(ctx, req.Version)
	if err != nil {
		return nil, err
	}

	defaultCardCollection, err = loadDefaultCardCollection(ctx, req.Version)
	if err != nil {
		return nil, err
	}

	defaultDecks, err = loadDefaultDecks(ctx, req.Version)
	if err != nil {
		return nil, err
	}

	aiDeckList, err = loadAIDecks(ctx, req.Version)
	if err != nil {
		return nil, err
	}

	overlordLevelingData, err = loadOverlordLevelingData(ctx, req.Version)
	if err != nil {
		return nil, err
	}

	var oracleAddress types.Address
	err = ctx.Get(oracleKey, &oracleAddress)
	if err != nil {
		return nil, err
	}

	return &zb_calls.GetInitResponse{
		InitData: &zb_data.InitData{
			Cards:             cardLibrary.Cards,
			Overlords:         overlordPrototypeList.Overlords,
			DefaultDecks:      defaultDecks.Decks,
			DefaultCollection: defaultCardCollection.Cards,
			AiDecks:           aiDeckList.Decks,
			OverlordLeveling:  overlordLevelingData,
			Version:           req.Version,
			Oracle:            &oracleAddress,
		},
	}, nil
}

// FIXME: duplicate of ListCardLibrary
func (z *ZombieBattleground) GetCardList(ctx contract.StaticContext, req *zb_calls.GetCardListRequest) (*zb_calls.GetCardListResponse, error) {
	if req.Version == "" {
		return nil, ErrVersionNotSet
	}

	cardList, err := loadCardLibrary(ctx, req.Version)
	if err != nil {
		return nil, err
	}

	return &zb_calls.GetCardListResponse{Cards: cardList.Cards}, nil
}

func (z *ZombieBattleground) GetAccount(ctx contract.StaticContext, req *zb_calls.GetAccountRequest) (*zb_data.Account, error) {
	var account zb_data.Account
	if err := ctx.Get(AccountKey(req.UserId), &account); err != nil {
		return nil, errors.Wrapf(err, "unable to retrieve account data for userId: %s", req.UserId)
	}
	return &account, nil
}

func (z *ZombieBattleground) UpdateAccount(ctx contract.Context, req *zb_calls.UpsertAccountRequest) (*zb_data.Account, error) {
	// Verify whether this privateKey associated with user
	if !isOwner(ctx, req.UserId) {
		return nil, ErrUserNotVerified
	}

	var account zb_data.Account
	accountKey := AccountKey(req.UserId)
	if err := ctx.Get(accountKey, &account); err != nil {
		return nil, errors.Wrapf(err, "unable to retrieve account data for userId: %s", req.UserId)
	}

	copyAccountInfo(&account, req)
	if err := ctx.Set(accountKey, &account); err != nil {
		return nil, errors.Wrapf(err, "error setting account information for userId: %s", req.UserId)
	}

	senderAddress := []byte(ctx.Message().Sender.Local)
	emitMsgJSON, err := prepareEmitMsgJSON(senderAddress, req.UserId, "updateaccount")
	if err == nil {
		ctx.EmitTopics(emitMsgJSON, TopicUpdateAccountEvent)
	}

	return &account, nil
}

func (z *ZombieBattleground) CreateAccount(ctx contract.Context, req *zb_calls.UpsertAccountRequest) error {
	if req.Version == "" {
		return ErrVersionNotSet
	}

	// confirm owner doesnt exist already
	if ctx.Has(AccountKey(req.UserId)) {
		ctx.Logger().Debug(fmt.Sprintf("user already exists -%s", req.UserId))
		return errors.New("user already exists")
	}

	var account zb_data.Account
	account.UserId = req.UserId
	account.Owner = ctx.Message().Sender.Bytes()
	copyAccountInfo(&account, req)

	if err := ctx.Set(AccountKey(req.UserId), &account); err != nil {
		return errors.Wrapf(err, "error setting account information for userId: %s", req.UserId)
	}
	ctx.GrantPermission([]byte(req.UserId), []string{OwnerRole})

	err := z.initializeUserDefaultCardCollection(ctx, req.Version, req.UserId)
	if err != nil {
		return errors.Wrap(err, "CreateAccount")
	}

	defaultDecks, err := z.initializeUserDefaultDecks(ctx, req.Version, req.UserId)
	if err != nil {
		return errors.Wrap(err, "CreateAccount")
	}

	//Emit CreateDeck event when creating new default decks for this new account
	for i := 0; i < len(defaultDecks.Decks); i++ {
		emitMsg := zb_calls.CreateDeckEvent{
			UserId:        req.UserId,
			SenderAddress: ctx.Message().Sender.Local.String(),
			Deck:          defaultDecks.Decks[i],
			Version:       req.Version,
		}

		data, err := proto.Marshal(&emitMsg)
		if err != nil {
			return errors.Wrap(err, "CreateAccount")
		}
		ctx.EmitTopics(data, TopicCreateDeckEvent)
	}

	err = saveAddressToUserIdLink(ctx, req.UserId, ctx.Message().Sender)
	if err != nil {
		return err
	}

	senderAddress := []byte(ctx.Message().Sender.Local)
	emitMsgJSON, err := prepareEmitMsgJSON(senderAddress, req.UserId, "createaccount")
	if err == nil {
		ctx.EmitTopics(emitMsgJSON, TopicCreateAccountEvent)
	}

	return nil
}

func (z *ZombieBattleground) Login(ctx contract.Context, req *zb_calls.LoginRequest) (*zb_calls.LoginResponse, error) {
	if !isOwner(ctx, req.UserId) {
		return nil, ErrUserNotVerified
	}

	if req.Version == "" {
		return nil, ErrVersionNotSet
	}

	err := saveAddressToUserIdLink(ctx, req.UserId, ctx.Message().Sender)
	if err != nil {
		return nil, err
	}

	userPersistentData, err := loadUserPersistentData(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	response := zb_calls.LoginResponse{}

	wipeExecuted, err := z.handleUserDataWipe(ctx, req.Version, req.UserId)
	if err != nil {
		return nil, err
	}

	response.DataWiped = wipeExecuted

	// apply any pending card collection changes that were pending because address to user id was not set
	userIdFound, err := z.applyPendingCardAmountChanges(ctx, ctx.Message().Sender)
	if err != nil {
		return nil, err
	}

	if !userIdFound {
		ctx.Logger().Warn("user id not found when doing applyPendingCardAmountChanges in Login", "userId", req.UserId, "userAddress", ctx.Message().Sender.String())
	}

	if userPersistentData.LastFullCardCollectionSyncPlasmachainBlockHeight == 0 {
		err = z.addGetUserFullCardCollectionOracleCommand(ctx, ctx.Message().Sender)
		if err != nil {
			return nil, err
		}

		// notify the client that a full collection sync is scheduled so it can wait a bit
		response.FullCardCollectionSyncExecuted = true
	}

	return &response, err
}

func (z *ZombieBattleground) UpdateUserElo(ctx contract.Context, req *zb_calls.UpdateUserEloRequest) error {
	// Verify whether this privateKey associated with user
	if !isOwner(ctx, req.UserId) {
		return ErrUserNotVerified
	}

	var account zb_data.Account
	accountKey := AccountKey(req.UserId)
	if err := ctx.Get(accountKey, &account); err != nil {
		return errors.Wrapf(err, "unable to retrieve account data for userId: %s", req.UserId)
	}

	// set elo score
	account.EloScore = req.EloScore

	if err := ctx.Set(accountKey, &account); err != nil {
		return errors.Wrapf(err, "error setting account elo score for userId: %s", req.UserId)
	}

	// emit event
	emitMsg := account
	data, err := proto.Marshal(&emitMsg)
	if err != nil {
		return err
	}
	ctx.EmitTopics(data, "zombiebattleground:update_elo")
	return nil
}

// CreateDeck appends the given deck to user's deck list
func (z *ZombieBattleground) CreateDeck(ctx contract.Context, req *zb_calls.CreateDeckRequest) (*zb_calls.CreateDeckResponse, error) {
	if req.Version == "" {
		return nil, ErrVersionNotSet
	}
	if req.Deck == nil {
		return nil, ErrDeckMustNotNil
	}
	if !isOwner(ctx, req.UserId) {
		return nil, ErrUserNotVerified
	}

	overlords, err := loadOverlordUserInstances(ctx, req.Version, req.UserId)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get overlord instances for userId: %s", req.UserId)
	}

	deckList, err := loadDecks(ctx, req.UserId, req.Version)
	if err != nil {
		return nil, errors.Wrap(err, "unable to load decks")
	}

	cardLibrary, err := loadCardLibrary(ctx, req.Version)
	if err != nil {
		return nil, err
	}

	userCardCollection, err := z.loadUserCardCollection(ctx, req.Version, req.UserId)
	if err != nil {
		return nil, err
	}

	err = validateDeck(false, cardLibrary, userCardCollection, req.Deck, deckList.Decks, overlords)
	if err != nil {
		return nil, errors.Wrap(err, "error validating deck")
	}

	// allocate new deck id starting from 1
	var newDeckID int64
	if len(deckList.Decks) != 0 {
		for _, deck := range deckList.Decks {
			if deck.Id > newDeckID {
				newDeckID = deck.Id
			}
		}
	}
	newDeckID++
	req.Deck.Id = newDeckID
	deckList.Decks = append(deckList.Decks, req.Deck)
	if err := saveDecks(ctx, req.Version, req.UserId, deckList); err != nil {
		return nil, err
	}

	senderAddress := ctx.Message().Sender.Local.String()
	emitMsg := zb_calls.CreateDeckEvent{
		UserId:        req.UserId,
		SenderAddress: senderAddress,
		Deck:          req.Deck,
		Version:       req.Version,
	}

	data, err := proto.Marshal(&emitMsg)
	if err != nil {
		return nil, err
	}
	ctx.EmitTopics(data, TopicCreateDeckEvent)

	return &zb_calls.CreateDeckResponse{DeckId: newDeckID}, nil
}

// EditDeck edits the deck by id
func (z *ZombieBattleground) EditDeck(ctx contract.Context, req *zb_calls.EditDeckRequest) error {
	if req.Version == "" {
		return ErrVersionNotSet
	}
	if req.Deck == nil {
		return fmt.Errorf("deck must not be nil")
	}
	if !isOwner(ctx, req.UserId) {
		return ErrUserNotVerified
	}

	overlords, err := loadOverlordUserInstances(ctx, req.Version, req.UserId)
	if err != nil {
		return errors.Wrapf(err, "unable to get overlord instances for userId: %s", req.UserId)
	}

	deckList, err := loadDecks(ctx, req.UserId, req.Version)
	if err != nil {
		return errors.Wrap(err, "unable to load decks")
	}

	cardLibrary, err := loadCardLibrary(ctx, req.Version)
	if err != nil {
		return err
	}

	userCardCollection, err := z.loadUserCardCollection(ctx, req.Version, req.UserId)
	if err != nil {
		return err
	}

	err = validateDeck(true, cardLibrary, userCardCollection, req.Deck, deckList.Decks, overlords)
	if err != nil {
		return errors.Wrap(err, "error validating deck")
	}

	existingDeck := getDeckByID(deckList.Decks, req.Deck.Id)
	if existingDeck == nil {
		return ErrNotFound
	}

	// update deck
	existingDeck.Reset()
	proto.Merge(existingDeck, req.Deck)

	// update decklist
	if err := saveDecks(ctx, req.Version, req.UserId, deckList); err != nil {
		return errors.Wrap(err, "error saving decks")
	}

	senderAddress := ctx.Message().Sender.Local.String()
	emitMsg := zb_calls.EditDeckEvent{
		UserId:        req.UserId,
		SenderAddress: senderAddress,
		Deck:          req.Deck,
		Version:       req.Version,
	}

	data, err := proto.Marshal(&emitMsg)
	if err != nil {
		return err
	}
	ctx.EmitTopics(data, TopicEditDeckEvent)

	return nil
}

// DeleteDeck deletes a user's deck by id
func (z *ZombieBattleground) DeleteDeck(ctx contract.Context, req *zb_calls.DeleteDeckRequest) error {
	if req.Version == "" {
		return ErrVersionNotSet
	}

	if !isOwner(ctx, req.UserId) {
		return ErrUserNotVerified
	}

	deckList, err := loadDecks(ctx, req.UserId, req.Version)
	if err != nil {
		return err
	}

	var deleted bool
	deckList.Decks, deleted = deleteDeckByID(deckList.Decks, req.DeckId)
	if !deleted {
		return fmt.Errorf("deck not found")
	}

	if err := saveDecks(ctx, req.Version, req.UserId, deckList); err != nil {
		return err
	}

	senderAddress := ctx.Message().Sender.Local.String()
	emitMsg := zb_calls.DeleteDeckEvent{
		UserId:        req.UserId,
		SenderAddress: senderAddress,
		DeckId:        req.DeckId,
	}

	data, err := proto.Marshal(&emitMsg)
	if err != nil {
		return err
	}
	ctx.EmitTopics(data, TopicDeleteDeckEvent)

	return nil
}

// ListDecks returns the user's decks
func (z *ZombieBattleground) ListDecks(ctx contract.Context, req *zb_calls.ListDecksRequest) (*zb_calls.ListDecksResponse, error) {
	if req.Version == "" {
		return nil, ErrVersionNotSet
	}

	deckList, err := loadDecks(ctx, req.UserId, req.Version)
	if err != nil {
		return nil, err
	}
	return &zb_calls.ListDecksResponse{
		Decks: deckList.Decks,
	}, nil
}

// GetDeck returns the deck by given id
func (z *ZombieBattleground) GetDeck(ctx contract.Context, req *zb_calls.GetDeckRequest) (*zb_calls.GetDeckResponse, error) {
	if req.Version == "" {
		return nil, ErrVersionNotSet
	}

	deckList, err := loadDecks(ctx, req.UserId, req.Version)
	if err != nil {
		return nil, err
	}
	deck := getDeckByID(deckList.Decks, req.DeckId)
	if deck == nil {
		return nil, contract.ErrNotFound
	}
	return &zb_calls.GetDeckResponse{Deck: deck}, nil
}

func (z *ZombieBattleground) GetAIDecks(ctx contract.StaticContext, req *zb_calls.GetAIDecksRequest) (*zb_calls.GetAIDecksResponse, error) {
	if req.Version == "" {
		return nil, ErrVersionNotSet
	}

	deckList, err := loadAIDecks(ctx, req.Version)
	if err != nil {
		return nil, err
	}

	return &zb_calls.GetAIDecksResponse{
		AiDecks: deckList.Decks,
	}, nil
}

// GetCollection returns the collection of the card own by the user
func (z *ZombieBattleground) GetCollection(ctx contract.Context, req *zb_calls.GetCollectionRequest) (*zb_calls.GetCollectionResponse, error) {
	if req.Version == "" {
		return nil, ErrVersionNotSet
	}

	collectionList, err := z.loadUserCardCollection(ctx, req.Version, req.UserId)
	if err != nil {
		return nil, err
	}
	return &zb_calls.GetCollectionResponse{Cards: collectionList}, nil
}

// ListCardLibrary list all the card library data
func (z *ZombieBattleground) ListCardLibrary(ctx contract.StaticContext, req *zb_calls.ListCardLibraryRequest) (*zb_calls.ListCardLibraryResponse, error) {
	if req.Version == "" {
		return nil, ErrVersionNotSet
	}

	cardList, err := loadCardLibrary(ctx, req.Version)
	if err != nil {
		return nil, err
	}

	return &zb_calls.ListCardLibraryResponse{Cards: cardList.Cards}, nil
}

func (z *ZombieBattleground) ListOverlordLibrary(ctx contract.StaticContext, req *zb_calls.ListOverlordLibraryRequest) (*zb_calls.ListOverlordLibraryResponse, error) {
	if req.Version == "" {
		return nil, ErrVersionNotSet
	}

	overlordPrototypeList, err := loadOverlordPrototypes(ctx, req.Version)
	if err != nil {
		return nil, err
	}
	return &zb_calls.ListOverlordLibraryResponse{Overlords: overlordPrototypeList.Overlords}, nil
}

func (z *ZombieBattleground) ListOverlordUserInstances(ctx contract.StaticContext, req *zb_calls.ListOverlordUserInstancesRequest) (*zb_calls.ListOverlordUserInstancesResponse, error) {
	if req.Version == "" {
		return nil, ErrVersionNotSet
	}

	overlordList, err := loadOverlordUserInstances(ctx, req.Version, req.UserId)
	if err != nil {
		return nil, err
	}
	return &zb_calls.ListOverlordUserInstancesResponse{Overlords: overlordList}, nil
}

func (z *ZombieBattleground) GetOverlordUserInstance(ctx contract.StaticContext, req *zb_calls.GetOverlordUserInstanceRequest) (*zb_calls.GetOverlordUserInstanceResponse, error) {
	if req.Version == "" {
		return nil, ErrVersionNotSet
	}

	overlordList, err := loadOverlordUserInstances(ctx, req.Version, req.UserId)
	if err != nil {
		return nil, err
	}
	overlord, found := getOverlordUserInstanceByPrototypeId(overlordList, req.OverlordId)
	if !found {
		return nil, fmt.Errorf("overlord with prototype id %d not found", req.OverlordId)
	}
	return &zb_calls.GetOverlordUserInstanceResponse{Overlord: overlord}, nil
}

func (z *ZombieBattleground) RegisterPlayerPool(ctx contract.Context, req *zb_calls.RegisterPlayerPoolRequest) (*zb_calls.RegisterPlayerPoolResponse, error) {
	// preparing user profile consisting of deck, score, ...
	_, err := getDeckWithRegistrationData(ctx, req.RegistrationData, req.RegistrationData.Version)
	if err != nil {
		return nil, err
	}

	if req.RegistrationData.Version == "" {
		return nil, fmt.Errorf("version not specified")
	}

	// sort tags
	if len(req.RegistrationData.Tags) > 0 {
		sort.Strings(req.RegistrationData.Tags)
	}

	profile := zb_data.PlayerProfile{
		RegistrationData: req.RegistrationData,
		UpdatedAt:        ctx.Now().Unix(),
	}

	var loadPlayerPoolFn func(contract.Context) (*zb_data.PlayerPool, error)
	var savePlayerPoolFn func(contract.Context, *zb_data.PlayerPool) error
	// if the tags is set, use tagged playerpool
	if len(profile.RegistrationData.Tags) > 0 {
		loadPlayerPoolFn = loadTaggedPlayerPool
		savePlayerPoolFn = saveTaggedPlayerPool
	} else {
		loadPlayerPoolFn = loadPlayerPool
		savePlayerPoolFn = savePlayerPool
	}

	// load player pool
	pool, err := loadPlayerPoolFn(ctx)
	if err != nil {
		return nil, err
	}

	match, _ := loadUserCurrentMatch(ctx, req.RegistrationData.UserId)
	if match != nil {
		return nil, errors.New("Player is already in a match")
	}

	targetProfile := findPlayerProfileByID(pool, req.RegistrationData.UserId)
	// if player is in the pool, remove the player from the pool first. otherwise, the profile won't get updated
	if targetProfile != nil {
		pool = removePlayerFromPool(pool, req.RegistrationData.UserId)
	}
	pool.PlayerProfiles = append(pool.PlayerProfiles, &profile)
	if err := savePlayerPoolFn(ctx, pool); err != nil {
		return nil, err
	}

	// prune the timed out player profile
	for _, pp := range pool.PlayerProfiles {
		updatedAt := time.Unix(pp.UpdatedAt, 0)
		if updatedAt.Add(MMTimeout).Before(ctx.Now()) {
			ctx.Logger().Info(fmt.Sprintf("Player profile %s timedout", pp.RegistrationData.UserId))
			// remove player from the pool
			pool = removePlayerFromPool(pool, pp.RegistrationData.UserId)
			if err := savePlayerPoolFn(ctx, pool); err != nil {
				return nil, err
			}
			// remove match
			match, _ := loadUserCurrentMatch(ctx, pp.RegistrationData.UserId)
			if match != nil {
				ctx.Delete(MatchKey(match.Id))
				match.Status = zb_data.Match_Timedout
				// remove player's match if existing
				ctx.Delete(UserMatchKey(pp.RegistrationData.UserId))
				// notify player
				emitMsg := zb_data.PlayerActionEvent{
					Match:            match,
					CreatedByBackend: true,
				}
				data, err := proto.Marshal(&emitMsg)
				if err != nil {
					return nil, err
				}
				ctx.EmitTopics(data, match.Topics...)
			}
		}
	}

	senderAddress := []byte(ctx.Message().Sender.Local)
	emitMsgJSON, err := prepareEmitMsgJSON(senderAddress, req.RegistrationData.UserId, "registerplayerpool")
	if err == nil {
		ctx.EmitTopics(emitMsgJSON, TopicRegisterPlayerPoolEvent)
	}

	return &zb_calls.RegisterPlayerPoolResponse{}, nil
}

func (z *ZombieBattleground) FindMatch(ctx contract.Context, req *zb_calls.FindMatchRequest) (*zb_calls.FindMatchResponse, error) {
	var loadPlayerPoolFn func(contract.Context) (*zb_data.PlayerPool, error)
	var savePlayerPoolFn func(contract.Context, *zb_data.PlayerPool) error
	// if the tags is set, use tagged playerpool
	if len(req.Tags) > 0 {
		loadPlayerPoolFn = loadTaggedPlayerPool
		savePlayerPoolFn = saveTaggedPlayerPool
	} else {
		loadPlayerPoolFn = loadPlayerPool
		savePlayerPoolFn = savePlayerPool
	}

	pool, err := loadPlayerPoolFn(ctx)
	if err != nil {
		return nil, err
	}
	match, _ := loadUserCurrentMatch(ctx, req.UserId)
	if match != nil {
		// timeout for matchmaking
		if match.Status == zb_data.Match_Matching {
			updatedAt := time.Unix(match.CreatedAt, 0)
			if updatedAt.Add(MMTimeout).Before(ctx.Now()) {
				ctx.Logger().Debug(fmt.Sprintf("Match %d timedout", match.Id))
				// remove match
				// ctx.Delete(MatchKey(match.Id))
				match.Status = zb_data.Match_Timedout
				if err := saveMatch(ctx, match); err != nil {
					return nil, err
				}
				// remove player's match if existing
				for _, player := range match.PlayerStates {
					ctx.Delete(UserMatchKey(player.Id))
				}
			}
		}
		// notify player
		emitMsg := zb_data.PlayerActionEvent{
			Match:            match,
			CreatedByBackend: true,
		}
		data, err := proto.Marshal(&emitMsg)
		if err != nil {
			return nil, err
		}

		topics := append(match.Topics, TopicFindMatchEvent)
		ctx.EmitTopics(data, topics...)

		return &zb_calls.FindMatchResponse{
			Match:      match,
			MatchFound: true,
		}, nil
	}
	playerProfile := findPlayerProfileByID(pool, req.UserId)
	if playerProfile == nil {
		return nil, errors.New("Player not found in player pool")
	}

	deck, err := getDeckWithRegistrationData(ctx, playerProfile.RegistrationData, playerProfile.RegistrationData.Version)
	if err != nil {
		return nil, err
	}

	// perform matchmaking function to calculate scores
	// steps:
	// 1. list all the candidates that has similar profiles TODO: match versions
	// 2. pick the most highest score
	// 3. if there is no candidate, sleep for MMWaitTime seconds
	retries := 0
	var matchedPlayerProfile *zb_data.PlayerProfile
	for retries < MMRetries {
		var playerScores []*PlayerScore
		for _, pp := range pool.PlayerProfiles {
			// skip the requesting player
			if pp.RegistrationData.UserId == req.UserId {
				continue
			}
			score := mmf(playerProfile, pp)
			// only non-negative score will be added
			if score > 0 {
				playerScores = append(playerScores, &PlayerScore{score: score, id: pp.RegistrationData.UserId})
			}
		}

		sortedPlayerScores := sortByPlayerScore(playerScores)
		if len(sortedPlayerScores) > 0 {
			matchedPlayerID := sortedPlayerScores[0].id
			matchedPlayerProfile = findPlayerProfileByID(pool, matchedPlayerID)
			// remove the match players from the pool
			pool = removePlayerFromPool(pool, matchedPlayerID)
			pool = removePlayerFromPool(pool, req.UserId)
			if err := savePlayerPoolFn(ctx, pool); err != nil {
				return nil, err
			}
			break
		}
		retries++
	}

	if matchedPlayerProfile == nil {
		return &zb_calls.FindMatchResponse{
			MatchFound: false,
		}, nil
	}

	// get matched player deck
	matchedDeck, err := getDeckWithRegistrationData(ctx, matchedPlayerProfile.RegistrationData, matchedPlayerProfile.RegistrationData.Version)
	if err != nil {
		return nil, err
	}

	// create match
	match = &zb_data.Match{
		Status: zb_data.Match_Matching,
		PlayerStates: []*zb_data.InitialPlayerState{
			&zb_data.InitialPlayerState{
				Id:            playerProfile.RegistrationData.UserId,
				Deck:          deck,
				MatchAccepted: false,
			},
			&zb_data.InitialPlayerState{
				Id:            matchedPlayerProfile.RegistrationData.UserId,
				Deck:          matchedDeck,
				MatchAccepted: false,
			},
		},
		Version: playerProfile.RegistrationData.Version, // TODO: match version of both players
		PlayerLastSeens: []*zb_data.PlayerTimestamp{
			&zb_data.PlayerTimestamp{
				Id:        playerProfile.RegistrationData.UserId,
				UpdatedAt: ctx.Now().Unix(),
			},
			&zb_data.PlayerTimestamp{
				Id:        matchedPlayerProfile.RegistrationData.UserId,
				UpdatedAt: ctx.Now().Unix(),
			},
		},
		PlayerDebugCheats: []*zb_data.DebugCheatsConfiguration{
			&playerProfile.RegistrationData.DebugCheats,
			&matchedPlayerProfile.RegistrationData.DebugCheats,
		},
	}

	if playerProfile.RegistrationData.DebugCheats.Enabled && playerProfile.RegistrationData.DebugCheats.UseCustomRandomSeed {
		match.RandomSeed = playerProfile.RegistrationData.DebugCheats.CustomRandomSeed
	} else {
		match.RandomSeed = ctx.Now().Unix()
	}

	match.CustomGameAddr = playerProfile.RegistrationData.CustomGame // TODO: make sure both players request same custom game?

	if err := createMatch(ctx, match, playerProfile.RegistrationData.UseBackendGameLogic); err != nil {
		return nil, err
	}

	// save user match
	if err := saveUserCurrentMatch(ctx, playerProfile.RegistrationData.UserId, match); err != nil {
		return nil, err
	}
	if err := saveUserCurrentMatch(ctx, matchedPlayerProfile.RegistrationData.UserId, match); err != nil {
		return nil, err
	}
	// save match
	// if err := saveMatch(ctx, match); err != nil {
	// 	return nil, err
	// }

	emitMsg := zb_data.PlayerActionEvent{
		Match: match,
	}
	data, err := proto.Marshal(&emitMsg)
	if err != nil {
		return nil, err
	}
	topics := append(match.Topics, TopicFindMatchEvent)
	ctx.EmitTopics(data, topics...)

	return &zb_calls.FindMatchResponse{
		Match:      match,
		MatchFound: true,
	}, nil
}

func (z *ZombieBattleground) AcceptMatch(ctx contract.Context, req *zb_calls.AcceptMatchRequest) (*zb_calls.AcceptMatchResponse, error) {
	match, err := loadUserCurrentMatch(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	if req.MatchId != match.Id {
		return nil, errors.New("match id not correct")
	}

	if match.Status != zb_data.Match_Matching {
		return nil, errors.New("Can't accept match, wrong status")
	}

	var opponentAccepted bool
	for _, playerState := range match.PlayerStates {
		if playerState.Id == req.UserId {
			playerState.MatchAccepted = true
		} else {
			opponentAccepted = playerState.MatchAccepted
		}
	}

	emitMsg := zb_data.PlayerActionEvent{
		Match:            match,
		CreatedByBackend: true,
	}

	if opponentAccepted {
		var customModeAddr loom.Address
		var customModeAddr2 *loom.Address
		var customModeAddrStr string
		//TODO cleanup how we do this parsing
		if match.CustomGameAddr != nil {
			customModeAddrStr = fmt.Sprintf("default:%s", match.CustomGameAddr.Local.String())
		}

		customModeAddr, err = loom.ParseAddress(customModeAddrStr)
		if err != nil {
			ctx.Logger().Debug(fmt.Sprintf("no custom game mode --%v\n", err))
		} else {
			customModeAddr2 = &customModeAddr
		}

		playerStates := []*zb_data.PlayerState{
			&zb_data.PlayerState{
				Id:    match.PlayerStates[0].Id,
				Deck:  match.PlayerStates[0].Deck,
				Index: -1,
			},
			&zb_data.PlayerState{
				Id:    match.PlayerStates[1].Id,
				Deck:  match.PlayerStates[1].Deck,
				Index: -1,
			},
		}

		gp, err := NewGamePlay(
			ctx,
			match.Id,
			match.Version,
			playerStates,
			match.RandomSeed,
			customModeAddr2,
			match.UseBackendGameLogic,
			match.PlayerDebugCheats,
		)
		if err != nil {
			return nil, err
		}
		if err := saveGameState(ctx, gp.State); err != nil {
			return nil, err
		}

		match.Status = zb_data.Match_Started

		emitMsg = zb_data.PlayerActionEvent{
			Match:            match,
			Block:            &zb_data.History{List: gp.history},
			CreatedByBackend: true,
		}
	}

	// save user match
	for _, playerState := range match.PlayerStates {
		if err := saveUserCurrentMatch(ctx, playerState.Id, match); err != nil {
			return nil, err
		}
	}
	// save match
	if err := saveMatch(ctx, match); err != nil {
		return nil, err
	}

	data, err := proto.Marshal(&emitMsg)
	if err != nil {
		return nil, err
	}
	topics := append(match.Topics, TopicAcceptMatchEvent)
	ctx.EmitTopics(data, topics...)

	return &zb_calls.AcceptMatchResponse{
		Match: match,
	}, nil
}

// TODO remove this
func (z *ZombieBattleground) GetPlayerPool(ctx contract.Context, req *zb_calls.PlayerPoolRequest) (*zb_calls.PlayerPoolResponse, error) {
	pool, err := loadPlayerPool(ctx)
	if err != nil {
		return nil, err
	}

	return &zb_calls.PlayerPoolResponse{
		Pool: pool,
	}, nil
}

// TODO remove this
func (z *ZombieBattleground) GetTaggedPlayerPool(ctx contract.Context, req *zb_calls.PlayerPoolRequest) (*zb_calls.PlayerPoolResponse, error) {
	pool, err := loadTaggedPlayerPool(ctx)
	if err != nil {
		return nil, err
	}

	return &zb_calls.PlayerPoolResponse{
		Pool: pool,
	}, nil
}

func (z *ZombieBattleground) CancelFindMatch(ctx contract.Context, req *zb_calls.CancelFindMatchRequest) (*zb_calls.CancelFindMatchResponse, error) {
	match, _ := loadUserCurrentMatch(ctx, req.UserId)

	if match != nil && match.Status != zb_data.Match_Ended {
		// remove current match
		for _, player := range match.PlayerStates {
			ctx.Delete(UserMatchKey(player.Id))
		}
		match.Status = zb_data.Match_Canceled
		if err := saveMatch(ctx, match); err != nil {
			return nil, err
		}
		// notify player
		emitMsg := zb_data.PlayerActionEvent{
			Match:            match,
			CreatedByBackend: true,
		}
		data, err := proto.Marshal(&emitMsg)
		if err != nil {
			return nil, err
		}
		if err == nil {
			ctx.EmitTopics(data, match.Topics...)
		}
	}

	var loadPlayerPoolFn func(contract.Context) (*zb_data.PlayerPool, error)
	var savePlayerPoolFn func(contract.Context, *zb_data.PlayerPool) error
	// if the tags is set, use tagged playerpool
	if len(req.Tags) > 0 {
		loadPlayerPoolFn = loadTaggedPlayerPool
		savePlayerPoolFn = saveTaggedPlayerPool
	} else {
		loadPlayerPoolFn = loadPlayerPool
		savePlayerPoolFn = savePlayerPool
	}

	// remove player from the player pool
	pool, err := loadPlayerPoolFn(ctx)
	if err != nil {
		return nil, err
	}
	pool = removePlayerFromPool(pool, req.UserId)
	if err := savePlayerPoolFn(ctx, pool); err != nil {
		return nil, err
	}

	return &zb_calls.CancelFindMatchResponse{}, nil
}

func (z *ZombieBattleground) GetMatch(ctx contract.StaticContext, req *zb_calls.GetMatchRequest) (*zb_calls.GetMatchResponse, error) {
	match, err := loadMatch(ctx, req.MatchId)
	if err != nil {
		return nil, err
	}

	return &zb_calls.GetMatchResponse{
		Match: match,
	}, nil
}

func (z *ZombieBattleground) GetGameState(ctx contract.StaticContext, req *zb_calls.GetGameStateRequest) (*zb_calls.GetGameStateResponse, error) {
	gameState, err := loadGameState(ctx, req.MatchId)
	if err != nil {
		return nil, err
	}

	return &zb_calls.GetGameStateResponse{
		GameState: gameState,
	}, nil
}

func (z *ZombieBattleground) GetInitialGameState(ctx contract.StaticContext, req *zb_calls.GetGameStateRequest) (*zb_calls.GetGameStateResponse, error) {
	initialGameState, err := loadInitialGameState(ctx, req.MatchId)
	if err != nil {
		return nil, err
	}

	return &zb_calls.GetGameStateResponse{
		GameState: initialGameState,
	}, nil
}

func (z *ZombieBattleground) AddSoloExperience(ctx contract.Context, req *zb_calls.AddSoloExperienceRequest) (*zb_calls.AddSoloExperienceResponse, error) {
	if req.Version == "" {
		return nil, fmt.Errorf("version not specified")
	}

	overlordLevelingData, err := loadOverlordLevelingData(ctx, req.Version)
	if err != nil {
		return nil, err
	}

	if err := applyExperience(ctx, req.Version, overlordLevelingData, req.UserId, parseUserIdToNumber(req.UserId), req.OverlordId, req.Experience, req.DeckId, req.IsWin); err != nil {
		return nil, err
	}

	return &zb_calls.AddSoloExperienceResponse{}, nil
}

func (z *ZombieBattleground) EndMatch(ctx contract.Context, req *zb_calls.EndMatchRequest) (*zb_calls.EndMatchResponse, error) {
	match, err := loadMatch(ctx, req.MatchId)
	if err != nil {
		return nil, err
	}

	match.Status = zb_data.Match_Ended
	if err := saveMatch(ctx, match); err != nil {
		return nil, err
	}

	// load game state
	gameState, err := loadGameState(ctx, req.MatchId)
	if err != nil {
		return nil, err
	}

	// save experience and level for both players
	overlordLevelingData, err := loadOverlordLevelingData(ctx, gameState.Version)
	if err != nil {
		return nil, err
	}

	for index, playerState := range match.PlayerStates {
		if err := applyExperience(
			ctx,
			match.Version,
			overlordLevelingData,
			playerState.Id,
			parseUserIdToNumber(playerState.Id),
			playerState.Deck.OverlordId,
			req.MatchExperiences[index],
			playerState.Deck.Id,
			req.WinnerId == playerState.Id,
		); err != nil {
			return nil, err
		}
	}

	// delete user match for both users
	for _, playerState := range match.PlayerStates {
		ctx.Delete(UserMatchKey(playerState.Id))
	}

	// update gamestate
	gamestate, err := loadGameState(ctx, match.Id)
	if err != nil {
		return nil, err
	}

	gp, err := GamePlayFrom(gamestate, match.UseBackendGameLogic, match.PlayerDebugCheats)
	if err != nil {
		return nil, err
	}
	gamestate.Winner = req.WinnerId
	gamestate.IsEnded = true
	if err := saveGameState(ctx, gamestate); err != nil {
		return nil, err
	}

	//TODO obviously this will need to change drastically once the logic is on the server
	gp.history = append(gp.history, &zb_data.HistoryData{
		Data: &zb_data.HistoryData_EndGame{
			EndGame: &zb_data.HistoryEndGame{
				UserId:   req.GetUserId(),
				MatchId:  req.MatchId,
				WinnerId: req.WinnerId,
			},
		},
	})
	// Don't think we need this since endgame should be emitted to match
	// match.Topics = append(match.Topics, "endgame")
	emitMsg := zb_data.PlayerActionEvent{
		Match:            match,
		Block:            &zb_data.History{List: gp.history},
		CreatedByBackend: true,
	}
	data, err := proto.Marshal(&emitMsg)
	if err != nil {
		return nil, err
	}
	ctx.EmitTopics(data, match.Topics...)

	return &zb_calls.EndMatchResponse{GameState: gamestate}, nil
}

func (z *ZombieBattleground) SendPlayerAction(ctx contract.Context, req *zb_calls.PlayerActionRequest) (*zb_calls.PlayerActionResponse, error) {
	match, err := loadMatch(ctx, req.MatchId)
	if err != nil {
		return nil, err
	}

	// check if the user is in the match
	found := false
	for _, player := range match.PlayerStates {
		if player.Id == req.PlayerAction.PlayerId {
			found = true
		}
	}
	if !found {
		return nil, errors.New("player not in the match")
	}

	gamestate, err := loadGameState(ctx, match.Id)
	if err != nil {
		return nil, err
	}
	gp, err := GamePlayFrom(gamestate, match.UseBackendGameLogic, match.PlayerDebugCheats)
	if err != nil {
		return nil, err
	}
	// TODO: change me. this is a bit hacky way to set card libarary
	cardlist, err := loadCardLibrary(ctx, gamestate.Version)
	if err != nil {
		return nil, err
	}
	gp.cardLibrary = cardlist
	gp.SetLogger(ctx.Logger())
	// add created timestamp
	req.PlayerAction.CreatedAt = ctx.Now().Unix()
	if err := gp.AddAction(req.PlayerAction); err != nil {
		return nil, err
	}

	req.PlayerAction.ActionOutcomes = gp.actionOutcomes
	gp.actionOutcomes = nil

	if req.PlayerAction.ActionOutcomes != nil && len(req.PlayerAction.ActionOutcomes) > 0 {
		ctx.Logger().Info(fmt.Sprintf("\n\nreq.PlayerAction.ActionOutcomes: %v\n\n", req.PlayerAction.ActionOutcomes))
	}

	if err := saveGameState(ctx, gamestate); err != nil {
		return nil, err
	}

	// update match status
	if match.Status == zb_data.Match_Started {
		match.Status = zb_data.Match_Playing
		if err := saveMatch(ctx, match); err != nil {
			return nil, err
		}
	}

	emitMsg := zb_data.PlayerActionEvent{
		PlayerAction:       req.PlayerAction,
		CurrentActionIndex: gamestate.CurrentActionIndex,
		Match:              match,
		Block:              &zb_data.History{List: gp.history},
		CreatedByBackend:   false,
	}

	data, err := proto.Marshal(&emitMsg)
	if err != nil {
		return nil, err
	}
	ctx.EmitTopics(data, match.Topics...)

	return &zb_calls.PlayerActionResponse{
		Match: match,
	}, nil
}

func (z *ZombieBattleground) SendBundlePlayerAction(ctx contract.Context, req *zb_calls.BundlePlayerActionRequest) (*zb_calls.BundlePlayerActionResponse, error) {
	match, err := loadMatch(ctx, req.MatchId)
	if err != nil {
		return nil, err
	}
	gamestate, err := loadGameState(ctx, match.Id)
	if err != nil {
		return nil, err
	}
	gp, err := GamePlayFrom(gamestate, match.UseBackendGameLogic, match.PlayerDebugCheats)
	if err != nil {
		return nil, err
	}
	gp.PrintState()
	if err := gp.AddBundleAction(req.PlayerActions...); err != nil {
		return nil, err
	}
	gp.PrintState()

	if err := saveGameState(ctx, gamestate); err != nil {
		return nil, err
	}

	// update match status
	if match.Status == zb_data.Match_Started {
		match.Status = zb_data.Match_Playing
		if err := saveMatch(ctx, match); err != nil {
			return nil, err
		}
	}

	return &zb_calls.BundlePlayerActionResponse{
		GameState: gamestate,
		Match:     match,
		History:   gp.history,
	}, nil
}

// ReplayGame simulate the game that has been created by initializing game from start and
// apply actions to from the current gamestate. ReplayGame does not save any gamestate.
func (z *ZombieBattleground) ReplayGame(ctx contract.Context, req *zb_calls.ReplayGameRequest) (*zb_calls.ReplayGameResponse, error) {
	match, err := loadMatch(ctx, req.MatchId)
	if err != nil {
		return nil, err
	}
	initGameState, err := loadInitialGameState(ctx, match.Id)
	if err != nil {
		return nil, err
	}
	gp, err := GamePlayFrom(initGameState, match.UseBackendGameLogic, match.PlayerDebugCheats)
	if err != nil {
		return nil, err
	}
	// TODO: change me. this is a bit hacky way to set card libarary
	cardlist, err := loadCardLibrary(ctx, initGameState.Version)
	if err != nil {
		return nil, err
	}
	gp.cardLibrary = cardlist
	gp.SetLogger(ctx.Logger())

	// get all actions from game states
	currentGameState, err := loadGameState(ctx, match.Id)
	if err != nil {
		return nil, err
	}

	actions := currentGameState.PlayerActions
	if req.StopAtActionIndex > -1 && int(req.StopAtActionIndex) < len(actions) {
		actions = actions[:int(req.StopAtActionIndex)]
	}

	if err := gp.AddBundleAction(actions...); err != nil {
		return nil, err
	}

	return &zb_calls.ReplayGameResponse{
		GameState:      initGameState,
		ActionOutcomes: gp.actionOutcomes,
	}, nil
}

func (z *ZombieBattleground) GetNotifications(ctx contract.StaticContext, req *zb_calls.GetNotificationsRequest) (*zb_calls.GetNotificationsResponse, error) {
	notificationList, err := loadUserNotifications(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &zb_calls.GetNotificationsResponse{
		Notifications: notificationList.Notifications,
	}, nil
}

func (z *ZombieBattleground) ClearNotifications(ctx contract.Context, req *zb_calls.ClearNotificationsRequest) (*zb_calls.ClearNotificationsResponse, error) {
	notificationList, err := loadUserNotifications(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	for _, id := range req.NotificationIds {
		notificationList.Notifications, err = removeNotification(notificationList.Notifications, id)
		if err != nil {
			return nil, err
		}
	}

	err = saveUserNotifications(ctx, req.UserId, notificationList)
	if err != nil {
		return nil, err
	}

	return &zb_calls.ClearNotificationsResponse{}, nil
}

func (z *ZombieBattleground) GetOverlordLevelingData(ctx contract.StaticContext, req *zb_calls.GetOverlordLevelingDataRequest) (*zb_calls.GetOverlordLevelingDataResponse, error) {
	if req.Version == "" {
		return nil, ErrVersionNotSet
	}

	overlordLevelingData, err := loadOverlordLevelingData(ctx, req.Version)
	if err != nil {
		return nil, err
	}

	return &zb_calls.GetOverlordLevelingDataResponse{
		OverlordLeveling: overlordLevelingData,
	}, nil
}

func (z *ZombieBattleground) KeepAlive(ctx contract.Context, req *zb_calls.KeepAliveRequest) (*zb_calls.KeepAliveResponse, error) {
	match, err := loadMatch(ctx, req.MatchId)
	if err != nil {
		return nil, err
	}

	var playerIndex = -1
	var playerID string
	for i, playerState := range match.PlayerStates {
		if playerState.Id == req.UserId {
			playerIndex = i
			playerID = playerState.Id
		}
	}
	if playerIndex < 0 {
		return nil, fmt.Errorf("player id %s not found", playerID)
	}

	if playerIndex > len(match.PlayerLastSeens)-1 {
		return nil, fmt.Errorf("player id %s not found", playerID)
	}

	var skipInitialChecking bool
	for _, lastseen := range match.PlayerLastSeens {
		if lastseen.UpdatedAt == 0 {
			skipInitialChecking = true
			break
		}
	}
	// init keepalive timestamp
	now := ctx.Now().Unix()
	if skipInitialChecking {
		for i := range match.PlayerLastSeens {
			match.PlayerLastSeens[i].UpdatedAt = now
		}
	}
	// update timestamp
	match.PlayerLastSeens[playerIndex].UpdatedAt = now
	if err := saveMatch(ctx, match); err != nil {
		return nil, err
	}

	if skipInitialChecking {
		return &zb_calls.KeepAliveResponse{}, nil
	}

	gamestate, err := loadGameState(ctx, match.Id)
	if err != nil {
		return nil, err
	}
	if gamestate.IsEnded {
		return &zb_calls.KeepAliveResponse{}, nil // just ignore for client check
	}

	gp, err := GamePlayFrom(gamestate, match.UseBackendGameLogic, match.PlayerDebugCheats)
	if err != nil {
		return nil, err
	}
	for _, lastseen := range match.PlayerLastSeens {
		lastSeenAt := time.Unix(lastseen.UpdatedAt, 0)
		if lastSeenAt.Add(KeepAliveTimeout).Before(ctx.Now()) {
			// create a leave match request and append to the game state
			leaveMatchAction := zb_data.PlayerAction{
				ActionType: zb_enums.PlayerActionType_LeaveMatch,
				PlayerId:   lastseen.Id,
				Action: &zb_data.PlayerAction_LeaveMatch{
					LeaveMatch: &zb_data.PlayerActionLeaveMatch{
						Reason: zb_data.PlayerActionLeaveMatch_KeepAliveTimeout,
					},
				},
				CreatedAt: ctx.Now().Unix(),
			}

			// ignore the error in case this method is called mutiple times
			if err := gp.AddAction(&leaveMatchAction); err == nil {
				if err := saveGameState(ctx, gamestate); err != nil {
					return nil, err
				}
			}
			// update match status
			match.Status = zb_data.Match_PlayerLeft
			if err := saveMatch(ctx, match); err != nil {
				return nil, err
			}
			// update winner
			leaveMatchReq := leaveMatchAction.GetLeaveMatch()
			leaveMatchReq.Winner = gp.State.Winner
			emitMsg := zb_data.PlayerActionEvent{
				PlayerAction:     &leaveMatchAction,
				Match:            match,
				Block:            &zb_data.History{List: gp.history},
				CreatedByBackend: true,
			}
			data, err := proto.Marshal(&emitMsg)
			if err != nil {
				return nil, err
			}
			ctx.EmitTopics(data, match.Topics...)
		}
	}

	return &zb_calls.KeepAliveResponse{}, nil
}

func (z *ZombieBattleground) GetContractState(ctx contract.StaticContext, req *zb_calls.EmptyRequest) (*zb_calls.GetContractStateResponse, error) {
	err := z.validateOracle(ctx)
	if err != nil {
		return nil, err
	}

	state, err := loadContractState(ctx)
	if err != nil {
		return nil, err
	}

	return &zb_calls.GetContractStateResponse{
		State: state,
	}, nil
}

func (z *ZombieBattleground) GetContractConfiguration(ctx contract.StaticContext, req *zb_calls.EmptyRequest) (*zb_calls.GetContractConfigurationResponse, error) {
	err := z.validateOracle(ctx)
	if err != nil {
		return nil, err
	}

	configuration, err := loadContractConfiguration(ctx)
	if err != nil {
		return nil, err
	}

	return &zb_calls.GetContractConfigurationResponse{
		Configuration: configuration,
	}, nil
}

func (z *ZombieBattleground) UpdateContractConfiguration(ctx contract.Context, req *zb_calls.UpdateContractConfigurationRequest) error {
	err := z.validateOracle(ctx)
	if err != nil {
		return err
	}

	configuration, err := loadContractConfiguration(ctx)
	if err != nil {
		if errors.Cause(err).Error() == ErrNotFound.Error() {
			configuration = &zb_data.ContractConfiguration{
				InitialFiatPurchaseTxId:     battleground_utility.MarshalBigIntProto(big.NewInt(0)),
				FiatPurchaseContractVersion: 0,
			}
		} else {
			return err
		}
	}

	state, err := loadContractState(ctx)
	if err != nil {
		if errors.Cause(err).Error() == ErrNotFound.Error() {
			state = &zb_data.ContractState{
				LastPlasmachainBlockNumber: 0,
				CurrentFiatPurchaseTxId:    battleground_utility.MarshalBigIntProto(big.NewInt(0)),
			}
		} else {
			return err
		}
	}

	changed := false
	if req.SetFiatPurchaseContractVersion {
		changed = true
		configuration.FiatPurchaseContractVersion = req.FiatPurchaseContractVersion
	}

	if req.SetInitialFiatPurchaseTxId {
		changed = true
		if req.InitialFiatPurchaseTxId == nil {
			return fmt.Errorf("InitialFiatPurchaseTxId == nil")
		}

		if req.InitialFiatPurchaseTxId.Value.Int.Cmp(configuration.InitialFiatPurchaseTxId.Value.Int) != 0 {
			configuration.InitialFiatPurchaseTxId = req.InitialFiatPurchaseTxId
			state.CurrentFiatPurchaseTxId = req.InitialFiatPurchaseTxId

			ctx.Logger().Info("txId reset", "configuration.InitialFiatPurchaseTxId", configuration.InitialFiatPurchaseTxId, "state.CurrentFiatPurchaseTxId", state.CurrentFiatPurchaseTxId)
		}
	}

	if req.SetUseCardLibraryAsUserCollection {
		changed = true
		configuration.UseCardLibraryAsUserCollection = req.UseCardLibraryAsUserCollection
	}

	if req.SetDataWipeConfiguration {
		changed = true
		found := false
		for _, existingDataWipeConfiguration := range configuration.DataWipeConfiguration {
			if existingDataWipeConfiguration.Version == req.DataWipeConfiguration.Version {
				existingDataWipeConfiguration.Reset()
				proto.Merge(existingDataWipeConfiguration, req.DataWipeConfiguration)
				found = true
				break
			}
		}

		if !found {
			configuration.DataWipeConfiguration = append(configuration.DataWipeConfiguration, req.DataWipeConfiguration)
		}
	}

	if req.SetCardCollectionSyncDataVersion {
		changed = true
		if req.CardCollectionSyncDataVersion == "" {
			return ErrVersionNotSet
		}

		configuration.CardCollectionSyncDataVersion = req.CardCollectionSyncDataVersion
	}

	if !changed {
		return fmt.Errorf("no configuration changes specified")
	}

	err = saveContractState(ctx, state)
	if err != nil {
		return err
	}

	err = saveContractConfiguration(ctx, configuration)
	if err != nil {
		return err
	}

	return nil
}

func (z *ZombieBattleground) UpdateOracle(ctx contract.Context, req *zb_calls.UpdateOracleRequest) error {
	var oldOraclePB types.Address
	err := ctx.Get(oracleKey, &oldOraclePB)
	if err != nil {
		return nil
	}

	oldOracle := loom.UnmarshalAddressPB(&oldOraclePB)
	newOraclePB := req.NewOracle
	newOracle := loom.UnmarshalAddressPB(newOraclePB)
	if ctx.Has(oracleKey) {
		if oldOracle.String() == newOracle.String() {
			return errors.New("cannot set new oracle to same address as old oracle")
		}
		if err := z.validateOracle(ctx); err != nil {
			return errors.Wrap(err, "sender is not the current oracle")
		}

		ctx.GrantPermissionTo(oldOracle, []byte(oldOraclePB.String()), "old-oracle")
	}
	ctx.GrantPermissionTo(newOracle, []byte(newOraclePB.String()), OracleRole)

	if err := ctx.Set(oracleKey, newOraclePB); err != nil {
		return errors.Wrap(err, "failed to set new oracle")
	}
	return nil
}

func (z *ZombieBattleground) GetGameMode(ctx contract.StaticContext, req *zb_calls.GetGameModeRequest) (*zb_data.GameMode, error) {
	gameModeList, err := loadGameModeList(ctx) // we get the game mode list first, because deleted modes won't be in there
	if err != nil {
		return nil, err
	}
	gameMode := getGameModeFromList(gameModeList, req.ID)
	if gameMode == nil {
		return nil, contract.ErrNotFound
	}

	return gameMode, nil
}

func (z *ZombieBattleground) CallCustomGameModeFunction(ctx contract.Context, req *zb_calls.CallCustomGameModeFunctionRequest) error {
	err := NewCustomGameMode(loom.Address{
		ChainID: req.Address.ChainId,
		Local:   req.Address.Local,
	}).CallFunction(ctx, req.CallData)

	if err != nil {
		return err
	}

	return nil
}

func (z *ZombieBattleground) GetGameModeCustomUi(ctx contract.StaticContext, req *zb_calls.GetCustomGameModeCustomUiRequest) (*zb_calls.GetCustomGameModeCustomUiResponse, error) {
	uiElements, err := NewCustomGameMode(loom.Address{
		ChainID: req.Address.ChainId,
		Local:   req.Address.Local,
	}).GetCustomUi(ctx)

	if err != nil {
		return nil, err
	}

	response := &zb_calls.GetCustomGameModeCustomUiResponse{
		UiElements: uiElements,
	}

	return response, nil
}

func (z *ZombieBattleground) ListGameModes(ctx contract.StaticContext, req *zb_calls.ListGameModesRequest) (*zb_data.GameModeList, error) {
	gameModeList, err := loadGameModeList(ctx)
	if err != nil {
		return nil, err
	}

	return gameModeList, nil
}

func validateGameModeReq(req *zb_calls.GameModeRequest) error {
	if req.Name == "" {
		return errors.New("GameMode name cannot be empty")
	}
	if utf8.RuneCountInString(req.Name) > MaxGameModeNameChar {
		return errors.New("GameMode name too long")
	}
	if req.Description == "" {
		return errors.New("GameMode Description cannot be empty")
	}
	if utf8.RuneCountInString(req.Description) > MaxGameModeDescriptionChar {
		return errors.New("GameMode Description too long")
	}
	if req.Version == "" {
		return errors.New("GameMode Version cannot be empty")
	}
	if utf8.RuneCountInString(req.Version) > MaxGameModeVersionChar {
		return errors.New("GameMode Version too long")
	}
	if req.Address == "" {
		return errors.New("GameMode address cannot be empty")
	}

	return nil
}

func (z *ZombieBattleground) AddGameMode(ctx contract.Context, req *zb_calls.GameModeRequest) (*zb_data.GameMode, error) {
	if err := validateGameModeReq(req); err != nil {
		return nil, err
	}

	// check if game mode with this name already exists
	gameModeList, err := loadGameModeList(ctx)
	if err != nil && err == contract.ErrNotFound {
		gameModeList = &zb_data.GameModeList{GameModes: []*zb_data.GameMode{}}
	}
	if gameMode := getGameModeFromListByName(gameModeList, req.Name); gameMode != nil {
		return nil, errors.New("A game mode with that name already exists")
	}

	// create a GUID from the hash of gameMode name and address
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(req.Name + req.Address))
	ID := hex.EncodeToString(h.Sum(nil))

	addr, err := loom.LocalAddressFromHexString(req.Address)
	if err != nil {
		return nil, err
	}

	owner := &types.Address{ChainId: ctx.ContractAddress().ChainID, Local: ctx.Message().Sender.Local}
	gameModeType := zb_data.GameModeType_Community

	// if request was made with a valid oracle, set type and owner to Loom
	var oldOraclePB types.Address
	err = ctx.Get(oracleKey, &oldOraclePB)
	oracleNotSet := err != nil && err.Error() == ErrNotFound.Error()

	if err := z.validateOracle(ctx); !oracleNotSet && err == nil {
		gameModeType = zb_data.GameModeType_Loom
		owner = loom.RootAddress(ctx.ContractAddress().ChainID).MarshalPB()
	}

	gameMode := &zb_data.GameMode{
		ID:           ID,
		Name:         req.Name,
		Description:  req.Description,
		Version:      req.Version,
		Address:      &types.Address{ChainId: ctx.ContractAddress().ChainID, Local: addr},
		Owner:        owner,
		GameModeType: gameModeType,
	}

	ctx.GrantPermission([]byte(ID), []string{OwnerRole})

	if err := addGameModeToList(ctx, gameMode); err != nil {
		return nil, err
	}

	return gameMode, nil
}

func (z *ZombieBattleground) UpdateGameMode(ctx contract.Context, req *zb_calls.UpdateGameModeRequest) (*zb_data.GameMode, error) {
	// Require either oracle or owner permission to delete a game mode
	err := z.validateOracle(ctx)
	if err == ErrOracleNotVerified {
		if ok, _ := ctx.HasPermission([]byte(req.ID), []string{OwnerRole}); !ok {
			return nil, ErrUserNotVerified
		}
	}

	gameModeList, err := loadGameModeList(ctx)
	if err != nil {
		return nil, err
	}
	gameMode := getGameModeFromList(gameModeList, req.ID)
	if gameMode == nil {
		return nil, contract.ErrNotFound
	}

	if req.Name != "" {
		if utf8.RuneCountInString(req.Name) > MaxGameModeNameChar {
			return nil, errors.New("GameMode name too long")
		}
		gameMode.Name = req.Name
	}
	if req.Description != "" {
		if utf8.RuneCountInString(req.Description) > MaxGameModeDescriptionChar {
			return nil, errors.New("GameMode Description too long")
		}
		gameMode.Description = req.Description
	}
	if req.Version != "" {
		if utf8.RuneCountInString(req.Version) > MaxGameModeVersionChar {
			return nil, errors.New("GameMode Version too long")
		}
		gameMode.Version = req.Version
	}
	if req.Address != "" {
		addr, err := loom.LocalAddressFromHexString(req.Address)
		if err != nil {
			return nil, err
		}
		gameMode.Address = &types.Address{ChainId: ctx.ContractAddress().ChainID, Local: addr}
	}

	if err = saveGameModeList(ctx, gameModeList); err != nil {
		return nil, err
	}

	return gameMode, nil
}

func (z *ZombieBattleground) DeleteGameMode(ctx contract.Context, req *zb_calls.DeleteGameModeRequest) error {
	// Require either oracle or owner permission to delete a game mode
	err := z.validateOracle(ctx)
	if err == ErrOracleNotVerified {
		if ok, _ := ctx.HasPermission([]byte(req.ID), []string{OwnerRole}); !ok {
			return ErrUserNotVerified
		}
	}

	gameModeList, err := loadGameModeList(ctx)
	if err != nil {
		return err
	}

	var deleted bool
	gameModeList, deleted = deleteGameMode(gameModeList, req.ID)
	if !deleted {
		return fmt.Errorf("game mode not found")
	}

	if err := saveGameModeList(ctx, gameModeList); err != nil {
		return err
	}

	return nil
}

func (z *ZombieBattleground) ProcessOracleEventBatch(ctx contract.Context, req *orctype.ProcessOracleEventBatchRequest) error {
	state, err := loadContractState(ctx)
	if err != nil {
		return err
	}

	eventCount := 0           // number of blocks that were actually processed in this batch
	lastEthBlock := uint64(0) // the last block processed in this batch

	addressToCardKeyToAmountChangesMap := addressToCardKeyToAmountChangesMap{}

	for _, ev := range req.Events {
		// Events in the batch are expected to be ordered by block, so a batch should contain
		// events from block N, followed by events from block N+1, any other order is invalid.
		if ev.EthBlock < lastEthBlock {
			ctx.Logger().Error("Oracle invalid event batch, block has already been processed",
				"block", ev.EthBlock)
			return ErrInvalidEventBatch
		}

		// Multiple validators might submit batches with overlapping block ranges because the
		// Gateway oracles will fetch events from Plasmachian at different times, with different
		// latencies, etc. Simply skip blocks that have already been processed.
		if ev.EthBlock <= state.LastPlasmachainBlockNumber {
			continue
		}

		switch payload := ev.Payload.(type) {
		case *orctype.PlasmachainEvent_Transfer:
			ctx.Logger().Info("got Transfer event")
			err := z.updateCardAmountChangeToAddressToCardAmountsChangeMap(
				addressToCardKeyToAmountChangesMap,
				payload.Transfer.From.Local,
				payload.Transfer.To.Local,
				cardKeyFromCardTokenId(payload.Transfer.TokenId.Value.Int64()),
				1,
				loom.UnmarshalAddressPB(req.ZbgCardContractAddress),
			)

			if err != nil {
				err = errors.Wrap(err, "error handling Transfer event")
				ctx.Logger().Error(err.Error())
				return err
			}
		case *orctype.PlasmachainEvent_TransferWithQuantity:
			ctx.Logger().Debug("got TransferWithQuantity event")
			err := z.updateCardAmountChangeToAddressToCardAmountsChangeMap(
				addressToCardKeyToAmountChangesMap,
				payload.TransferWithQuantity.From.Local,
				payload.TransferWithQuantity.To.Local,
				cardKeyFromCardTokenId(payload.TransferWithQuantity.TokenId.Value.Int64()),
				uint(payload.TransferWithQuantity.Amount.Value.Uint64()),
				loom.UnmarshalAddressPB(req.ZbgCardContractAddress),
			)

			if err != nil {
				err = errors.Wrap(err, "error handling TransferWithQuantity event")
				ctx.Logger().Error(err.Error())
				return err
			}
		case *orctype.PlasmachainEvent_BatchTransfer:
			ctx.Logger().Info("got TransferWithQuantity event")

			for index, cardTokenId := range payload.BatchTransfer.TokenIds {
				amount := payload.BatchTransfer.Amounts[index]
				err := z.updateCardAmountChangeToAddressToCardAmountsChangeMap(
					addressToCardKeyToAmountChangesMap,
					payload.BatchTransfer.From.Local,
					payload.BatchTransfer.To.Local,
					cardKeyFromCardTokenId(cardTokenId.Value.Int64()),
					uint(amount.Value.Uint64()),
					loom.UnmarshalAddressPB(req.ZbgCardContractAddress),
				)

				if err != nil {
					err = errors.Wrap(err, "error handling TransferWithQuantity event")
					ctx.Logger().Error(err.Error())
					return err
				}
			}
		case nil:
			ctx.Logger().Error("Oracle missing event payload")
			continue
		default:
			ctx.Logger().Error("Oracle unknown event payload type %T", payload)
			continue
		}

		if ev.EthBlock > lastEthBlock {
			lastEthBlock = ev.EthBlock
		}

		eventCount++
	}

	if debugEnabled {
		for address, cardAmountsChangeMap := range addressToCardKeyToAmountChangesMap {
			fmt.Printf("Address %s card amount changes:\n", address)
			for cardKey, amountChange := range cardAmountsChangeMap {
				fmt.Printf("   (%s) = %d\n", cardKey.String(), amountChange)
			}
		}
	}

	// If there are no new events in this batch return an error so that the batch tx isn't
	// propagated to the other nodes.
	if eventCount == 0 {
		return fmt.Errorf("no new events found in the batch")
	}

	err = z.applyAddressToCardAmountsChangeMapDelta(ctx, addressToCardKeyToAmountChangesMap)
	if err != nil {
		return errors.Wrap(err, "failed to apply address to card amounts change map")
	}

	ctx.Logger().Debug("setting last Plasmachain block", "LastPlasmachainBlockNumber", req.LastPlasmachainBlockNumber)

	state.LastPlasmachainBlockNumber = req.LastPlasmachainBlockNumber
	return saveContractState(ctx, state)
}

func (z *ZombieBattleground) SetLastPlasmaBlockNumber(ctx contract.Context, req *zb_calls.SetLastPlasmaBlockNumberRequest) error {
	err := z.validateOracle(ctx)
	if err != nil {
		return err
	}

	state, err := loadContractState(ctx)
	if err != nil {
		return err
	}

	ctx.Logger().Debug("setting last Plasmachain block", "LastPlasmachainBlockNumber", req.LastPlasmachainBlockNumber)
	state.LastPlasmachainBlockNumber = req.LastPlasmachainBlockNumber
	return saveContractState(ctx, state)
}

func (z *ZombieBattleground) GetContractBuildMetadata(ctx contract.StaticContext, req *zb_calls.GetContractBuildMetadataRequest) (*zb_calls.GetContractBuildMetadataResponse, error) {
	return &zb_calls.GetContractBuildMetadataResponse{
		Date:   BuildDate,
		GitSha: BuildGitSha,
		Build:  BuildNumber,
	}, nil
}

func (z *ZombieBattleground) GetPendingMintingTransactionReceipts(ctx contract.StaticContext, req *zb_calls.GetPendingMintingTransactionReceiptsRequest) (*zb_calls.GetPendingMintingTransactionReceiptsResponse, error) {
	if !isOwner(ctx, req.UserId) {
		return nil, ErrUserNotVerified
	}

	receiptCollection, err := loadPendingMintingTransactionReceipts(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &zb_calls.GetPendingMintingTransactionReceiptsResponse{
		ReceiptCollection: receiptCollection,
	}, nil
}

func (z *ZombieBattleground) ConfirmPendingMintingTransactionReceipt(ctx contract.Context, req *zb_calls.ConfirmPendingMintingTransactionReceiptRequest) error {
	if !isOwner(ctx, req.UserId) {
		return ErrUserNotVerified
	}

	receiptCollection, err := loadPendingMintingTransactionReceipts(ctx, req.UserId)
	if err != nil {
		return err
	}

	found := false
	for i, receipt := range receiptCollection.Receipts {
		if receipt.TxId.Value.Int.Cmp(req.TxId.Value.Int) == 0 {
			receiptCollection.Receipts = append(receiptCollection.Receipts[:i], receiptCollection.Receipts[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("receipt for txId %s not found", req.TxId.Value.String())
	}

	err = savePendingMintingTransactionReceipts(ctx, req.UserId, receiptCollection)
	if err != nil {
		return err
	}

	return nil
}

func (z *ZombieBattleground) DebugMintBoosterPackReceipt(ctx contract.Context, req *zb_calls.DebugMintBoosterPackReceiptRequest) (*zb_calls.DebugMintBoosterPackReceiptResponse, error) {
	if !debugEnabled {
		return nil, ErrDebugNotEnabled
	}

	userIdString := defaultUserIdPrefix + req.UserId.Value.Int.String()
	receipt, err := mintBoosterPacksAndSave(ctx, userIdString, req.UserId.Value.Int, uint(req.BoosterAmount))
	if err != nil {
		return nil, err
	}

	receiptJson, err := json.Marshal(receipt)
	if err != nil {
		return nil, err
	}

	return &zb_calls.DebugMintBoosterPackReceiptResponse{
		ReceiptJson: string(receiptJson),
		Receipt:     receipt.MarshalPB(),
	}, nil
}

func (z *ZombieBattleground) DebugGetUserIdByAddress(ctx contract.StaticContext, req *zb_calls.DebugGetUserIdByAddressRequest) (*zb_data.UserIdContainer, error) {
	found, userId, err := loadUserIdByAddress(ctx, loom.UnmarshalAddressPB(req.Address))
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, ErrNotFound
	}

	return &zb_data.UserIdContainer{
		UserId: userId,
	}, nil
}

func (z *ZombieBattleground) GetOracleCommandRequestList(ctx contract.StaticContext, req *orctype.GetOracleCommandRequestListRequest) (*orctype.GetOracleCommandRequestListResponse, error) {
	err := z.validateOracle(ctx)
	if err != nil {
		return nil, err
	}

	list, err := loadOracleCommandRequestList(ctx)
	if err != nil {
		return nil, err
	}

	return &orctype.GetOracleCommandRequestListResponse{
		CommandRequests: list.Commands,
	}, nil
}

func (z *ZombieBattleground) ProcessOracleCommandResponseBatch(ctx contract.Context, req *orctype.ProcessOracleCommandResponseBatchRequest) (*zb_calls.EmptyResponse, error) {
	err := z.validateOracle(ctx)
	if err != nil {
		return nil, err
	}

	err = z.processOracleCommandResponseBatchInternal(ctx, req.CommandResponses)
	if err != nil {
		return nil, err
	}

	return &zb_calls.EmptyResponse{}, nil
}

func (z *ZombieBattleground) RequestUserFullCardCollectionSync(ctx contract.Context, req *zb_calls.RequestUserFullCardCollectionSyncRequest) (*zb_calls.EmptyResponse, error) {
	err := z.isOwnerOrOracle(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	err = z.addGetUserFullCardCollectionOracleCommand(ctx, ctx.Message().Sender)
	if err != nil {
		return nil, err
	}

	return &zb_calls.EmptyResponse{}, nil
}

func (z *ZombieBattleground) DebugGetPendingCardAmountChangeItems(ctx contract.StaticContext, req *zb_calls.DebugGetPendingCardAmountChangeItemsRequest) (*zb_calls.DebugGetPendingCardAmountChangeItemsResponse, error) {
	container, err := loadPendingCardAmountChangesContainerByAddress(ctx, loom.UnmarshalAddressPB(req.Address))
	if err != nil {
		return nil, err
	}

	return &zb_calls.DebugGetPendingCardAmountChangeItemsResponse{
		Container: container,
	}, nil
}

func (z *ZombieBattleground) validateOracle(ctx contract.StaticContext) error {
	var oldOraclePB types.Address
	err := ctx.Get(oracleKey, &oldOraclePB)
	if err != nil && err.Error() == ErrNotFound.Error() {
		return nil
	}

	if ok, _ := ctx.HasPermissionFor(ctx.Message().Sender, []byte(ctx.Message().Sender.MarshalPB().String()), []string{OracleRole}); !ok {
		return ErrOracleNotVerified
	}

	if ok, _ := ctx.HasPermissionFor(ctx.Message().Sender, []byte(ctx.Message().Sender.MarshalPB().String()), []string{"old-oracle"}); ok {
		return errors.New("This oracle is expired. Please use latest oracle")
	}
	return nil
}

func (z *ZombieBattleground) initializeUserDefaultCardCollection(ctx contract.Context, version string, userId string) error {
	// add default collection list
	defaultCardCollection, err := loadDefaultCardCollection(ctx, version)
	if err != nil {
		return errors.Wrap(err, "error initializing user default card collection")
	}

	if err := saveUserCardCollection(ctx, userId, defaultCardCollection); err != nil {
		return errors.Wrap(err, "error initializing user default card collection")
	}

	return nil
}

func (z *ZombieBattleground) initializeUserDefaultDecks(ctx contract.Context, version string, userId string) (decks *zb_data.DeckList, err error) {
	decks, err = loadDefaultDecks(ctx, version)
	if err != nil {
		return nil, errors.Wrap(err, "error initializing user default decks")
	}

	// update default deck with none-zero id
	for i := 0; i < len(decks.Decks); i++ {
		decks.Decks[i].Id = int64(i + 1)
	}
	if err := saveDecks(ctx, version, userId, decks); err != nil {
		return nil, errors.Wrap(err, "error initializing user default decks")
	}

	return decks, nil
}

func (z *ZombieBattleground) handleUserDataWipe(ctx contract.Context, version string, userId string) (wipeExecuted bool, err error) {
	wipeExecuted = false
	configuration, err := loadContractConfiguration(ctx)
	if err != nil {
		return false, errors.Wrap(err, "error handling user data wipe")
	}

	userPersistentData, err := loadUserPersistentData(ctx, userId)
	if err != nil {
		return false, errors.Wrap(err, "error handling user data wipe")
	}

	var matchingDataWipeConfiguration *zb_data.DataWipeConfiguration
loop:
	for _, dataWipeItem := range configuration.DataWipeConfiguration {
		// Check if a wipe is configured for this data version
		if dataWipeItem.Version == version {
			// Check if wipe for this version was already executed for this user
			for _, alreadyExecutedWipeVersion := range userPersistentData.ExecutedDataWipesVersions {
				if alreadyExecutedWipeVersion == version {
					break loop
				}
			}

			matchingDataWipeConfiguration = dataWipeItem
		}
	}

	if matchingDataWipeConfiguration != nil {
		wipeExecuted = true

		if matchingDataWipeConfiguration.WipeDecks {
			_, err = z.initializeUserDefaultDecks(ctx, version, userId)
			if err != nil {
				return false, errors.Wrap(err, "error wiping user decks")
			}
		}

		if matchingDataWipeConfiguration.WipeOverlordUserInstances {
			emptyOverlordUserInstances := zb_data.OverlordUserDataList{
				OverlordsUserData: []*zb_data.OverlordUserData{},
			}
			err = saveOverlordUserDataList(ctx, userId, &emptyOverlordUserInstances)
			if err != nil {
				return false, errors.Wrap(err, "error wiping overlord user data")
			}
		}

		userPersistentData.ExecutedDataWipesVersions =
			append(userPersistentData.ExecutedDataWipesVersions, matchingDataWipeConfiguration.Version)
		err = saveUserPersistentData(ctx, userId, userPersistentData)
		if err != nil {
			return false, errors.Wrap(err, "error handling user data wipe")
		}
	}

	return wipeExecuted, nil
}

func (z *ZombieBattleground) isOwnerOrOracle(ctx contract.StaticContext, userId string) error {
	if !isOwner(ctx, userId) {
		err := z.validateOracle(ctx)
		if err != nil {
			return errors.Wrap(ErrNotOwnerOrOracleNotVerified, err.Error())
		}
		return ErrUserNotVerified
	}

	return nil
}

func applyExperience(
	ctx contract.Context,
	version string,
	overlordLevelingData *zb_data.OverlordLevelingData,
	userId string,
	userIdInt *big.Int,
	overlordId int64,
	experience int64,
	deckId int64,
	isWin bool,
) error {
	overlordUserInstances, err := loadOverlordUserInstances(ctx, version, userId)
	if err != nil {
		return err
	}

	overlord, found := getOverlordUserInstanceByPrototypeId(overlordUserInstances, overlordId)
	if !found {
		return fmt.Errorf("overlord with prototype id %d not found", overlordId)
	}

	if err := applyExperienceInternal(ctx, userId, userIdInt, overlordLevelingData, overlordUserInstances, overlord, experience, deckId, isWin); err != nil {
		return errors.Wrap(err, "failed to apply experience")
	}

	return nil
}

func applyExperienceInternal(
	ctx contract.Context,
	userId string,
	userIdInt *big.Int,
	overlordLevelingData *zb_data.OverlordLevelingData,
	overlordUserInstances []*zb_data.OverlordUserInstance,
	targetOverlordUserInstance *zb_data.OverlordUserInstance,
	matchExperience int64,
	deckId int64,
	isWin bool,
) error {
	oldExperience := targetOverlordUserInstance.UserData.Experience
	oldLevel := int32(targetOverlordUserInstance.UserData.Level)

	targetOverlordUserInstance.UserData.Experience += matchExperience
	newLevel := calculateOverlordLevel(overlordLevelingData, targetOverlordUserInstance.UserData)
	levelRewards := make([]*zb_data.LevelReward, 0)
	if newLevel > int32(targetOverlordUserInstance.UserData.Level) {
		targetOverlordUserInstance.UserData.Level = int64(newLevel)

		// Get rewards for all in-between levels
		for level := oldLevel + 1; level <= newLevel; level++ {
			levelReward := getLevelReward(overlordLevelingData, level)
			if levelReward != nil {
				levelRewards = append(levelRewards, levelReward)
			}
		}

		for i := 0; i < len(levelRewards); i++ {
			// skill rewards
			switch reward := levelRewards[i].Reward.(type) {
			case *zb_data.LevelReward_SkillReward:
				skillReward := reward.SkillReward
				var skillToUnlock *zb_data.OverlordSkillPrototype = nil
				for j := 0; j < len(targetOverlordUserInstance.Prototype.Skills); j++ {
					if j == int(skillReward.SkillIndex) {
						skillToUnlock = targetOverlordUserInstance.Prototype.Skills[j]
						break
					}
				}

				if skillToUnlock == nil {
					return fmt.Errorf("failed to find skill for reward")
				}

				alreadyUnlocked := false
				for _, unlockedSkillId := range targetOverlordUserInstance.UserData.UnlockedSkillIds {
					if unlockedSkillId == skillToUnlock.Id {
						alreadyUnlocked = true
						break
					}
				}

				if !alreadyUnlocked {
					targetOverlordUserInstance.UserData.UnlockedSkillIds = append(targetOverlordUserInstance.UserData.UnlockedSkillIds, skillToUnlock.Id)
				}

				// TODO: Update decks with no skills for players convenience?
			case *zb_data.LevelReward_BoosterPackReward:
				boosterPackReward := reward.BoosterPackReward

				_, err := mintBoosterPacksAndSave(ctx, userId, userIdInt, uint(boosterPackReward.Amount))
				if err != nil {
					return err
				}
			}
		}
	}

	overlordUserDataList := &zb_data.OverlordUserDataList{
		OverlordsUserData: getOverlordsUserDataFromOverlordUserInstances(overlordUserInstances),
	}
	if err := saveOverlordUserDataList(ctx, userId, overlordUserDataList); err != nil {
		return err
	}

	// Set the notification
	notifications, err := loadUserNotifications(ctx, userId)
	if err != nil {
		return err
	}

	//  There can only be one EndMatch notification at any time
loop:
	for _, notification := range notifications.Notifications {
		switch notification.Type {
		case zb_data.NotificationType_EndMatch:
			notifications.Notifications, err = removeNotification(notifications.Notifications, notification.Id)
			if err != nil {
				return err
			}

			break loop
		}
	}

	notification := createBaseNotification(ctx, notifications.Notifications, zb_data.NotificationType_EndMatch)
	notification.Notification = &zb_data.Notification_EndMatch{
		EndMatch: &zb_data.NotificationEndMatch{
			OverlordId:    targetOverlordUserInstance.Prototype.Id,
			OldExperience: oldExperience,
			OldLevel:      oldLevel,
			NewExperience: targetOverlordUserInstance.UserData.Experience,
			NewLevel:      int32(targetOverlordUserInstance.UserData.Level),
			Rewards:       levelRewards,
			IsWin:         isWin,
			DeckId:        deckId,
		},
	}

	notifications.Notifications = append(notifications.Notifications, notification)
	if err := saveUserNotifications(ctx, userId, notifications); err != nil {
		return err
	}

	return nil
}

var Contract plugin.Contract = contract.MakePluginContract(&ZombieBattleground{})

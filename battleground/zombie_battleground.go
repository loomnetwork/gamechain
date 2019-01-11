package battleground

import (
	"crypto/ecdsa"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gogo/protobuf/proto"
	orctype "github.com/loomnetwork/gamechain/types/oracle"
	"github.com/loomnetwork/gamechain/types/zb"
	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/plugin"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	"github.com/loomnetwork/go-loom/types"
	solsha3 "github.com/miguelmota/go-solidity-sha3"
	"github.com/pkg/errors"
)

const (
	MaxGameModeNameChar           = 48
	MaxGameModeDescriptionChar    = 255
	MaxGameModeVersionChar        = 16
	TurnTimeout                   = 120 * time.Second
	KeepAliveTimeout              = 60 * time.Second // client keeps sending keepalive every 30 second. have to make sure we have some buffer for network delays
	RewardTypeTutorialCompleted   = "tutorial-completed"
	TutorialRewardContractVersion = 1
)

const (
	TopicCreateAccountEvent      = "zombiebattleground:createaccount"
	TopicUpdateAccountEvent      = "zombiebattleground:updateaccount"
	TopicCreateDeckEvent         = "zombiebattleground:createdeck"
	TopicEditDeckEvent           = "zombiebattleground:editdeck"
	TopicDeleteDeckEvent         = "zombiebattleground:deletedeck"
	TopicAddHeroExpEvent         = "zombiebattleground:addheroexperience"
	TopicRegisterPlayerPoolEvent = "zombiebattleground:registerplayerpool"
	TopicFindMatchEvent          = "zombiebattleground:findmatch"
	TopicAcceptMatchEvent        = "zombiebattleground:acceptmatch"
	// match pattern match:id e.g. match:1, match:2, ...
	TopicMatchEventPrefix = "match:"
)

const (
	OracleRole = "oracle"
	OwnerRole  = "user"
)

var (
	// secret
	secret string
	// privateKey to sign reward
	privateKeyStr string
	// Error list
	ErrOracleNotSpecified = errors.New("oracle not specified")
	ErrInvalidEventBatch  = errors.New("invalid event batch")
)

type ZombieBattleground struct {
}

func (z *ZombieBattleground) Meta() (plugin.Meta, error) {
	return plugin.Meta{
		Name:    "ZombieBattleground",
		Version: "1.0.0",
	}, nil
}

func (z *ZombieBattleground) Init(ctx contract.Context, req *zb.InitRequest) error {
	secret = os.Getenv("SECRET_KEY")
	if secret == "" {
		secret = "justsowecantestwithoutenvvar"
	}

	privateKeyStr = os.Getenv("GAMECHAIN_PRIVATE_KEY")

	if req.Oracle != nil {
		ctx.GrantPermissionTo(loom.UnmarshalAddressPB(req.Oracle), []byte(req.Oracle.String()), OracleRole)
		if err := ctx.Set(oracleKey, req.Oracle); err != nil {
			return errors.Wrap(err, "Error setting oracle")
		}
	}

	// init state
	state := zb.GamechainState{
		LastPlasmachainBlockNum: 1,
	}
	if err := saveState(ctx, &state); err != nil {
		return err
	}

	// initialize card library
	cardList := zb.CardList{
		Cards: req.Cards,
	}
	if err := ctx.Set(MakeVersionedKey(req.Version, cardListKey), &cardList); err != nil {
		return err
	}

	// initialize heroes
	heroList := zb.HeroList{
		Heroes: req.Heroes,
	}
	if err := ctx.Set(MakeVersionedKey(req.Version, heroListKey), &heroList); err != nil {
		return err
	}

	// initialize card collection
	cardCollectionList := zb.CardCollectionList{
		Cards: req.DefaultCollection,
	}
	if err := ctx.Set(MakeVersionedKey(req.Version, defaultCollectionKey), &cardCollectionList); err != nil {
		return err
	}

	// initialize default deck
	deckList := zb.DeckList{
		Decks: req.DefaultDecks,
	}
	if err := ctx.Set(MakeVersionedKey(req.Version, defaultDeckKey), &deckList); err != nil {
		return err
	}

	// initialize default heroes
	defaultHeroList := zb.HeroList{
		Heroes: req.Heroes,
	}
	if err := ctx.Set(MakeVersionedKey(req.Version, defaultHeroesKey), &defaultHeroList); err != nil {
		return err
	}

	// initialize AI decks
	aiDeckList := zb.AIDeckList{
		Decks: req.AiDecks,
	}
	if err := ctx.Set(MakeVersionedKey(req.Version, aiDecksKey), &aiDeckList); err != nil {
		return err
	}

	return nil
}

func (z *ZombieBattleground) UpdateInit(ctx contract.Context, req *zb.UpdateInitRequest) error {
	var heroList zb.HeroList
	var defaultHeroList zb.HeroList
	var cardList zb.CardList
	var cardCollectionList zb.CardCollectionList
	var deckList zb.DeckList

	// initialize card library
	cardList.Cards = req.Cards
	if req.Cards == nil {
		if req.OldVersion != "" {
			if err := ctx.Get(MakeVersionedKey(req.OldVersion, cardListKey), &cardList); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("'cards' key missing, old version not specified")
		}
	}
	if err := ctx.Set(MakeVersionedKey(req.Version, cardListKey), &cardList); err != nil {
		return err
	}

	// initialize heroes
	heroList.Heroes = req.Heroes
	defaultHeroList.Heroes = req.Heroes
	if req.Heroes == nil {
		if req.OldVersion != "" {
			if err := ctx.Get(MakeVersionedKey(req.OldVersion, heroListKey), &heroList); err != nil {
				return err
			}
			if err := ctx.Get(MakeVersionedKey(req.OldVersion, defaultHeroesKey), &defaultHeroList); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("'heroes' key missing, old version not specified")
		}

	}
	if err := ctx.Set(MakeVersionedKey(req.Version, heroListKey), &heroList); err != nil {
		return err
	}
	if err := ctx.Set(MakeVersionedKey(req.Version, defaultHeroesKey), &defaultHeroList); err != nil {
		return err
	}

	// initialize default collection
	cardCollectionList.Cards = req.DefaultCollection
	if req.DefaultCollection == nil {
		if req.OldVersion != "" {
			if err := ctx.Get(MakeVersionedKey(req.OldVersion, defaultCollectionKey), &cardCollectionList); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("'default_collection' key missing, old version not specified")
		}

	}
	if err := ctx.Set(MakeVersionedKey(req.Version, defaultCollectionKey), &cardCollectionList); err != nil {
		return err
	}

	// initialize default deck
	deckList.Decks = req.DefaultDecks
	if req.DefaultDecks == nil {
		if req.OldVersion != "" {
			if err := ctx.Get(MakeVersionedKey(req.OldVersion, defaultDeckKey), &deckList); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("'default_decks' key missing, old version not specified")
		}
	}
	if err := ctx.Set(MakeVersionedKey(req.Version, defaultDeckKey), &deckList); err != nil {
		return err
	}

	// initialize AI decks
	aiDeckList := zb.AIDeckList{
		Decks: req.AiDecks,
	}
	if req.AiDecks == nil {
		if req.OldVersion != "" {
			if err := ctx.Get(MakeVersionedKey(req.OldVersion, aiDecksKey), &deckList); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("'ai_decks' key missing, old version not specified")
		}
	}

	if err := ctx.Set(MakeVersionedKey(req.Version, aiDecksKey), &aiDeckList); err != nil {
		return err
	}

	return nil
}

func (z *ZombieBattleground) GetInit(ctx contract.StaticContext, req *zb.GetInitRequest) (*zb.GetInitResponse, error) {
	var cardList zb.CardList
	var heroList zb.HeroList
	var defaultHeroList zb.HeroList
	var cardCollectionList zb.CardCollectionList
	var deckList zb.DeckList
	var aiDeckList zb.AIDeckList

	if err := ctx.Get(MakeVersionedKey(req.Version, cardListKey), &cardList); err != nil {
		return nil, errors.Wrap(err, "error getting cardList")
	}

	if err := ctx.Get(MakeVersionedKey(req.Version, heroListKey), &heroList); err != nil {
		return nil, errors.Wrap(err, "error getting heroList")
	}

	if err := ctx.Get(MakeVersionedKey(req.Version, defaultHeroesKey), &defaultHeroList); err != nil {
		return nil, errors.Wrap(err, "error getting default heroList")
	}

	if err := ctx.Get(MakeVersionedKey(req.Version, defaultCollectionKey), &cardCollectionList); err != nil {
		return nil, errors.Wrap(err, "error getting default collectionList")
	}

	if err := ctx.Get(MakeVersionedKey(req.Version, defaultDeckKey), &deckList); err != nil {
		return nil, errors.Wrap(err, "error getting default deckList")
	}

	if err := ctx.Get(MakeVersionedKey(req.Version, aiDecksKey), &aiDeckList); err != nil {
		return nil, errors.Wrap(err, "error getting aiDeckList")
	}

	return &zb.GetInitResponse{
		Cards:             cardList.Cards,
		Heroes:            heroList.Heroes,
		DefaultHeroes:     defaultHeroList.Heroes,
		DefaultDecks:      deckList.Decks,
		DefaultCollection: cardCollectionList.Cards,
		AiDecks:           aiDeckList.Decks,
		Version:           req.Version,
	}, nil
}

func (z *ZombieBattleground) UpdateCardList(ctx contract.Context, req *zb.UpdateCardListRequest) error {
	cardList := zb.CardList{
		Cards: req.Cards,
	}
	return saveCardList(ctx, req.Version, &cardList)
}

func (z *ZombieBattleground) GetCardList(ctx contract.Context, req *zb.GetCardListRequest) (*zb.GetCardListResponse, error) {
	cardlist, err := loadCardList(ctx, req.Version)
	if err != nil {
		return nil, err
	}
	return &zb.GetCardListResponse{Cards: cardlist.Cards}, nil
}

func (z *ZombieBattleground) GetAccount(ctx contract.StaticContext, req *zb.GetAccountRequest) (*zb.Account, error) {
	var account zb.Account
	if err := ctx.Get(AccountKey(req.UserId), &account); err != nil {
		return nil, errors.Wrapf(err, "unable to retrieve account data for userId: %s", req.UserId)
	}
	return &account, nil
}

func (z *ZombieBattleground) UpdateAccount(ctx contract.Context, req *zb.UpsertAccountRequest) (*zb.Account, error) {
	// Verify whether this privateKey associated with user
	if !isOwner(ctx, req.UserId) {
		return nil, ErrUserNotVerified
	}

	var account zb.Account
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

func (z *ZombieBattleground) CreateAccount(ctx contract.Context, req *zb.UpsertAccountRequest) error {
	// confirm owner doesnt exist already
	if ctx.Has(AccountKey(req.UserId)) {
		return errors.New("user already exists")
	}

	var account zb.Account
	account.UserId = req.UserId
	account.Owner = ctx.Message().Sender.Bytes()
	copyAccountInfo(&account, req)

	if err := ctx.Set(AccountKey(req.UserId), &account); err != nil {
		return errors.Wrapf(err, "error setting account information for userId: %s", req.UserId)
	}
	ctx.GrantPermission([]byte(req.UserId), []string{OwnerRole})

	// add default collection list
	var collectionList zb.CardCollectionList
	if err := ctx.Get(MakeVersionedKey(req.Version, defaultCollectionKey), &collectionList); err != nil {
		return errors.Wrapf(err, "unable to get default collectionlist")
	}
	if err := ctx.Set(CardCollectionKey(req.UserId), &collectionList); err != nil {
		return errors.Wrapf(err, "unable to save card collection for userId: %s", req.UserId)
	}

	var deckList zb.DeckList
	if err := ctx.Get(MakeVersionedKey(req.Version, defaultDeckKey), &deckList); err != nil {
		return errors.Wrapf(err, "unable to get default decks")
	}
	// update default deck with none-zero id
	for i := 0; i < len(deckList.Decks); i++ {
		deckList.Decks[i].Id = int64(i + 1)
	}
	if err := ctx.Set(DecksKey(req.UserId), &deckList); err != nil {
		return errors.Wrapf(err, "unable to save decks for userId: %s", req.UserId)
	}

	var heroes zb.HeroList
	if err := ctx.Get(MakeVersionedKey(req.Version, defaultHeroesKey), &heroes); err != nil {
		return errors.Wrapf(err, "unable to get default hero")
	}
	if err := ctx.Set(HeroesKey(req.UserId), &heroes); err != nil {
		return errors.Wrapf(err, "unable to save heroes for userId: %s", req.UserId)
	}

	senderAddress := []byte(ctx.Message().Sender.Local)
	emitMsgJSON, err := prepareEmitMsgJSON(senderAddress, req.UserId, "createaccount")
	if err == nil {
		ctx.EmitTopics(emitMsgJSON, TopicCreateAccountEvent)
	}

	return nil
}

func (z *ZombieBattleground) UpdateUserElo(ctx contract.Context, req *zb.UpdateUserEloRequest) error {
	// Verify whether this privateKey associated with user
	if !isOwner(ctx, req.UserId) {
		return ErrUserNotVerified
	}

	var account zb.Account
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
	ctx.EmitTopics([]byte(data), "zombiebattleground:update_elo")
	return nil
}

// CreateDeck appends the given deck to user's deck list
func (z *ZombieBattleground) CreateDeck(ctx contract.Context, req *zb.CreateDeckRequest) (*zb.CreateDeckResponse, error) {
	if req.Deck == nil {
		return nil, ErrDeckMustNotNil
	}
	if !isOwner(ctx, req.UserId) {
		return nil, ErrUserNotVerified
	}
	// validate hero
	heroes, err := loadHeroes(ctx, req.UserId)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get hero data for userId: %s", req.UserId)
	}
	if err := validateDeckHero(heroes.Heroes, req.Deck.HeroId); err != nil {
		return nil, err
	}
	// validate version on card library
	cardLibrary, err := loadCardList(ctx, req.Version)
	if err != nil {
		return nil, err
	}
	if err := validateCardLibrary(cardLibrary.Cards, req.Deck.Cards); err != nil {
		return nil, err
	}

	// Since the server side does not have any knowleadge on user's collection, we skip this logic on the server side for now.
	// TODO: Turn on the check when the server side knows user's collection
	// validating against default card collection
	// var defaultCollection zb.CardCollectionList
	// if err := ctx.Get(MakeVersionedKey(req.Version, defaultCollectionKey), &defaultCollection); err != nil {
	// 	return nil, errors.Wrapf(err, "unable to get default collectionlist")
	// }
	// // make sure the given cards and amount must be a subset of user's cards
	// if err := validateDeckCollections(defaultCollection.Cards, req.Deck.Cards); err != nil {
	// 	return nil, err
	// }

	deckList, err := loadDecks(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	// check duplicate name
	if existing := getDeckByName(deckList.Decks, req.Deck.Name); existing != nil {
		return nil, ErrDeckNameExists
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
	if err := saveDecks(ctx, req.UserId, deckList); err != nil {
		return nil, err
	}

	senderAddress := ctx.Message().Sender.Local.String()
	emitMsg := zb.CreateDeckEvent{
		UserId:        req.UserId,
		SenderAddress: senderAddress,
		Deck:          req.Deck,
		Version:       req.Version,
	}

	data, err := proto.Marshal(&emitMsg)
	if err != nil {
		return nil, err
	}
	ctx.EmitTopics([]byte(data), TopicCreateDeckEvent)

	return &zb.CreateDeckResponse{DeckId: newDeckID}, nil
}

// EditDeck edits the deck by id
func (z *ZombieBattleground) EditDeck(ctx contract.Context, req *zb.EditDeckRequest) error {
	if req.Deck == nil {
		return fmt.Errorf("deck must not be nil")
	}
	if !isOwner(ctx, req.UserId) {
		return ErrUserNotVerified
	}
	// validate hero
	heroes, err := loadHeroes(ctx, req.UserId)
	if err := validateDeckHero(heroes.Heroes, req.Deck.HeroId); err != nil {
		return err
	}
	// validate version on card library
	cardLibrary, err := loadCardList(ctx, req.Version)
	if err != nil {
		return err
	}
	if err := validateCardLibrary(cardLibrary.Cards, req.Deck.Cards); err != nil {
		return err
	}

	// Since the server side does not have any knowleadge on user's collection, we skip this logic on the server side for now.
	// TODO: Turn on the check when the server side knows user's collection
	// validating against default card collection
	// var defaultCollection zb.CardCollectionList
	// if err := ctx.Get(MakeVersionedKey(req.Version, defaultCollectionKey), &defaultCollection); err != nil {
	// 	return nil, errors.Wrapf(err, "unable to get default collectionlist")
	// }
	// // make sure the given cards and amount must be a subset of user's cards
	// if err := validateDeckCollections(defaultCollection.Cards, req.Deck.Cards); err != nil {
	// 	return nil, err
	// }

	// validate deck
	deckList, err := loadDecks(ctx, req.UserId)
	if err != nil {
		return err
	}
	// TODO: check if this still valid
	// The deck name should be validated on the client side, not server
	if err := validateDeckName(deckList.Decks, req.Deck); err != nil {
		return err
	}

	deckID := req.Deck.Id
	existingDeck := getDeckByID(deckList.Decks, deckID)
	if existingDeck == nil {
		return ErrNotfound
	}
	// update deck
	existingDeck.Name = req.Deck.Name
	existingDeck.Cards = req.Deck.Cards
	existingDeck.HeroId = req.Deck.HeroId
	existingDeck.PrimarySkill = req.Deck.PrimarySkill
	existingDeck.SecondarySkill = req.Deck.SecondarySkill

	// update decklist
	if err := saveDecks(ctx, req.UserId, deckList); err != nil {
		return err
	}

	senderAddress := ctx.Message().Sender.Local.String()
	emitMsg := zb.EditDeckEvent{
		UserId:        req.UserId,
		SenderAddress: senderAddress,
		Deck:          req.Deck,
		Version:       req.Version,
	}

	data, err := proto.Marshal(&emitMsg)
	if err != nil {
		return err
	}
	ctx.EmitTopics([]byte(data), TopicEditDeckEvent)

	return nil
}

// DeleteDeck deletes a user's deck by id
func (z *ZombieBattleground) DeleteDeck(ctx contract.Context, req *zb.DeleteDeckRequest) error {
	if !isOwner(ctx, req.UserId) {
		return ErrUserNotVerified
	}

	deckList, err := loadDecks(ctx, req.UserId)
	if err != nil {
		return err
	}

	var deleted bool
	deckList.Decks, deleted = deleteDeckByID(deckList.Decks, req.DeckId)
	if !deleted {
		return fmt.Errorf("deck not found")
	}

	if err := saveDecks(ctx, req.UserId, deckList); err != nil {
		return err
	}

	senderAddress := ctx.Message().Sender.Local.String()
	emitMsg := zb.DeleteDeckEvent{
		UserId:        req.UserId,
		SenderAddress: senderAddress,
		DeckId:        req.DeckId,
	}

	data, err := proto.Marshal(&emitMsg)
	if err != nil {
		return err
	}
	ctx.EmitTopics([]byte(data), TopicDeleteDeckEvent)

	return nil
}

// ListDecks returns the user's decks
func (z *ZombieBattleground) ListDecks(ctx contract.StaticContext, req *zb.ListDecksRequest) (*zb.ListDecksResponse, error) {
	deckList, err := loadDecks(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &zb.ListDecksResponse{
		Decks: deckList.Decks,
	}, nil
}

// GetDeck returns the deck by given id
func (z *ZombieBattleground) GetDeck(ctx contract.StaticContext, req *zb.GetDeckRequest) (*zb.GetDeckResponse, error) {
	deckList, err := loadDecks(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	deck := getDeckByID(deckList.Decks, req.DeckId)
	if deck == nil {
		return nil, contract.ErrNotFound
	}
	return &zb.GetDeckResponse{Deck: deck}, nil
}

func (z *ZombieBattleground) SetAIDecks(ctx contract.Context, req *zb.SetAIDecksRequest) error {
	deckList := zb.AIDeckList{
		Decks: req.Decks,
	}
	return saveAIDecks(ctx, req.Version, &deckList)
}

func (z *ZombieBattleground) GetAIDecks(ctx contract.StaticContext, req *zb.GetAIDecksRequest) (*zb.GetAIDecksResponse, error) {
	deckList, err := loadAIDecks(ctx, req.Version)
	if err != nil {
		return nil, err
	}
	return &zb.GetAIDecksResponse{
		Decks: deckList.Decks,
	}, nil
}

// GetCollection returns the collection of the card own by the user
func (z *ZombieBattleground) GetCollection(ctx contract.StaticContext, req *zb.GetCollectionRequest) (*zb.GetCollectionResponse, error) {
	collectionList, err := loadCardCollection(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &zb.GetCollectionResponse{Cards: collectionList.Cards}, nil
}

// ListCardLibrary list all the card library data
func (z *ZombieBattleground) ListCardLibrary(ctx contract.StaticContext, req *zb.ListCardLibraryRequest) (*zb.ListCardLibraryResponse, error) {
	var cardList zb.CardList
	if err := ctx.Get(MakeVersionedKey(req.Version, cardListKey), &cardList); err != nil {
		return nil, err
	}

	return &zb.ListCardLibraryResponse{Cards: cardList.Cards}, nil
}

func (z *ZombieBattleground) ListHeroLibrary(ctx contract.StaticContext, req *zb.ListHeroLibraryRequest) (*zb.ListHeroLibraryResponse, error) {
	var heroList zb.HeroList
	if err := ctx.Get(MakeVersionedKey(req.Version, heroListKey), &heroList); err != nil {
		return nil, err
	}
	return &zb.ListHeroLibraryResponse{Heroes: heroList.Heroes}, nil
}

func (z *ZombieBattleground) UpdateHeroLibrary(ctx contract.Context, req *zb.UpdateHeroLibraryRequest) (*zb.UpdateHeroLibraryResponse, error) {
	var heroList = zb.HeroList{
		Heroes: req.Heroes,
	}
	if err := ctx.Set(MakeVersionedKey(req.Version, heroListKey), &heroList); err != nil {
		return nil, err
	}
	return &zb.UpdateHeroLibraryResponse{}, nil
}

func (z *ZombieBattleground) ListHeroes(ctx contract.StaticContext, req *zb.ListHeroesRequest) (*zb.ListHeroesResponse, error) {
	heroList, err := loadHeroes(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &zb.ListHeroesResponse{Heroes: heroList.Heroes}, nil
}

func (z *ZombieBattleground) GetHero(ctx contract.StaticContext, req *zb.GetHeroRequest) (*zb.GetHeroResponse, error) {
	heroList, err := loadHeroes(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	hero := getHeroById(heroList.Heroes, req.HeroId)
	if hero == nil {
		return nil, contract.ErrNotFound
	}
	return &zb.GetHeroResponse{Hero: hero}, nil
}

func (z *ZombieBattleground) SetHero(ctx contract.Context, req *zb.SetHeroRequest) (*zb.SetHeroResponse, error) {
	if req.Hero == nil {
		return nil, fmt.Errorf("Hero is null")
	}

	heroList, err := loadHeroes(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	hero := getHeroById(heroList.Heroes, req.HeroId)
	if hero == nil {
		return nil, contract.ErrNotFound
	}
	hero = proto.Clone(req.Hero).(*zb.Hero)

	// make sure we don't override hero id
	hero.HeroId = req.HeroId

	if err := saveHeroes(ctx, req.UserId, heroList); err != nil {
		return nil, err
	}

	senderAddress := []byte(ctx.Message().Sender.Local)
	emitMsgJSON, err := prepareEmitMsgJSON(senderAddress, req.UserId, "setHero")
	if err == nil {
		ctx.EmitTopics(emitMsgJSON, "zombiebattleground:sethero")
	}

	return &zb.SetHeroResponse{Hero: hero}, nil
}

func (z *ZombieBattleground) AddHeroExperience(ctx contract.Context, req *zb.AddHeroExperienceRequest) (*zb.AddHeroExperienceResponse, error) {
	if req.Experience <= 0 {
		return nil, fmt.Errorf("experience needs to be greater than zero")
	}
	if !isOwner(ctx, req.UserId) {
		return nil, ErrUserNotVerified
	}

	heroList, err := loadHeroes(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	hero := getHeroById(heroList.Heroes, req.HeroId)
	if hero == nil {
		return nil, contract.ErrNotFound
	}
	hero.Experience += req.Experience

	if err := saveHeroes(ctx, req.UserId, heroList); err != nil {
		return nil, err
	}

	senderAddress := []byte(ctx.Message().Sender.Local)
	emitMsgJSON, err := prepareEmitMsgJSON(senderAddress, req.UserId, "addHeroExperience")
	if err == nil {
		ctx.EmitTopics(emitMsgJSON, TopicAddHeroExpEvent)
	}

	return &zb.AddHeroExperienceResponse{HeroId: hero.HeroId, Experience: hero.Experience}, nil
}

func (z *ZombieBattleground) SetHeroExperience(ctx contract.Context, req *zb.SetHeroExperienceRequest) (*zb.SetHeroExperienceResponse, error) {
	if req.Experience <= 0 {
		return nil, fmt.Errorf("experience needs to be greater than zero")
	}

	heroList, err := loadHeroes(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	hero := getHeroById(heroList.Heroes, req.HeroId)
	if hero == nil {
		return nil, contract.ErrNotFound
	}
	hero.Experience = req.Experience

	if err := saveHeroes(ctx, req.UserId, heroList); err != nil {
		return nil, err
	}

	senderAddress := []byte(ctx.Message().Sender.Local)
	emitMsgJSON, err := prepareEmitMsgJSON(senderAddress, req.UserId, "setHeroExperience")
	if err == nil {
		ctx.EmitTopics(emitMsgJSON, "zombiebattleground:setheroexperience")
	}

	return &zb.SetHeroExperienceResponse{HeroId: hero.HeroId, Experience: hero.Experience}, nil
}

func (z *ZombieBattleground) SetHeroLevel(ctx contract.Context, req *zb.SetHeroLevelRequest) (*zb.SetHeroLevelResponse, error) {
	if req.Level <= 0 {
		return nil, fmt.Errorf("level needs to be greater than zero")
	}

	heroList, err := loadHeroes(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	hero := getHeroById(heroList.Heroes, req.HeroId)
	if hero == nil {
		return nil, contract.ErrNotFound
	}
	hero.Level = req.Level

	if err := saveHeroes(ctx, req.UserId, heroList); err != nil {
		return nil, err
	}

	senderAddress := []byte(ctx.Message().Sender.Local)
	emitMsgJSON, err := prepareEmitMsgJSON(senderAddress, req.UserId, "setHeroLevel")
	if err == nil {
		ctx.EmitTopics(emitMsgJSON, "zombiebattleground:setherolevel")
	}

	return &zb.SetHeroLevelResponse{HeroId: hero.HeroId, Level: hero.Level}, nil
}

func (z *ZombieBattleground) GetHeroSkills(ctx contract.StaticContext, req *zb.GetHeroSkillsRequest) (*zb.GetHeroSkillsResponse, error) {
	heroList, err := loadHeroes(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	hero := getHeroById(heroList.Heroes, req.HeroId)
	if hero == nil {
		return nil, contract.ErrNotFound
	}
	return &zb.GetHeroSkillsResponse{HeroId: hero.HeroId, Skills: hero.Skills}, nil
}

func (z *ZombieBattleground) RegisterPlayerPool(ctx contract.Context, req *zb.RegisterPlayerPoolRequest) (*zb.RegisterPlayerPoolResponse, error) {
	// preparing user profile consisting of deck, score, ...
	_, err := getDeckWithRegistrationData(ctx, req.RegistrationData)
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

	profile := zb.PlayerProfile{
		RegistrationData: req.RegistrationData,
		UpdatedAt:        ctx.Now().Unix(),
	}

	fmt.Printf("RegisterPlayerPool: %+v\n", profile)

	var loadPlayerPoolFn func(contract.StaticContext) (*zb.PlayerPool, error)
	var savePlayerPoolFn func(contract.Context, *zb.PlayerPool) error
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
				match.Status = zb.Match_Timedout
				// remove player's match if existing
				ctx.Delete(UserMatchKey(pp.RegistrationData.UserId))
				// notify player
				emitMsg := zb.PlayerActionEvent{
					Match: match,
				}
				data, err := proto.Marshal(&emitMsg)
				if err != nil {
					return nil, err
				}
				ctx.EmitTopics([]byte(data), match.Topics...)
			}
		}
	}

	senderAddress := []byte(ctx.Message().Sender.Local)
	emitMsgJSON, err := prepareEmitMsgJSON(senderAddress, req.RegistrationData.UserId, "registerplayerpool")
	if err == nil {
		ctx.EmitTopics(emitMsgJSON, TopicRegisterPlayerPoolEvent)
	}

	return &zb.RegisterPlayerPoolResponse{}, nil
}

func (z *ZombieBattleground) FindMatch(ctx contract.Context, req *zb.FindMatchRequest) (*zb.FindMatchResponse, error) {
	var loadPlayerPoolFn func(contract.StaticContext) (*zb.PlayerPool, error)
	var savePlayerPoolFn func(contract.Context, *zb.PlayerPool) error
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
		if match.Status == zb.Match_Matching {
			updatedAt := time.Unix(match.CreatedAt, 0)
			if updatedAt.Add(MMTimeout).Before(ctx.Now()) {
				ctx.Logger().Debug(fmt.Sprintf("Match %d timedout", match.Id))
				// remove match
				// ctx.Delete(MatchKey(match.Id))
				match.Status = zb.Match_Timedout
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
		emitMsg := zb.PlayerActionEvent{
			Match: match,
		}
		data, err := proto.Marshal(&emitMsg)
		if err != nil {
			return nil, err
		}
		if err == nil {
			topics := append(match.Topics, TopicFindMatchEvent)
			ctx.EmitTopics([]byte(data), topics...)
		}

		return &zb.FindMatchResponse{
			Match:      match,
			MatchFound: true,
		}, nil
	}
	playerProfile := findPlayerProfileByID(pool, req.UserId)
	if playerProfile == nil {
		return nil, errors.New("Player not found in player pool")
	}

	deck, err := getDeckWithRegistrationData(ctx, playerProfile.RegistrationData)
	if err != nil {
		return nil, err
	}

	// perform matchmaking function to calculate scores
	// steps:
	// 1. list all the candidates that has similar profiles TODO: match versions
	// 2. pick the most highest score
	// 3. if there is no candidate, sleep for MMWaitTime seconds
	retries := 0
	var matchedPlayerProfile *zb.PlayerProfile
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
		return &zb.FindMatchResponse{
			MatchFound: false,
		}, nil
	}

	// get matched player deck
	matchedDeck, err := getDeckWithRegistrationData(ctx, matchedPlayerProfile.RegistrationData)
	if err != nil {
		return nil, err
	}

	// create match
	match = &zb.Match{
		Status: zb.Match_Matching,
		PlayerStates: []*zb.InitialPlayerState{
			&zb.InitialPlayerState{
				Id:            playerProfile.RegistrationData.UserId,
				Deck:          deck,
				MatchAccepted: false,
			},
			&zb.InitialPlayerState{
				Id:            matchedPlayerProfile.RegistrationData.UserId,
				Deck:          matchedDeck,
				MatchAccepted: false,
			},
		},
		Version: playerProfile.RegistrationData.Version, // TODO: match version of both players
		PlayerLastSeens: []*zb.PlayerTimestamp{
			&zb.PlayerTimestamp{
				Id:        playerProfile.RegistrationData.UserId,
				UpdatedAt: ctx.Now().Unix(),
			},
			&zb.PlayerTimestamp{
				Id:        matchedPlayerProfile.RegistrationData.UserId,
				UpdatedAt: ctx.Now().Unix(),
			},
		},
		PlayerDebugCheats: []*zb.DebugCheatsConfiguration{
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

	emitMsg := zb.PlayerActionEvent{
		Match: match,
	}
	data, err := proto.Marshal(&emitMsg)
	if err != nil {
		return nil, err
	}
	topics := append(match.Topics, TopicFindMatchEvent)
	ctx.EmitTopics([]byte(data), topics...)

	return &zb.FindMatchResponse{
		Match:      match,
		MatchFound: true,
	}, nil
}

func (z *ZombieBattleground) AcceptMatch(ctx contract.Context, req *zb.AcceptMatchRequest) (*zb.AcceptMatchResponse, error) {
	match, err := loadUserCurrentMatch(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	if req.MatchId != match.Id {
		return nil, errors.New("match id not correct")
	}

	if match.Status != zb.Match_Matching {
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

	emitMsg := zb.PlayerActionEvent{
		Match: match,
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

		playerStates := []*zb.PlayerState{
			&zb.PlayerState{
				Id:   match.PlayerStates[0].Id,
				Deck: match.PlayerStates[0].Deck,
			},
			&zb.PlayerState{
				Id:   match.PlayerStates[1].Id,
				Deck: match.PlayerStates[1].Deck,
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

		match.Status = zb.Match_Started

		emitMsg = zb.PlayerActionEvent{
			Match: match,
			Block: &zb.History{List: gp.history},
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
	ctx.EmitTopics([]byte(data), topics...)

	return &zb.AcceptMatchResponse{
		Match: match,
	}, nil
}

// TODO remove this
func (z *ZombieBattleground) GetPlayerPool(ctx contract.StaticContext, req *zb.PlayerPoolRequest) (*zb.PlayerPoolResponse, error) {
	pool, err := loadPlayerPool(ctx)
	if err != nil {
		return nil, err
	}

	return &zb.PlayerPoolResponse{
		Pool: pool,
	}, nil
}

func (z *ZombieBattleground) CancelFindMatch(ctx contract.Context, req *zb.CancelFindMatchRequest) (*zb.CancelFindMatchResponse, error) {
	match, _ := loadUserCurrentMatch(ctx, req.UserId)

	if match != nil {
		// remove current match
		for _, player := range match.PlayerStates {
			ctx.Delete(UserMatchKey(player.Id))
		}
		match.Status = zb.Match_Canceled
		if err := saveMatch(ctx, match); err != nil {
			return nil, err
		}
		// notify player
		emitMsg := zb.PlayerActionEvent{
			Match: match,
		}
		data, err := proto.Marshal(&emitMsg)
		if err != nil {
			return nil, err
		}
		if err == nil {
			ctx.EmitTopics([]byte(data), match.Topics...)
		}
	}

	var loadPlayerPoolFn func(contract.StaticContext) (*zb.PlayerPool, error)
	var savePlayerPoolFn func(contract.Context, *zb.PlayerPool) error
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

	return &zb.CancelFindMatchResponse{}, nil
}

func (z *ZombieBattleground) GetMatch(ctx contract.StaticContext, req *zb.GetMatchRequest) (*zb.GetMatchResponse, error) {
	match, err := loadMatch(ctx, req.MatchId)
	if err != nil {
		return nil, err
	}

	return &zb.GetMatchResponse{
		Match: match,
	}, nil
}

func (z *ZombieBattleground) GetGameState(ctx contract.StaticContext, req *zb.GetGameStateRequest) (*zb.GetGameStateResponse, error) {
	gameState, err := loadGameState(ctx, req.MatchId)
	if err != nil {
		return nil, err
	}

	return &zb.GetGameStateResponse{
		GameState: gameState,
	}, nil
}

func (z *ZombieBattleground) EndMatch(ctx contract.Context, req *zb.EndMatchRequest) (*zb.EndMatchResponse, error) {
	match, err := loadMatch(ctx, req.MatchId)
	if err != nil {
		return nil, err
	}

	match.Status = zb.Match_Ended
	if err := saveMatch(ctx, match); err != nil {
		return nil, err
	}
	// delete user match
	ctx.Delete(UserMatchKey(req.UserId))

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
	gp.history = append(gp.history, &zb.HistoryData{
		Data: &zb.HistoryData_EndGame{
			EndGame: &zb.HistoryEndGame{
				UserId:   req.GetUserId(),
				MatchId:  req.MatchId,
				WinnerId: req.WinnerId,
			},
		},
	})
	// Don't think we need this since endgame should be emitted to match
	// match.Topics = append(match.Topics, "endgame")
	emitMsg := zb.PlayerActionEvent{
		Match: match,
		Block: &zb.History{List: gp.history},
	}
	data, err := proto.Marshal(&emitMsg)
	if err != nil {
		return nil, err
	}
	ctx.EmitTopics([]byte(data), match.Topics...)

	return &zb.EndMatchResponse{GameState: gamestate}, nil
}

func (z *ZombieBattleground) CheckGameStatus(ctx contract.Context, req *zb.CheckGameStatusRequest) (*zb.CheckGameStatusResponse, error) {
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
	// Check if the current player is gone for more than timeout.
	// If there is no action added, we check the gamestate created time.
	var createdAt time.Time
	latestAction := gp.current()
	if latestAction == nil {
		createdAt = time.Unix(gamestate.CreatedAt, 0)
	} else {
		createdAt = time.Unix(latestAction.CreatedAt, 0)
	}
	activePlayer := gp.activePlayer()
	if createdAt.Add(TurnTimeout).Before(ctx.Now()) {
		// create a leave match request and append to the game state
		leaveMatchAction := zb.PlayerAction{
			ActionType: zb.PlayerActionType_LeaveMatch,
			PlayerId:   activePlayer.Id,
			Action: &zb.PlayerAction_LeaveMatch{
				LeaveMatch: &zb.PlayerActionLeaveMatch{},
			},
			CreatedAt: ctx.Now().Unix(),
		}
		err := gp.AddAction(&leaveMatchAction)
		// ignore the error in case this method is called mutiple times
		if err == nil {
			if err := saveGameState(ctx, gamestate); err != nil {
				return nil, err
			}
		}
		// update match status
		match.Status = zb.Match_PlayerLeft
		if err := saveMatch(ctx, match); err != nil {
			return nil, err
		}
		// update winner
		leaveMatchReq := leaveMatchAction.GetLeaveMatch()
		leaveMatchReq.Winner = gp.State.Winner
		emitMsg := zb.PlayerActionEvent{
			PlayerAction: &leaveMatchAction,
			Match:        match,
			Block:        &zb.History{List: gp.history},
		}
		data, err := proto.Marshal(&emitMsg)
		if err != nil {
			return nil, err
		}
		ctx.EmitTopics([]byte(data), match.Topics...)
	}

	return &zb.CheckGameStatusResponse{}, nil
}

func (z *ZombieBattleground) SendPlayerAction(ctx contract.Context, req *zb.PlayerActionRequest) (*zb.PlayerActionResponse, error) {
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
	if match.Status == zb.Match_Started {
		match.Status = zb.Match_Playing
		if err := saveMatch(ctx, match); err != nil {
			return nil, err
		}
	}

	emitMsg := zb.PlayerActionEvent{
		PlayerAction: req.PlayerAction,
		Match:        match,
		Block:        &zb.History{List: gp.history},
	}

	data, err := proto.Marshal(&emitMsg)
	if err != nil {
		return nil, err
	}
	ctx.EmitTopics([]byte(data), match.Topics...)

	return &zb.PlayerActionResponse{
		Match: match,
	}, nil
}

func (z *ZombieBattleground) SendBundlePlayerAction(ctx contract.Context, req *zb.BundlePlayerActionRequest) (*zb.BundlePlayerActionResponse, error) {
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
	if match.Status == zb.Match_Started {
		match.Status = zb.Match_Playing
		if err := saveMatch(ctx, match); err != nil {
			return nil, err
		}
	}

	return &zb.BundlePlayerActionResponse{
		GameState: gamestate,
		Match:     match,
		History:   gp.history,
	}, nil
}

func (z *ZombieBattleground) KeepAlive(ctx contract.Context, req *zb.KeepAliveRequest) (*zb.KeepAliveResponse, error) {
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
		return &zb.KeepAliveResponse{}, nil
	}

	gamestate, err := loadGameState(ctx, match.Id)
	if err != nil {
		return nil, err
	}
	if gamestate.IsEnded {
		return &zb.KeepAliveResponse{}, nil // just ignore for client check
	}

	gp, err := GamePlayFrom(gamestate, match.UseBackendGameLogic, match.PlayerDebugCheats)
	if err != nil {
		return nil, err
	}
	for _, lastseen := range match.PlayerLastSeens {
		lastSeenAt := time.Unix(lastseen.UpdatedAt, 0)
		if lastSeenAt.Add(KeepAliveTimeout).Before(ctx.Now()) {
			// create a leave match request and append to the game state
			leaveMatchAction := zb.PlayerAction{
				ActionType: zb.PlayerActionType_LeaveMatch,
				PlayerId:   lastseen.Id,
				Action: &zb.PlayerAction_LeaveMatch{
					LeaveMatch: &zb.PlayerActionLeaveMatch{},
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
			match.Status = zb.Match_PlayerLeft
			if err := saveMatch(ctx, match); err != nil {
				return nil, err
			}
			// update winner
			leaveMatchReq := leaveMatchAction.GetLeaveMatch()
			leaveMatchReq.Winner = gp.State.Winner
			emitMsg := zb.PlayerActionEvent{
				PlayerAction: &leaveMatchAction,
				Match:        match,
				Block:        &zb.History{List: gp.history},
			}
			data, err := proto.Marshal(&emitMsg)
			if err != nil {
				return nil, err
			}
			ctx.EmitTopics([]byte(data), match.Topics...)
		}
	}

	return &zb.KeepAliveResponse{}, nil
}

func (z *ZombieBattleground) GetState(ctx contract.StaticContext, req *zb.GetGamechainStateRequest) (*zb.GetGamechainStateResponse, error) {
	state, err := loadState(ctx)
	if err != nil {
		return nil, err
	}
	return &zb.GetGamechainStateResponse{
		State: state,
	}, nil
}

func (z *ZombieBattleground) InitState(ctx contract.Context, req *zb.InitGamechainStateRequest) error {
	state, err := loadState(ctx)
	if err != nil && err != contract.ErrNotFound {
		return err
	}
	if state != nil {
		return fmt.Errorf("state already inilialized")
	}
	if req.Oracle == nil {
		return ErrOracleNotSpecified
	}

	if err := z.validateOracle(ctx, req.Oracle); err != nil {
		return err
	}
	state = &zb.GamechainState{
		LastPlasmachainBlockNum: 1,
	}
	return saveState(ctx, state)
}

func (z *ZombieBattleground) UpdateVersions(ctx contract.Context, req *zb.UpdateVersionsRequest) error {
	var err error
	if req.ContentVersion != "" {
		err = ctx.Set(contentVersionKey, &zb.ContentVersion{
			ContentVersion: req.ContentVersion,
		})
		if err != nil {
			return err
		}
	}

	if req.PvpVersion != "" {
		err = ctx.Set(pvpVersionKey, &zb.PvpVersion{
			PvpVersion: req.PvpVersion,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (z *ZombieBattleground) GetVersions(ctx contract.StaticContext, req *zb.GetVersionsRequest) (*zb.GetVersionsResponse, error) {
	var contentVersion zb.ContentVersion
	var pvpVersion zb.PvpVersion
	var err error
	err = ctx.Get(contentVersionKey, &contentVersion)
	if err != nil {
		return nil, err
	}

	err = ctx.Get(pvpVersionKey, &pvpVersion)
	if err != nil {
		return nil, err
	}

	return &zb.GetVersionsResponse{
		ContentVersion: contentVersion.ContentVersion,
		PvpVersion:     pvpVersion.PvpVersion,
	}, nil
}

func (z *ZombieBattleground) UpdateOracle(ctx contract.Context, params *zb.UpdateOracle) error {
	if ctx.Has(oracleKey) {
		if params.OldOracle.String() == params.NewOracle.String() {
			return errors.New("Cannot set new oracle to same address as old oracle")
		}
		if err := z.validateOracle(ctx, params.OldOracle); err != nil {
			return errors.Wrap(err, "validating oracle")
		}
		ctx.GrantPermission([]byte(params.OldOracle.String()), []string{"old-oracle"})
	}
	ctx.GrantPermission([]byte(params.NewOracle.String()), []string{OracleRole})

	if err := ctx.Set(oracleKey, params.NewOracle); err != nil {
		return errors.Wrap(err, "setting new oracle")
	}
	return nil
}

func (z *ZombieBattleground) validateOracle(ctx contract.Context, zo *types.Address) error {
	if ok, _ := ctx.HasPermission([]byte(zo.String()), []string{OracleRole}); !ok {
		return errors.New("Oracle unverified")
	}

	if ok, _ := ctx.HasPermission([]byte(zo.String()), []string{"old-oracle"}); ok {
		return errors.New("This oracle is expired. Please use latest oracle")
	}
	return nil
}

func (z *ZombieBattleground) GetGameMode(ctx contract.StaticContext, req *zb.GetGameModeRequest) (*zb.GameMode, error) {
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

func (z *ZombieBattleground) CallCustomGameModeFunction(ctx contract.Context, req *zb.CallCustomGameModeFunctionRequest) error {
	err := NewCustomGameMode(loom.Address{
		ChainID: req.Address.ChainId,
		Local:   req.Address.Local,
	}).CallFunction(ctx, req.CallData)

	if err != nil {
		return err
	}

	return nil
}

func (z *ZombieBattleground) GetGameModeCustomUi(ctx contract.StaticContext, req *zb.GetCustomGameModeCustomUiRequest) (*zb.GetCustomGameModeCustomUiResponse, error) {
	uiElements, err := NewCustomGameMode(loom.Address{
		ChainID: req.Address.ChainId,
		Local:   req.Address.Local,
	}).GetCustomUi(ctx)

	if err != nil {
		return nil, err
	}

	response := &zb.GetCustomGameModeCustomUiResponse{
		UiElements: uiElements,
	}

	return response, nil
}

func (z *ZombieBattleground) ListGameModes(ctx contract.StaticContext, req *zb.ListGameModesRequest) (*zb.GameModeList, error) {
	gameModeList, err := loadGameModeList(ctx)
	if err != nil {
		return nil, err
	}

	return gameModeList, nil
}

func validateGameModeReq(req *zb.GameModeRequest) error {
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

func (z *ZombieBattleground) AddGameMode(ctx contract.Context, req *zb.GameModeRequest) (*zb.GameMode, error) {
	if err := validateGameModeReq(req); err != nil {
		return nil, err
	}

	// check if game mode with this name already exists
	gameModeList, err := loadGameModeList(ctx)
	if err != nil && err == contract.ErrNotFound {
		gameModeList = &zb.GameModeList{GameModes: []*zb.GameMode{}}
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

	gameModeType := zb.GameModeType_Community
	owner := &types.Address{ChainId: ctx.ContractAddress().ChainID, Local: ctx.Message().Sender.Local}
	// if request was made with a valid oracle, set type and owner to Loom
	if req.Oracle != "" {
		oracleLocal, err := loom.LocalAddressFromHexString(req.Oracle)
		if err != nil {
			return nil, err
		}

		oracleAddr := &types.Address{ChainId: ctx.ContractAddress().ChainID, Local: oracleLocal}

		if err := z.validateOracle(ctx, oracleAddr); err != nil {
			return nil, err
		}

		gameModeType = zb.GameModeType_Loom
		owner = loom.RootAddress(ctx.ContractAddress().ChainID).MarshalPB()
	}

	gameMode := &zb.GameMode{
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

func (z *ZombieBattleground) UpdateGameMode(ctx contract.Context, req *zb.UpdateGameModeRequest) (*zb.GameMode, error) {
	// Require either oracle or owner permission to update a game mode
	if req.Oracle != "" {
		oracleLocal, err := loom.LocalAddressFromHexString(req.Oracle)
		if err != nil {
			return nil, err
		}

		oracleAddr := &types.Address{ChainId: ctx.ContractAddress().ChainID, Local: oracleLocal}

		if err := z.validateOracle(ctx, oracleAddr); err != nil {
			return nil, err
		}
	} else if ok, _ := ctx.HasPermission([]byte(req.ID), []string{OwnerRole}); !ok {
		return nil, ErrUserNotVerified
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

func (z *ZombieBattleground) DeleteGameMode(ctx contract.Context, req *zb.DeleteGameModeRequest) error {
	// Require either oracle or owner permission to delete a game mode
	if req.Oracle != "" {
		oracleLocal, err := loom.LocalAddressFromHexString(req.Oracle)
		if err != nil {
			return err
		}

		oracleAddr := &types.Address{ChainId: ctx.ContractAddress().ChainID, Local: oracleLocal}

		if err := z.validateOracle(ctx, oracleAddr); err != nil {
			return err
		}
	} else if ok, _ := ctx.HasPermission([]byte(req.ID), []string{OwnerRole}); !ok {
		return ErrUserNotVerified
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

func (z *ZombieBattleground) RewardTutorialCompleted(ctx contract.Context, req *zb.RewardTutorialCompletedRequest) (*zb.RewardTutorialCompletedResponse, error) {
	if !isOwner(ctx, req.UserId) {
		return nil, ErrUserNotVerified
	}

	rewardClaimed, err := getRewardClaimed(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	if rewardClaimed != nil {
		if rewardClaimed.RewardType == RewardTypeTutorialCompleted {
			return nil, fmt.Errorf("reward already claimed")
		}
	}

	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		return nil, fmt.Errorf("error reading private key")
	}

	/*
		nonce, err := getNonce(ctx)
		if err != nil {
			return nil, err
		}
	*/

	//rewardType := RewardTypeTutorialCompleted

	// assign rewards
	var minionPack uint64
	minionPack = 5

	userIDUint, err := getOrGenerateUserIDUint(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	verifySignResult, err := generateVerifyHash(userIDUint, minionPack, TutorialRewardContractVersion, privateKey)
	if err != nil {
		return nil, err
	}

	if len(verifySignResult.Signature) != 132 {
		return nil, fmt.Errorf("signature length invalid")
	}

	r := verifySignResult.Signature[0:66]
	s := "0x" + verifySignResult.Signature[66:130]
	vStr := verifySignResult.Signature[130:132]
	v, err := strconv.ParseUint(vStr, 16, 8)
	if err != nil {
		return nil, err
	}

	return &zb.RewardTutorialCompletedResponse{
		UserId:     req.UserId,
		UserIdUint: userIDUint,
		RewardType: RewardTypeTutorialCompleted,
		//	Nonce:      nonce,
		Hash:       verifySignResult.Hash,
		R:          r,
		S:          s,
		V:          v,
		MinionPack: minionPack,
	}, nil
}

func (z *ZombieBattleground) ConfirmRewardClaimed(ctx contract.Context, req *zb.ConfirmRewardClaimedRequest) error {
	if !isOwner(ctx, req.UserId) {
		return ErrUserNotVerified
	}

	userIDUint, err := getUserIDUint(ctx, req.UserId)
	if err != nil {
		return err
	}

	// check whether the userIDUint in the request is same as the one generated by the RewardTutorialCompleted call
	// the uint user id is generated only in the RewardTutorialCompleted call, so it can be used for safety check
	// this is to prevent confirmation of reward claim without actually generating a reward via RewardTutorialCompleted
	if req.UserIdUint != userIDUint {
		return fmt.Errorf("invalid request")
	}
	// TODO: add rewardType checks?
	err = setRewardClaimed(ctx, req.UserId, req.RewardType)
	return err
}

type verifySignResult struct {
	Hash      string
	Signature string
}

func generateVerifyHash(userID uint64, minionPack uint64, tutorialRewardContractVersion uint64, privKey *ecdsa.PrivateKey) (*verifySignResult, error) {

	hash, err := createHash(userID, minionPack, tutorialRewardContractVersion)

	if err != nil {
		return nil, err
	}

	sig, err := soliditySign(hash, privKey)

	if err != nil {
		return nil, err
	}

	return &verifySignResult{
		Hash:      "0x" + hex.EncodeToString(hash),
		Signature: "0x" + hex.EncodeToString(sig),
	}, nil
}

func createHash(userID uint64, minionPack uint64, tutorialRewardContractVersion uint64) ([]byte, error) {

	hash := solsha3.SoliditySHA3(
		solsha3.Uint256(strconv.FormatUint(userID, 10)),
		solsha3.Uint256(strconv.FormatUint(minionPack, 10)),
		solsha3.Uint256(strconv.FormatUint(tutorialRewardContractVersion, 10)),
	)

	if len(hash) == 0 {
		return nil, errors.New("Failed to generate hash")
	}

	return hash, nil
}

func soliditySign(data []byte, privKey *ecdsa.PrivateKey) ([]byte, error) {
	sig, err := crypto.Sign(data, privKey)
	if err != nil {
		return nil, err
	}

	v := sig[len(sig)-1]
	sig[len(sig)-1] = v + 27
	return sig, nil
}

func (z *ZombieBattleground) ProcessEventBatch(ctx contract.Context, req *orctype.ProcessEventBatchRequest) error {
	state, err := loadState(ctx)
	if err != nil {
		return err
	}

	blockCount := 0           // number of blocks that were actually processed in this batch
	lastEthBlock := uint64(0) // the last block processed in this batch

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
		if ev.EthBlock <= state.LastPlasmachainBlockNum {
			continue
		}

		switch payload := ev.Payload.(type) {
		case *orctype.PlasmachainEvent_Card:
			if err := validateGeneratedCard(payload.Card); err != nil {
				return err
			}
			userID := string(payload.Card.Owner.Local) // should be bytes that represents address
			cardID := payload.Card.CardID.Value.Int64()
			amount := payload.Card.Amount.Value.Int64()
			err := z.syncCardToCollection(ctx, userID, cardID, amount, req.CardVersion)
			if err != nil {
				ctx.Logger().Error("Oracle failed to add card to user collection", "err", err)
				return err
			}
			ctx.Logger().Info(fmt.Sprintf("add card: %+v", payload.Card))
		case nil:
			ctx.Logger().Error("Oracle missing event payload")
			continue

		default:
			ctx.Logger().Error("Oracle unknown event payload type %T", payload)
			continue
		}

		if ev.EthBlock > lastEthBlock {
			blockCount++
			lastEthBlock = ev.EthBlock
		}
	}

	// If there are no new events in this batch return an error so that the batch tx isn't
	// propagated to the other nodes.
	if blockCount == 0 {
		return fmt.Errorf("no new events found in the batch")
	}

	state.LastPlasmachainBlockNum = lastEthBlock

	return saveState(ctx, state)
}

func validateGeneratedCard(card *orctype.PlasmachainGeneratedCard) error {
	if card == nil {
		return errors.New("card is nil")
	}
	if card.CardID == nil {
		return errors.New("card id is nil")
	}
	if card.Amount == nil {
		return errors.New("card amount is nil")
	}
	if card.Owner == nil {
		return errors.New("card owner is nil")
	}
	if card.Contract == nil {
		return errors.New("card contract is nil")
	}
	return nil
}

func (z *ZombieBattleground) syncCardToCollection(ctx contract.Context, userID string, cardID int64, amount int64, version string) error {
	cardCollection, err := loadCardCollection(ctx, userID)
	// the user key might not be created at the time we fetch data from plasmachain
	// we need to create a key to make sure gamechain does not return error
	if err != nil && err == contract.ErrNotFound {
		req := zb.UpsertAccountRequest{
			UserId:  userID,
			Version: version,
		}
		if err := z.CreateAccount(ctx, &req); err != nil {
			return err
		}
	}
	cardlib, err := loadCardList(ctx, version)
	if err != nil {
		return err
	}

	// Map from cardID to mouldID
	// the formular is cardID = mouldID + x
	// for example cardID 250 = 25 + 0
	//   or 161 = 16 + 1
	mouldID := cardID / 10
	card, ok := findCardByMouldID(cardlib, mouldID)
	if !ok {
		return fmt.Errorf("card mould id %d not found in card library", mouldID)
	}

	// add to collection
	found := false
	for i := range cardCollection.Cards {
		if cardCollection.Cards[i].CardName == card.Name {
			cardCollection.Cards[i].Amount += amount
			found = true
			break
		}
	}
	if !found {
		cardCollection.Cards = append(cardCollection.Cards, &zb.CardCollectionCard{
			CardName: card.Name,
			Amount:   amount,
		})
	}

	return saveCardCollection(ctx, userID, cardCollection)
}

func (z *ZombieBattleground) SetLastPlasmaBlockNum(ctx contract.Context, req *zb.SetLastPlasmaBlockNumRequest) error {
	state, err := loadState(ctx)
	if err != nil && err != contract.ErrNotFound {
		return err
	}
	// TODO: Need to validate oracle
	// if req.Oracle == nil {
	// 	return ErrOracleNotSpecified
	// }

	// if err := z.validateOracle(ctx, req.Oracle); err != nil {
	// 	return err
	// }
	state = &zb.GamechainState{
		LastPlasmachainBlockNum: req.LastBlockNum,
	}
	return saveState(ctx, state)
}

var Contract plugin.Contract = contract.MakePluginContract(&ZombieBattleground{})

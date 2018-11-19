package battleground

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"sort"
	"time"
	"unicode/utf8"

	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/plugin"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	"github.com/loomnetwork/go-loom/types"
	"github.com/pkg/errors"
)

type ZombieBattleground struct {
	ClientSideRuleOverride bool //disables all checks to ensure the client can work before server is fully implemented
}

const (
	MaxGameModeNameChar        = 48
	MaxGameModeDescriptionChar = 255
	MaxGameModeVersionChar     = 16
	TurnTimeout                = 120 * time.Second
)

var secret string

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
	disableClientSideOverride := os.Getenv("DISABLE_CLIENT_SIDE_OVERRIDE")
	if disableClientSideOverride == "false" {
		z.ClientSideRuleOverride = false
	} else {
		z.ClientSideRuleOverride = true
	}

	if req.Oracle != nil {
		ctx.GrantPermissionTo(loom.UnmarshalAddressPB(req.Oracle), []byte(req.Oracle.String()), "oracle")
		if err := ctx.Set(oracleKey, req.Oracle); err != nil {
			return errors.Wrap(err, "Error setting oracle")
		}
	}

	// initialize card library
	cardList := zb.CardList{
		Cards: req.Cards,
	}
	if err := ctx.Set(MakeVersionedKey(req.Version, cardListKey), &cardList); err != nil {
		return err
	}
	// initialize heros
	heroList := zb.HeroList{
		Heroes: req.Heroes,
	}
	if err := ctx.Set(MakeVersionedKey(req.Version, heroListKey), &heroList); err != nil {
		return err
	}

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

	defaultHeroList := zb.HeroList{
		Heroes: req.Heroes,
	}
	if err := ctx.Set(MakeVersionedKey(req.Version, defaultHeroesKey), &defaultHeroList); err != nil {
		return err
	}

	// initialize default AI decks
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

	// initialize heros
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

	// initialize default AI decks
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

	return &zb.GetInitResponse{
		Cards:             cardList.Cards,
		Heroes:            heroList.Heroes,
		DefaultHeroes:     defaultHeroList.Heroes,
		DefaultDecks:      deckList.Decks,
		DefaultCollection: cardCollectionList.Cards,
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
		ctx.EmitTopics(emitMsgJSON, "zombiebattleground:updateaccount")
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
		ctx.EmitTopics(emitMsgJSON, "zombiebattleground:createaccount")
	}

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

	senderAddress := []byte(ctx.Message().Sender.Local)
	emitMsgJSON, err := prepareEmitMsgJSON(senderAddress, req.UserId, "createdeck")
	if err == nil {
		ctx.EmitTopics(emitMsgJSON, "zombiebattleground:createdeck")
	}
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

	// update decklist
	if err := saveDecks(ctx, req.UserId, deckList); err != nil {
		return err
	}

	senderAddress := []byte(ctx.Message().Sender.Local)
	emitMsgJSON, err := prepareEmitMsgJSON(senderAddress, req.UserId, "editdeck")
	if err == nil {
		ctx.EmitTopics(emitMsgJSON, "zombiebattleground:editdeck")
	}
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
		ctx.EmitTopics(emitMsgJSON, "zombiebattleground:addheroexperience")
	}

	return &zb.AddHeroExperienceResponse{HeroId: hero.HeroId, Experience: hero.Experience}, nil
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
	dl, err := loadDecks(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	deck := getDeckByID(dl.Decks, req.DeckId)
	if deck == nil {
		return nil, fmt.Errorf("deck id %d not found", req.DeckId)
	}

	if req.Version == "" {
		return nil, fmt.Errorf("version not specified")
	}

	// sort tags
	if len(req.Tags) > 0 {
		sort.Strings(req.Tags)
	}

	profile := zb.PlayerProfile{
		UserId:     req.UserId,
		DeckId:     deck.Id,
		UpdatedAt:  ctx.Now().Unix(),
		Version:    req.Version,
		RandomSeed: req.RandomSeed,
		CustomGame: req.CustomGame,
		Tags:       req.Tags,
	}

	var loadPlayerPoolFn func(contract.StaticContext) (*zb.PlayerPool, error)
	var savePlayerPoolFn func(contract.Context, *zb.PlayerPool) error
	// if the tags is set, use tagged playerpool
	if len(profile.Tags) > 0 {
		loadPlayerPoolFn = loadTaggedPlayerPool
		savePlayerPoolFn = saveTaggedPlayerPool
	} else {
		loadPlayerPoolFn = loadPlayerPool
		savePlayerPoolFn = savePlayerPool
	}

	ctx.Logger().Info(fmt.Sprintf("Register Player Pool user=%s", req.UserId))
	// load player pool
	pool, err := loadPlayerPoolFn(ctx)
	if err != nil {
		return nil, err
	}

	match, _ := loadUserCurrentMatch(ctx, req.UserId)
	if match != nil {
		return nil, errors.New("Player is already in a match")
	}

	// if player is not in pool, add them
	if targetProfile := findPlayerProfileByID(pool, req.UserId); targetProfile == nil {
		pool.PlayerProfiles = append(pool.PlayerProfiles, &profile)
		if err := savePlayerPoolFn(ctx, pool); err != nil {
			return nil, err
		}
	}

	// prune the timed out player profile
	for _, pp := range pool.PlayerProfiles {
		updatedAt := time.Unix(pp.UpdatedAt, 0)
		if updatedAt.Add(MMTimeout).Before(ctx.Now()) {
			ctx.Logger().Debug(fmt.Sprintf("Player profile %s timedout", pp.UserId))
			// remove player from the pool
			pool = removePlayerFromPool(pool, pp.UserId)
			// remove match
			match, _ := loadUserCurrentMatch(ctx, pp.UserId)
			if match != nil {
				ctx.Delete(MatchKey(match.Id))
				match.Status = zb.Match_Timedout
				// remove player's match if existing
				ctx.Delete(UserMatchKey(pp.UserId))
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

	ctx.Logger().Info(fmt.Sprintf("Player Pool %+v", pool))

	senderAddress := []byte(ctx.Message().Sender.Local)
	emitMsgJSON, err := prepareEmitMsgJSON(senderAddress, req.UserId, "registerplayerpool")
	if err == nil {
		ctx.EmitTopics(emitMsgJSON, "zombiebattleground:registerplayerpool")
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
			ctx.EmitTopics([]byte(data), match.Topics...)
		}

		return &zb.FindMatchResponse{
			Match: match,
		}, nil
	}
	playerProfile := findPlayerProfileByID(pool, req.UserId)
	if playerProfile == nil {
		return nil, errors.New("Player not found in player pool")
	}

	dl, err := loadDecks(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	deck := getDeckByID(dl.Decks, playerProfile.DeckId)
	if deck == nil {
		return nil, fmt.Errorf("deck id %d not found while matchmaking", playerProfile.DeckId)
	}

	// perform matchmaking function to calculate scores
	// steps:
	// 1. list all the candidates that has similar profiles TODO: match versions
	// 2. pick the most highest score
	// 3. if there is no candidate, sleep for MMWaitTime seconds
	retries := 0
	var matchedPlayerProfile *zb.PlayerProfile
	for retries < MMRetries {
		ctx.Logger().Debug(fmt.Sprintf("Matchmaking for user=%s retries=%d", req.UserId, retries))

		ctx.Logger().Debug(fmt.Sprintf("PlayerPool size before running matchmaking: %d", len(pool.PlayerProfiles)))
		var playerScores []*PlayerScore
		for _, pp := range pool.PlayerProfiles {
			// skip the requesting player
			if pp.UserId == req.UserId {
				continue
			}
			score := mmf(playerProfile, pp)
			// only non-negative score will be added
			if score > 0 {
				playerScores = append(playerScores, &PlayerScore{score: score, id: pp.UserId})
			}
		}

		sortedPlayerScores := sortByPlayerScore(playerScores)
		if len(sortedPlayerScores) > 0 {
			matchedPlayerID := sortedPlayerScores[0].id
			ctx.Logger().Debug(fmt.Sprintf("PlayerPool matched player %s vs %s", req.UserId, matchedPlayerID))
			matchedPlayerProfile = findPlayerProfileByID(pool, matchedPlayerID)
			// remove the match players from the pool
			pool = removePlayerFromPool(pool, matchedPlayerID)
			pool = removePlayerFromPool(pool, req.UserId)
			if err := savePlayerPoolFn(ctx, pool); err != nil {
				return nil, err
			}
			ctx.Logger().Debug(fmt.Sprintf("PlayerPool size after running matchmaking: %d", len(pool.PlayerProfiles)))
			break
		}
		retries++
	}

	if matchedPlayerProfile == nil {
		return nil, errors.New("Matchmaking failed, couldnt get matchedPlayerProfile")
	}

	// get matched player deck
	matchedDl, err := loadDecks(ctx, matchedPlayerProfile.UserId)
	if err != nil {
		return nil, err
	}
	matchedDeck := getDeckByID(matchedDl.Decks, matchedPlayerProfile.DeckId)
	if matchedDeck == nil {
		return nil, fmt.Errorf("deck id %d not found", matchedPlayerProfile.DeckId)
	}

	// create match
	match = &zb.Match{
		Status: zb.Match_Matching,
		PlayerStates: []*zb.PlayerState{
			&zb.PlayerState{
				Id:            playerProfile.UserId,
				Deck:          deck,
				MatchAccepted: false,
			},
			&zb.PlayerState{
				Id:            matchedPlayerProfile.UserId,
				Deck:          matchedDeck,
				MatchAccepted: false,
			},
		},
		Version: playerProfile.Version, // TODO: match version of both players
	}

	match.RandomSeed = playerProfile.RandomSeed //TODO: seed should really come from somewhere else
	if match.RandomSeed == 0 {
		match.RandomSeed = ctx.Now().Unix()
	}

	match.CustomGameAddr = playerProfile.CustomGame // TODO: make sure both players request same custom game?

	if err := createMatch(ctx, match); err != nil {
		return nil, err
	}

	// save user match
	if err := saveUserCurrentMatch(ctx, playerProfile.UserId, match); err != nil {
		return nil, err
	}
	if err := saveUserCurrentMatch(ctx, matchedPlayerProfile.UserId, match); err != nil {
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
	ctx.EmitTopics([]byte(data), match.Topics...)

	return &zb.FindMatchResponse{
		Match: match,
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
		var addr loom.Address
		var addr2 *loom.Address
		var addrStr string
		//TODO cleanup how we do this parsing
		if match.CustomGameAddr != nil {
			addrStr = fmt.Sprintf("default:%s", match.CustomGameAddr.Local.String())
		}

		addr, err = loom.ParseAddress(addrStr)
		if err != nil {
			ctx.Logger().Debug(fmt.Sprintf("no custom game mode --%v\n", err))
		} else {
			addr2 = &addr
		}

		ctx.Logger().Info(fmt.Sprintf("NewGamePlay-clientSideRuleOverride-%t\n", z.ClientSideRuleOverride))
		gp, err := NewGamePlay(ctx, match.Id, match.Version, match.PlayerStates, match.RandomSeed, addr2, z.ClientSideRuleOverride)
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
	ctx.EmitTopics([]byte(data), match.Topics...)

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

func (z *ZombieBattleground) DebugFindMatch(ctx contract.Context, req *zb.DebugFindMatchRequest) (*zb.FindMatchResponse, error) {
	// preparing user profile consisting of deck, score, ...
	dl, err := loadDecks(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	var deck *zb.Deck
	if req.DeckId > 0 {
		deck = getDeckByID(dl.Decks, req.DeckId)
	}
	if req.Deck != nil {
		deck = req.Deck
	}

	if deck == nil {
		return nil, fmt.Errorf("deck not set")
	}

	// sort tags
	if len(req.Tags) > 0 {
		sort.Strings(req.Tags)
	}

	profile := zb.PlayerProfile{
		UserId:    req.UserId,
		DeckId:    deck.Id,
		UpdatedAt: ctx.Now().Unix(),
		Tags:      req.Tags,
	}

	var loadPlayerPoolFn func(contract.StaticContext) (*zb.PlayerPool, error)
	var savePlayerPoolFn func(contract.Context, *zb.PlayerPool) error
	// if the tags is set, use tagged playerpool
	if len(profile.Tags) > 0 {
		loadPlayerPoolFn = loadTaggedPlayerPool
		savePlayerPoolFn = saveTaggedPlayerPool
	} else {
		loadPlayerPoolFn = loadPlayerPool
		savePlayerPoolFn = savePlayerPool
	}

	// perform matchmaking function to calculate scores
	// steps:
	// 1. list all the candidates that has similar profiles
	// 2. pick the most highest score
	// 3. if there is no candidate, sleep for MMWaitTime seconds
	retries := 0
	var matchedPlayerProfile *zb.PlayerProfile
	for retries < MMRetries {
		ctx.Logger().Debug(fmt.Sprintf("Matchmaking for user=%s retires=%d", req.UserId, retries))
		// load player pool
		pool, err := loadPlayerPoolFn(ctx)
		if err != nil {
			return nil, err
		}

		// prune the timed out player profile
		for _, pp := range pool.PlayerProfiles {
			updatedAt := time.Unix(pp.UpdatedAt, 0)
			if updatedAt.Add(MMTimeout).Before(ctx.Now()) {
				ctx.Logger().Debug(fmt.Sprintf("Player profile %s timedout", pp.UserId))
				// remove player from the pool
				pool = removePlayerFromPool(pool, pp.UserId)
				// remove match
				match, _ := loadUserCurrentMatch(ctx, pp.UserId)
				if match != nil {
					ctx.Delete(MatchKey(match.Id))
					match.Status = zb.Match_Timedout
					// remove player's match if existing
					ctx.Delete(UserMatchKey(pp.UserId))
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

		var playerScores []*PlayerScore
		for _, pp := range pool.PlayerProfiles {
			// skip the requesting player
			if pp.UserId == req.UserId {
				continue
			}
			score := mmf(&profile, pp)
			// only non-negtive score will be added
			if score > 0 {
				playerScores = append(playerScores, &PlayerScore{score: score, id: pp.UserId})
			}
		}

		sortedPlayerScores := sortByPlayerScore(playerScores)
		if len(sortedPlayerScores) == 0 {
			pool.PlayerProfiles = append(pool.PlayerProfiles, &profile)
			if err := savePlayerPoolFn(ctx, pool); err != nil {
				return nil, err
			}
			// create match
			match := &zb.Match{
				Status: zb.Match_Matching,
				PlayerStates: []*zb.PlayerState{
					&zb.PlayerState{
						Id:   req.UserId,
						Deck: deck,
					},
				},
				Version: req.Version,
			}

			if err := createMatch(ctx, match); err != nil {
				return nil, err
			}
			// save user match
			if err := saveUserCurrentMatch(ctx, req.UserId, match); err != nil {
				return nil, err
			}
			return &zb.FindMatchResponse{
				Match: match,
			}, nil
		}

		if len(sortedPlayerScores) > 0 {
			matchedPlayerID := sortedPlayerScores[0].id
			ctx.Logger().Debug(fmt.Sprintf("PlayerPool matched player %s vs %s", req.UserId, matchedPlayerID))
			matchedPlayerProfile = findPlayerProfileByID(pool, matchedPlayerID)
			// remove the match players from the pool
			pool = removePlayerFromPool(pool, matchedPlayerID)
			pool = removePlayerFromPool(pool, req.UserId)
			if err := savePlayerPoolFn(ctx, pool); err != nil {
				return nil, err
			}
			ctx.Logger().Debug(fmt.Sprintf("PlayerPool size after running matchmaking: %d", len(pool.PlayerProfiles)))
			break
		}
		retries++
	}

	// get and update the match
	player1 := matchedPlayerProfile.UserId
	match, err := loadUserCurrentMatch(ctx, player1)
	if err != nil && err != contract.ErrNotFound {
		return nil, err
	}

	match.PlayerStates = append(match.PlayerStates, &zb.PlayerState{
		Id:   req.UserId,
		Deck: deck,
	})
	match.Status = zb.Match_Started

	// save user match
	if err := saveUserCurrentMatch(ctx, req.UserId, match); err != nil {
		return nil, err
	}
	if err := saveUserCurrentMatch(ctx, player1, match); err != nil {
		return nil, err
	}
	// save match
	if err := saveMatch(ctx, match); err != nil {
		return nil, err
	}

	// create game state
	match.RandomSeed = req.RandomSeed
	if match.RandomSeed == 0 {
		match.RandomSeed = ctx.Now().Unix()
	}

	var addr loom.Address
	var addr2 *loom.Address
	var addrStr string
	//TODO cleanup how we do this parsing
	if req.CustomGame != nil {
		addrStr = fmt.Sprintf("default:%s", req.CustomGame.Local.String())
	}

	addr, err = loom.ParseAddress(addrStr)
	if err != nil {
		ctx.Logger().Debug(fmt.Sprintf("no custom game mode --%v\n", err))
	} else {
		addr2 = &addr
	}

	ctx.Logger().Log(fmt.Sprintf("NewGamePlay-clientSideRuleOverride-%t\n", z.ClientSideRuleOverride))
	gp, err := NewGamePlay(ctx, match.Id, req.Version, match.PlayerStates, match.RandomSeed, addr2, z.ClientSideRuleOverride)
	if err != nil {
		return nil, err
	}
	if err := saveGameState(ctx, gp.State); err != nil {
		return nil, err
	}

	// accept match
	emitMsg := zb.PlayerActionEvent{
		Match: match,
		Block: &zb.History{List: gp.history},
	}
	data, err := proto.Marshal(&emitMsg)
	if err != nil {
		return nil, err
	}
	ctx.EmitTopics([]byte(data), match.Topics...)

	return &zb.FindMatchResponse{
		Match: match,
	}, nil
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

	gp, err := GamePlayFrom(gamestate, z.ClientSideRuleOverride)
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
	match.Topics = append(match.Topics, "endgame")
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
	gp, err := GamePlayFrom(gamestate, z.ClientSideRuleOverride)
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
	gp, err := GamePlayFrom(gamestate, z.ClientSideRuleOverride)
	if err != nil {
		return nil, err
	}
	// add created timestamp
	req.PlayerAction.CreatedAt = ctx.Now().Unix()
	if err := gp.AddAction(req.PlayerAction); err != nil {
		return nil, err
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
		GameState: gamestate,
		Match:     match,
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
	gp, err := GamePlayFrom(gamestate, z.ClientSideRuleOverride)
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
	ctx.GrantPermission([]byte(params.NewOracle.String()), []string{"oracle"})

	if err := ctx.Set(oracleKey, params.NewOracle); err != nil {
		return errors.Wrap(err, "setting new oracle")
	}
	return nil
}

func (z *ZombieBattleground) validateOracle(ctx contract.Context, zo *types.Address) error {
	if ok, _ := ctx.HasPermission([]byte(zo.String()), []string{"oracle"}); !ok {
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

var Contract plugin.Contract = contract.MakePluginContract(&ZombieBattleground{ClientSideRuleOverride: true})

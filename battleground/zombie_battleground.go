package battleground

import (
	"fmt"
	"sort"
	"strings"

	"github.com/loomnetwork/go-loom/plugin"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/pkg/errors"
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
	// initialize card library
	cardList := zb.CardList{
		Cards: req.Cards,
	}
	if err := ctx.Set(cardListKey, &cardList); err != nil {
		return err
	}
	// initialize heros
	heroList := zb.HeroList{
		Heroes: req.Heroes,
	}
	if err := ctx.Set(heroListKey, &heroList); err != nil {
		return err
	}

	cardCollectionList := zb.CardCollectionList{
		Cards: req.DefaultCollection,
	}
	if err := ctx.Set(defaultCollectionKey, &cardCollectionList); err != nil {
		return err
	}

	// initialize default deck
	deckList := zb.DeckList{
		Decks: req.DefaultDecks,
	}
	if err := ctx.Set(defaultDeckKey, &deckList); err != nil {
		return err
	}

	defaultHeroList := zb.HeroList{
		Heroes: req.Heroes,
	}
	if err := ctx.Set(defaultHeroesKey, &defaultHeroList); err != nil {
		return err
	}

	return nil
}

func (z *ZombieBattleground) GetAccount(ctx contract.StaticContext, req *zb.GetAccountRequest) (*zb.Account, error) {
	var account zb.Account
	userKeySpace := NewUserKeySpace(req.UserId)

	if err := ctx.Get(userKeySpace.AccountKey(), &account); err != nil {
		return nil, errors.Wrapf(err, "unable to retrieve account data for userId: %s", req.UserId)
	}

	return &account, nil
}

func (z *ZombieBattleground) UpdateAccount(ctx contract.Context, req *zb.UpsertAccountRequest) (*zb.Account, error) {
	var account zb.Account

	senderAddress := []byte(ctx.Message().Sender.Local)
	userKeySpace := NewUserKeySpace(req.UserId)

	// Verify whether this privateKey associated with user
	if !isUser(ctx, req.UserId) {
		return nil, fmt.Errorf("userId: %s is not verified", req.UserId)
	}

	if err := ctx.Get(userKeySpace.AccountKey(), &account); err != nil {
		return nil, errors.Wrapf(err, "unable to retrieve account data for userId: %s", req.UserId)
	}

	copyAccountInfo(&account, req)
	if err := ctx.Set(userKeySpace.AccountKey(), &account); err != nil {
		return nil, errors.Wrapf(err, "error setting account information for userId: %s", req.UserId)
	}

	emitMsgJSON, err := prepareEmitMsgJSON(senderAddress, req.UserId, "updateaccount")
	if err == nil {
		ctx.EmitTopics(emitMsgJSON, "zombiebattleground:updateaccount")
	}

	return &account, nil
}

func (z *ZombieBattleground) CreateAccount(ctx contract.Context, req *zb.UpsertAccountRequest) error {
	senderAddress := []byte(ctx.Message().Sender.Local)
	userKeySpace := NewUserKeySpace(req.UserId)
	var deckList zb.DeckList

	// confirm owner doesnt exist already
	if ctx.Has(userKeySpace.AccountKey()) {
		return errors.New("user already exists")
	}

	var account zb.Account
	account.UserId = req.UserId
	account.Owner = ctx.Message().Sender.Bytes()

	copyAccountInfo(&account, req)

	if err := ctx.Set(userKeySpace.AccountKey(), &account); err != nil {
		return errors.Wrapf(err, "error setting account information for userId: %s", req.UserId)
	}

	ctx.GrantPermission([]byte(req.UserId), []string{"user"})

	// add default collection list
	var collectionList zb.CardCollectionList
	if err := ctx.Get(defaultCollectionKey, &collectionList); err != nil {
		return errors.Wrapf(err, "unable to get default collectionlist")
	}
	if err := ctx.Set(userKeySpace.CardCollectionKey(), &collectionList); err != nil {
		return errors.Wrapf(err, "unable to save card collection for userId: %s", req.UserId)
	}

	if err := ctx.Get(defaultDeckKey, &deckList); err != nil {
		return errors.Wrapf(err, "unable to get default decks")
	}
	if err := ctx.Set(userKeySpace.DecksKey(), &deckList); err != nil {
		return errors.Wrapf(err, "unable to save decks for userId: %s", req.UserId)
	}

	var heroes zb.HeroList
	if err := ctx.Get(defaultHeroesKey, &heroes); err != nil {
		return errors.Wrapf(err, "unable to get default hero")
	}
	if err := ctx.Set(userKeySpace.HeroesKey(), &heroes); err != nil {
		return errors.Wrapf(err, "unable to save heroinfo for userId: %s", req.UserId)
	}

	emitMsgJSON, err := prepareEmitMsgJSON(senderAddress, req.UserId, "createaccount")
	if err == nil {
		ctx.EmitTopics(emitMsgJSON, "zombiebattleground:createaccount")
	}

	return nil
}

// CreateDeck appends the given deck to user's deck list
func (z *ZombieBattleground) CreateDeck(ctx contract.Context, req *zb.CreateDeckRequest) (*zb.CreateDeckResponse, error) {
	userID := strings.TrimSpace(req.UserId)
	userKeySpace := NewUserKeySpace(userID)

	if req.Deck == nil {
		return nil, fmt.Errorf("deck must not be nil")
	}

	if !isUser(ctx, userID) {
		return nil, fmt.Errorf("user is not verified")
	}

	var userCollection zb.CardCollectionList
	if err := ctx.Get(userKeySpace.CardCollectionKey(), &userCollection); err != nil {
		return nil, errors.Wrapf(err, "unable to get collection data for userId: %s", req.UserId)
	}

	if err := validateDeckCollections(userCollection.Cards, req.Deck.Cards); err != nil {
		return nil, err
	}

	var heroes zb.HeroList
	if err := ctx.Get(userKeySpace.HeroesKey(), &heroes); err != nil {
		return nil, err
	}

	if err := validateDeckHero(heroes.Heroes, req.Deck.HeroId); err != nil {
		return nil, err
	}

	var deckList zb.DeckList
	err := ctx.Get(userKeySpace.DecksKey(), &deckList)
	if err != nil && err != contract.ErrNotFound {
		return nil, err
	}

	if err := validateDeckName(deckList.Decks, req.Deck); err != nil {
		return nil, err
	}

	// allocate new deck id
	var newDeckId int64 = 0
	if len(deckList.Decks) != 0 {
		for _, deck := range deckList.Decks {
			if deck.Id > newDeckId {
				newDeckId = deck.Id
			}
		}

		newDeckId++
	}

	req.Deck.Id = newDeckId

	deckList.Decks = mergeDeckSets(deckList.Decks, []*zb.Deck{req.Deck})
	deckList.LastModificationTimestamp = req.LastModificationTimestamp;

	if err := ctx.Set(userKeySpace.DecksKey(), &deckList); err != nil {
		return nil, err
	}

	senderAddress := []byte(ctx.Message().Sender.Local)
	emitMsgJSON, err := prepareEmitMsgJSON(senderAddress, userID, "createdeck")
	if err == nil {
		ctx.EmitTopics(emitMsgJSON, "zombiebattleground:createdeck")
	}

	return &zb.CreateDeckResponse{DeckId: newDeckId}, nil
}

// EditDeck edits the deck by id
func (z *ZombieBattleground) EditDeck(ctx contract.Context, req *zb.EditDeckRequest) error {
	if req.Deck == nil {
		return fmt.Errorf("deck must not be nil")
	}
	if !isUser(ctx, req.UserId) {
		return fmt.Errorf("user is not verified")
	}

	userKeySpace := NewUserKeySpace(req.UserId)
	var userCollection zb.CardCollectionList
	if err := ctx.Get(userKeySpace.CardCollectionKey(), &userCollection); err != nil {
		return errors.Wrapf(err, "unable to get collection data for userId: %s", req.UserId)
	}

	var deckList zb.DeckList
	err := ctx.Get(userKeySpace.DecksKey(), &deckList)
	if err != nil && err != contract.ErrNotFound {
		return err
	}
	if err := validateDeckCollections(userCollection.Cards, req.Deck.Cards); err != nil {
		return err
	}

	var heroes zb.HeroList
	if err := ctx.Get(userKeySpace.HeroesKey(), &heroes); err != nil {
		return err
	}
	if err := validateDeckHero(heroes.Heroes, req.Deck.HeroId); err != nil {
		return err
	}

	if err := editDeck(deckList.Decks, req.Deck); err != nil {
		return err
	}

	deckList.LastModificationTimestamp = req.LastModificationTimestamp
	if err := ctx.Set(userKeySpace.DecksKey(), &deckList); err != nil {
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
	userID := strings.TrimSpace(req.UserId)
	userKeySpace := NewUserKeySpace(userID)

	if !isUser(ctx, userID) {
		return fmt.Errorf("user is not verified")
	}

	var deckList zb.DeckList
	err := ctx.Get(userKeySpace.DecksKey(), &deckList)
	if err == contract.ErrNotFound {
		return err
	}
	if err != nil {
		return err
	}

	var deleted bool
	deckList.Decks, deleted = deleteDeckById(deckList.Decks, req.DeckId)
	if !deleted {
		return fmt.Errorf("deck not found")
	}

	deckList.LastModificationTimestamp = req.LastModificationTimestamp
	if err := ctx.Set(userKeySpace.DecksKey(), &deckList); err != nil {
		return err
	}
	return nil
}

// ListDecks returns the user's decks
func (z *ZombieBattleground) ListDecks(ctx contract.StaticContext, req *zb.ListDecksRequest) (*zb.ListDecksResponse, error) {
	userID := strings.TrimSpace(req.UserId)
	userKeySpace := NewUserKeySpace(userID)

	var deckList zb.DeckList
	err := ctx.Get(userKeySpace.DecksKey(), &deckList)

	if err != nil && err != contract.ErrNotFound {
		return nil, err
	}

	return &zb.ListDecksResponse{
		Decks: deckList.Decks,
		LastModificationTimestamp: deckList.LastModificationTimestamp,
	}, nil
}

// GetDeck returns the deck by given id
func (z *ZombieBattleground) GetDeck(ctx contract.StaticContext, req *zb.GetDeckRequest) (*zb.GetDeckResponse, error) {
	userID := strings.TrimSpace(req.UserId)
	userKeySpace := NewUserKeySpace(userID)

	var deckList zb.DeckList
	err := ctx.Get(userKeySpace.DecksKey(), &deckList)

	if err != nil && err != contract.ErrNotFound {
		return nil, err
	}

	deck := getDeckById(deckList.Decks, req.DeckId)
	if deck == nil {
		return nil, contract.ErrNotFound
	}
	return &zb.GetDeckResponse{Deck: deck}, nil
}

// GetCollection returns the collection of the card own by the user
func (z *ZombieBattleground) GetCollection(ctx contract.StaticContext, req *zb.GetCollectionRequest) (*zb.GetCollectionResponse, error) {
	userID := strings.TrimSpace(req.UserId)
	userKeySpace := NewUserKeySpace(userID)

	var collectionList zb.CardCollectionList
	err := ctx.Get(userKeySpace.CardCollectionKey(), &collectionList)
	if err == contract.ErrNotFound {
		return &zb.GetCollectionResponse{}, err
	}
	if err != nil {
		return nil, err
	}

	return &zb.GetCollectionResponse{Cards: collectionList.Cards}, nil
}

// ListCardLibrary list all the card library data
func (z *ZombieBattleground) ListCardLibrary(ctx contract.StaticContext, req *zb.ListCardLibraryRequest) (*zb.ListCardLibraryResponse, error) {
	var cardList zb.CardList
	if err := ctx.Get(cardListKey, &cardList); err != nil {
		return nil, err
	}
	// convert to card list to card library view grouped by set
	var category = make(map[string][]*zb.Card)
	for _, card := range cardList.Cards {
		if _, ok := category[card.Set]; !ok {
			category[card.Set] = make([]*zb.Card, 0)
		}
		category[card.Set] = append(category[card.Set], card)
	}
	// order sets by name
	var setNames []string
	for k := range category {
		setNames = append(setNames, k)
	}
	sort.Strings(setNames)

	var sets []*zb.CardSet
	for _, setName := range setNames {
		cards, ok := category[setName]
		if !ok {
			continue
		}
		set := &zb.CardSet{
			Name:  setName,
			Cards: cards,
		}
		sets = append(sets, set)
	}

	return &zb.ListCardLibraryResponse{Sets: sets}, nil
}

// ListHero return all the heros
func (z *ZombieBattleground) ListHeroLibrary(ctx contract.StaticContext, req *zb.ListHeroLibraryRequest) (*zb.ListHeroLibraryResponse, error) {
	var heroList zb.HeroList
	if err := ctx.Get(heroListKey, &heroList); err != nil {
		return nil, err
	}

	return &zb.ListHeroLibraryResponse{Heroes: heroList.Heroes}, nil
}

func (z *ZombieBattleground) ListHeroes(ctx contract.StaticContext, req *zb.ListHeroesRequest) (*zb.ListHeroesResponse, error) {
	userKeySpace := NewUserKeySpace(req.UserId)
	var heroes zb.HeroList

	if err := ctx.Get(userKeySpace.HeroesKey(), &heroes); err != nil {
		return nil, err
	}

	return &zb.ListHeroesResponse{Heroes: heroes.Heroes}, nil
}

func (z *ZombieBattleground) GetHero(ctx contract.StaticContext, req *zb.GetHeroRequest) (*zb.GetHeroResponse, error) {
	userKeySpace := NewUserKeySpace(req.UserId)
	var heroes zb.HeroList

	if err := ctx.Get(userKeySpace.HeroesKey(), &heroes); err != nil {
		return nil, err
	}

	hero := getHeroById(heroes.Heroes, req.HeroId)
	if hero == nil {
		return nil, contract.ErrNotFound
	}

	return &zb.GetHeroResponse{Hero: hero}, nil

}

func (z *ZombieBattleground) AddHeroExperience(ctx contract.Context, req *zb.AddHeroExperienceRequest) (*zb.AddHeroExperienceResponse, error) {
	userKeySpace := NewUserKeySpace(req.UserId)
	var heroes zb.HeroList

	if req.Experience <= 0 {
		return nil, fmt.Errorf("experience needs to be greater than zero.")
	}

	if !isUser(ctx, strings.TrimSpace(req.UserId)) {
		return nil, fmt.Errorf("user is not verified")
	}

	if err := ctx.Get(userKeySpace.HeroesKey(), &heroes); err != nil {
		return nil, err
	}

	hero := getHeroById(heroes.Heroes, req.HeroId)
	if hero == nil {
		return nil, contract.ErrNotFound
	}

	hero.Experience += req.Experience

	if err := ctx.Set(userKeySpace.HeroesKey(), &heroes); err != nil {
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
	userKeySpace := NewUserKeySpace(req.UserId)
	var heroes zb.HeroList

	if err := ctx.Get(userKeySpace.HeroesKey(), &heroes); err != nil {
		return nil, err
	}

	hero := getHeroById(heroes.Heroes, req.HeroId)
	if hero == nil {
		return nil, contract.ErrNotFound
	}

	return &zb.GetHeroSkillsResponse{HeroId: hero.HeroId, Skills: hero.Skills}, nil
}

var Contract plugin.Contract = contract.MakePluginContract(&ZombieBattleground{})

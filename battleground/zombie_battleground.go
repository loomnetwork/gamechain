package battleground

import (
	"encoding/json"
	"fmt"
	"sort"

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
	if err := ctx.Get(defaultCollectionKey, &collectionList); err != nil {
		return errors.Wrapf(err, "unable to get default collectionlist")
	}
	if err := ctx.Set(CardCollectionKey(req.UserId), &collectionList); err != nil {
		return errors.Wrapf(err, "unable to save card collection for userId: %s", req.UserId)
	}

	var deckList zb.DeckList
	if err := ctx.Get(defaultDeckKey, &deckList); err != nil {
		return errors.Wrapf(err, "unable to get default decks")
	}
	if err := ctx.Set(DecksKey(req.UserId), &deckList); err != nil {
		return errors.Wrapf(err, "unable to save decks for userId: %s", req.UserId)
	}

	var heroes zb.HeroList
	if err := ctx.Get(defaultHeroesKey, &heroes); err != nil {
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
	// validate user card collection
	userCollection, err := loadCardCollection(ctx, req.UserId)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get collection data for userId: %s", req.UserId)
	}
	// make sure the given cards and amount must be a subset of user's cards
	if err := validateDeckCollections(userCollection.Cards, req.Deck.Cards); err != nil {
		return nil, err
	}

	deckList, err := loadDecks(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	// check duplicate name
	if existing := getDeckByName(deckList.Decks, req.Deck.Name); existing != nil {
		return nil, ErrDeckNameExists
	}
	// allocate new deck id
	// TODO: check if this won't cause nondeterministic result
	var newDeckID int64
	if len(deckList.Decks) != 0 {
		for _, deck := range deckList.Decks {
			if deck.Id > newDeckID {
				newDeckID = deck.Id
			}
		}
		newDeckID++
	}
	req.Deck.Id = newDeckID
	deckList.Decks = append(deckList.Decks, req.Deck)
	deckList.LastModificationTimestamp = req.LastModificationTimestamp
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
	// validate user card collection
	userCollection, err := loadCardCollection(ctx, req.UserId)
	if err != nil {
		return errors.Wrapf(err, "unable to get collection data for userId: %s", req.UserId)
	}
	if err := validateDeckCollections(userCollection.Cards, req.Deck.Cards); err != nil {
		return err
	}
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
	deckList.LastModificationTimestamp = req.LastModificationTimestamp
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

	deckList.LastModificationTimestamp = req.LastModificationTimestamp
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
		LastModificationTimestamp: deckList.LastModificationTimestamp,
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

func (z *ZombieBattleground) ListHeroLibrary(ctx contract.StaticContext, req *zb.ListHeroLibraryRequest) (*zb.ListHeroLibraryResponse, error) {
	var heroList zb.HeroList
	if err := ctx.Get(heroListKey, &heroList); err != nil {
		return nil, err
	}
	return &zb.ListHeroLibraryResponse{Heroes: heroList.Heroes}, nil
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

func (z *ZombieBattleground) FindMatch(ctx contract.Context, req *zb.FindMatchRequest) (*zb.FindMatchResponse, error) {
	// find the room available for the user to be filled in; otherwise, create a new one
	roomlist, err := loadRoomList(ctx)
	if err != nil {
		return nil, err
	}
	if len(roomlist.Rooms) == 0 {
		room := zb.Room{
			Id:     1, // fixed for now
			Topics: []string{fmt.Sprintf("room:%d", 1)},
			Status: zb.Room_Waiting,
			PlayerStates: []*zb.PlayerState{
				&zb.PlayerState{
					Id:     req.UserId,
					Status: zb.PlayerState_Ready,
				},
			},
		}
		roomlist.Rooms = append(roomlist.Rooms, &room)
		if err := saveLoomList(ctx, roomlist); err != nil {
			return nil, err
		}
		return &zb.FindMatchResponse{
			Room: &room,
		}, nil
	}

	room := roomlist.Rooms[0]
	if req.UserId != room.PlayerStates[0].Id {
		room.PlayerStates = append(room.PlayerStates, &zb.PlayerState{
			Id:     req.UserId,
			Status: zb.PlayerState_Ready,
		})
	}

	if err := saveLoomList(ctx, roomlist); err != nil {
		return nil, err
	}

	// the return result should include the topic for the client to subscribe to
	return &zb.FindMatchResponse{
		Room: room,
	}, nil
}

func (z *ZombieBattleground) StartMatch(ctx contract.Context, req *zb.StartMatchRequest) error {
	// find the room available for the user to be filled in; otherwise, create a new one
	roomlist, err := loadRoomList(ctx)
	if err != nil {
		return err
	}
	if len(roomlist.Rooms) == 0 {
		return contract.ErrNotFound
	}

	room := roomlist.Rooms[0]
	room.Status = zb.Room_Ready
	emitMsg := zb.EventRoom{
		Room: room,
	}

	data, err := json.Marshal(emitMsg)
	if err != nil {
		return err
	}
	if err == nil {
		ctx.EmitTopics(data, room.Topics[0])
	}

	return nil
}

func (z *ZombieBattleground) SendAction(ctx contract.Context, req *zb.ActionRequest) error {
	roomlist, err := loadRoomList(ctx)
	if err != nil {
		return err
	}
	if len(roomlist.Rooms) == 0 {
		return contract.ErrNotFound
	}

	room := roomlist.Rooms[0]

	emitMsg := zb.ActionEvent{
		RoomId:  room.Id,
		UserId:  req.UserId,
		Message: req.Message,
	}
	data, err := json.Marshal(emitMsg)
	if err != nil {
		return err
	}
	if err == nil {
		ctx.EmitTopics(data, room.Topics[0])
	}

	return nil
}

var Contract plugin.Contract = contract.MakePluginContract(&ZombieBattleground{})

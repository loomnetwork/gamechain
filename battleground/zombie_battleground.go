package battleground

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"sort"
	"unicode/utf8"

	"github.com/golang/protobuf/jsonpb"
	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/plugin"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	"github.com/loomnetwork/go-loom/types"
	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/pkg/errors"
)

type ZombieBattleground struct {
}

const (
	MaxGameModeNameChar        = 48
	MaxGameModeDescriptionChar = 255
	MaxGameModeVersionChar     = 16
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

	return nil
}

func (z *ZombieBattleground) UpdateInit(ctx contract.Context, req *zb.UpdateInitRequest) error {

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

	return nil
}

func (z *ZombieBattleground) UpdateCardList(ctx contract.Context, req *zb.UpdateCardListRequest) error {
	cardList := zb.CardList{
		Cards: req.Cards,
	}
	if err := ctx.Set(MakeVersionedKey(req.Version, cardListKey), &cardList); err != nil {
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
	if err := ctx.Get(MakeVersionedKey(req.Version, cardListKey), &cardList); err != nil {
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
	if err := ctx.Get(MakeVersionedKey(req.Version, heroListKey), &heroList); err != nil {
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
	// load deck id
	dl, err := loadDecks(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	deck := getDeckByID(dl.Decks, req.DeckId)
	if deck == nil {
		return nil, fmt.Errorf("deck id %d not found", req.DeckId)
	}

	// register the user to match making pool
	// TODO: chan ge to scan users in matchmakings
	infos, err := loadMatchMakingInfoList(ctx)
	if err != nil {
		return nil, err
	}

	var info *zb.MatchMakingInfo
	for _, inf := range infos.Infos {
		if inf.UserId == req.UserId {
			continue
		}
		info = inf
	}

	if info == nil {
		// save user info
		info = &zb.MatchMakingInfo{
			UserId: req.UserId,
			Deck:   deck,
		}
		infos.Infos = append(infos.Infos, info)
		if err := saveMatchMakingInfoList(ctx, infos); err != nil {
			return nil, err
		}

		match := &zb.Match{
			Status: zb.Match_Matching,
			PlayerStates: []*zb.PlayerState{
				&zb.PlayerState{
					Id:   req.UserId,
					Deck: deck,
				},
			},
		}

		if err := createMatch(ctx, match); err != nil {
			return nil, err
		}
		// save user match
		// TODO: clean up the previous match?
		if err := saveUserMatch(ctx, req.UserId, match); err != nil {
			return nil, err
		}
		return &zb.FindMatchResponse{
			Match: match,
		}, nil
	}

	// get and update the match
	opponentID := info.UserId
	match, err := loadUserMatch(ctx, opponentID)
	if err != nil && err != contract.ErrNotFound {
		return nil, err
	}
	match.PlayerStates = append(match.PlayerStates, &zb.PlayerState{
		Id:   req.UserId,
		Deck: deck,
	})
	match.Status = zb.Match_Started

	// save user match
	// TODO: clean up the previous match?
	if err := saveUserMatch(ctx, req.UserId, match); err != nil {
		return nil, err
	}

	// remove info from match making list by making sure that only second player remove it once
	newinfos := make([]*zb.MatchMakingInfo, 0)
	for _, inf := range infos.Infos {
		if inf.UserId == opponentID {
			continue
		}
		newinfos = append(newinfos, inf)
	}
	infos.Infos = newinfos
	if err := saveMatchMakingInfoList(ctx, infos); err != nil {
		return nil, err
	}

	if err := saveMatch(ctx, match); err != nil {
		return nil, err
	}

	// create game state
	seed := ctx.Now().Unix()
	gp, err := NewGamePlay(match.Id, match.PlayerStates, seed)
	if err != nil {
		return nil, err
	}
	if err := saveGameState(ctx, gp.State); err != nil {
		return nil, err
	}

	// accept match
	emitMsg := zb.PlayerActionEvent{
		Match:     match,
		GameState: gp.State,
	}
	data, err := new(jsonpb.Marshaler).MarshalToString(&emitMsg)
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

func (z *ZombieBattleground) GetMatch(ctx contract.Context, req *zb.GetMatchRequest) (*zb.GetMatchResponse, error) {
	match, err := loadMatch(ctx, req.MatchId)
	if err != nil {
		return nil, err
	}
	gameState, _ := loadGameState(ctx, req.MatchId)

	return &zb.GetMatchResponse{
		Match:     match,
		GameState: gameState,
	}, nil
}

func (z *ZombieBattleground) GetGameState(ctx contract.Context, req *zb.GetGameStateRequest) (*zb.GetGameStateResponse, error) {
	gameState, err := loadGameState(ctx, req.MatchId)
	if err != nil {
		return nil, err
	}

	return &zb.GetGameStateResponse{
		GameState: gameState,
	}, nil
}

func (z *ZombieBattleground) SetGameState(ctx contract.Context, req *zb.SetGameStateRequest) (*zb.SetGameStateResponse, error) {
	err := saveGameState(ctx, req.GameState)
	if err != nil {
		return nil, err
	}

	return &zb.SetGameStateResponse{}, nil
}

func (z *ZombieBattleground) LeaveMatch(ctx contract.Context, req *zb.LeaveMatchRequest) (*zb.LeaveMatchResponse, error) {
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

	// TODO: Change on gamestate

	emitMsg := zb.PlayerActionEvent{
		UserId: req.UserId,
		Match:  match,
	}
	data, err := new(jsonpb.Marshaler).MarshalToString(&emitMsg)
	if err != nil {
		return nil, err
	}
	if err == nil {
		ctx.EmitTopics([]byte(data), match.Topics...)
	}

	return &zb.LeaveMatchResponse{}, nil
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
	gp, err := GamePlayFrom(gamestate)
	if err != nil {
		return nil, err
	}
	gp.PrintState()
	if err := gp.AddAction(req.PlayerAction); err != nil {
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

	emitMsg := zb.PlayerActionEvent{
		PlayerActionType: req.PlayerAction.ActionType,
		UserId:           req.PlayerAction.PlayerId,
		PlayerAction:     req.PlayerAction,
		Match:            match,
		GameState:        gamestate,
	}
	data, err := new(jsonpb.Marshaler).MarshalToString(&emitMsg)
	if err != nil {
		return nil, err
	}
	if err == nil {
		ctx.EmitTopics([]byte(data), match.Topics...)
	}

	return &zb.PlayerActionResponse{
		GameState: gamestate,
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

var Contract plugin.Contract = contract.MakePluginContract(&ZombieBattleground{})

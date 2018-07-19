package battleground

import (
	"encoding/json"
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
		Heros: req.Heros,
	}
	if err := ctx.Set(heroListKey, &heroList); err != nil {
		return err
	}
	// initialize default card collection
	collectionList := zb.CardCollectionList{
		Cards: req.DefaultCollection,
	}
	if err := ctx.Set(defaultCollectionKey, &collectionList); err != nil {
		return err
	}
	// initialize default deck
	deckList := zb.DeckList{
		Decks: req.DefaultDecks,
	}
	if err := ctx.Set(defaultDeckKey, &deckList); err != nil {
		return err
	}
	return nil
}

func (z *ZombieBattleground) getDecks(deckSet []*zb.Deck, decksToQuery []string) []*zb.Deck {
	deckMap := make(map[string]*zb.Deck)
	decks := make([]*zb.Deck, len(decksToQuery))

	for _, deck := range deckSet {
		deckMap[deck.Name] = deck
	}

	i := 0
	for _, deckName := range decksToQuery {
		if _, ok := deckMap[deckName]; !ok {
			continue
		}

		decks[i] = deckMap[deckName]
		i++
	}

	return decks[:i]
}

func (z *ZombieBattleground) deleteDecks(deckSet []*zb.Deck, decksToDelete []string) ([]*zb.Deck, bool, error) {
	deckMap := make(map[string]*zb.Deck)

	for _, deck := range deckSet {
		deckMap[deck.Name] = deck
	}

	for _, deckName := range decksToDelete {
		delete(deckMap, deckName)
	}

	newArray := make([]*zb.Deck, len(deckMap))

	if len(newArray) == 0 {
		return nil, false, errors.New("cannot delete only deck available")
	}

	i := 0
	for _, deck := range deckSet {
		if _, ok := deckMap[deck.Name]; !ok {
			continue
		}

		newArray[i] = deck
		i++
	}

	return newArray, len(newArray) == len(deckSet), nil
}

func (z *ZombieBattleground) isUser(ctx contract.Context, userId string) bool {
	ok, _ := ctx.HasPermission([]byte(userId), []string{"user"})
	return ok
}

func (z *ZombieBattleground) prepareEmitMsgJSON(address []byte, owner, method string) ([]byte, error) {
	emitMsg := struct {
		Owner  string
		Method string
		Addr   []byte
	}{owner, method, address}

	return json.Marshal(emitMsg)
}

func (z *ZombieBattleground) copyAccountInfo(account *zb.Account, req *zb.UpsertAccountRequest) {
	account.PhoneNumberVerified = req.PhoneNumberVerified
	account.RewardRedeemed = req.RewardRedeemed
	account.IsKickstarter = req.IsKickstarter
	account.Image = req.Image
	account.EmailNotification = req.EmailNotification
	account.EloScore = req.EloScore
	account.CurrentTier = req.CurrentTier
	account.GameMembershipTier = req.GameMembershipTier
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
	userId := strings.TrimSpace(req.UserId)
	userKeySpace := NewUserKeySpace(userId)

	// Verify whether this privateKey associated with user
	if !z.isUser(ctx, userId) {
		return nil, fmt.Errorf("userId: %s is not verified", req.UserId)
	}

	if err := ctx.Get(userKeySpace.AccountKey(), &account); err != nil {
		return nil, errors.Wrapf(err, "unable to retrieve account data for userId: %s", req.UserId)
	}

	z.copyAccountInfo(&account, req)
	if err := ctx.Set(userKeySpace.AccountKey(), &account); err != nil {
		return nil, errors.Wrapf(err, "error setting account information for userId: %s", req.UserId)
	}

	ctx.Logger().Info("updated zombiebattleground account", "user_id", userId, "address", senderAddress)

	emitMsgJSON, err := z.prepareEmitMsgJSON(senderAddress, userId, "updateaccount")
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("error marshalling emit message for userId:%s. Error:%s", req.UserId, err))
	} else {
		ctx.EmitTopics(emitMsgJSON, "zombiebattleground:updateaccount")
	}

	return &account, nil
}

func (z *ZombieBattleground) CreateAccount(ctx contract.Context, req *zb.UpsertAccountRequest) error {
	userId := strings.TrimSpace(req.UserId)
	senderAddress := []byte(ctx.Message().Sender.Local)
	userKeySpace := NewUserKeySpace(userId)

	// confirm owner doesnt exist already
	if ctx.Has(userKeySpace.AccountKey()) {
		return errors.New("user already exists.\n")
	}

	var account zb.Account
	account.UserId = userId
	account.Owner = ctx.Message().Sender.Bytes()

	z.copyAccountInfo(&account, req)

	if err := ctx.Set(userKeySpace.AccountKey(), &account); err != nil {
		return errors.Wrapf(err, "error setting account information for userId: %s", req.UserId)
	}

	ctx.GrantPermission([]byte(userId), []string{"user"})

	// add default collection list
	var collectionList zb.CardCollectionList
	if err := ctx.Get(defaultCollectionKey, &collectionList); err != nil {
		return err
	}
	if err := ctx.Set(userKeySpace.CardCollectionKey(), &collectionList); err != nil {
		return err
	}

	var deckList zb.DeckList
	if err := ctx.Get(defaultDeckKey, &deckList); err != nil {
		return err
	}
	if err := ctx.Set(userKeySpace.DecksKey(), &deckList); err != nil {
		return err
	}

	ctx.Logger().Info("created zombiebattleground account", "userId", userId, "address", senderAddress)

	emitMsgJSON, err := z.prepareEmitMsgJSON(senderAddress, userId, "createaccount")
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("error marshalling emit message for userId:%s. Error:%s", req.UserId, err))
	} else {
		ctx.EmitTopics(emitMsgJSON, "zombiebattleground:createaccount")
	}

	return nil
}

// Deck related functions
func (z *ZombieBattleground) GetDecks(ctx contract.StaticContext, req *zb.GetDecksRequest) (*zb.DeckList, error) {
	var DeckList zb.DeckList

	userId := strings.TrimSpace(req.UserId)
	userKeySpace := NewUserKeySpace(userId)

	if err := ctx.Get(userKeySpace.DecksKey(), &DeckList); err != nil {
		return nil, errors.Wrapf(err, "unable to get decks for userId: %s", userId)
	}

	return &DeckList, nil
}

func (z *ZombieBattleground) GetDeck(ctx contract.StaticContext, req *zb.GetDeckRequest) (*zb.Deck, error) {
	var DeckList zb.DeckList

	userId := strings.TrimSpace(req.UserId)
	userKeySpace := NewUserKeySpace(userId)
	deckId := strings.TrimSpace(req.DeckId)

	if err := ctx.Get(userKeySpace.DecksKey(), &DeckList); err != nil {
		return nil, errors.Wrapf(err, "unable to get decks for userId: %s", userId)
	}

	result := z.getDecks(DeckList.Decks, []string{deckId})

	if len(result) != 0 {
		return result[0], nil
	} else {
		return nil, fmt.Errorf("deck: %s not found.", deckId)
	}
}

func (z *ZombieBattleground) DeleteDeck(ctx contract.Context, req *zb.DeleteDeckRequest) error {
	var DeckList zb.DeckList
	var err error
	var deleted bool

	userId := strings.TrimSpace(req.UserId)
	senderAddress := []byte(ctx.Message().Sender.Local)
	userKeySpace := NewUserKeySpace(userId)
	deckId := strings.TrimSpace(req.DeckId)

	if !z.isUser(ctx, userId) {
		return fmt.Errorf("userId: %s is not verified", req.UserId)
	}

	if err = ctx.Get(userKeySpace.DecksKey(), &DeckList); err != nil {
		return errors.Wrapf(err, "unable to get decks for userId: %s", userId)
	}

	if DeckList.Decks, deleted, err = z.deleteDecks(DeckList.Decks, []string{deckId}); err != nil {
		return errors.Wrapf(err, "unable to delete deck: %s", deckId)
	}

	if err = ctx.Set(userKeySpace.DecksKey(), &DeckList); err != nil {
		return errors.Wrapf(err, "unable to save decks for userId: %s", userId)
	}

	if deleted {
		ctx.Logger().Info("Deleted zombiebattleground deck", "userId", userId, "deckId", deckId, "address", senderAddress)

		emitMsgJSON, err := z.prepareEmitMsgJSON(senderAddress, userId, "deletedeck")
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("error marshalling emit message for userId:%s. Error:%s", req.UserId, err))
		} else {
			ctx.EmitTopics(emitMsgJSON, "zombiebattleground:deletedeck")
		}
	}

	return nil
}

func (z *ZombieBattleground) AddDeck(ctx contract.Context, req *zb.AddDeckRequest) error {
	userId := strings.TrimSpace(req.UserId)
	senderAddress := []byte(ctx.Message().Sender.Local)
	userKeySpace := NewUserKeySpace(userId)

	if !z.isUser(ctx, userId) {
		return fmt.Errorf("userId: %s is not verified", req.UserId)
	}

	var collectionList zb.CardCollectionList
	if err := ctx.Get(userKeySpace.CardCollectionKey(), &collectionList); err != nil {
		return errors.Wrapf(err, "unable to get collection data for userId: %s", req.UserId)
	}

	var deckList zb.DeckList
	if err := ctx.Get(userKeySpace.DecksKey(), &deckList); err != nil {
		return errors.Wrapf(err, "unable to get decks for userId: %s", userId)
	}

	// TODO Check if req.Deck is nil
	if err := validateDeckAddition(collectionList.Cards, req.Deck.Cards); err != nil {
		return errors.Wrapf(err, "unable to validate deck data for userId: %s", req.UserId)
	}

	deckList.Decks = mergeDeckSets(deckList.Decks, []*zb.Deck{req.Deck})

	if err := ctx.Set(userKeySpace.DecksKey(), &deckList); err != nil {
		return errors.Wrapf(err, "unable to save decks for userId: %s", userId)
	}

	ctx.Logger().Info("Created zombiebattleground deck", "userId", userId, "deckId", req.Deck.Name, "address", senderAddress)

	emitMsgJSON, err := z.prepareEmitMsgJSON(senderAddress, userId, "adddeck")
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("error marshalling emit message for userId:%s. Error:%s", req.UserId, err))
	} else {
		ctx.EmitTopics(emitMsgJSON, "zombiebattleground:adddeck")
	}

	return nil

}

func (z *ZombieBattleground) ListCardLibrary(ctx contract.StaticContext, req *zb.ListCardLibraryRequest) (*zb.ListCardLibraryResponse, error) {
	var cardList zb.CardList
	if err := ctx.Get(cardListKey, &cardList); err != nil {
		return nil, err
	}
	// convert to card list to card library view grouped by element
	var category = make(map[string][]*zb.Card)
	for _, card := range cardList.Cards {
		if _, ok := category[card.Element]; !ok {
			category[card.Element] = make([]*zb.Card, 0)
		}
		category[card.Element] = append(category[card.Element], card)
	}
	// order the element by name
	var elements []string
	for k := range category {
		elements = append(elements, k)
	}
	sort.Strings(elements)

	var sets []*zb.CardSet
	for _, elem := range elements {
		cards, ok := category[elem]
		if !ok {
			continue
		}
		set := &zb.CardSet{
			Name:  elem,
			Cards: cards,
		}
		sets = append(sets, set)
	}

	return &zb.ListCardLibraryResponse{Sets: sets}, nil
}

func (z *ZombieBattleground) ListHero(ctx contract.StaticContext, req *zb.ListHeroRequest) (*zb.ListHeroResponse, error) {
	var heroList zb.HeroList
	if err := ctx.Get(heroListKey, &heroList); err != nil {
		return nil, err
	}

	return &zb.ListHeroResponse{Heros: heroList.Heros}, nil
}

var Contract plugin.Contract = contract.MakePluginContract(&ZombieBattleground{})

package battleground

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/loomnetwork/go-loom/plugin"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/pkg/errors"
)

type ZombieBattleground struct {
}

func (z *ZombieBattleground) getDecks(deckSet []*zb.ZBDeck, decksToQuery []string) []*zb.ZBDeck {
	deckMap := make(map[string]*zb.ZBDeck)
	decks := make([]*zb.ZBDeck, len(decksToQuery))

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

func (z *ZombieBattleground) mergeDeckSets(deckSet1 []*zb.ZBDeck, deckSet2 []*zb.ZBDeck) []*zb.ZBDeck {
	deckMap := make(map[string]*zb.ZBDeck)

	for _, deck := range deckSet1 {
		deckMap[deck.Name] = deck
	}

	for _, deck := range deckSet2 {
		deckMap[deck.Name] = deck
	}

	newArray := make([]*zb.ZBDeck, len(deckMap))

	i := 0
	for _, deck := range deckMap {
		newArray[i] = deck
		i++
	}

	return newArray
}

func (z *ZombieBattleground) deleteDecks(deckSet []*zb.ZBDeck, decksToDelete []string) ([]*zb.ZBDeck, error) {
	deckMap := make(map[string]*zb.ZBDeck)

	for _, deck := range deckSet {
		deckMap[deck.Name] = deck
	}

	for _, deckName := range decksToDelete {
		delete(deckMap, deckName)
	}

	newArray := make([]*zb.ZBDeck, len(deckMap))

	i := 0
	for _, deck := range deckMap {
		newArray[i] = deck
		i++
	}

	if len(newArray) == 0 {
		return nil, errors.New("Cannot delete only deck available")
	}

	return newArray, nil
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

func (z *ZombieBattleground) copyAccountInfo(account *zb.ZBAccount, req *zb.UpsertAccountRequest) {
	account.PhoneNumberVerified = req.PhoneNumberVerified
	account.RewardRedeemed = req.RewardRedeemed
	account.IsKickstarter = req.IsKickstarter
	account.Image = req.Image
	account.EmailNotification = req.EmailNotification
	account.EloScore = req.EloScore
	account.CurrentTier = req.CurrentTier
	account.GameMembershipTier = req.GameMembershipTier
}

func (z *ZombieBattleground) Meta() (plugin.Meta, error) {
	return plugin.Meta{
		Name:    "ZombieBattleground",
		Version: "1.0.0",
	}, nil
}

func (z *ZombieBattleground) Init(ctx contract.Context, req *zb.InitRequest) error {
	ctx.Set(InitDataKey(), req)
	return nil
}

func (z *ZombieBattleground) GetAccount(ctx contract.StaticContext, req *zb.GetAccountRequest) (*zb.ZBAccount, error) {
	var account zb.ZBAccount
	userKeySpace := NewUserKeySpace(req.UserId)

	if err := ctx.Get(userKeySpace.AccountKey(), &account); err != nil {
		return nil, errors.Wrapf(err, "Unable to retrieve account data for userId: %s", req.UserId)
	}

	return &account, nil
}

func (z *ZombieBattleground) UpdateAccount(ctx contract.Context, req *zb.UpsertAccountRequest) (*zb.ZBAccount, error) {
	var account zb.ZBAccount

	senderAddress := []byte(ctx.Message().Sender.Local)
	userId := strings.TrimSpace(req.UserId)
	userKeySpace := NewUserKeySpace(userId)

	// Verify whether this privateKey associated with user
	if !z.isUser(ctx, userId) {
		return nil, fmt.Errorf("UserId: %s is not verified", req.UserId)
	}

	if err := ctx.Get(userKeySpace.AccountKey(), &account); err != nil {
		return nil, errors.Wrapf(err, "Unable to retrieve account data for userId: %s", req.UserId)
	}

	z.copyAccountInfo(&account, req)
	if err := ctx.Set(userKeySpace.AccountKey(), &account); err != nil {
		return nil, errors.Wrapf(err, "Error setting account information for userId: %s", req.UserId)
	}

	ctx.Logger().Info("Updated zombiebattleground account", "user_id", userId, "address", senderAddress)

	emitMsgJSON, err := z.prepareEmitMsgJSON(senderAddress, userId, "updateaccount")
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("Error marshalling emit message for userId:%s. Error:%s", req.UserId, err))
	} else {
		ctx.EmitTopics(emitMsgJSON, "zombiebattleground:updateaccount")
	}

	return &account, nil
}

func (z *ZombieBattleground) CreateAccount(ctx contract.Context, req *zb.UpsertAccountRequest) error {
	var account zb.ZBAccount

	userId := strings.TrimSpace(req.UserId)
	senderAddress := []byte(ctx.Message().Sender.Local)
	userKeySpace := NewUserKeySpace(userId)

	var initData zb.InitRequest
	if err := ctx.Get(InitDataKey(), &initData); err != nil {
		return errors.Wrapf(err, "Unable to retrieve initdata.")
	}

	// confirm owner doesnt exist already
	if ctx.Has(userKeySpace.AccountKey()) {
		return errors.New("User already exists.\n")
	}

	account.UserId = userId
	account.Owner = ctx.Message().Sender.Bytes()

	z.copyAccountInfo(&account, req)

	if err := ctx.Set(userKeySpace.AccountKey(), &account); err != nil {
		return errors.Wrapf(err, "Error setting account information for userId: %s", req.UserId)
	}

	ctx.GrantPermission([]byte(userId), []string{"user"})

	ctx.Set(userKeySpace.DecksKey(), &zb.UserDecks{Decks: initData.DefaultDecks})

	ctx.Logger().Info("Created zombiebattleground account", "userId", userId, "address", senderAddress)

	emitMsgJSON, err := z.prepareEmitMsgJSON(senderAddress, userId, "createaccount")
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("Error marshalling emit message for userId:%s. Error:%s", req.UserId, err))
	} else {
		ctx.EmitTopics(emitMsgJSON, "zombiebattleground:createaccount")
	}

	return nil
}

// Deck related functions
func (z *ZombieBattleground) GetDecks(ctx contract.StaticContext, req *zb.GetDecksRequest) (*zb.UserDecks, error) {
	var userDecks zb.UserDecks

	userId := strings.TrimSpace(req.UserId)
	userKeySpace := NewUserKeySpace(userId)

	if err := ctx.Get(userKeySpace.DecksKey(), &userDecks); err != nil {
		return nil, errors.Wrapf(err, "Unable to get decks for userId: %s", userId)
	}

	return &userDecks, nil
}

func (z *ZombieBattleground) GetDeck(ctx contract.StaticContext, req *zb.GetDeckRequest) (*zb.ZBDeck, error) {
	var userDecks zb.UserDecks

	userId := strings.TrimSpace(req.UserId)
	userKeySpace := NewUserKeySpace(userId)
	deckId := strings.TrimSpace(req.DeckId)

	if err := ctx.Get(userKeySpace.DecksKey(), &userDecks); err != nil {
		return nil, errors.Wrapf(err, "Unable to get decks for userId: %s", userId)
	}

	result := z.getDecks(userDecks.Decks, []string{deckId})

	if len(result) != 0 {
		return result[0], nil
	} else {
		return nil, fmt.Errorf("Deck: %s not found.", deckId)
	}
}

func (z *ZombieBattleground) DeleteDeck(ctx contract.Context, req *zb.DeleteDeckRequest) error {
	var userDecks zb.UserDecks
	var err error

	userId := strings.TrimSpace(req.UserId)
	senderAddress := []byte(ctx.Message().Sender.Local)
	userKeySpace := NewUserKeySpace(userId)
	deckId := strings.TrimSpace(req.DeckId)

	if err = ctx.Get(userKeySpace.DecksKey(), &userDecks); err != nil {
		return errors.Wrapf(err, "Unable to get decks for userId: %s", userId)
	}

	if userDecks.Decks, err = z.deleteDecks(userDecks.Decks, []string{deckId}); err != nil {
		return errors.Wrapf(err, "Unable to delete deck: %s", deckId)
	}

	if err = ctx.Set(userKeySpace.DecksKey(), &userDecks); err != nil {
		return errors.Wrapf(err, "Unable to save decks for userId: %s", userId)
	}

	// TODO: Only emit events and log in case something actually got deleted
	ctx.Logger().Info("Deleted zombiebattleground deck", "userId", userId, "deckId", deckId, "address", senderAddress)

	emitMsgJSON, err := z.prepareEmitMsgJSON(senderAddress, userId, "deletedeck")
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("Error marshalling emit message for userId:%s. Error:%s", req.UserId, err))
	} else {
		ctx.EmitTopics(emitMsgJSON, "zombiebattleground:deletedeck")
	}

	return nil
}

var Contract plugin.Contract = contract.MakePluginContract(&ZombieBattleground{})

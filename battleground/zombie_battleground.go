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

func (z *ZombieBattleground) isUser(ctx contract.Context, userId string) bool {
	ok, _ := ctx.HasPermission([]byte(userId), []string{"user"})
	return ok
}

func (z *ZombieBattleground) constructUserKey(userId string) []byte {
	return []byte("user:" + userId)
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

func (z *ZombieBattleground) Meta() (plugin.Meta, error) {
	return plugin.Meta{
		Name:    "ZombieBattleground",
		Version: "1.0.0",
	}, nil
}

func (z *ZombieBattleground) Init(ctx contract.Context, req *plugin.Request) error {
	return nil
}

func (z *ZombieBattleground) GetAccount(ctx contract.StaticContext, req *zb.GetAccountRequest) (*zb.Account, error) {
	var account zb.Account

	if err := ctx.Get(z.constructUserKey(strings.TrimSpace(req.UserId)), &account); err != nil {
		return nil, errors.Wrapf(err, "Unable to retrieve account data for userId: %s", req.UserId)
	}

	return &account, nil
}

func (z *ZombieBattleground) UpdateAccount(ctx contract.Context, req *zb.UpsertAccountRequest) (*zb.Account, error) {
	var account zb.Account
	senderAddress := []byte(ctx.Message().Sender.Local)
	userId := strings.TrimSpace(req.UserId)

	// Verify whether this privateKey associated with user
	if !z.isUser(ctx, userId) {
		return nil, fmt.Errorf("UserId: %s is not verified", req.UserId)
	}

	if err := ctx.Get(z.constructUserKey(userId), &account); err != nil {
		return nil, errors.Wrapf(err, "Unable to retrieve account data for userId: %s", req.UserId)
	}

	z.copyAccountInfo(&account, req)
	if err := ctx.Set(z.constructUserKey(userId), &account); err != nil {
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
	var account zb.Account

	userId := strings.TrimSpace(req.UserId)
	senderAddress := []byte(ctx.Message().Sender.Local)

	// confirm owner doesnt exist already
	if ctx.Has(z.constructUserKey(userId)) {
		return errors.New("User already exists.\n")
	}

	account.UserId = userId
	account.Owner = ctx.Message().Sender.MarshalPB()

	z.copyAccountInfo(&account, req)

	if err := ctx.Set(z.constructUserKey(userId), &account); err != nil {
		return errors.Wrapf(err, "Error setting account information for userId: %s", req.UserId)
	}

	ctx.GrantPermission([]byte(userId), []string{"user"})

	ctx.Logger().Info("Created zombiebattleground account", "userId", userId, "address", senderAddress)

	emitMsgJSON, err := z.prepareEmitMsgJSON(senderAddress, userId, "createaccount")
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("Error marshalling emit message for userId:%s. Error:%s", req.UserId, err))
	} else {
		ctx.EmitTopics(emitMsgJSON, "zombiebattleground:createaccount")
	}

	return nil
}

var Contract plugin.Contract = contract.MakePluginContract(&ZombieBattleground{})

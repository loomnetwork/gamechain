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

func (z *ZombieBattleground) isOwner(ctx contract.Context, username string) bool {
	ok, _ := ctx.HasPermission([]byte(username), []string{"owner"})
	return ok
}

func (z *ZombieBattleground) constructOwnerKey(owner string) []byte {
	return []byte("owner:" + owner)
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
	owner := strings.TrimSpace(req.Username)

	if err := ctx.Get(z.constructOwnerKey(owner), &account); err != nil {
		return nil, errors.Wrapf(err, "Unable to retrieve account data for username: %s", req.Username)
	}

	return &account, nil
}

func (z *ZombieBattleground) UpdateAccount(ctx contract.Context, req *zb.UpsertAccountRequest) (*zb.Account, error) {
	var account zb.Account
	senderAddress := []byte(ctx.Message().Sender.Local)

	// Verify whether user is actually an owner or not
	if !z.isOwner(ctx, strings.TrimSpace(req.Username)) {
		return nil, fmt.Errorf("Username: %s is not verified", req.Username)
	}

	owner := strings.TrimSpace(req.Username)

	if err := ctx.Get(z.constructOwnerKey(owner), &account); err != nil {
		return nil, errors.Wrapf(err, "Unable to retrieve account data for username: %s", req.Username)
	}

	z.copyAccountInfo(&account, req)
	if err := ctx.Set(z.constructOwnerKey(owner), &account); err != nil {
		return nil, errors.Wrapf(err, "Error setting account information for user: %s", req.Username)
	}

	ctx.Logger().Info("Updated zombiebattleground account", "owner", owner, "address", senderAddress)

	emitMsgJSON, err := z.prepareEmitMsgJSON(senderAddress, owner, "updateaccount")
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("Error marshalling emit message for username:%s. Error:%s", req.Username, err))
	} else {
		ctx.EmitTopics(emitMsgJSON, "zombiebattleground:updateaccount")
	}

	return &account, nil
}

func (z *ZombieBattleground) CreateAccount(ctx contract.Context, req *zb.UpsertAccountRequest) error {
	var account zb.Account
	var uuid string

	senderAddress := []byte(ctx.Message().Sender.Local)

	// confirm owner doesnt exist already
	if ctx.Has(z.constructOwnerKey(strings.TrimSpace(req.Username))) {
		return errors.New("Owner already exists.\n")
	}

	owner := strings.TrimSpace(req.Username)

	uuid, err := generateUUID()
	if err != nil {
		return errors.Wrapf(err, "Unable to generate Account Id for user: %s", req.Username)
	}

	account.Id = uuid
	account.Owner = ctx.Message().Sender.MarshalPB()
	account.Username = req.Username

	z.copyAccountInfo(&account, req)

	if err := ctx.Set(z.constructOwnerKey(owner), &account); err != nil {
		return errors.Wrapf(err, "Error setting account information for user: %s", req.Username)
	}

	ctx.GrantPermission([]byte(owner), []string{"owner"})

	ctx.Logger().Info("Created zombiebattleground account", "owner", owner, "address", senderAddress)

	emitMsgJSON, err := z.prepareEmitMsgJSON(senderAddress, owner, "createaccount")
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("Error marshalling emit message for username:%s. Error:%s", req.Username, err))
	} else {
		ctx.EmitTopics(emitMsgJSON, "zombiebattleground:createaccount")
	}

	return nil
}

var Contract plugin.Contract = contract.MakePluginContract(&ZombieBattleground{})

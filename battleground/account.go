package battleground

import (
	"fmt"
	"strings"

	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/pkg/errors"
)

func copyAccountInfo(account *zb.Account, req *zb.UpsertAccountRequest) {
	account.PhoneNumberVerified = req.PhoneNumberVerified
	account.RewardRedeemed = req.RewardRedeemed
	account.IsKickstarter = req.IsKickstarter
	account.Image = req.Image
	account.EmailNotification = req.EmailNotification
	account.EloScore = req.EloScore
	account.CurrentTier = req.CurrentTier
	account.GameMembershipTier = req.GameMembershipTier
}

func GetAccount(ctx contract.StaticContext, req *zb.GetAccountRequest) (*zb.Account, error) {
	var account zb.Account
	owner := strings.TrimSpace(req.Username)

	if err := ctx.Get(OwnerKey(owner), &account); err != nil {
		return nil, errors.Wrapf(err, "Unable to retrieve account data for username: %s", req.Username)
	}

	return &account, nil
}

func UpdateAccount(ctx contract.Context, req *zb.UpsertAccountRequest) (*zb.Account, error) {
	var account zb.Account
	senderAddress := []byte(ctx.Message().Sender.Local)

	// Verify whether user is actually an owner or not
	if !isOwner(ctx, strings.TrimSpace(req.Username)) {
		return nil, fmt.Errorf("Username: %s is not verified", req.Username)
	}

	owner := strings.TrimSpace(req.Username)

	if err := ctx.Get(OwnerKey(owner), &account); err != nil {
		return nil, errors.Wrapf(err, "Unable to retrieve account data for username: %s", req.Username)
	}

	copyAccountInfo(&account, req)
	if err := ctx.Set(OwnerKey(owner), &account); err != nil {
		return nil, errors.Wrapf(err, "Error setting account information for user: %s", req.Username)
	}

	ctx.Logger().Info("Updated zombiebattleground account", "owner", owner, "address", senderAddress)

	emitMsgJSON, err := prepareEmitMsgJSON(senderAddress, owner, "updateaccount")
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("Error marshalling emit message for username:%s. Error:%s", req.Username, err))
	} else {
		ctx.EmitTopics(emitMsgJSON, "zombiebattleground:updateaccount")
	}

	return &account, nil
}

func CreateAccount(ctx contract.Context, req *zb.UpsertAccountRequest) error {
	var account zb.Account
	var uuid string

	senderAddress := []byte(ctx.Message().Sender.Local)

	// confirm owner doesnt exist already
	if ctx.Has(OwnerKey(strings.TrimSpace(req.Username))) {
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

	copyAccountInfo(&account, req)

	if err := ctx.Set(OwnerKey(owner), &account); err != nil {
		return errors.Wrapf(err, "Error setting account information for user: %s", req.Username)
	}

	ctx.GrantPermission([]byte(owner), []string{"owner"})

	ctx.Logger().Info("Created zombiebattleground account", "owner", owner, "address", senderAddress)

	emitMsgJSON, err := prepareEmitMsgJSON(senderAddress, owner, "createaccount")
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("Error marshalling emit message for username:%s. Error:%s", req.Username, err))
	} else {
		ctx.EmitTopics(emitMsgJSON, "zombiebattleground:createaccount")
	}

	return nil
}

package battleground

import (
	"fmt"
	orctype "github.com/loomnetwork/gamechain/types/oracle"
	"github.com/loomnetwork/go-loom"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	"github.com/pkg/errors"
)

func (z *ZombieBattleground) processOracleCommandResponseBatchInternal(ctx contract.Context, commandResponses []*orctype.OracleCommandResponse) error {
	commandRequestList, err := loadOracleCommandRequestList(ctx)
	if err != nil {
		return errors.Wrap(err, "processOracleCommandResponseBatchInternal")
	}

	newCommandRequests := make([]*orctype.OracleCommandRequest, 0)
	for _, commandResponseOneOf := range commandResponses {
		// check if a command with such id even exists
		var matchingCommandRequest *orctype.OracleCommandRequest = nil
		for _, existingCommandRequest := range commandRequestList.Commands {
			if existingCommandRequest.CommandId == commandResponseOneOf.CommandId {
				matchingCommandRequest = existingCommandRequest
				break
			}
		}

		if matchingCommandRequest == nil {
			ctx.Logger().Warn(fmt.Errorf("unknown oracle command with id %d", commandResponseOneOf.CommandId).Error())
			continue
		}

		err = nil
		switch commandResponse := commandResponseOneOf.Command.(type) {
		case *orctype.OracleCommandResponse_GetUserFullCardCollection:
			err = z.processOracleCommandResponseGetUserFullCardCollection(
				ctx,
				loom.UnmarshalAddressPB(commandResponse.GetUserFullCardCollection.UserAddress),
				commandResponse.GetUserFullCardCollection.OwnedCards,
				commandResponse.GetUserFullCardCollection.BlockHeight,
			)
			break
		}

		// We allow single commands to fail
		// Just keep them unconfirmed until something is fixed for it to become processed and confirmed.
		if err != nil {
			newCommandRequests = append(newCommandRequests, matchingCommandRequest)
			ctx.Logger().Error(errors.Wrap(err, "error processing oracle command response").Error())
		}
	}

	commandRequestList.Commands = newCommandRequests
	err = saveOracleCommandRequestList(ctx, commandRequestList)
	if err != nil {
		return errors.Wrap(err, "processOracleCommandResponseBatchInternal")
	}

	return nil
}

func (z *ZombieBattleground) saveOracleCommandRequestToList(
	ctx contract.Context,
	commandRequest *orctype.OracleCommandRequest,
	filterSameTypePredicate func(request *orctype.OracleCommandRequest) (mustRemove bool),
) error {
	// load data
	commandRequestList, err := loadOracleCommandRequestList(ctx)
	if err != nil {
		return err
	}

	contractState, err := loadContractState(ctx)
	if err != nil {
		return err
	}

	// filter our commands. Useful for commands that must be unique in list
	filteredCommandsRequests := make([]*orctype.OracleCommandRequest, 0)
	for _, commandRequest := range commandRequestList.Commands {
		if !filterSameTypePredicate(commandRequest) {
			filteredCommandsRequests = append(filteredCommandsRequests, commandRequest)
		}
	}

	filteredCommandsRequests = append(filteredCommandsRequests, commandRequest)

	//save data
	contractState.CurrentOracleCommandId++
	err = saveContractState(ctx, contractState)
	if err != nil {
		return err
	}

	commandRequestList.Commands = filteredCommandsRequests
	err = saveOracleCommandRequestList(ctx, commandRequestList)
	if err != nil {
		return err
	}

	return nil
}

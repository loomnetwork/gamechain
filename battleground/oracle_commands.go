package battleground

import (
	"fmt"
	orctype "github.com/loomnetwork/gamechain/types/oracle"
	"github.com/loomnetwork/go-loom"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	"github.com/pkg/errors"
)

func createBaseOracleCommand(ctx contract.StaticContext) (*orctype.OracleCommandRequest, error) {
	contractState, err := loadContractState(ctx)
	if err != nil {
		return nil, err
	}

	return &orctype.OracleCommandRequest{
		CommandId: contractState.CurrentOracleCommandId,
	}, nil
}

func (z *ZombieBattleground) addGetUserFullCardCollectionOracleCommand(ctx contract.Context, address loom.Address) error {
	configuration, err := loadContractConfiguration(ctx)
	if err != nil {
		return errors.Wrap(err, "addGetUserFullCardCollectionOracleCommand")
	}

	if configuration.CardCollectionSyncDataVersion == "" {
		return errors.New("configuration.CardCollectionSyncDataVersion is not set")
	}

	command, err := createBaseOracleCommand(ctx)
	if err != nil {
		return errors.Wrap(err, "addGetUserFullCardCollectionOracleCommand")
	}

	command.Command = &orctype.OracleCommandRequest_GetUserFullCardCollection{
		GetUserFullCardCollection: &orctype.OracleCommandRequest_GetUserFullCardCollectionCommandRequest{
			UserAddress: address.MarshalPB(),
		},
	}

	err = z.saveOracleCommandRequestToList(ctx, command, func(request *orctype.OracleCommandRequest) (mustRemove bool) {
		// Remove any other request for the same address
		return request.GetGetUserFullCardCollection() != nil && loom.UnmarshalAddressPB(request.GetGetUserFullCardCollection().UserAddress).Compare(address) == 0
	})

	if err != nil {
		return errors.Wrap(err, "addGetUserFullCardCollectionOracleCommand")
	}

	return nil
}

func (z *ZombieBattleground) processOracleCommandResponseBatchInternal(ctx contract.Context, commandResponses []*orctype.OracleCommandResponse) error {
	commandRequestList, err := loadOracleCommandRequestList(ctx)
	if err != nil {
		return errors.Wrap(err, "processOracleCommandResponseBatchInternal")
	}

	unhandledCommandRequests := make([]*orctype.OracleCommandRequest, 0)
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
			if err != nil {
				err = errors.Wrap(err, "processOracleCommandResponseGetUserFullCardCollection failed")
			}
			break
		}

		// We allow single commands to fail
		// Just keep them unconfirmed until something is fixed for it to become processed and confirmed.
		if err != nil {
			unhandledCommandRequests = append(unhandledCommandRequests, matchingCommandRequest)
			ctx.Logger().Error(errors.Wrap(err, "error processing oracle command response").Error())
		}
	}

	commandRequestList.Commands = unhandledCommandRequests
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

	// filter out commands. Useful for commands that must be unique in list
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

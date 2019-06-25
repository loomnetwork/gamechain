package battleground

import (
	"github.com/loomnetwork/gamechain/types/oracle"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
)

func (z *ZombieBattleground) saveOracleCommandRequest(
	ctx contract.Context,
	commandRequest *oracle.OracleCommandRequest,
	filterSameTypePredicate func(request *oracle.OracleCommandRequest) (remove bool),
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
	filteredCommandsRequests := make([]*oracle.OracleCommandRequest, 0)
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

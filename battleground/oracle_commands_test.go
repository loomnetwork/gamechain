package battleground

import (
	"github.com/loomnetwork/gamechain/types/oracle"
	"github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/common"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	"github.com/loomnetwork/go-loom/types"
	assert "github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestProcessEventBatch(t *testing.T) {
	c := &ZombieBattleground{}
	var pubKeyHexString = "3866f776276246e4f9998aa90632931d89b0d3a5930e804e02299533f55b39e1"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)

	req := &oracle.ProcessOracleEventBatchRequest{
		LastPlasmachainBlockNumber: 123,
		Events: []*oracle.PlasmachainEvent{
			{
				EthBlock: 120,
				Payload: &oracle.PlasmachainEvent_TransferWithQuantity{
					TransferWithQuantity: &oracle.PlasmachainEventTransferWithQuantity{
						Amount: &types.BigUInt{Value: common.BigUInt{Int: big.NewInt(3)}},
						TokenId: &types.BigUInt{Value: common.BigUInt{Int: big.NewInt(100)}},
						From: loom.MustParseAddress("default:0x0000000000000000000000000000000000000001").MarshalPB(),
						To: loom.MustParseAddress("default:0x0000000000000000000000000000000000000002").MarshalPB(),
					},
				},
			},
			{
				EthBlock: 120,
				Payload: &oracle.PlasmachainEvent_TransferWithQuantity{
					TransferWithQuantity: &oracle.PlasmachainEventTransferWithQuantity{
						Amount: &types.BigUInt{Value: common.BigUInt{Int: big.NewInt(7)}},
						TokenId: &types.BigUInt{Value: common.BigUInt{Int: big.NewInt(100)}},
						From: loom.MustParseAddress("default:0x0000000000000000000000000000000000000001").MarshalPB(),
						To: loom.MustParseAddress("default:0x0000000000000000000000000000000000000002").MarshalPB(),
					},
				},
			},
			{
				EthBlock: 120,
				Payload: &oracle.PlasmachainEvent_TransferWithQuantity{
					TransferWithQuantity: &oracle.PlasmachainEventTransferWithQuantity{
						Amount: &types.BigUInt{Value: common.BigUInt{Int: big.NewInt(2)}},
						TokenId: &types.BigUInt{Value: common.BigUInt{Int: big.NewInt(100)}},
						From: loom.MustParseAddress("default:0x0000000000000000000000000000000000000002").MarshalPB(),
						To: loom.MustParseAddress("default:0x0000000000000000000000000000000000000001").MarshalPB(),
					},
				},
			},
		},
	}
	err := c.ProcessOracleEventBatch(ctx, req)
	assert.Nil(t, err)
}

func TestCreateOracleCommandRequest(t *testing.T) {
	c := &ZombieBattleground{}
	var pubKeyHexString = "3866f776276246e4f9998aa90632931d89b0d3a5930e804e02299533f55b39e1"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)

	t.Run("No commands initially", func(t *testing.T) {
		commandRequestList, err := loadOracleCommandRequestList(ctx)
		assert.Nil(t, err)
		assert.NotNil(t, commandRequestList)
		assert.NotNil(t, commandRequestList.Commands)
		assert.Equal(t, 0, len(commandRequestList.Commands))

		state, err := loadContractState(ctx)
		assert.Nil(t, err)
		assert.NotNil(t, state)
		assert.Equal(t, uint64(0), state.CurrentOracleCommandId)
	})

	t.Run("Add basic command", func(t *testing.T) {
		command := &oracle.OracleCommandRequest{}
		err := c.saveOracleCommandRequest(ctx, command, func(request *oracle.OracleCommandRequest) (remove bool) {
			return false
		})

		assert.Nil(t, err)

		commandRequestList, err := loadOracleCommandRequestList(ctx)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(commandRequestList.Commands))

		state, err := loadContractState(ctx)
		assert.Nil(t, err)
		assert.NotNil(t, state)
		assert.Equal(t, uint64(1), state.CurrentOracleCommandId)
	})

	t.Run("Add command with data", func(t *testing.T) {
		command := &oracle.OracleCommandRequest{
			Command: &oracle.OracleCommandRequest_GetUserFullCardCollection{
				GetUserFullCardCollection:	&oracle.OracleCommandRequest_GetUserFullCardCollectionCommandRequest {
					UserId: &types.BigUInt{Value: common.BigUInt{Int: big.NewInt(373)}},
				},
			},
		}
		err := c.saveOracleCommandRequest(ctx, command, func(request *oracle.OracleCommandRequest) (remove bool) {
			return false
		})

		assert.Nil(t, err)

		commandRequestList, err := loadOracleCommandRequestList(ctx)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(commandRequestList.Commands))
		assert.NotNil(t, commandRequestList.Commands[1].GetGetUserFullCardCollection())
		assert.Equal(t, int64(373), commandRequestList.Commands[1].GetGetUserFullCardCollection().UserId.Value.Int64())

		state, err := loadContractState(ctx)
		assert.Nil(t, err)
		assert.NotNil(t, state)
		assert.Equal(t, uint64(2), state.CurrentOracleCommandId)
	})

	t.Run("Add unique command", func(t *testing.T) {
		command := &oracle.OracleCommandRequest{
			Command: &oracle.OracleCommandRequest_GetUserFullCardCollection{
				GetUserFullCardCollection:	&oracle.OracleCommandRequest_GetUserFullCardCollectionCommandRequest {
					UserId: &types.BigUInt{Value: common.BigUInt{Int: big.NewInt(374)}},
				},
			},
		}
		err := c.saveOracleCommandRequest(ctx, command, func(request *oracle.OracleCommandRequest) (remove bool) {
			return request.GetGetUserFullCardCollection() != nil
		})

		assert.Nil(t, err)

		commandRequestList, err := loadOracleCommandRequestList(ctx)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(commandRequestList.Commands))
		assert.NotNil(t, commandRequestList.Commands[1].GetGetUserFullCardCollection())
		assert.Equal(t, int64(374), commandRequestList.Commands[1].GetGetUserFullCardCollection().UserId.Value.Int64())

		state, err := loadContractState(ctx)
		assert.Nil(t, err)
		assert.NotNil(t, state)
		assert.Equal(t, uint64(3), state.CurrentOracleCommandId)
	})
}
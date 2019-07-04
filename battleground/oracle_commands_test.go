package battleground

import (
	"github.com/loomnetwork/gamechain/tools/battleground_utility"
	orctype "github.com/loomnetwork/gamechain/types/oracle"
	"github.com/loomnetwork/go-loom"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
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

	req := &orctype.ProcessOracleEventBatchRequest{
		LastPlasmachainBlockNumber: 123,
		ZbgCardContractAddress: loom.MustParseAddress("default:0x0000000000000000000000000000000000000099").MarshalPB(),
		Events: []*orctype.PlasmachainEvent{
			{
				EthBlock: 120,
				Payload: &orctype.PlasmachainEvent_TransferWithQuantity{
					TransferWithQuantity: &orctype.PlasmachainEventTransferWithQuantity{
						Amount:  battleground_utility.MarshalBigIntProto(big.NewInt(3)),
						TokenId: battleground_utility.MarshalBigIntProto(big.NewInt(100)),
						From:    loom.MustParseAddress("default:0x0000000000000000000000000000000000000001").MarshalPB(),
						To:      loom.MustParseAddress("default:0x0000000000000000000000000000000000000002").MarshalPB(),
					},
				},
			},
			{
				EthBlock: 120,
				Payload: &orctype.PlasmachainEvent_TransferWithQuantity{
					TransferWithQuantity: &orctype.PlasmachainEventTransferWithQuantity{
						Amount:  battleground_utility.MarshalBigIntProto(big.NewInt(7)),
						TokenId: battleground_utility.MarshalBigIntProto(big.NewInt(100)),
						From:    loom.MustParseAddress("default:0x0000000000000000000000000000000000000001").MarshalPB(),
						To:      loom.MustParseAddress("default:0x0000000000000000000000000000000000000002").MarshalPB(),
					},
				},
			},
			{
				EthBlock: 120,
				Payload: &orctype.PlasmachainEvent_TransferWithQuantity{
					TransferWithQuantity: &orctype.PlasmachainEventTransferWithQuantity{
						Amount:  battleground_utility.MarshalBigIntProto(big.NewInt(2)),
						TokenId: battleground_utility.MarshalBigIntProto(big.NewInt(100)),
						From:    loom.MustParseAddress("default:0x0000000000000000000000000000000000000002").MarshalPB(),
						To:      loom.MustParseAddress("default:0x0000000000000000000000000000000000000001").MarshalPB(),
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
	user1Address := loom.MustParseAddress("default:0x0000000000000000000000000000000000000003")

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
		command, err := createBaseOracleCommand(ctx)
		assert.Nil(t, err)
		err = c.saveOracleCommandRequestToList(ctx, command, func(request *orctype.OracleCommandRequest) (mustRemove bool) {
			return false
		})

		assert.Nil(t, err)

		commandRequestList, err := loadOracleCommandRequestList(ctx)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(commandRequestList.Commands))
		assert.Equal(t, uint64(0), commandRequestList.Commands[0].CommandId)

		state, err := loadContractState(ctx)
		assert.Nil(t, err)
		assert.NotNil(t, state)
		assert.Equal(t, uint64(1), state.CurrentOracleCommandId)
	})

	t.Run("Add command with data", func(t *testing.T) {
		command, err := createBaseOracleCommand(ctx)
		assert.Nil(t, err)
		command.Command = &orctype.OracleCommandRequest_GetUserFullCardCollection{
			GetUserFullCardCollection: &orctype.OracleCommandRequest_GetUserFullCardCollectionCommandRequest{
				UserAddress: user1Address.MarshalPB(),
			},
		}

		err = c.saveOracleCommandRequestToList(ctx, command, func(request *orctype.OracleCommandRequest) (mustRemove bool) {
			return false
		})

		assert.Nil(t, err)

		commandRequestList, err := loadOracleCommandRequestList(ctx)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(commandRequestList.Commands))
		assert.NotNil(t, commandRequestList.Commands[1].GetGetUserFullCardCollection())
		assert.Equal(t, user1Address.String(),loom.UnmarshalAddressPB(commandRequestList.Commands[1].GetGetUserFullCardCollection().UserAddress).String())
		assert.Equal(t, uint64(1), commandRequestList.Commands[1].CommandId)

		state, err := loadContractState(ctx)
		assert.Nil(t, err)
		assert.NotNil(t, state)
		assert.Equal(t, uint64(2), state.CurrentOracleCommandId)
	})

	t.Run("Add unique type per-user command", func(t *testing.T) {
		command, err := createBaseOracleCommand(ctx)
		assert.Nil(t, err)
		command.Command = &orctype.OracleCommandRequest_GetUserFullCardCollection{
			GetUserFullCardCollection: &orctype.OracleCommandRequest_GetUserFullCardCollectionCommandRequest{
				UserAddress: user1Address.MarshalPB(),
			},
		}

		err = c.saveOracleCommandRequestToList(ctx, command, func(request *orctype.OracleCommandRequest) (mustRemove bool) {
			return request.GetGetUserFullCardCollection() != nil && loom.UnmarshalAddressPB(request.GetGetUserFullCardCollection().UserAddress).Compare(user1Address) == 0
		})

		assert.Nil(t, err)

		commandRequestList, err := loadOracleCommandRequestList(ctx)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(commandRequestList.Commands))
		assert.NotNil(t, commandRequestList.Commands[1].GetGetUserFullCardCollection())
		assert.Equal(t, user1Address.String(),loom.UnmarshalAddressPB(commandRequestList.Commands[1].GetGetUserFullCardCollection().UserAddress).String())

		state, err := loadContractState(ctx)
		assert.Nil(t, err)
		assert.NotNil(t, state)
		assert.Equal(t, uint64(3), state.CurrentOracleCommandId)
	})
}

func TestGetUserFullCardCollectionCommandRequest(t *testing.T) {
	c := &ZombieBattleground{}
	var pubKeyHexString = "3866f776276246e4f9998aa90632931d89b0d3a5930e804e02299533f55b39e1"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)
	user1Address := loom.MustParseAddress("default:0x0000000000000000000000000000000000000003")
	user2Address := loom.MustParseAddress("default:0x0000000000000000000000000000000000000004")

	t.Run("GetUserFullCardCollectionCommandRequest works", func(t *testing.T) {
		err := c.addGetUserFullCardCollectionOracleCommand(ctx, user1Address)
		assert.Nil(t, err)
		err = c.addGetUserFullCardCollectionOracleCommand(ctx, user2Address)
		assert.Nil(t, err)

		commandRequestList, err := loadOracleCommandRequestList(ctx)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(commandRequestList.Commands))
		assert.Equal(t, uint64(0), commandRequestList.Commands[0].CommandId)
		assert.Equal(t, uint64(1), commandRequestList.Commands[1].CommandId)
		assert.Equal(t, user1Address.String(), loom.UnmarshalAddressPB(commandRequestList.Commands[0].GetGetUserFullCardCollection().UserAddress).String())
		assert.Equal(t, user2Address.String(), loom.UnmarshalAddressPB(commandRequestList.Commands[1].GetGetUserFullCardCollection().UserAddress).String())

		state, err := loadContractState(ctx)
		assert.Nil(t, err)
		assert.NotNil(t, state)
		assert.Equal(t, uint64(2), state.CurrentOracleCommandId)
	})

	t.Run("GetUserFullCardCollectionCommandRequest is unique per-user", func(t *testing.T) {
		err := c.addGetUserFullCardCollectionOracleCommand(ctx, user1Address)
		assert.Nil(t, err)

		commandRequestList, err := loadOracleCommandRequestList(ctx)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(commandRequestList.Commands))
		assert.Equal(t, uint64(1), commandRequestList.Commands[0].CommandId)
		assert.Equal(t, uint64(2), commandRequestList.Commands[1].CommandId)
		assert.Equal(t, user2Address.String(), loom.UnmarshalAddressPB(commandRequestList.Commands[0].GetGetUserFullCardCollection().UserAddress).String())
		assert.Equal(t, user1Address.String(), loom.UnmarshalAddressPB(commandRequestList.Commands[1].GetGetUserFullCardCollection().UserAddress).String())

		state, err := loadContractState(ctx)
		assert.Nil(t, err)
		assert.NotNil(t, state)
		assert.Equal(t, uint64(3), state.CurrentOracleCommandId)
	})
}

package battleground

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/gamechain/types/common"
	"github.com/loomnetwork/go-loom"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
)

// We will create one instance of this per deployed game mode
type CustomGameMode struct {
	// Address of game mode contract deployed to Loom EVM.
	tokenAddr   loom.Address
	contractABI *abi.ABI
}

func (c *CustomGameMode) UpdateInitialPlayerGameState(ctx contract.Context, gameState *zb.GameState) error {
	serializedGameState, err := c.serializeGameState(gameState)
	if err != nil {
		return err
	}

	serializedGameStateChangeActions, err := c.callOnMatchStarting(ctx, serializedGameState)
	if err != nil {
		return err
	}

	err = c.deserializeAndApplyGameStateChangeActions(gameState, serializedGameStateChangeActions)
	if err != nil {
		return err
	}

	return nil
}

func NewCustomGameMode(tokenAddr loom.Address) *CustomGameMode {
	erc20ABI, err := abi.JSON(strings.NewReader(zbGameModeABI))
	if err != nil {
		panic(err)
	}
	return &CustomGameMode{
		tokenAddr:   tokenAddr,
		contractABI: &erc20ABI,
	}
}

func (c *CustomGameMode) staticCallEVM(ctx contract.Context, method string, result interface{}, params ...interface{}) error {
	input, err := c.contractABI.Pack(method, params...)
	if err != nil {
		return err
	}
	var output []byte
	if err := contract.StaticCallEVM(ctx, c.tokenAddr, input, &output); err != nil {
		return err
	}
	return c.contractABI.Unpack(result, method, output)
}

func (c *CustomGameMode) callEVM(ctx contract.Context, method string, params ...interface{}) ([]byte, error) {
	input, err := c.contractABI.Pack(method, params...)
	if err != nil {
		return nil, err
	}
	var evmOut []byte
	return evmOut, contract.CallEVM(ctx, c.tokenAddr, input, &evmOut)
}

func (c *CustomGameMode) callOnMatchStarting(ctx contract.Context, serializedGameState []byte) (serializedGameStateChangeActions []byte, err error) {
	ctx.Logger().Info(fmt.Sprintf("serializedGameState----------------%v\n", serializedGameState))
	if err := c.staticCallEVM(ctx, "onMatchStarting", &serializedGameStateChangeActions, serializedGameState); err != nil {
		return nil, err
	}

	ctx.Logger().Info(fmt.Sprintf("serializedGameStateChangeActions----------------%v\n", serializedGameStateChangeActions))
	return serializedGameStateChangeActions, nil
}

func (c *CustomGameMode) serializeGameState(state *zb.GameState) (bytes []byte, err error) {
	rb := NewReverseBuffer(make([]byte, 256))
	if err = binary.Write(rb, binary.BigEndian, int64(state.Id)); err != nil {
		return nil, err
	}
	if err = binary.Write(rb, binary.BigEndian, byte(state.CurrentPlayerIndex)); err != nil {
		return nil, err
	}
	for _, playerState := range state.PlayerStates {
		if err = binary.Write(rb, binary.BigEndian, byte(playerState.Hp)); err != nil {
			return nil, err
		}
		if err = binary.Write(rb, binary.BigEndian, byte(playerState.Mana)); err != nil {
			return nil, err
		}
	}

	return rb.GetFilledSlice(), nil
}

func (c *CustomGameMode) deserializeAndApplyGameStateChangeActions(state *zb.GameState, serializedActions []byte) (err error) {
	rb := NewReverseBuffer(serializedActions)
	for {
		var action battleground.GameStateChangeAction
		if err = binary.Read(rb, binary.BigEndian, &action); err != nil {
			return
		}

		mustBreak := false
		switch action {
		case battleground.GameStateChangeAction_None:
			mustBreak = true
		case battleground.GameStateChangeAction_SetPlayerDefense:
			var playerIndex byte
			if err = binary.Read(rb, binary.BigEndian, &playerIndex); err != nil {
				return
			}

			var newDefense byte
			if err = binary.Read(rb, binary.BigEndian, &newDefense); err != nil {
				return
			}

			state.PlayerStates[playerIndex].Hp = int32(newDefense)
		case battleground.GameStateChangeAction_SetPlayerGoo:
			var playerIndex byte
			if err = binary.Read(rb, binary.BigEndian, &playerIndex); err != nil {
				return
			}

			var newGoo byte
			if err = binary.Read(rb, binary.BigEndian, &newGoo); err != nil {
				return
			}

			state.PlayerStates[playerIndex].Mana = int32(newGoo)
		default:
			return errors.New(fmt.Sprintf("Unknown game state change action %d", action))
		}

		if mustBreak {
			return nil
		}
	}
}

// From Zombiebattleground game mode repo
const zbGameModeABI = `
[
	{
		"constant": true,
		"inputs": [],
		"name": "name",
		"outputs": [
			{
				"name": "",
				"type": "string"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "gameState",
				"type": "bytes"
			}
		],
		"name": "onMatchStarting",
		"outputs": [
			{
				"name": "",
				"type": "bytes"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "getStaticConfigs",
		"outputs": [
			{
				"name": "",
				"type": "uint256[]"
			},
			{
				"name": "",
				"type": "uint256[]"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"name": "_from",
				"type": "address"
			}
		],
		"name": "MatchedStarted",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"name": "player1Addr",
				"type": "address"
			},
			{
				"indexed": false,
				"name": "player2Addr",
				"type": "address"
			},
			{
				"indexed": false,
				"name": "winner",
				"type": "uint256"
			}
		],
		"name": "MatchFinished",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"name": "to",
				"type": "address"
			},
			{
				"indexed": false,
				"name": "tokens",
				"type": "uint256"
			}
		],
		"name": "AwardTokens",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"name": "to",
				"type": "address"
			},
			{
				"indexed": false,
				"name": "cardID",
				"type": "uint256"
			}
		],
		"name": "AwardCard",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"name": "to",
				"type": "address"
			},
			{
				"indexed": false,
				"name": "packCount",
				"type": "uint256"
			},
			{
				"indexed": false,
				"name": "packType",
				"type": "uint256"
			}
		],
		"name": "AwardPack",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"name": "_from",
				"type": "address"
			}
		],
		"name": "UserRegistered",
		"type": "event"
	}
]
`

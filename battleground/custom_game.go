package battleground

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
)

// We will create one instance of this per deployed game mode
type CustomGameMode struct {
	// Address of game mode contract deployed to Loom EVM.
	tokenAddr   loom.Address
	contractABI *abi.ABI
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

func (c *CustomGameMode) UpdateInitialPlayerGameState(ctx contract.Context, gameState *zb.GameState) (err error) {
	serializedGameState, err := c.serializeGameState(gameState)
	if err != nil {
		return
	}

	serializedGameStateChangeActions, err := c.callOnMatchStarting(ctx, serializedGameState)
	if err != nil {
		return
	}

	err = c.deserializeAndApplyGameStateChangeActions(gameState, serializedGameStateChangeActions)
	if err != nil {
		return
	}

	return
}

func (c *CustomGameMode) GetCustomUi(ctx contract.StaticContext) (uiElements []*zb.CustomGameModeCustomUiElement, err error) {
	serializedCustomUi, err := c.callGetCustomUi(ctx)
	if err != nil {
		return
	}

	uiElements, err = c.deserializeCustomUi(serializedCustomUi)
	if err != nil {
		return
	}

	return uiElements, nil
}

func (c *CustomGameMode) CallFunction(ctx contract.Context, method string) (err error) {
	// crude way to call a function with no inputs and outputs without an ABI
	input := crypto.Keccak256([]byte(method + "()"))[:4]

	ctx.Logger().Info(fmt.Sprintf("methodCallAbi ----------------%v\n", input))

	var evmOut []byte
	return contract.CallEVM(ctx, c.tokenAddr, input, &evmOut)
}

func (c *CustomGameMode) staticCallEVM(ctx contract.StaticContext, method string, result interface{}, params ...interface{}) error {
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
	if err := c.staticCallEVM(ctx, "onMatchStarting", &serializedGameStateChangeActions, serializedGameState); err != nil {
		return nil, err
	}

	return serializedGameStateChangeActions, nil
}

func (c *CustomGameMode) callGetCustomUi(ctx contract.StaticContext) (serializedCustomUi []byte, err error) {
	if err := c.staticCallEVM(ctx, "getCustomUi", &serializedCustomUi); err != nil {
		return nil, err
	}

	return serializedCustomUi, nil
}

// From Zombiebattleground game mode repo
const zbGameModeABI = `
[
    {
      "inputs": [],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "constructor"
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
    },
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
      "constant": true,
      "inputs": [],
      "name": "GameStart",
      "outputs": [
        {
          "name": "",
          "type": "uint256"
        }
      ],
      "payable": false,
      "stateMutability": "pure",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "name": "",
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
      "stateMutability": "pure",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "getCustomUi",
      "outputs": [
        {
          "name": "",
          "type": "bytes"
        }
      ],
      "payable": false,
      "stateMutability": "pure",
      "type": "function"
    }
  ]
`

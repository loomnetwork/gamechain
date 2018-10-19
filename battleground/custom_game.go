package battleground

import (
	"errors"
	"fmt"
	"github.com/loomnetwork/go-ethereum/common/hexutil"
	"runtime"
	"runtime/debug"
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

func (c *CustomGameMode) CallOnMatchStartingBeforeInitialDraw(ctx contract.Context, gameState *zb.GameState) (err error) {
	defer func() {
		if r := recover(); r != nil {
			stackBuffer := debug.Stack()
			runtime.Stack(stackBuffer, false)

			err = errors.New(fmt.Sprintf("! Error in custom mode hook - CallOnMatchStartingBeforeInitialDraw\n%v\n%s", r, string(stackBuffer)))
		}
	}()

	return c.callAndApplyMatchHook(ctx, "onMatchStartingBeforeInitialDraw", gameState)
}

func (c *CustomGameMode) CallOnMatchStartingAfterInitialDraw(ctx contract.Context, gameState *zb.GameState) (err error) {
	defer func() {
		if r := recover(); r != nil {
			stackBuffer := debug.Stack()
			runtime.Stack(stackBuffer, false)

			err = errors.New(fmt.Sprintf("! Error in custom mode hook - CallOnMatchStartingAfterInitialDraw\n%v\n%s", r, string(stackBuffer)))
		}
	}()

	return c.callAndApplyMatchHook(ctx, "onMatchStartingAfterInitialDraw", gameState)
}

func (c *CustomGameMode) GetCustomUi(ctx contract.StaticContext) (uiElements []*zb.CustomGameModeCustomUiElement, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("! Error in custom mode hook - GetCustomUi\n", ))
			debug.PrintStack()
		}
	}()

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

func (c *CustomGameMode) CallFunction(ctx contract.Context, abiInput []byte) (err error) {
	_, e := c.callEVMRaw(ctx, abiInput)
	return e
}

func (c *CustomGameMode) callAndApplyMatchHook(ctx contract.Context, matchHookName string, gameState *zb.GameState) (err error) {
	serializedGameState, err := c.serializeGameState(gameState)
	if err != nil {
		return
	}

	fmt.Printf("------- serializedGameState %s\n%s\n%s\n--------\n", matchHookName, hexutil.Encode(serializedGameState))

	serializedGameStateChangeActions, err := c.callMatchHook(ctx, matchHookName, serializedGameState)
	if err != nil {
		fmt.Printf("-- Error in match hook %s:\n%v",  matchHookName, err)
		return
	}

	fmt.Printf("------- serializedGameStateChangeActions %s\n%s\n%s\n--------\n", matchHookName, hexutil.Encode(serializedGameStateChangeActions))

	err = c.deserializeAndApplyGameStateChangeActions(gameState, serializedGameStateChangeActions)
	if err != nil {
		return
	}

	return
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

func (c *CustomGameMode) staticCallEVMRaw(ctx contract.StaticContext, abiInput []byte) ([]byte, error) {
	var output []byte
	if err := contract.StaticCallEVM(ctx, c.tokenAddr, abiInput, &output); err != nil {
		return nil, err
	}
	return output, nil
}

func (c *CustomGameMode) callEVM(ctx contract.Context, method string, params ...interface{}) ([]byte, error) {
	input, err := c.contractABI.Pack(method, params...)
	if err != nil {
		return nil, err
	}

	return c.callEVMRaw(ctx, input)
}

func (c *CustomGameMode) callEVMRaw(ctx contract.Context, abiInput []byte) ([]byte, error) {
	var evmOut []byte
	return evmOut, contract.CallEVM(ctx, c.tokenAddr, abiInput, &evmOut)
}

func (c *CustomGameMode) callMatchHook(ctx contract.Context, matchHookName string, serializedGameState []byte) (serializedGameStateChangeActions []byte, err error) {
/*	if serializedGameStateChangeActions, err = c.callEVM(ctx, matchHookName, serializedGameState); err != nil {
		return nil, err
	}*/

	if err  = c.staticCallEVM(ctx, matchHookName, &serializedGameStateChangeActions, serializedGameState); err != nil {
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
      "inputs": [
        {
          "name": "",
          "type": "bytes"
        }
      ],
      "name": "onMatchStartingBeforeInitialDraw",
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
      "inputs": [
        {
          "name": "serializedGameState",
          "type": "bytes"
        }
      ],
      "name": "onMatchStartingAfterInitialDraw",
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
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "name": "val",
          "type": "int32"
        }
      ],
      "name": "incrementCounter",
      "outputs": [],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [],
      "name": "incrementCounter",
      "outputs": [],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    }
  ]
`

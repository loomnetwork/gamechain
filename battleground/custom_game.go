package battleground

import (
	"errors"
	"fmt"
	"github.com/loomnetwork/gamechain/types/zb/zb_data"
	"github.com/loomnetwork/go-loom/plugin/types"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/loomnetwork/go-loom"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
)

// We will create one instance of this per deployed game mode
type CustomGameMode struct {
	// Address of game mode contract deployed to Loom EVM.
	tokenAddr   loom.Address
}

var zbCustomGameModeABI *abi.ABI
var zbGameStateChangesEventTopic string

func init() {
	gameModeAbi, err := abi.JSON(strings.NewReader(zbGameModeABI))
	if err != nil {
		panic(err)
	}

	zbCustomGameModeABI = &gameModeAbi
	zbGameStateChangesEventTopic = zbCustomGameModeABI.Events[zbGameStateChangesEventName].Id().String()
}

func NewCustomGameMode(tokenAddr loom.Address) *CustomGameMode {
	return &CustomGameMode{
		tokenAddr:   tokenAddr,
	}
}

func (c *CustomGameMode) CallHookBeforeMatchStart(ctx contract.Context, gameplay *Gameplay) (err error) {
	return c.callAndApplyMatchHook(ctx, "beforeMatchStart", gameplay)
}

func (c *CustomGameMode) CallHookAfterInitialDraw(ctx contract.Context, gameplay *Gameplay) (err error) {
	return c.callAndApplyMatchHook(ctx, "afterInitialDraw", gameplay)
}

func (c *CustomGameMode) GetCustomUi(ctx contract.StaticContext) (uiElements []*zb_data.CustomGameModeCustomUiElement, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = recoverFromHook(err, "getCustomUi", r)
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

func (c *CustomGameMode) callAndApplyMatchHook(ctx contract.Context, matchHookName string, gameplay *Gameplay) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = recoverFromHook(err, matchHookName, r)
		}
	}()

	serializedGameState, err := c.serializeGameState(gameplay.State)
	if err != nil {
		return
	}

	serializedGameStateChangeActions, err := c.callMatchHook(ctx, matchHookName, serializedGameState)
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("Error in match hook %s: \n%v\n", matchHookName, err))
		return
	}

	err = c.deserializeAndApplyGameStateChangeActions(ctx, gameplay, serializedGameStateChangeActions)
	if err != nil {
		return
	}

	return
}

func recoverFromHook(err error, hookName string, r interface{}) error {
	stackBuffer := debug.Stack()
	runtime.Stack(stackBuffer, false)
	err = errors.New(fmt.Sprintf("! Error in custom mode hook - %s\n%v\n%s", hookName, r, string(stackBuffer)))
	return err
}

func (c *CustomGameMode) staticCallEVM(ctx contract.StaticContext, method string, result interface{}, params ...interface{}) error {
	input, err := zbCustomGameModeABI.Pack(method, params...)
	if err != nil {
		return err
	}
	var output []byte
	if err := contract.StaticCallEVM(ctx, c.tokenAddr, input, &output); err != nil {
		return err
	}
	return zbCustomGameModeABI.Unpack(result, method, output)
}

func (c *CustomGameMode) staticCallEVMRaw(ctx contract.StaticContext, abiInput []byte) ([]byte, error) {
	var output []byte
	if err := contract.StaticCallEVM(ctx, c.tokenAddr, abiInput, &output); err != nil {
		return nil, err
	}
	return output, nil
}

func (c *CustomGameMode) callEVM(ctx contract.Context, method string, params ...interface{}) (*types.EvmTxReceipt, error) {
	input, err := zbCustomGameModeABI.Pack(method, params...)
	if err != nil {
		return nil, err
	}

	return c.callEVMRaw(ctx, input)
}

func (c *CustomGameMode) callEVMRaw(ctx contract.Context, abiInput []byte) (*types.EvmTxReceipt, error) {
	var evmOut []byte
	err := contract.CallEVM(ctx, c.tokenAddr, abiInput, &evmOut)
	if err != nil {
		return nil, err
	}

	receipt, err := ctx.GetEvmTxReceipt(evmOut)
	if err != nil {
		return nil, err
	}
	return &receipt, nil
}

func (c *CustomGameMode) callMatchHook(ctx contract.Context, matchHookName string, serializedGameState []byte) (serializedGameStateChangeActions []byte, err error) {
	txReceipt, err := c.callEVM(ctx, matchHookName, serializedGameState)
	if err != nil {
		return nil, err
	}

	// ZBGameMode contract has default implementations of the hooks that don't emit any events.
	// So it is fine if the receipt has no logs, but if there are some, at least one of them must be GameStateChanges.
	if len(txReceipt.Logs) == 0 {
		return make([]byte, 0, 0), nil
	}

	var gameStateChangesLog *types.EventData = nil
	logsLoop:
	for _, log := range txReceipt.Logs {
		for _, logTopic := range log.Topics {
			if logTopic == zbGameStateChangesEventTopic {
				gameStateChangesLog = log
				break logsLoop
			}
		}
	}

	if gameStateChangesLog == nil {
		return nil, errors.New(fmt.Sprintf("Expected event %s", zbGameStateChangesEventName))
	}

	return gameStateChangesLog.EncodedBody, nil
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
          "indexed": false,
          "name": "serializedChanges",
          "type": "bytes"
        }
      ],
      "name": "GameStateChanges",
      "type": "event"
    },
    {
      "constant": false,
      "inputs": [
        {
          "name": "",
          "type": "bytes"
        }
      ],
      "name": "beforeMatchStart",
      "outputs": [],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "name": "",
          "type": "bytes"
        }
      ],
      "name": "afterInitialDraw",
      "outputs": [],
      "payable": false,
      "stateMutability": "nonpayable",
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
    }
  ]
`
const zbGameStateChangesEventName string = "GameStateChanges"
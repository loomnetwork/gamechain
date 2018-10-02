package battleground

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	loom "github.com/loomnetwork/go-loom"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	"github.com/loomnetwork/zombie_battleground/types/zb"
)

// We will create one instance of this per deployed game mode
type CustomGameMode struct {
	// Address of game mode contract deployed to Loom EVM.
	tokenAddr   loom.Address
	contractABI *abi.ABI
}

func (c *CustomGameMode) UpdateInitialPlayerGameState(ctx contract.Context, players []*zb.PlayerState) error {
	hpVal, err := c.GetStaticConfigs(ctx)
	if err != nil {
		return err
	}
	for _, v := range players {
		v.Hp = int32(hpVal)
	}
	return nil
}

//Todo simple example to get bootstrapped
func (c *CustomGameMode) GetStaticConfigs(ctx contract.Context) (int64, error) {

	//  GetStaticConfigs()

	var (
		ret0 = new([]*big.Int)
		ret1 = new([]*big.Int)
	)
	out := []interface{}{
		ret0,
		ret1,
	}
	if err := c.staticCallEVM(ctx, "getStaticConfigs", &out); err != nil {
		return 0, err
	}
	for k, v := range *ret1 {
		if k == 0 {
			ctx.Logger().Info(fmt.Sprintf("health----------------%v\n", v))
			return v.Int64(), nil
		}
	}
	return 0, nil
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

// From Zombiebattleground game mode repo
const zbGameModeABI = `
[
    {
      "constant": true,
      "inputs": [
        {
          "name": "_interfaceId",
          "type": "bytes4"
        }
      ],
      "name": "supportsInterface",
      "outputs": [
        {
          "name": "",
          "type": "bool"
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
          "type": "uint256"
        }
      ],
      "name": "userAccts",
      "outputs": [
        {
          "name": "",
          "type": "address"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "costToEnter",
      "outputs": [
        {
          "name": "",
          "type": "uint256"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "InterfaceId_ERC165",
      "outputs": [
        {
          "name": "",
          "type": "bytes4"
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
          "type": "address"
        }
      ],
      "name": "userGames",
      "outputs": [
        {
          "name": "status",
          "type": "uint256"
        },
        {
          "name": "gamesCount",
          "type": "uint256"
        },
        {
          "name": "wins",
          "type": "uint256"
        },
        {
          "name": "loses",
          "type": "uint256"
        }
      ],
      "payable": false,
      "stateMutability": "view",
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
      "name": "activeUsers",
      "outputs": [
        {
          "name": "",
          "type": "uint256"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
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
      "constant": false,
      "inputs": [
        {
          "name": "useraddr",
          "type": "address"
        },
        {
          "name": "ticketId",
          "type": "uint256"
        },
        {
          "name": "gameId",
          "type": "uint256"
        },
        {
          "name": "_v",
          "type": "uint8[]"
        },
        {
          "name": "_r",
          "type": "bytes32[]"
        },
        {
          "name": "_s",
          "type": "bytes32[]"
        }
      ],
      "name": "RegisterGame",
      "outputs": [],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "name": "useraddr1",
          "type": "address"
        },
        {
          "name": "useraddr2",
          "type": "address"
        }
      ],
      "name": "GameStart",
      "outputs": [],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "name": "player1Addr",
          "type": "address"
        },
        {
          "name": "player2Addr",
          "type": "address"
        },
        {
          "name": "winner",
          "type": "uint256"
        }
      ],
      "name": "GameFinished",
      "outputs": [],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    }
  ]
`

// +build evm

package battleground

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"

	//	"github.com/loomnetwork/go-loom/plugin"

	"github.com/loomnetwork/loomchain"
	"github.com/loomnetwork/loomchain/eth/subs"
	lvm "github.com/loomnetwork/loomchain/vm"
	"github.com/stretchr/testify/assert"
)

// Implements loomchain.EventHandler interface
type fakeEventHandler struct {
}

func (eh *fakeEventHandler) Post(height uint64, e *loomchain.EventData) error {
	return nil
}

func (eh *fakeEventHandler) EmitBlockTx(height uint64) error {
	return nil
}

func (eh *fakeEventHandler) SubscriptionSet() *loomchain.SubscriptionSet {
	return nil
}

func (eh *fakeEventHandler) EthSubscriptionSet() *subs.EthSubscriptionSet {
	return nil
}

func deployEVMContract(vm lvm.VM, contractBinary string, caller loom.Address) (loom.Address, *abi.ABI, error) {
	contractAddr := loom.Address{}
	//hexByteCode, err := ioutil.ReadFile("testdata/" + filename + ".bin")
	//if err != nil {
	//	return contractAddr, nil, err
	//}
	hexByteCode := contractBinary
	//	abiBytes, err := ioutil.ReadFile("testdata/" + filename + ".abi")
	//	if err != nil {
	//		return contractAddr, nil, err
	//	}
	abiBytes := zbGameModeABI
	contractABI, err := abi.JSON(strings.NewReader(string(abiBytes)))
	if err != nil {
		return contractAddr, nil, err
	}
	byteCode := common.FromHex(string(hexByteCode))
	fmt.Printf("vm-%v\n", vm)
	fmt.Printf("caller-%v\n", caller)
	_, contractAddr, err = vm.Create(caller, byteCode, loom.NewBigUIntFromInt(0))
	if err != nil {
		return contractAddr, nil, err
	}
	return contractAddr, &contractABI, nil
}

/*func TestCustomGameMode(t *testing.T) {

	var pubKeyHexString = "e4008e26428a9bca87465e8de3a8d0e9c37a56ca619d3d6202b0567528786618"
	pubKey, _ := hex.DecodeString(pubKeyHexString)

	addr := loom.Address{
		//ChainID: "default",
		Local: loom.LocalAddressFromPublicKey(pubKey),
	}

	owner := loom.RootAddress("chain")
	fakeCtx := plugin.CreateFakeContextWithEVM(addr, loom.RootAddress("chain"))

	evm := levm.NewLoomVm(fakeCtx.State, nil, nil, nil)
	evmContractAddr, evmContractABI, err := deployEVMContract(evm, "conquermode", owner)
	require.NoError(t, err)

	fmt.Printf("deployed contract -%v -%v \n", evmContractAddr, evmContractABI)

	gwCtx := contract.WrapPluginContext(fakeCtx.WithAddress(addr))

	cg := NewCustomGameMode(evmContractAddr)
	fmt.Printf("getsomething from addr-%v\n", addr)
	res, err := cg.GetStaticConfigs(gwCtx) //addr)

	assert.Nil(t, err)
	//if err != nil {
	//	assert.FailNow(t, "error reading data in contract")
	//}
	assert.Equal(t, res, int64(30))
}*/

func createSimpleGame(t *testing.T) *Gameplay {
	// fakeCtx := plugin.CreateFakeContextWithEVM(loom.RootAddress("chain"), loom.RootAddress("chain"))
	// gwCtx := contract.WrapPluginContext(fakeCtx.WithAddress(loom.RootAddress("chain")))
	var c *ZombieBattleground
	var pubKeyHexString = "e4008e26428a9bca87465e8de3a8d0e9c37a56ca619d3d6202b0567528786618"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)

	var deckList zb.DeckList
	err := ctx.Get(MakeVersionedKey("v1", defaultDeckKey), &deckList)
	assert.Nil(t, err)
	player1 := "player-1"
	player2 := "player-2"
	players := []*zb.PlayerState{
		{Id: player1, Deck: deckList.Decks[0]},
		{Id: player2, Deck: deckList.Decks[0]},
	}
	seed := int64(0)
	gp, err := NewGamePlay(ctx, 5, "v1", players, seed, nil, false)
	gp.createGame(nil)

	assert.Nil(t, err)

	return gp
}

func TestCustomGameMode_DeserializeGameStateChangeActions(t *testing.T) {
	gp := createSimpleGame(t)
	cgm := NewCustomGameMode(loom.RootAddress("chain"))
	buffer := common.FromHex("0x00000000000000000000000000000000000000000000000000000000000000000000000000000000080100000002070000000002060100000001050000000001")

	err := cgm.deserializeAndApplyGameStateChangeActions(*gp.ctx, gp.State, buffer)
	assert.Nil(t, err)

	assert.Equal(t, int32(5), gp.State.PlayerStates[0].Defense)
	assert.Equal(t, int32(6), gp.State.PlayerStates[1].Defense)
	assert.Equal(t, int32(7), gp.State.PlayerStates[0].CurrentGoo)
	assert.Equal(t, int32(8), gp.State.PlayerStates[1].CurrentGoo)
}

func TestCustomGameMode_DeserializeGameStateChangeActionsUnknownAction(t *testing.T) {
	gp := createSimpleGame(t)
	cgm := NewCustomGameMode(loom.RootAddress("chain"))
	buffer := common.FromHex("0x000000000000000000000000000000000000000000000000000000000000000000000000000000000801000000020700000000020601000000010500000000F9")

	err := cgm.deserializeAndApplyGameStateChangeActions(*gp.ctx, gp.State, buffer)
	assert.NotEqual(t, err, nil)
}

func TestCustomGameMode_DeserializeCustomUiElements(t *testing.T) {
	cgm := NewCustomGameMode(loom.RootAddress("chain"))
	buffer := common.FromHex("00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000436c69636b204d650000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000008000000960000012c0000012c000002a300000002536f6d65205665727920436f6f6c207465787421000000000000000000000000000000000000000000000000000000000000000000000000000000000000001400000096000000c8000000e60000001900000001")

	uiElements, err := cgm.deserializeCustomUi(buffer)
	assert.Nil(t, err)

	assert.Equal(
		t,
		zb.Rect{
			Position: &zb.Vector2Int{X: 25, Y: 230},
			Size_:    &zb.Vector2Int{X: 200, Y: 150},
		},
		*uiElements[0].Rect)
	label := uiElements[0].UiElement.(*zb.CustomGameModeCustomUiElement_Label)
	assert.Equal(t, "Some Very Cool text!", label.Label.Text)

	assert.Equal(
		t,
		zb.Rect{
			Position: &zb.Vector2Int{X: 675, Y: 300},
			Size_:    &zb.Vector2Int{X: 300, Y: 150},
		},
		*uiElements[1].Rect)

	button := uiElements[1].UiElement.(*zb.CustomGameModeCustomUiElement_Button)
	assert.Equal(t, "Click Me", button.Button.Title)
}

func TestCustomGameMode_SerializeGameState(t *testing.T) {
	gp := createSimpleGame(t)
	cgm := NewCustomGameMode(loom.RootAddress("chain"))
	bytes, err := cgm.serializeGameState(gp.State)
	assert.Nil(t, err)

	bytesHex := hexutil.Encode(bytes)

	assert.Equal(
		t,
		"0x0a0a06030000140100000000010000000301000000035075736868680000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000600000015010000000001000000020100000002507566666572000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000060000000e010000000001000000010100000002576865657a7900000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000006000000110100000000010000000101000000015768697a70617200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000700000013010000000001000000010100000001536f6f7468736179657200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a0000001000000005010000000001000000010100000001417a7572617a000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000060000000b010000000001000000010100000001417a7572617a000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000060000000c010000000001000000010100000002576865657a79000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000060000001200000003000000000000000244656661756c740000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000070000000000000000706c617965722d3200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000080a0a0603000014010000000001000000010100000001536f6f7468736179657200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a00000004010000000001000000010100000001536f6f7468736179657200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a00000005010000000001000000020100000003426f756e63657200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000700000009010000000001000000030100000003507573686868000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000060000000a010000000001000000010100000001417a7572617a000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000060000000000000005010000000001000000010100000001417a7572617a0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000600000001010000000001000000010100000002576865657a790000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000600000007010000000001000000010100000002576865657a79000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000060000000600000003000000000000000244656661756c740000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000070000000000000000706c617965722d310000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000008000000000000000005",
		bytesHex)
}

func TestCustomGameMode_DeserializeStrings(t *testing.T) {
	buffer := common.FromHex("0x436f6f6c20427574746f6e000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000b")
	rb := NewReverseBuffer(buffer)
	str, err := deserializeString(rb)
	assert.Nil(t, err)
	assert.Equal(t, str, "Cool Button")

	buffer = common.FromHex("0x000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000006f6c20427574746f6e20300000000000000000000000000000000000000000003720436f6f6c20427574746f6e203820436f6f6c20427574746f6e203920436f746f6e203520436f6f6c20427574746f6e203620436f6f6c20427574746f6e2020427574746f6e203320436f6f6c20427574746f6e203420436f6f6c20427574436f6f6c20427574746f6e203120436f6f6c20427574746f6e203220436f6f6c000000000000000000000000000000000000000000000000000000000000008b")
	rb = NewReverseBuffer(buffer)
	str, err = deserializeString(rb)
	assert.Nil(t, err)
	assert.Equal(t, str, "Cool Button 1 Cool Button 2 Cool Button 3 Cool Button 4 Cool Button 5 Cool Button 6 Cool Button 7 Cool Button 8 Cool Button 9 Cool Button 0")
}

func TestCustomGameMode_SerializeStrings(t *testing.T) {
	rb := NewReverseBuffer(make([]byte, 256))
	err := serializeString(rb, "Cool Button")
	assert.Nil(t, err)
	_, err = rb.Seek(0, io.SeekStart)
	assert.Nil(t, err)
	str, err := deserializeString(rb)
	assert.Nil(t, err)
	assert.Equal(t, str, "Cool Button")

	rb = NewReverseBuffer(make([]byte, 256))
	err = serializeString(rb, "Cool Button 1 Cool Button 2 Cool Button 3 Cool Button 4 Cool Button 5 Cool Button 6 Cool Button 7 Cool Button 8 Cool Button 9 Cool Button 0")
	assert.Nil(t, err)
	_, err = rb.Seek(0, io.SeekStart)
	assert.Nil(t, err)
	str, err = deserializeString(rb)
	assert.Nil(t, err)
	assert.Equal(t, str, "Cool Button 1 Cool Button 2 Cool Button 3 Cool Button 4 Cool Button 5 Cool Button 6 Cool Button 7 Cool Button 8 Cool Button 9 Cool Button 0")
}

// From Zombiebattleground game mode repo
const zbGameModeBIN = `0x6080604052601960065534801561001557600080fd5b506100487f01ffc9a700000000000000000000000000000000000000000000000000000000640100000000610093810204565b604080518082019091526001815260026020820181905261006991816100ff565b5060408051808201909152601e81526001602082015261008d90600390600261014a565b506101a7565b7fffffffff0000000000000000000000000000000000000000000000000000000080821614156100c257600080fd5b7fffffffff00000000000000000000000000000000000000000000000000000000166000908152602081905260409020805460ff19166001179055565b82805482825590600052602060002090810192821561013a579160200282015b8281111561013a57825182559160200191906001019061011f565b5061014692915061018a565b5090565b82805482825590600052602060002090810192821561013a579160200282015b8281111561013a578251829060ff1690559160200191906001019061016a565b6101a491905b808211156101465760008155600101610190565b90565b6108d4806101b66000396000f3006080604052600436106100ae5763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166301ffc9a781146100b357806306fdde03146100fe5780630ab2bb131461018857806314155419146101bc57806319fa8f50146101e357806346c84f521461022a57806374aa34de146102565780638ec585391461026b5780638fd70787146102b2578063919c5417146102d9578063d6a6cc4d14610326575b600080fd5b3480156100bf57600080fd5b506100ea7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19600435166103d4565b604080519115158252519081900360200190f35b34801561010a57600080fd5b50610113610408565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561014d578181015183820152602001610135565b50505050905090810190601f16801561017a5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561019457600080fd5b506101a060043561043f565b60408051600160a060020a039092168252519081900360200190f35b3480156101c857600080fd5b506101d1610467565b60408051918252519081900360200190f35b3480156101ef57600080fd5b506101f861046d565b604080517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff199092168252519081900360200190f35b34801561023657600080fd5b50610254600160a060020a0360043581169060243516604435610491565b005b34801561026257600080fd5b506101d161051f565b34801561027757600080fd5b5061028c600160a060020a0360043516610525565b604080519485526020850193909352838301919091526060830152519081900360800190f35b3480156102be57600080fd5b50610254600160a060020a036004358116906024351661054c565b3480156102e557600080fd5b5061025460048035600160a060020a03169060248035916044359160643580820192908101359160843580820192908101359160a4359081019101356105fb565b34801561033257600080fd5b5061033b6106af565b604051808060200180602001838103835285818151815260200191508051906020019060200280838360005b8381101561037f578181015183820152602001610367565b50505050905001838103825284818151815260200191508051906020019060200280838360005b838110156103be5781810151838201526020016103a6565b5050505090500194505050505060405180910390f35b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff191660009081526020819052604090205460ff1690565b60408051808201909152600b81527f436f6e717565724d6f6465000000000000000000000000000000000000000000602082015290565b600580548290811061044d57fe5b600091825260209091200154600160a060020a0316905081565b60065481565b7f01ffc9a70000000000000000000000000000000000000000000000000000000081565b60038060018314156104a957506001905060006104ba565b82600214156104ba57506000905060015b6104c4858361079b565b6104ce848261079b565b60408051600160a060020a0380881682528616602082015280820185905290517f5f7ceea33fec204ed93d7f1122512181f7ffa7f57ded7f75390cb195cfbbf49d9181900360600190a15050505050565b60055490565b60046020526000908152604090208054600182015460028301546003909301549192909184565b600160a060020a038083166000908152600460205260408082209284168252902060028254141561057957fe5b60028154141561058557fe5b6001825560018155604051600160a060020a038516907fa7b0e02d36afaa61a94171ae1d80245349ae39d9a7c0d2f35dd100c9288b674b90600090a2604051600160a060020a038416907fa7b0e02d36afaa61a94171ae1d80245349ae39d9a7c0d2f35dd100c9288b674b90600090a250505050565b600160a060020a03891660009081526004602052604090206001881461061d57fe5b600081556005805460018101825560009182527f036b6384b5eca791c62761152d0c79bb0604c104a5fb6f4eb0703f3154bb3db001805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a038d1690811790915560405190917f54db7a5cb4735e1aac1f53db512d3390390bb6637bd30ad4bf9fc98667d9b9b991a250505050505050505050565b60608060608060006002805490506040519080825280602002602001820160405280156106e6578160200160208202803883390190505b506003546040805182815260208084028201019091529194508015610715578160200160208202803883390190505b509150600090505b60025481101561079157600280548290811061073557fe5b9060005260206000200154838281518110151561074e57fe5b60209081029091010152600380548290811061076657fe5b9060005260206000200154828281518110151561077f57fe5b6020908102909101015260010161071d565b5090939092509050565b600160a060020a03821660009081526004602052604090206002815414156107bf57fe5b6001818101805482019055821415610883576002810180546001019081905560071415610828576006546040805191825251600160a060020a038516917fc409fe35c244d795f363c16fcfb864119423794fe7c1d3de8abcad79056c1cd8919081900360200190a25b8060020154600c141561087e576040805160018152600060208201528151600160a060020a038616927f6de83f319e4e49f88b6ce3aa22ea20cc78fafa95706834fc295aed209d279079928290030190a2600281555b6108a3565b8115156108a357600380820180546001019081905514156108a357600281555b5050505600a165627a7a723058206040d7185a7cd9f238802da3fd8d8c8660e5f7bf0a7854f9fb810b088f4468e10029`

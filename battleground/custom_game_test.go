// +build evm

package battleground

import (
	"encoding/hex"
	"fmt"
	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/require"
	"io"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/loomnetwork/go-loom"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	//	"github.com/loomnetwork/go-loom/plugin"

	"github.com/loomnetwork/loomchain"
	"github.com/loomnetwork/loomchain/eth/subs"
	"github.com/loomnetwork/loomchain/plugin"
	lvm "github.com/loomnetwork/loomchain/vm"
	levm "github.com/loomnetwork/loomchain/evm"
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
	fakeCtx := plugin.CreateFakeContextWithEVM(loom.RootAddress("chain"), loom.RootAddress("chain"))
	gwCtx := contract.WrapPluginContext(fakeCtx.WithAddress(loom.RootAddress("chain")))

	player1 := "player-1"
	player2 := "player-2"
	players := []*zb.PlayerState{
		{Id: player1, Deck: &defaultDeck1},
		{Id: player2, Deck: &defaultDeck2},
	}
	seed := int64(0)
	gp, err := NewGamePlay(gwCtx, 5, players, seed, nil)

	assert.Nil(t, err)

	return gp
}

func TestDeserializeGameStateChangeActions(t *testing.T) {
	gp := createSimpleGame(t)
	cgm := NewCustomGameMode(loom.RootAddress("chain"))
	buffer := common.FromHex("0x00000000000000000000000000000000000000000000000000000000000000000000000000000000080100000002070000000002060100000001050000000001")

	err := cgm.deserializeAndApplyGameStateChangeActions(gp.State, buffer)
	assert.Nil(t, err)

	assert.Equal(t, int32(5), gp.State.PlayerStates[0].Hp)
	assert.Equal(t, int32(6), gp.State.PlayerStates[1].Hp)
	assert.Equal(t, int32(7), gp.State.PlayerStates[0].Mana)
	assert.Equal(t, int32(8), gp.State.PlayerStates[1].Mana)
}

func TestDeserializeGameStateChangeActions2(t *testing.T) {
	gp := createSimpleGame(t)
	cgm := NewCustomGameMode(loom.RootAddress("chain"))
	buffer := common.FromHex("0x000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000044765726d73000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000500000000000000035a68616d70696f6e000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000802010000000300000000000000044765726d73000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000500000000000000035a68616d70696f6e0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000008020000000003")

	err := cgm.deserializeAndApplyGameStateChangeActions(gp.State, buffer)
	assert.Nil(t, err)

	assert.Equal(t, "Zhampion", gp.State.PlayerStates[0].Deck.Cards[0].CardName)
	assert.Equal(t, int64(3), gp.State.PlayerStates[0].Deck.Cards[0].Amount)
	assert.Equal(t, "Germs", gp.State.PlayerStates[0].Deck.Cards[1].CardName)
	assert.Equal(t, int64(4), gp.State.PlayerStates[0].Deck.Cards[1].Amount)

	assert.Equal(t, "Zhampion", gp.State.PlayerStates[1].Deck.Cards[0].CardName)
	assert.Equal(t, int64(3), gp.State.PlayerStates[1].Deck.Cards[0].Amount)
	assert.Equal(t, "Germs", gp.State.PlayerStates[1].Deck.Cards[1].CardName)
	assert.Equal(t, int64(4), gp.State.PlayerStates[1].Deck.Cards[1].Amount)
}

func TestDeserializeGameStateChangeActionsUnknownAction(t *testing.T) {
	gp := createSimpleGame(t)
	cgm := NewCustomGameMode(loom.RootAddress("chain"))
	buffer := common.FromHex("0x000000000000000000000000000000000000000000000000000000000000000000000000000000000801000000020700000000020601000000010500000000F9")

	err := cgm.deserializeAndApplyGameStateChangeActions(gp.State, buffer)
	assert.NotEqual(t, err, nil)
}

func TestDeserializeCustomUiElements(t *testing.T) {
	cgm := NewCustomGameMode(loom.RootAddress("chain"))
	buffer := common.FromHex("0x000000000000000000000000000000000000000000000000736f6d6546756e6374696f6e0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c436c69636b204d65000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000800000096000000c80000001e0000001900000002536f6d65205665727920436f6f6c207465787421000000000000000000000000000000000000000000000000000000000000000000000000000000000000001400000096000000c8000000e60000001900000001")

	uiElements, err := cgm.deserializeCustomUi(buffer)
	assert.Nil(t, err)

	assert.Equal(
		t,
		zb.Rect{
			Position: &zb.Vector2Int {	X: 25,	Y: 230},
			Size_: &zb.Vector2Int {	X: 200,	Y: 150},
		},
		*uiElements[0].Rect)
	label := uiElements[0].UiElement.(*zb.CustomGameModeCustomUiElement_Label)
	assert.Equal(t, "Some Very Cool text!", label.Label.Text)

	assert.Equal(
		t,
		zb.Rect{
			Position: &zb.Vector2Int {	X: 25,	Y: 30},
			Size_: &zb.Vector2Int {	X: 200,	Y: 150},
		},
		*uiElements[1].Rect)

	button := uiElements[1].UiElement.(*zb.CustomGameModeCustomUiElement_Button)
	assert.Equal(t, "Click Me", button.Button.Title)
	assert.Equal(t, "someFunction", button.Button.OnClickFunctionName)
}

func TestSerializeGameState(t *testing.T) {
	gp := createSimpleGame(t)
	cgm := NewCustomGameMode(loom.RootAddress("chain"))
	bytes, err := cgm.serializeGameState(gp.State)
	assert.Nil(t, err)

	bytesHex := hexutil.Encode(bytes)

	assert.Equal(t, "0x00000000000000024765797a6572000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000060000000000000004467265657a7a6565000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000800000000000000044a65747465720000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000600000000000000044f7a6d6f7a697a00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000700000000000000045a6e6f776d616e0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000070000000000000004497a7a65000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000000000000443657262657275730000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000008000000000000000447617267616e7475610000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000908000000000000000244656661756c743200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000080000000000000002011400000000000000025a68616d70696f6e00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000080000000000000004466972652d4d6177000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000800000000000000044d6f646f0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000040000000000000004576572657a6f6d620000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000008000000000000000443796e6465726d616e00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000009000000000000000442757272726e6e00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000700000000000000045175617a69000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000500000000000000045079726f6d617a00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000708000000000000000144656661756c7431000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000800000000000000010114000000000000000005", bytesHex)
}

func TestDeserializeStrings(t *testing.T) {
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

func TestSerializeStrings(t *testing.T) {
	rb := NewReverseBuffer(make([]byte, 256))
	err := serializeString(rb, "Cool Button")
	assert.Nil(t, err)
	_, err = rb.Seek(0, io.SeekStart)
	assert.Nil(t, err)
	fmt.Println(rb.remainingBytes)
	str, err := deserializeString(rb)
	fmt.Println(rb.buffer)
	assert.Nil(t, err)
	assert.Equal(t, str, "Cool Button")

	rb = NewReverseBuffer(make([]byte, 256))
	err = serializeString(rb, "Cool Button 1 Cool Button 2 Cool Button 3 Cool Button 4 Cool Button 5 Cool Button 6 Cool Button 7 Cool Button 8 Cool Button 9 Cool Button 0")
	assert.Nil(t, err)
	_, err = rb.Seek(0, io.SeekStart)
	assert.Nil(t, err)
	fmt.Println(rb.remainingBytes)
	str, err = deserializeString(rb)
	fmt.Println(rb.buffer)
	assert.Nil(t, err)
	assert.Equal(t, str, "Cool Button 1 Cool Button 2 Cool Button 3 Cool Button 4 Cool Button 5 Cool Button 6 Cool Button 7 Cool Button 8 Cool Button 9 Cool Button 0")
}

func TestMemoryLeak(t *testing.T) {
	var pubKeyHexString = "e4008e26428a9bca87465e8de3a8d0e9c37a56ca619d3d6202b0567528786618"
	pubKey, _ := hex.DecodeString(pubKeyHexString)

	addr := loom.Address{
		//ChainID: "default",
		Local: loom.LocalAddressFromPublicKey(pubKey),
	}

	owner := loom.RootAddress("chain")
	fakeCtx := plugin.CreateFakeContextWithEVM(addr, loom.RootAddress("chain"))

	evm := levm.NewLoomVm(fakeCtx.State, nil, nil, nil)
	evmContractAddr, evmContractABI, err := deployEVMContract(evm, testMemoryLeakBIN, owner)
	require.NoError(t, err)

	fmt.Printf("deployed contract -%v -%v \n", evmContractAddr, evmContractABI)

	gwCtx := contract.WrapPluginContext(fakeCtx.WithAddress(addr))


	player1 := "player-1"
	player2 := "player-2"
	players := []*zb.PlayerState{
		{Id: player1, Deck: &defaultDeck1},
		{Id: player2, Deck: &defaultDeck2},
	}
	seed := int64(0)
	gp, err := NewGamePlay(gwCtx, 5, players, seed, &evmContractAddr)
	assert.Nil(t, err)

	fmt.Printf("getsomething from addr-%v\n", addr)
	err = gp.customGameMode.UpdateInitialPlayerGameState(gwCtx, gp.State)
}

// From Zombiebattleground game mode repo
const zbGameModeBIN = `0x6080604052601960065534801561001557600080fd5b506100487f01ffc9a700000000000000000000000000000000000000000000000000000000640100000000610093810204565b604080518082019091526001815260026020820181905261006991816100ff565b5060408051808201909152601e81526001602082015261008d90600390600261014a565b506101a7565b7fffffffff0000000000000000000000000000000000000000000000000000000080821614156100c257600080fd5b7fffffffff00000000000000000000000000000000000000000000000000000000166000908152602081905260409020805460ff19166001179055565b82805482825590600052602060002090810192821561013a579160200282015b8281111561013a57825182559160200191906001019061011f565b5061014692915061018a565b5090565b82805482825590600052602060002090810192821561013a579160200282015b8281111561013a578251829060ff1690559160200191906001019061016a565b6101a491905b808211156101465760008155600101610190565b90565b6108d4806101b66000396000f3006080604052600436106100ae5763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166301ffc9a781146100b357806306fdde03146100fe5780630ab2bb131461018857806314155419146101bc57806319fa8f50146101e357806346c84f521461022a57806374aa34de146102565780638ec585391461026b5780638fd70787146102b2578063919c5417146102d9578063d6a6cc4d14610326575b600080fd5b3480156100bf57600080fd5b506100ea7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19600435166103d4565b604080519115158252519081900360200190f35b34801561010a57600080fd5b50610113610408565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561014d578181015183820152602001610135565b50505050905090810190601f16801561017a5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561019457600080fd5b506101a060043561043f565b60408051600160a060020a039092168252519081900360200190f35b3480156101c857600080fd5b506101d1610467565b60408051918252519081900360200190f35b3480156101ef57600080fd5b506101f861046d565b604080517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff199092168252519081900360200190f35b34801561023657600080fd5b50610254600160a060020a0360043581169060243516604435610491565b005b34801561026257600080fd5b506101d161051f565b34801561027757600080fd5b5061028c600160a060020a0360043516610525565b604080519485526020850193909352838301919091526060830152519081900360800190f35b3480156102be57600080fd5b50610254600160a060020a036004358116906024351661054c565b3480156102e557600080fd5b5061025460048035600160a060020a03169060248035916044359160643580820192908101359160843580820192908101359160a4359081019101356105fb565b34801561033257600080fd5b5061033b6106af565b604051808060200180602001838103835285818151815260200191508051906020019060200280838360005b8381101561037f578181015183820152602001610367565b50505050905001838103825284818151815260200191508051906020019060200280838360005b838110156103be5781810151838201526020016103a6565b5050505090500194505050505060405180910390f35b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff191660009081526020819052604090205460ff1690565b60408051808201909152600b81527f436f6e717565724d6f6465000000000000000000000000000000000000000000602082015290565b600580548290811061044d57fe5b600091825260209091200154600160a060020a0316905081565b60065481565b7f01ffc9a70000000000000000000000000000000000000000000000000000000081565b60038060018314156104a957506001905060006104ba565b82600214156104ba57506000905060015b6104c4858361079b565b6104ce848261079b565b60408051600160a060020a0380881682528616602082015280820185905290517f5f7ceea33fec204ed93d7f1122512181f7ffa7f57ded7f75390cb195cfbbf49d9181900360600190a15050505050565b60055490565b60046020526000908152604090208054600182015460028301546003909301549192909184565b600160a060020a038083166000908152600460205260408082209284168252902060028254141561057957fe5b60028154141561058557fe5b6001825560018155604051600160a060020a038516907fa7b0e02d36afaa61a94171ae1d80245349ae39d9a7c0d2f35dd100c9288b674b90600090a2604051600160a060020a038416907fa7b0e02d36afaa61a94171ae1d80245349ae39d9a7c0d2f35dd100c9288b674b90600090a250505050565b600160a060020a03891660009081526004602052604090206001881461061d57fe5b600081556005805460018101825560009182527f036b6384b5eca791c62761152d0c79bb0604c104a5fb6f4eb0703f3154bb3db001805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a038d1690811790915560405190917f54db7a5cb4735e1aac1f53db512d3390390bb6637bd30ad4bf9fc98667d9b9b991a250505050505050505050565b60608060608060006002805490506040519080825280602002602001820160405280156106e6578160200160208202803883390190505b506003546040805182815260208084028201019091529194508015610715578160200160208202803883390190505b509150600090505b60025481101561079157600280548290811061073557fe5b9060005260206000200154838281518110151561074e57fe5b60209081029091010152600380548290811061076657fe5b9060005260206000200154828281518110151561077f57fe5b6020908102909101015260010161071d565b5090939092509050565b600160a060020a03821660009081526004602052604090206002815414156107bf57fe5b6001818101805482019055821415610883576002810180546001019081905560071415610828576006546040805191825251600160a060020a038516917fc409fe35c244d795f363c16fcfb864119423794fe7c1d3de8abcad79056c1cd8919081900360200190a25b8060020154600c141561087e576040805160018152600060208201528151600160a060020a038616927f6de83f319e4e49f88b6ce3aa22ea20cc78fafa95706834fc295aed209d279079928290030190a2600281555b6108a3565b8115156108a357600380820180546001019081905514156108a357600281555b5050505600a165627a7a723058206040d7185a7cd9f238802da3fd8d8c8660e5f7bf0a7854f9fb810b088f4468e10029`
const testMemoryLeakBIN = `0x6080604052600060015534801561001557600080fd5b506040805160c08101825260018152600260208201526003918101919091526004606082015260056080820152600660a082018190526100579160009161005d565b506100ca565b82805482825590600052602060002090810192821561009d579160200282015b8281111561009d578251829060ff1690559160200191906001019061007d565b506100a99291506100ad565b5090565b6100c791905b808211156100a957600081556001016100b3565b90565b610e4f806100d96000396000f30060806040526004361061006c5763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166306fdde0381146100715780630f60f773146100fb5780634cf2e2dc146101545780635b34b9661461017b578063f72db61e14610192575b600080fd5b34801561007d57600080fd5b506100866101a7565b6040805160208082528351818301528351919283929083019185019080838360005b838110156100c05781810151838201526020016100a8565b50505050905090810190601f1680156100ed5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561010757600080fd5b506040805160206004803580820135601f81018490048402850184019095528484526100869436949293602493928401919081908401838280828437509497506101de9650505050505050565b34801561016057600080fd5b50610169610283565b60408051918252519081900360200190f35b34801561018757600080fd5b50610190610288565b005b34801561019e57600080fd5b50610086610292565b60408051808201909152600b81527f4578616d706c6547616d65000000000000000000000000000000000000000000602082015290565b60606101e8610da2565b6101f0610dc1565b604080516161a88082526161e0820190925260009182919060208201620c35008038833901905050945061027a565b83604001515182101561026e575060005b604084015180518390811061024157fe5b9060200190602002015160400151606001515181101561026357600101610230565b60019091019061021f565b610277836107ef565b94505b50505050919050565b600390565b6001805481019055565b606061029c610dc1565b6102ae8161040063ffffffff6107b616565b60408051608081018252601981830190815261012c60608301819052908252825180840184529081526096602082810191909152808301919091528251808401909352600f83527f436f756e7465722056616c75653a2000000000000000000000000000000000009083015261032c9183919063ffffffff61080516565b6040805160808101825261014581830190815261012c6060830181905290825282518084019093528252609660208381019190915281019190915260015461038691906103789061082d565b83919063ffffffff61080516565b604080516080810182526019818301908152601e60608301528152815180830183526101f4815260c86020828101919091528083019190915282518084018452600881527f436c69636b204d65000000000000000000000000000000000000000000000000818301528351808501909452600c84527f736f6d6546756e6374696f6e000000000000000000000000000000000000000091840191909152610436928492919063ffffffff61095616565b6040805160808101825261026c818301908152601e606083015281528151808301835261012c815260c86020828101919091528083019190915282518084018452600281527f2b31000000000000000000000000000000000000000000000000000000000000818301528351808501909452601084527f696e6372656d656e74436f756e74657200000000000000000000000000000000918401919091526104e7928492919063ffffffff61095616565b6104f0816107ef565b91505090565b6104fe610dff565b8152602001906001900390816104f65750506040890152600093505b8760400151518410156107ac5761053186886109b6565b604089015180518690811061054257fe5b6020908102909101015160ff909116905261055d60086109bb565b8603955061056b86886109b6565b604089015180518690811061057c57fe5b602090810290910181015160ff90921691015261059960086109bb565b860395506105a786886109b6565b600790810b900b83526105ba60406109bb565b860395506105c88688610bea565b9450846040519080825280601f01601f1916602001820160405280156105f8578160200160208202803883390190505b506020840181905261060d9087908990610c0f565b848603955061061c86886109b6565b836040019060070b908160070b8152505061063760406109bb565b8603955061064586886109b6565b915061065160086109bb565b860395508160ff1660405190808252806020026020018201604052801561069257816020015b61067f610dc1565b8152602001906001900390816106775790505b5060608401525060005b826060015151811015610776576106b38688610bea565b9450846040519080825280601f01601f1916602001820160405280156106e3578160200160208202803883390190505b5060608401518051839081106106f557fe5b60209081029091010151526060830151805161072a9188918a91908590811061071a57fe5b6020908102909101015151610c0f565b848603955061073986886109b6565b606084015180518390811061074a57fe5b6020908102909101810151600792830b90920b91015261076a60406109bb565b9095039460010161069c565b8288604001518581518110151561078957fe5b602090810290910101516040015261079f610dd9565b600190940193925061051a565b5050505050505050565b806040519080825280601f01601f1916602001820160405280156107e4578160200160208202803883390190505b508252602090910152565b5190565b60209093018051939093039092525050565b61081183600184610c88565b6108248360200151828560000151610cb2565b6107f381610cf8565b60608160008083818415156108775760408051808201909152600181527f30000000000000000000000000000000000000000000000000000000000000006020820152955061094c565b8493505b831561089257600190920191600a8404935061087b565b826040519080825280601f01601f1916602001820160405280156108c0578160200160208202803883390190505b5091505060001982015b84156109485781516000198201917f01000000000000000000000000000000000000000000000000000000000000006030600a8906010291849190811061090d57fe5b9060200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350600a850494506108ca565b8195505b5050505050919050565b61096284600285610c88565b6109758460200151838660000151610cb2565b61097e82610cf8565b6020850180519190910390819052845161099a91908390610cb2565b6109a381610cf8565b6020909401805194909403909352505050565b015190565b60008160088114610ac85760108114610ad15760188114610ada5760208114610ae35760288114610aec5760308114610af55760388114610afe5760408114610b075760488114610b105760508114610b195760588114610b225760608114610b2b5760688114610b345760708114610b3d5760788114610b465760808114610b4f5760888114610b585760908114610b615760988114610b6a5760a08114610b735760a88114610b7c5760b08114610b855760b88114610b8e5760c08114610b975760c88114610ba05760d08114610ba95760d88114610bb25760e08114610bbb5760e88114610bc45760f08114610bcd5760f88114610bd6576101008114610bdf5760209150610be4565b60019150610be4565b60029150610be4565b60039150610be4565b60049150610be4565b60059150610be4565b60069150610be4565b60079150610be4565b60089150610be4565b60099150610be4565b600a9150610be4565b600b9150610be4565b600c9150610be4565b600d9150610be4565b600e9150610be4565b600f9150610be4565b60109150610be4565b60119150610be4565b60129150610be4565b60139150610be4565b60149150610be4565b60159150610be4565b60169150610be4565b60179150610be4565b60189150610be4565b60199150610be4565b601a9150610be4565b601b9150610be4565b601c9150610be4565b601d9150610be4565b601e9150610be4565b601f9150610be4565b602091505b50919050565b80820151602081046001016000601f83161115610c05576001015b6020029392505050565b81830151602081046001016000601f83161115610c2a576001015b60005b81811015610c4f578486015160208202850152601f1990950194600101610c2d565b505050505050565b63ffffffff168460000151610c82565b610c7160206109bb565b602090920180519290920390915250565b90910152565b610c928383610d15565b610ca58360000151846020015183610d29565b6020909301929092525050565b815160208082049160009190061115610cc9576001015b60010160005b81811015610cf1576020810284015183860152601f1990940193600101610ccf565b5050505050565b80516020808204910615610d0a576001015b600101602002919050565b610c678260200151826002811115610c5757fe5b600080839050610d3e85828560000151610d58565b9050610d4f85828560200151610d58565b95945050505050565b80516000908390610d6e90829060030b87610c82565b610d7860206109bb565b81039050610d8e81846020015160030b87610c82565b610d9860206109bb565b9003949350505050565b6040805160608181018352600080835260208301529181019190915290565b60408051808201909152606081526000602082015290565b604080516080810182526000808252606060208301819052928201528181019190915290565b6040805160c08101825260008082526020820152908101610e1e610dd9565b9052905600a165627a7a72305820e0d8e44b9a66fa5beca7eddf0684e76cb7a4d628ad364e7442421b6b7d4664870029`
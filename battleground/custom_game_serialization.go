package battleground

import (
	"encoding/binary"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/gamechain/types/common"
	"github.com/loomnetwork/gamechain/types/zb"
	"io"

	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
)

func (c *CustomGameMode) serializeGameState(state *zb.GameState) (bytes []byte, err error) {
	rb := newPanicReaderWriterProxy(NewReverseBuffer(make([]byte, 8192)))
	binary.Write(rb, binary.BigEndian, int64(state.Id))
	binary.Write(rb, binary.BigEndian, byte(state.CurrentPlayerIndex))

	for _, playerState := range state.PlayerStates {
		serializeString(rb, playerState.Id)
		c.serializeDeck(rb, playerState.Deck)

		simpleCardsInHand := newSimpleCardInstanceArrayFromCardInstanceArray(playerState.CardsInHand)
		simpleCardsInDeck := newSimpleCardInstanceArrayFromCardInstanceArray(playerState.CardsInDeck)

		c.serializeSimpleCardInstanceArray(rb, simpleCardsInHand)
		c.serializeSimpleCardInstanceArray(rb, simpleCardsInDeck)
		binary.Write(rb, binary.BigEndian, byte(playerState.Defense))
		binary.Write(rb, binary.BigEndian, byte(playerState.CurrentGoo))
		binary.Write(rb, binary.BigEndian, byte(playerState.GooVials))
		binary.Write(rb, binary.BigEndian, uint32(playerState.TurnTime))
		binary.Write(rb, binary.BigEndian, byte(playerState.InitialCardsInHandCount))
		binary.Write(rb, binary.BigEndian, byte(playerState.MaxCardsInPlay))
		binary.Write(rb, binary.BigEndian, byte(playerState.MaxCardsInHand))
		binary.Write(rb, binary.BigEndian, byte(playerState.MaxGooVials))
	}

	return rb.readWriter.(*ReverseBuffer).GetFilledSlice(), nil
}

func newSimpleCardInstanceArrayFromCardInstanceArray(cards []*zb.CardInstance) ([]*SimpleCardInstance) {
	simpleCards := make([]*SimpleCardInstance, len(cards))
	for i, card := range cards {
		simpleCards[i] = newSimpleCardInstanceFromCardInstance(card)
	}

	return simpleCards
}

func (c *CustomGameMode) serializeDeck(writer io.Writer, deck *zb.Deck) (err error) {
	binary.Write(writer, binary.BigEndian, int64(deck.Id))
	serializeString(writer, deck.Name)
	binary.Write(writer, binary.BigEndian, int64(deck.HeroId))

	return nil
}

func (c *CustomGameMode) serializeSimpleCardInstance(writer io.Writer, card *SimpleCardInstance) (err error) {
	binary.Write(writer, binary.BigEndian, int32(card.instanceId))
	serializeString(writer, card.mouldName)
	binary.Write(writer, binary.BigEndian, int32(card.defense))
	binary.Write(writer, binary.BigEndian, bool(card.defenseInherited))
	binary.Write(writer, binary.BigEndian, int32(card.attack))
	binary.Write(writer, binary.BigEndian, bool(card.attackInherited))
	binary.Write(writer, binary.BigEndian, int32(card.gooCost))
	binary.Write(writer, binary.BigEndian, bool(card.gooCostInherited))

	return nil
}

func (c *CustomGameMode) deserializeSimpleCardInstance(reader io.Reader) (simpleCard *SimpleCardInstance, err error) {
	simpleCard = &SimpleCardInstance{}

	binary.Read(reader, binary.BigEndian, &simpleCard.instanceId)
	simpleCard.mouldName, err = deserializeString(reader)
	binary.Read(reader, binary.BigEndian, &simpleCard.defense)
	binary.Read(reader, binary.BigEndian, &simpleCard.defenseInherited)
	binary.Read(reader, binary.BigEndian, &simpleCard.attack)
	binary.Read(reader, binary.BigEndian, &simpleCard.attackInherited)
	binary.Read(reader, binary.BigEndian, &simpleCard.gooCost)
	binary.Read(reader, binary.BigEndian, &simpleCard.gooCostInherited)

	return
}

func (c *CustomGameMode) serializeSimpleCardInstanceArray(writer io.Writer, cards []*SimpleCardInstance) (err error) {
	binary.Write(writer, binary.BigEndian, uint32(len(cards)))

	for _, card := range cards {
		c.serializeSimpleCardInstance(writer, card)
	}

	return nil
}

func (c *CustomGameMode) deserializeSimpleCardInstanceArray(reader io.Reader) (cards []*SimpleCardInstance, err error) {
	var cardCount uint32
	binary.Read(reader, binary.BigEndian, &cardCount)

	cards = make([]*SimpleCardInstance, cardCount)
	for i := uint32(0); i < cardCount; i++ {
		cards[i], _ = c.deserializeSimpleCardInstance(reader)
	}

	return cards, nil
}

func (c *CustomGameMode) updateCardFromSimpleCard(ctx contract.Context, card *zb.CardInstance, simpleCard *SimpleCardInstance, cardLibraryCard *zb.Card) (*zb.CardInstance, error) {
	newCard := newCardInstanceFromCardDetails(
		cardLibraryCard,
		card.InstanceId,
		card.Owner,
	)

	newCard.Prototype = proto.Clone(newCard.Prototype).(*zb.CardPrototype)

	if !simpleCard.defenseInherited {
		newCard.Defense = simpleCard.defense
	}

	if !simpleCard.attackInherited {
		newCard.Attack = simpleCard.attack
	}

	if !simpleCard.gooCostInherited {
		newCard.GooCost = simpleCard.gooCost
	}

	return newCard, nil
}

func (c *CustomGameMode) updateCardsFromSimpleCards(
	ctx contract.Context,
	gameplay *Gameplay,
	cards []*zb.CardInstance,
	simpleCards []*SimpleCardInstance,
) (newCards []*zb.CardInstance, err error) {
	for _, simpleCard := range simpleCards {
		var newCard *zb.CardInstance
		isMatchingInstanceIdFound := false
		for _, card := range cards {
			if simpleCard.instanceId == card.InstanceId {
				cardLibraryCard, err := getCardDetails(gameplay.cardLibrary, simpleCard.mouldName)
				if err != nil {
					return nil, err
				}

				newCard, err = c.updateCardFromSimpleCard(ctx, card, simpleCard, cardLibraryCard)

				isMatchingInstanceIdFound = true
				break
			}
		}

		if !isMatchingInstanceIdFound {
			return nil, fmt.Errorf("card with instance ID %d not found", simpleCard.instanceId)
		}

		newCards = append(newCards, newCard)
	}

	return newCards, nil
}

func (c *CustomGameMode) deserializeAndApplyGameStateChangeActions(ctx contract.Context, gameplay *Gameplay, serializedActions []byte) (err error) {
	if len(serializedActions) == 0 {
		return nil
	}

	reader := newPanicReaderWriterProxy(NewReverseBuffer(serializedActions))
	for {
		var action battleground.GameStateChangeAction
		binary.Read(reader, binary.BigEndian, &action)

		mustBreak := false
		switch action {
		case battleground.GameStateChangeAction_None:
			mustBreak = true
		case battleground.GameStateChangeAction_SetPlayerDefense:
			var playerIndex byte
			binary.Read(reader, binary.BigEndian, &playerIndex)

			var newValue byte
			binary.Read(reader, binary.BigEndian, &newValue)

			gameplay.State.PlayerStates[playerIndex].Defense = int32(newValue)
		case battleground.GameStateChangeAction_SetPlayerCurrentGoo:
			var playerIndex byte
			binary.Read(reader, binary.BigEndian, &playerIndex)

			var newValue byte
			binary.Read(reader, binary.BigEndian, &newValue)

			gameplay.State.PlayerStates[playerIndex].CurrentGoo = int32(newValue)
		case battleground.GameStateChangeAction_SetPlayerGooVials:
			var playerIndex byte
			binary.Read(reader, binary.BigEndian, &playerIndex)

			var newValue byte
			binary.Read(reader, binary.BigEndian, &newValue)

			gameplay.State.PlayerStates[playerIndex].GooVials = int32(newValue)
		case battleground.GameStateChangeAction_SetPlayerCardsInDeck:
			var playerIndex byte
			binary.Read(reader, binary.BigEndian, &playerIndex)

			simpleCards, _ := c.deserializeSimpleCardInstanceArray(reader)
			gameplay.State.PlayerStates[playerIndex].CardsInDeck, err =
				c.updateCardsFromSimpleCards(
					ctx,
					gameplay,
					gameplay.State.PlayerStates[playerIndex].CardsInDeck,
					simpleCards,
				)

			if err != nil {
				return
			}
		case battleground.GameStateChangeAction_SetPlayerCardsInHand:
			var playerIndex byte
			binary.Read(reader, binary.BigEndian, &playerIndex)

			simpleCards, _ := c.deserializeSimpleCardInstanceArray(reader)
			gameplay.State.PlayerStates[playerIndex].CardsInHand, err =
				c.updateCardsFromSimpleCards(
					ctx,
					gameplay,
					gameplay.State.PlayerStates[playerIndex].CardsInHand,
					simpleCards,
				)
			if err != nil {
				return
			}
		case battleground.GameStateChangeAction_SetPlayerInitialCardsInHandCount:
			var playerIndex byte
			binary.Read(reader, binary.BigEndian, &playerIndex)

			var newValue byte
			binary.Read(reader, binary.BigEndian, &newValue)

			gameplay.State.PlayerStates[playerIndex].InitialCardsInHandCount = int32(newValue)
		case battleground.GameStateChangeAction_SetPlayerMaxCardsInPlay:
			var playerIndex byte
			binary.Read(reader, binary.BigEndian, &playerIndex)

			var newValue byte
			binary.Read(reader, binary.BigEndian, &newValue)

			gameplay.State.PlayerStates[playerIndex].MaxCardsInPlay = int32(newValue)
		case battleground.GameStateChangeAction_SetPlayerMaxCardsInHand:
			var playerIndex byte
			binary.Read(reader, binary.BigEndian, &playerIndex)

			var newValue byte
			binary.Read(reader, binary.BigEndian, &newValue)

			gameplay.State.PlayerStates[playerIndex].MaxCardsInHand = int32(newValue)
		case battleground.GameStateChangeAction_SetPlayerMaxGooVials:
			var playerIndex byte
			binary.Read(reader, binary.BigEndian, &playerIndex)

			var newValue byte
			binary.Read(reader, binary.BigEndian, &newValue)

			gameplay.State.PlayerStates[playerIndex].MaxGooVials = int32(newValue)
		case battleground.GameStateChangeAction_SetPlayerTurnTime:
			var playerIndex byte
			binary.Read(reader, binary.BigEndian, &playerIndex)

			var newValue uint32
			binary.Read(reader, binary.BigEndian, &newValue)

			gameplay.State.PlayerStates[playerIndex].TurnTime = int32(newValue)
		default:
			return fmt.Errorf("unknown game state change action %d", action)
		}

		if mustBreak {
			return nil
		}
	}

	return gameplay.validateGameState()
}

func (c *CustomGameMode) deserializeCustomUi(serializedCustomUi []byte) (uiElements []*zb.CustomGameModeCustomUiElement, err error) {
	if len(serializedCustomUi) == 0 {
		return make([]*zb.CustomGameModeCustomUiElement, 0), nil
	}

	rb := newPanicReaderWriterProxy(NewReverseBuffer(serializedCustomUi))
	for {
		var elementType battleground.CustomUiElement
		binary.Read(rb, binary.BigEndian, &elementType)

		mustBreak := false
		switch elementType {
		case battleground.CustomUiElement_None:
			mustBreak = true
		case battleground.CustomUiElement_Label:
			var element zb.CustomGameModeCustomUiElement
			var label zb.CustomGameModeCustomUiLabel

			rect, _ := deserializeRect(rb)
			element.Rect = &rect
			label.Text, _ = deserializeString(rb)
			element.UiElement = &zb.CustomGameModeCustomUiElement_Label{Label: &label}

			uiElements = append(uiElements, &element)
		case battleground.CustomUiElement_Button:
			var element zb.CustomGameModeCustomUiElement
			var button zb.CustomGameModeCustomUiButton

			rect, _ := deserializeRect(rb)
			element.Rect = &rect
			button.Title, _ = deserializeString(rb)
			callDataStr, _ := deserializeString(rb)
			button.CallData = []byte(callDataStr)
			element.UiElement = &zb.CustomGameModeCustomUiElement_Button{Button: &button}

			uiElements = append(uiElements, &element)
		default:
			return nil, fmt.Errorf("unknown custom UI element type %d", elementType)
		}

		if mustBreak {
			return
		}
	}
}

type SimpleCardInstance struct {
	instanceId       int32
	mouldName        string
	defense          int32
	defenseInherited bool
	attack           int32
	attackInherited  bool
	gooCost          int32
	gooCostInherited bool
}

func newSimpleCardInstanceFromCardInstance(card *zb.CardInstance) *SimpleCardInstance {
	return &SimpleCardInstance{
		instanceId:       card.InstanceId,
		mouldName:        card.Prototype.Name,
		attack:           card.Attack,
		attackInherited:  true,
		defense:          card.Defense,
		defenseInherited: true,
		gooCost:          card.Prototype.GooCost,
		gooCostInherited: true,
	}
}

type PanicReaderWriterProxy struct {
	readWriter io.ReadWriter
}

func newPanicReaderWriterProxy(readWriter io.ReadWriter) *PanicReaderWriterProxy {
	prw := new(PanicReaderWriterProxy)
	prw.readWriter = readWriter
	return prw
}

func (prw *PanicReaderWriterProxy) Read(p []byte) (n int, err error) {
	n, err = prw.readWriter.Read(p)
	if err != nil {
		panic(err)
	}

	return n, nil
}

func (prw *PanicReaderWriterProxy) Write(p []byte) (n int, err error) {
	n, err = prw.readWriter.Write(p)
	if err != nil {
		panic(err)
	}

	return n, nil
}

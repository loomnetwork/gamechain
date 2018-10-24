package battleground

import (
	"encoding/binary"
	"fmt"
	"github.com/loomnetwork/gamechain/types/common"
	"github.com/loomnetwork/gamechain/types/zb"
	"io"

	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
)

func (c *CustomGameMode) serializeGameState(state *zb.GameState) (bytes []byte, err error) {
	rb := NewPanicReaderWriterProxy(NewReverseBuffer(make([]byte, 8192)))
	binary.Write(rb, binary.BigEndian, int64(state.Id))
	binary.Write(rb, binary.BigEndian, byte(state.CurrentPlayerIndex))

	for _, playerState := range state.PlayerStates {
		serializeString(rb, playerState.Id)
		c.serializeDeck(rb, playerState.Deck)

		simpleCardsInHand := convertCardInstanceArrayToSimpleCardInstanceArray(playerState.CardsInHand)
		simpleCardsInDeck := convertCardInstanceArrayToSimpleCardInstanceArray(playerState.CardsInDeck)

		c.serializeSimpleCardInstanceArray(rb, simpleCardsInHand)
		c.serializeSimpleCardInstanceArray(rb, simpleCardsInDeck)
		binary.Write(rb, binary.BigEndian, byte(playerState.Defense))
		binary.Write(rb, binary.BigEndian, byte(playerState.CurrentGoo))
		binary.Write(rb, binary.BigEndian, byte(playerState.GooVials))
		binary.Write(rb, binary.BigEndian, byte(playerState.InitialCardsInHandCount))
		binary.Write(rb, binary.BigEndian, byte(playerState.MaxCardsInPlay))
		binary.Write(rb, binary.BigEndian, byte(playerState.MaxCardsInHand))
		binary.Write(rb, binary.BigEndian, byte(playerState.MaxGooVials))
	}

	return rb.readWriter.(*ReverseBuffer).GetFilledSlice(), nil
}

func convertCardInstanceArrayToSimpleCardInstanceArray(cards []*zb.CardInstance) ([]*SimpleCardInstance) {
	simpleCards := make([]*SimpleCardInstance, len(cards))
	for i, card := range cards {
		simpleCards[i] = &SimpleCardInstance{
			instanceId: card.InstanceId,
			name:       card.Prototype.Name,
		}
	}

	return simpleCards
}

func (c *CustomGameMode) serializeDeck(writer io.Writer, deck *zb.Deck) (err error) {
	binary.Write(writer, binary.BigEndian, int64(deck.Id))
	serializeString(writer, deck.Name)
	binary.Write(writer, binary.BigEndian, int64(deck.HeroId))

	return nil
}

func (c *CustomGameMode) serializeCardPrototype(writer io.Writer, card *zb.CardPrototype) (err error) {
	serializeString(writer, card.Name)
	binary.Write(writer, binary.BigEndian, uint8(card.GooCost))

	return nil
}

func (c *CustomGameMode) deserializeCardPrototype(reader io.Reader) (card *zb.CardPrototype, err error) {
	name, err := deserializeString(reader)

	var gooCost uint8
	binary.Read(reader, binary.BigEndian, &gooCost)

	return &zb.CardPrototype{
		Name:    name,
		GooCost: int32(gooCost),
	}, nil
}

func (c *CustomGameMode) serializeCardInstance(writer io.Writer, card *zb.CardInstance) (err error) {
	binary.Write(writer, binary.BigEndian, int32(card.InstanceId))
	c.serializeCardPrototype(writer, card.Prototype)
	binary.Write(writer, binary.BigEndian, int32(card.Defense))
	binary.Write(writer, binary.BigEndian, int32(card.Attack))
	serializeString(writer, card.Owner)

	return nil
}

func (c *CustomGameMode) deserializeCardInstance(reader io.Reader) (card *zb.CardInstance, err error) {
	var instanceId int32
	binary.Read(reader, binary.BigEndian, &instanceId)

	cardPrototype, _ := c.deserializeCardPrototype(reader)

	var defense int32
	binary.Read(reader, binary.BigEndian, &defense)

	var attack int32
	binary.Read(reader, binary.BigEndian, &attack)

	owner, _ := deserializeString(reader)

	return &zb.CardInstance{
		InstanceId: instanceId,
		Prototype:  cardPrototype,
		Defense:    defense,
		Attack:     attack,
		Owner:      owner,
	}, nil
}

func (c *CustomGameMode) serializeSimpleCardInstance(writer io.Writer, card *SimpleCardInstance) (err error) {
	binary.Write(writer, binary.BigEndian, int32(card.instanceId))
	serializeString(writer, card.name)

	return nil
}

func (c *CustomGameMode) deserializeSimpleCardInstance(reader io.Reader) (card *SimpleCardInstance, err error) {
	var instanceId int32
	binary.Read(reader, binary.BigEndian, &instanceId)

	name, err := deserializeString(reader)

	return &SimpleCardInstance{
		instanceId: instanceId,
		name:       name,
	}, nil
}

func (c *CustomGameMode) serializeCardInstanceArray(writer io.Writer, cards []*zb.CardInstance) (err error) {
	binary.Write(writer, binary.BigEndian, uint32(len(cards)))

	for _, card := range cards {
		c.serializeCardInstance(writer, card)
	}

	return nil
}

func (c *CustomGameMode) deserializeCardInstanceArray(reader io.Reader) (cards []*zb.CardInstance, err error) {
	var cardCount uint32
	binary.Read(reader, binary.BigEndian, &cardCount)

	cards = make([]*zb.CardInstance, cardCount)
	for i := uint32(0); i < cardCount; i++ {
		cards[i], _ = c.deserializeCardInstance(reader)
	}

	return cards, nil
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

func updateFromSimpleCards(ctx contract.Context, state *zb.GameState, cards []*zb.CardInstance, simpleCards []*SimpleCardInstance) (error) {
	cardLibrary, err := getCardLibrary(ctx, state.Version)
	if err != nil {
		return err
	}

	for _, simpleCard := range simpleCards {
		isMatchingInstanceIdFound := false
		for cardIndex, card := range cards {
			if simpleCard.instanceId == card.InstanceId {
				cardLibraryCard, err := getCardDetails(cardLibrary, simpleCard.name)
				if err != nil {
					return err
				}

				cards[cardIndex] =
					newCardInstanceFromCardDetails(
						cardLibraryCard,
						card.InstanceId,
						card.Owner,
					)

				isMatchingInstanceIdFound = true
				break
			}
		}

		if !isMatchingInstanceIdFound {
			return fmt.Errorf("card with instance ID %d not found", simpleCard.instanceId)
		}
	}

	return nil
}

func (c *CustomGameMode) deserializeAndApplyGameStateChangeActions(ctx contract.Context, state *zb.GameState, serializedActions []byte) (err error) {
	if len(serializedActions) == 0 {
		return nil
	}

	reader := NewPanicReaderWriterProxy(NewReverseBuffer(serializedActions))
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

			state.PlayerStates[playerIndex].Defense = int32(newValue)
		case battleground.GameStateChangeAction_SetPlayerCurrentGoo:
			var playerIndex byte
			binary.Read(reader, binary.BigEndian, &playerIndex)

			var newValue byte
			binary.Read(reader, binary.BigEndian, &newValue)

			state.PlayerStates[playerIndex].CurrentGoo = int32(newValue)
		case battleground.GameStateChangeAction_SetPlayerGooVials:
			var playerIndex byte
			binary.Read(reader, binary.BigEndian, &playerIndex)

			var newValue byte
			binary.Read(reader, binary.BigEndian, &newValue)

			state.PlayerStates[playerIndex].GooVials = int32(newValue)
		case battleground.GameStateChangeAction_SetPlayerCardsInDeck:
			var playerIndex byte
			binary.Read(reader, binary.BigEndian, &playerIndex)

			simpleCards, _ := c.deserializeSimpleCardInstanceArray(reader)
			err := updateFromSimpleCards(ctx, state, state.PlayerStates[playerIndex].CardsInDeck, simpleCards)

			if err != nil {
				return err
			}
		case battleground.GameStateChangeAction_SetPlayerCardsInHand:
			var playerIndex byte
			binary.Read(reader, binary.BigEndian, &playerIndex)

			simpleCards, _ := c.deserializeSimpleCardInstanceArray(reader)
			err := updateFromSimpleCards(ctx, state, state.PlayerStates[playerIndex].CardsInHand, simpleCards)

			if err != nil {
				return err
			}
		case battleground.GameStateChangeAction_SetPlayerInitialCardsInHandCount:
			var playerIndex byte
			binary.Read(reader, binary.BigEndian, &playerIndex)

			var newValue byte
			binary.Read(reader, binary.BigEndian, &newValue)

			state.PlayerStates[playerIndex].InitialCardsInHandCount = int32(newValue)
		case battleground.GameStateChangeAction_SetPlayerMaxCardsInPlay:
			var playerIndex byte
			binary.Read(reader, binary.BigEndian, &playerIndex)

			var newValue byte
			binary.Read(reader, binary.BigEndian, &newValue)

			state.PlayerStates[playerIndex].MaxCardsInPlay = int32(newValue)
		case battleground.GameStateChangeAction_SetPlayerMaxCardsInHand:
			var playerIndex byte
			binary.Read(reader, binary.BigEndian, &playerIndex)

			var newValue byte
			binary.Read(reader, binary.BigEndian, &newValue)

			state.PlayerStates[playerIndex].MaxCardsInHand = int32(newValue)
		case battleground.GameStateChangeAction_SetPlayerMaxGooVials:
			var playerIndex byte
			binary.Read(reader, binary.BigEndian, &playerIndex)

			var newValue byte
			binary.Read(reader, binary.BigEndian, &newValue)

			state.PlayerStates[playerIndex].MaxGooVials = int32(newValue)
		default:
			return fmt.Errorf("unknown game state change action %d", action)
		}

		if mustBreak {
			return nil
		}
	}
}

func (c *CustomGameMode) deserializeCustomUi(serializedCustomUi []byte) (uiElements []*zb.CustomGameModeCustomUiElement, err error) {
	if len(serializedCustomUi) == 0 {
		return make([]*zb.CustomGameModeCustomUiElement, 0), nil
	}

	rb := NewPanicReaderWriterProxy(NewReverseBuffer(serializedCustomUi))
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
	instanceId int32
	name       string
}

type PanicReaderWriterProxy struct {
	readWriter io.ReadWriter
}

func NewPanicReaderWriterProxy(readWriter io.ReadWriter) *PanicReaderWriterProxy {
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

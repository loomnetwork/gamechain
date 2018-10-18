package battleground

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/loomnetwork/gamechain/types/common"
	"github.com/loomnetwork/gamechain/types/zb"
)

func (c *CustomGameMode) serializeGameState(state *zb.GameState) (bytes []byte, err error) {
	rb := NewReverseBuffer(make([]byte, 8192))
	if err = binary.Write(rb, binary.BigEndian, int64(state.Id)); err != nil {
		return nil, err
	}
	if err = binary.Write(rb, binary.BigEndian, byte(state.CurrentPlayerIndex)); err != nil {
		return nil, err
	}
	for _, playerState := range state.PlayerStates {
		// Generic state
		if err = binary.Write(rb, binary.BigEndian, byte(playerState.Defense)); err != nil {
			return nil, err
		}

		if err = binary.Write(rb, binary.BigEndian, byte(playerState.CurrentGoo)); err != nil {
			return nil, err
		}

		if err = binary.Write(rb, binary.BigEndian, byte(playerState.GooVials)); err != nil {
			return nil, err
		}

		// Deck info
		if err = binary.Write(rb, binary.BigEndian, int64(playerState.Deck.Id)); err != nil {
			return nil, err
		}

		if err = serializeString(rb, playerState.Deck.Name); err != nil {
			return nil, err
		}

		if err = binary.Write(rb, binary.BigEndian, int64(playerState.Deck.HeroId)); err != nil {
			return nil, err
		}

		// Deck cards
		if err = c.serializeCardCollectionArray(rb, playerState.Deck.Cards); err != nil {
			return nil, err
		}

		if err = c.serializeCardInstanceArray(rb, playerState.CardsInDeck); err != nil {
			return nil, err
		}

		if err = c.serializeCardInstanceArray(rb, playerState.CardsInHand); err != nil {
			return nil, err
		}
	}

	return rb.GetFilledSlice(), nil
}

func (c *CustomGameMode) serializeCardCollection(rb *ReverseBuffer, card *zb.CardCollection) (err error) {
	if err = serializeString(rb, card.CardName); err != nil {
		return err
	}

	if err = binary.Write(rb, binary.BigEndian, uint64(card.Amount)); err != nil {
		return err
	}

	return nil
}

func (c *CustomGameMode) deserializeCardCollection(rb *ReverseBuffer) (card *zb.CardCollection, err error) {
	name, err := deserializeString(rb)
	if err != nil {
		return nil, err
	}

	var amount int64
	if err = binary.Read(rb, binary.BigEndian, &amount); err != nil {
		return nil, err
	}

	return &zb.CardCollection{
		CardName: name,
		Amount: amount,
	}, nil
}

func (c *CustomGameMode) serializeCardCollectionArray(rb *ReverseBuffer, cards []*zb.CardCollection) (err error) {
	if err = binary.Write(rb, binary.BigEndian, uint64(len(cards))); err != nil {
		return err
	}

	for _, card := range cards {
		if err = c.serializeCardCollection(rb, card); err != nil {
			return err
		}
	}

	return nil
}

func (c *CustomGameMode) deserializeCardCollectionArray(rb *ReverseBuffer) (cards []*zb.CardCollection, err error) {
	var cardCount uint64
	if err = binary.Read(rb, binary.BigEndian, &cardCount); err != nil {
		return
	}

	cards = make([]*zb.CardCollection, cardCount)
	for i := uint64(0); i < cardCount; i++ {
		cards[i], err = c.deserializeCardCollection(rb)
		if err != nil {
			return nil, err
		}
	}

	return cards, nil
}

func (c *CustomGameMode) serializeCardPrototype(rb *ReverseBuffer, card *zb.CardPrototype) (err error) {
	if err = serializeString(rb, card.Name); err != nil {
		return err
	}

	if err = binary.Write(rb, binary.BigEndian, uint8(card.GooCost)); err != nil {
		return err
	}

	return nil
}

func (c *CustomGameMode) deserializeCardPrototype(rb *ReverseBuffer) (card *zb.CardPrototype, err error) {
	name, err := deserializeString(rb)
	if err != nil {
		return nil, err
	}

	var gooCost uint8
	if err = binary.Read(rb, binary.BigEndian, &gooCost); err != nil {
		return nil, err
	}

	return &zb.CardPrototype{
		Name: name,
		GooCost: int32(gooCost),
	}, nil
}

func (c *CustomGameMode) serializeCardInstance(rb *ReverseBuffer, card *zb.CardInstance) (err error) {
	if err = binary.Write(rb, binary.BigEndian, int32(card.InstanceId)); err != nil {
		return err
	}

	if err = c.serializeCardPrototype(rb, card.Prototype); err != nil {
		return err
	}

	if err = binary.Write(rb, binary.BigEndian, int32(card.Defense)); err != nil {
		return err
	}

	if err = binary.Write(rb, binary.BigEndian, int32(card.Attack)); err != nil {
		return err
	}

	if err = serializeString(rb, card.Owner); err != nil {
		return err
	}

	return nil
}

func (c *CustomGameMode) deserializeCardInstance(rb *ReverseBuffer) (card *zb.CardInstance, err error) {
	var instanceId int32
	if err = binary.Read(rb, binary.BigEndian, &instanceId); err != nil {
		return nil, err
	}

	var cardPrototype *zb.CardPrototype
	if cardPrototype, err = c.deserializeCardPrototype(rb); err != nil {
		return nil, err
	}

	var defense int32
	if err = binary.Read(rb, binary.BigEndian, &defense); err != nil {
		return nil, err
	}

	var attack int32
	if err = binary.Read(rb, binary.BigEndian, &attack); err != nil {
		return nil, err
	}

	owner, err := deserializeString(rb)
	if err != nil {
		return nil, err
	}

	return &zb.CardInstance{
		InstanceId: instanceId,
		Prototype: cardPrototype,
		Defense: defense,
		Attack: attack,
		Owner: owner,
	}, nil
}

func (c *CustomGameMode) serializeCardInstanceArray(rb *ReverseBuffer, cards []*zb.CardInstance) (err error) {
	if err = binary.Write(rb, binary.BigEndian, uint64(len(cards))); err != nil {
		return err
	}

	for _, card := range cards {
		if err = c.serializeCardInstance(rb, card); err != nil {
			return err
		}
	}

	return nil
}

func (c *CustomGameMode) deserializeCardInstanceArray(rb *ReverseBuffer) (cards []*zb.CardInstance, err error) {
	var cardCount uint64
	if err = binary.Read(rb, binary.BigEndian, &cardCount); err != nil {
		return
	}

	cards = make([]*zb.CardInstance, cardCount)
	for i := uint64(0); i < cardCount; i++ {
		cards[i], err = c.deserializeCardInstance(rb)
		if err != nil {
			return nil, err
		}
	}

	return cards, nil
}

func (c *CustomGameMode) deserializeAndApplyGameStateChangeActions(state *zb.GameState, serializedActions []byte) (err error) {
	if len(serializedActions) == 0 {
		return nil
	}

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

			state.PlayerStates[playerIndex].Defense = int32(newDefense)
		case battleground.GameStateChangeAction_SetPlayerCurrentGoo:
			var playerIndex byte
			if err = binary.Read(rb, binary.BigEndian, &playerIndex); err != nil {
				return
			}

			var newCurrentGoo byte
			if err = binary.Read(rb, binary.BigEndian, &newCurrentGoo); err != nil {
				return
			}

			state.PlayerStates[playerIndex].CurrentGoo = int32(newCurrentGoo)
		case battleground.GameStateChangeAction_SetPlayerGooVials:
			var playerIndex byte
			if err = binary.Read(rb, binary.BigEndian, &playerIndex); err != nil {
				return
			}

			var newGooVials byte
			if err = binary.Read(rb, binary.BigEndian, &newGooVials); err != nil {
				return
			}

			state.PlayerStates[playerIndex].GooVials = int32(newGooVials)
		case battleground.GameStateChangeAction_SetPlayerInitialDeckCards:
			var playerIndex byte
			if err = binary.Read(rb, binary.BigEndian, &playerIndex); err != nil {
				return
			}

			state.PlayerStates[playerIndex].Deck.Cards, err = c.deserializeCardCollectionArray(rb)
			if err != nil {
				return err
			}
		case battleground.GameStateChangeAction_SetPlayerCardsInDeck:
			var playerIndex byte
			if err = binary.Read(rb, binary.BigEndian, &playerIndex); err != nil {
				return
			}

			state.PlayerStates[playerIndex].CardsInDeck, err = c.deserializeCardInstanceArray(rb)
			if err != nil {
				return err
			}
		case battleground.GameStateChangeAction_SetPlayerCardsInHand:
			var playerIndex byte
			if err = binary.Read(rb, binary.BigEndian, &playerIndex); err != nil {
				return
			}

			state.PlayerStates[playerIndex].CardsInHand, err = c.deserializeCardInstanceArray(rb)
			if err != nil {
				return err
			}
		default:
			return errors.New(fmt.Sprintf("Unknown game state change action %d", action))
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

	rb := NewReverseBuffer(serializedCustomUi)
	for {
		var elementType battleground.CustomUiElement
		if err = binary.Read(rb, binary.BigEndian, &elementType); err != nil {
			return
		}

		mustBreak := false
		switch elementType {
		case battleground.CustomUiElement_None:
			mustBreak = true
		case battleground.CustomUiElement_Label:
			var element zb.CustomGameModeCustomUiElement
			var label zb.CustomGameModeCustomUiLabel

			rect, err := deserializeRect(rb)
			if err != nil {
				return nil, err
			}
			element.Rect = &rect

			if label.Text, err = deserializeString(rb); err != nil {
				return nil, err
			}

			element.UiElement = &zb.CustomGameModeCustomUiElement_Label { Label: &label }

			uiElements = append(uiElements, &element)
		case battleground.CustomUiElement_Button:
			var element zb.CustomGameModeCustomUiElement
			var button zb.CustomGameModeCustomUiButton

			rect, err := deserializeRect(rb)
			if err != nil {
				return nil, err
			}
			element.Rect = &rect

			if button.Title, err = deserializeString(rb); err != nil {
				return nil, err
			}

			if button.OnClickFunctionName, err = deserializeString(rb); err != nil {
				return nil, err
			}

			element.UiElement = &zb.CustomGameModeCustomUiElement_Button { Button: &button }

			uiElements = append(uiElements, &element)
		default:
			return nil, errors.New(fmt.Sprintf("Unknown custom UI element type %d", elementType))
		}

		if mustBreak {
			return
		}
	}
}

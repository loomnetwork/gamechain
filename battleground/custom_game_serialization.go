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
		if err = binary.Write(rb, binary.BigEndian, byte(playerState.Hp)); err != nil {
			return nil, err
		}
		if err = binary.Write(rb, binary.BigEndian, byte(playerState.Mana)); err != nil {
			return nil, err
		}

		// Deck
		if err = binary.Write(rb, binary.BigEndian, int64(playerState.Deck.Id)); err != nil {
			return nil, err
		}

		if err = serializeString(rb, playerState.Deck.Name); err != nil {
			return nil, err
		}

		if err = binary.Write(rb, binary.BigEndian, int64(playerState.Deck.HeroId)); err != nil {
			return nil, err
		}

		if err = binary.Write(rb, binary.BigEndian, uint8(len(playerState.Deck.Cards))); err != nil {
			return nil, err
		}

		for _, card := range playerState.Deck.Cards {
			if err = serializeString(rb, card.CardName); err != nil {
				return nil, err
			}

			if err = binary.Write(rb, binary.BigEndian, int64(card.Amount)); err != nil {
				return nil, err
			}
		}
	}

	return rb.GetFilledSlice(), nil
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
		case battleground.GameStateChangeAction_SetPlayerDeckCards:
			var playerIndex byte
			if err = binary.Read(rb, binary.BigEndian, &playerIndex); err != nil {
				return
			}

			var cardCount byte
			if err = binary.Read(rb, binary.BigEndian, &cardCount); err != nil {
				return
			}

			cards := make([]*zb.CardCollection, cardCount)
			for i := byte(0); i < cardCount; i++ {
				name, err := deserializeString(rb)
				if err != nil {
					return err
				}

				var amount int64
				if err = binary.Read(rb, binary.BigEndian, &amount); err != nil {
					return err
				}

				cards[i] = &zb.CardCollection{
					CardName: name,
					Amount: amount,
				}
			}

			state.PlayerStates[playerIndex].Deck.Cards = cards
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

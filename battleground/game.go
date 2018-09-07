package battleground

import (
	"fmt"

	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/pkg/errors"
)

var (
	errInvalidActivePlayer = errors.New("invalid active player")
)

type stateFn func(*gameplay) stateFn

type gameplay struct {
	State   zb.GameState
	stateFn stateFn
	errc    chan error
}

func InitGameplay(state zb.GameState) *gameplay {
	// TODO: validate players

	// TODO: shuffle cards
	// state.Player1.CardInDeck = cardListFromDeck(state.Player1.Deck)
	// state.Player2.CardInDeck = cardListFromDeck(state.Player2.Deck)

	return &gameplay{
		State: state,
		errc:  make(chan error),
	}
}

func (g *gameplay) PlayAction(action *zb.PlayerAction) error {
	if err := g.checkCurrentPlayer(action); err != nil {
		return err
	}
	g.State.PlayerActions = append(g.State.PlayerActions, action)
	return nil
}

func (g *gameplay) checkCurrentPlayer(action *zb.PlayerAction) error {
	activePlayer := g.activePlayer()
	if activePlayer.Id != action.PlayerId {
		return errInvalidActivePlayer
	}
	return nil
}

func (g *gameplay) play() {
	for g.stateFn = gameStart; g.stateFn != nil; {
		g.stateFn = g.stateFn(g)
	}
}

func (g *gameplay) next() *zb.PlayerAction {
	if g.State.CurrentActionIndex+1 > int64(len(g.State.PlayerActions)-1) {
		return nil
	}
	action := g.State.PlayerActions[g.State.CurrentActionIndex+1]
	g.State.CurrentActionIndex++
	return action
}

func (g *gameplay) peek() *zb.PlayerAction {
	if g.State.CurrentActionIndex+1 > int64(len(g.State.PlayerActions)) {
		return nil
	}

	action := g.State.PlayerActions[g.State.CurrentActionIndex+1]
	return action
}

func (g *gameplay) current() *zb.PlayerAction {
	action := g.State.PlayerActions[g.State.CurrentActionIndex]
	return action
}

func (g *gameplay) activePlayer() *zb.PlayerState {
	return g.State.PlayerStates[g.State.CurrentPlayerIndex]
}

func (g *gameplay) printState() {
	state := g.State
	fmt.Printf("============State Info=============\n")
	fmt.Printf("Current Player: %v\n", state.CurrentPlayerIndex)

	for i, player := range g.State.PlayerStates {
		if g.State.CurrentPlayerIndex == int32(i) {
			fmt.Printf("Player1: %s 🧟\n", player.Id)
		} else {
			fmt.Printf("Player1: %s\n", player.Id)
		}
		fmt.Printf("\thp: %v\n", player.Hp)
		fmt.Printf("\tmana: %v\n", player.Mana)
		// fmt.Printf("\tdeck: %v\n", state.Player1.Deck)
		fmt.Printf("\tcard in hand (%d): %v\n", len(player.CardInHand), player.CardInHand)
		fmt.Printf("\tcard on field (%d): %v\n", len(player.CardOnBoard), player.CardOnBoard)
		fmt.Printf("\tcard in deck (%d): %v\n", len(player.CardInDeck), player.CardInDeck)
	}

	fmt.Printf("Actions: count %v\n", len(state.PlayerActions))
	for i, action := range state.PlayerActions {
		if int64(i) == state.CurrentActionIndex {
			fmt.Printf("\t>>> %v\n", action)
		} else {
			fmt.Printf("\t%v\n", action)
		}
	}
	fmt.Printf("Current Action Index: %v\n", state.CurrentActionIndex)
	fmt.Printf("=========================\n")
}

func gameStart(g *gameplay) stateFn {
	fmt.Printf("state: %v\n", zb.ActionType_START_GAME)
	if isEnded(g) {
		return nil
	}

	// determine the next action
	g.printState()
	next := g.next()
	if next == nil {
		return nil
	}

	switch next.ActionType {
	case zb.ActionType_DRAW_CARD:
		return actionDrawCard
	default:
		return nil
	}
}

func actionDrawCard(g *gameplay) stateFn {
	fmt.Printf("state: %v\n", zb.ActionType_DRAW_CARD)
	if isEnded(g) {
		return nil
	}
	current := g.current()
	if current == nil {
		return nil
	}

	// check player turn
	if err := g.checkCurrentPlayer(current); err != nil {
		return nil
	}
	// drawcard to the active player

	// player := g.activePlayer()
	// cards, deck := drawFromCardList(player.CardInDeck, 1)
	// TODO: persist card to state
	// fmt.Printf("--> drawn card: %v\n", cards)
	// fmt.Printf("--> card in deck: %v\n", deck)

	// determine the next action
	g.printState()
	next := g.next()
	if next == nil {
		return nil
	}

	switch next.ActionType {
	case zb.ActionType_END_TURN:
		return actionEndTurn
	case zb.ActionType_DRAW_CARD:
		return actionDrawCard
	default:
		return nil
	}
}

func actionCardAttack(g *gameplay) stateFn {
	fmt.Printf("state: %v\n", zb.ActionType_CARD_ATTACK)
	if isEnded(g) {
		return nil
	}
	// current action
	current := g.current()
	if current == nil {
		return nil
	}

	// check player turn
	if err := g.checkCurrentPlayer(current); err != nil {
		return nil
	}

	// determine the next action
	g.printState()
	next := g.next()
	if next == nil {
		return nil
	}

	switch next.ActionType {
	case zb.ActionType_END_TURN:
		return actionEndTurn
	default:
		return nil
	}
}

func actionEndTurn(g *gameplay) stateFn {
	fmt.Printf("state: %v\n", zb.ActionType_END_TURN)
	if isEnded(g) {
		return nil
	}
	// current action
	current := g.current()
	if current == nil {
		return nil
	}
	// check player turn
	if err := g.checkCurrentPlayer(current); err != nil {
		return nil
	}
	g.switchPlayerTurn()

	// determine the next action
	g.printState()
	next := g.next()
	if next == nil {
		return nil
	}

	switch next.ActionType {
	case zb.ActionType_END_TURN:
		return actionEndTurn
	case zb.ActionType_DRAW_CARD:
		return actionDrawCard
	default:
		return nil
	}
}

func (g *gameplay) switchPlayerTurn() {
	g.State.CurrentPlayerIndex = (g.State.CurrentPlayerIndex + 1) % int32(len(g.State.PlayerStates))
}

// // drawCardFromDeck draw cards from active user decks and update user state
// func (g *gameplay) drawCardFromDeck(player *zb.ActivePlayer, n int) (cards []*zb.Card) {
// 	var player *zb.Player
// 	if g.State.ActivePlayer == zb.ActivePlayer_player1 {
// 		player = g.State.Player1
// 	} else {
// 		player = g.State.Player2
// 	}

// 	return
// }

func isEnded(g *gameplay) bool {
	for _, player := range g.State.PlayerStates {
		if player.Hp <= 0 {
			return true
		}
	}
	return false
}

package battleground

import (
	"fmt"

	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/pkg/errors"
)

var (
	errInvalidPlayer         = errors.New("invalid player")
	errCurrentActionNotfound = errors.New("current action not found")
	errInvalidAction         = errors.New("invalid action")
)

type stateFn func(*gameplay) stateFn

type gameplay struct {
	State   zb.GameState
	stateFn stateFn
	errc    chan error
}

func NewGameplay(state zb.GameState) *gameplay {
	// TODO: validate players

	// TODO: shuffle cards
	for i := range state.PlayerStates {
		state.PlayerStates[i].CardsInDeck = cardInstanceFromDeck(state.PlayerStates[i].Deck)
	}

	g := &gameplay{
		State: state,
		errc:  make(chan error),
	}
	go g.run()
	return g
}

// Wait blocks until errc returns error or get closed
func (g *gameplay) Wait() error {
	select {
	case err := <-g.errc:
		return err
	}
}

func (g *gameplay) AddAction(action *zb.PlayerAction) error {
	if err := g.checkCurrentPlayer(action); err != nil {
		return err
	}
	g.State.PlayerActions = append(g.State.PlayerActions, action)
	// reset errc to make sure that game does not block
	g.errc = make(chan error)
	// resume the gameplay
	go g.resume()
	select {
	case err := <-g.errc:
		return err
	}
}

func (g *gameplay) checkCurrentPlayer(action *zb.PlayerAction) error {
	activePlayer := g.activePlayer()
	if activePlayer.Id != action.PlayerId {
		return errInvalidPlayer
	}
	return nil
}

func (g *gameplay) run() {
	for g.stateFn = gameStart; g.stateFn != nil; {
		g.stateFn = g.stateFn(g)
	}
	close(g.errc)
	fmt.Printf("gameplay stopped at action index %d\n", g.State.CurrentActionIndex)
}

func (g *gameplay) resume() {
	// get the current state
	next := g.next()
	if next == nil {
		g.actionError(errCurrentActionNotfound)
		return
	}
	var state stateFn
	switch next.ActionType {
	case zb.PlayerActionType_CardAttack:
		state = actionCardAttack
	case zb.PlayerActionType_DrawCardPlayer:
		state = actionDrawCard
	case zb.PlayerActionType_PlayCard:
		state = actionPlayCard
	case zb.PlayerActionType_EndTurn:
		state = actionEndTurn
	default:
		g.actionError(errInvalidAction)
		return
	}

	fmt.Printf("gameplay resumed at action index %d\n", g.State.CurrentActionIndex)
	for g.stateFn = state; g.stateFn != nil; {
		g.stateFn = g.stateFn(g)
	}
	close(g.errc)
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

func (g *gameplay) changePlayerTurn() {
	g.State.CurrentPlayerIndex = (g.State.CurrentPlayerIndex + 1) % int32(len(g.State.PlayerStates))
}

func (g *gameplay) actionError(err error) stateFn {
	g.errc <- err
	return nil
}

func (g *gameplay) isEnded() bool {
	for _, player := range g.State.PlayerStates {
		if player.Hp <= 0 {
			return true
		}
	}
	return false
}

func (g *gameplay) printState() {
	state := g.State
	fmt.Printf("============StateInfo=============\n")
	fmt.Printf("Current Player Index: %v\n", state.CurrentPlayerIndex)

	for i, player := range g.State.PlayerStates {
		if g.State.CurrentPlayerIndex == int32(i) {
			fmt.Printf("Player%d: %s 🧟\n", i, player.Id)
		} else {
			fmt.Printf("Player%d: %s\n", i, player.Id)
		}
		fmt.Printf("\thp: %v\n", player.Hp)
		fmt.Printf("\tmana: %v\n", player.Mana)
		// fmt.Printf("\tdeck: %v\n", state.Player1.Deck)
		fmt.Printf("\tcard in hand (%d): %v\n", len(player.CardsInHand), player.CardsInHand)
		fmt.Printf("\tcard on board (%d): %v\n", len(player.CardsOnBoard), player.CardsOnBoard)
		fmt.Printf("\tcard in deck (%d): %v\n", len(player.CardsInDeck), player.CardsInDeck)
	}

	fmt.Printf("Actions: count %v\n", len(state.PlayerActions))
	for i, action := range state.PlayerActions {
		if int64(i) == state.CurrentActionIndex {
			fmt.Printf("   -->> [%d] %v\n", i, action)
		} else {
			fmt.Printf("\t[%d] %v\n", i, action)
		}
	}
	fmt.Printf("Current Action Index: %v\n", state.CurrentActionIndex)
	fmt.Printf("=========================\n")
}

func gameStart(g *gameplay) stateFn {
	fmt.Printf("state: %v\n", zb.PlayerActionType_StartGame)
	if g.isEnded() {
		return nil
	}

	// determine the next action
	g.printState()
	next := g.next()
	if next == nil {
		return nil
	}

	switch next.ActionType {
	case zb.PlayerActionType_DrawCardPlayer:
		return actionDrawCard
	default:
		return nil
	}
}

func actionDrawCard(g *gameplay) stateFn {
	fmt.Printf("state: %v\n", zb.PlayerActionType_DrawCardPlayer)
	if g.isEnded() {
		return nil
	}
	current := g.current()
	if current == nil {
		return nil
	}

	// check player turn
	if err := g.checkCurrentPlayer(current); err != nil {
		return g.actionError(err)
	}

	// draw card
	// TODO: handle card limit in hand
	// TODO: handle empty deck
	if len(g.activePlayer().CardsInDeck) > 0 {
		card := g.activePlayer().CardsInDeck[0]
		g.activePlayer().CardsInHand = append(g.activePlayer().CardsInHand, card)
		// remove card from CardsInDeck
		g.activePlayer().CardsInDeck = g.activePlayer().CardsInDeck[1:]
	}

	// determine the next action
	g.printState()
	next := g.next()
	if next == nil {
		return nil
	}

	switch next.ActionType {
	case zb.PlayerActionType_EndTurn:
		return actionEndTurn
	case zb.PlayerActionType_DrawCardPlayer:
		return actionDrawCard
	case zb.PlayerActionType_PlayCard:
		return actionPlayCard
	case zb.PlayerActionType_CardAttack:
		return actionCardAttack
	default:
		return nil
	}
}

func actionPlayCard(g *gameplay) stateFn {
	fmt.Printf("state: %v\n", zb.PlayerActionType_PlayCard)
	if g.isEnded() {
		return nil
	}
	current := g.current()
	if current == nil {
		return nil
	}

	// check player turn
	if err := g.checkCurrentPlayer(current); err != nil {
		return g.actionError(err)
	}

	// draw card
	// TODO: handle card limit on board
	if len(g.activePlayer().CardsInHand) > 0 {
		card := g.activePlayer().CardsInHand[0]
		g.activePlayer().CardsOnBoard = append(g.activePlayer().CardsOnBoard, card)
		g.activePlayer().CardsInHand = g.activePlayer().CardsInHand[1:]
	}

	// determine the next action
	g.printState()
	next := g.next()
	if next == nil {
		return nil
	}

	switch next.ActionType {
	case zb.PlayerActionType_EndTurn:
		return actionEndTurn
	case zb.PlayerActionType_DrawCardPlayer:
		return actionDrawCard
	case zb.PlayerActionType_PlayCard:
		return actionPlayCard
	case zb.PlayerActionType_CardAttack:
		return actionCardAttack
	default:
		return nil
	}
}

func actionCardAttack(g *gameplay) stateFn {
	fmt.Printf("state: %v\n", zb.PlayerActionType_CardAttack)
	if g.isEnded() {
		return nil
	}
	// current action
	current := g.current()
	if current == nil {
		return nil
	}

	// check player turn
	if err := g.checkCurrentPlayer(current); err != nil {
		return g.actionError(err)
	}

	// TODO: card attack

	// determine the next action
	g.printState()
	next := g.next()
	if next == nil {
		return nil
	}

	switch next.ActionType {
	case zb.PlayerActionType_EndTurn:
		return actionEndTurn
	case zb.PlayerActionType_DrawCardPlayer:
		return actionDrawCard
	case zb.PlayerActionType_PlayCard:
		return actionPlayCard
	case zb.PlayerActionType_CardAttack:
		return actionCardAttack
	default:
		return nil
	}
}

func actionEndTurn(g *gameplay) stateFn {
	fmt.Printf("state: %v\n", zb.PlayerActionType_EndTurn)
	if g.isEnded() {
		return nil
	}
	// current action
	current := g.current()
	if current == nil {
		return nil
	}
	// check player turn
	if err := g.checkCurrentPlayer(current); err != nil {
		return g.actionError(err)
	}
	// change player turn
	g.changePlayerTurn()

	// determine the next action
	g.printState()
	next := g.next()
	if next == nil {
		return nil
	}

	switch next.ActionType {
	case zb.PlayerActionType_EndTurn:
		return actionEndTurn
	case zb.PlayerActionType_DrawCardPlayer:
		return actionDrawCard
	case zb.PlayerActionType_PlayCard:
		return actionPlayCard
	case zb.PlayerActionType_CardAttack:
		return actionCardAttack
	default:
		return nil
	}
}

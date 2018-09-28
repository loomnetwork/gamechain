package battleground

import (
	"fmt"
	"math/rand"

	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/pkg/errors"
)

var (
	errInvalidPlayer         = errors.New("invalid player")
	errCurrentActionNotfound = errors.New("current action not found")
	errInvalidAction         = errors.New("invalid action")
	errNotEnoughPlayer       = errors.New("not enough players")
	errAlreadyTossCoin       = errors.New("already tossed coin")
	errNoCurrentPlayer       = errors.New("no current player")
)

type Gameplay struct {
	State   *zb.GameState
	stateFn stateFn
	err     error
}

type stateFn func(*Gameplay) stateFn

// NewGamePlay initializes GamePlay with default game state and run to the  latest state
func NewGamePlay(id int64, players []*zb.PlayerState) (*Gameplay, error) {
	state := &zb.GameState{
		Id:                 id,
		CurrentActionIndex: -1, // use -1 to avoid confict with default value
		PlayerStates:       players,
		CurrentPlayerIndex: -1, // use -1 to avoid confict with default value
	}
	return GamePlayFrom(state)
}

// GamePlayFrom initializes and run game to the latest state
func GamePlayFrom(state *zb.GameState) (*Gameplay, error) {
	g := &Gameplay{State: state}
	if err := g.run(); err != nil {
		return nil, err
	}
	return g, nil
}

// TossCoin decides who the first player is
func (g *Gameplay) TossCoin(seed int64) error {
	if len(g.State.PlayerStates) == 0 {
		return errNotEnoughPlayer
	}
	// prevent modifiying already-init state
	if g.State.CurrentPlayerIndex != -1 {
		return errAlreadyTossCoin
	}

	// TODO: test this function on multinode validators
	r := rand.New(rand.NewSource(seed))
	n := r.Int31n(int32(len(g.State.PlayerStates)))
	g.State.CurrentPlayerIndex = n
	return nil
}

// InitPlayers suffle cards in players' decks
func (g *Gameplay) InitPlayers(seed int64) error {
	if g.State.CurrentPlayerIndex < 0 {
		return errNoCurrentPlayer
	}
	for i := 0; i < len(g.State.PlayerStates); i++ {
		deck := g.State.PlayerStates[i].Deck
		g.State.PlayerStates[i].CardsInDeck = shuffleCardInDeck(deck, seed)
		g.State.PlayerStates[i].Hp = 20
		g.State.PlayerStates[i].Mana = 1

		// draw cards 3 card for mulligan
		n := 3
		g.State.PlayerStates[i].CardsInHand = g.State.PlayerStates[i].CardsInDeck[:n]
		g.State.PlayerStates[i].CardsInDeck = g.State.PlayerStates[i].CardsInDeck[n:]
	}
	return nil
}

func (g *Gameplay) AddAction(action *zb.PlayerAction) error {
	if err := g.checkCurrentPlayer(action); err != nil {
		return err
	}
	g.State.PlayerActions = append(g.State.PlayerActions, action)
	// resume the Gameplay
	return g.resume()
}

func (g *Gameplay) checkCurrentPlayer(action *zb.PlayerAction) error {
	activePlayer := g.activePlayer()
	if activePlayer.Id != action.PlayerId {
		return errInvalidPlayer
	}
	return nil
}

func (g *Gameplay) run() error {
	for g.stateFn = gameStart; g.stateFn != nil; {
		g.stateFn = g.stateFn(g)
	}
	fmt.Printf("Gameplay stopped at action index %d, err=%v\n", g.State.CurrentActionIndex, g.err)
	return g.err
}

func (g *Gameplay) resume() error {
	// get the current state
	next := g.next()
	if next == nil {
		return errCurrentActionNotfound
	}
	var state stateFn
	switch next.ActionType {
	case zb.PlayerActionType_CardAttack:
		state = actionCardAttack
	case zb.PlayerActionType_DrawCard:
		state = actionDrawCard
	case zb.PlayerActionType_CardPlay:
		state = actionCardPlay
	case zb.PlayerActionType_EndTurn:
		state = actionEndTurn
	default:
		return errInvalidAction
	}

	fmt.Printf("Gameplay resumed at action index %d\n", g.State.CurrentActionIndex)
	for g.stateFn = state; g.stateFn != nil; {
		g.stateFn = g.stateFn(g)
	}
	return g.err
}

func (g *Gameplay) next() *zb.PlayerAction {
	if g.State.CurrentActionIndex+1 > int64(len(g.State.PlayerActions)-1) {
		return nil
	}
	action := g.State.PlayerActions[g.State.CurrentActionIndex+1]
	g.State.CurrentActionIndex++
	return action
}

func (g *Gameplay) peek() *zb.PlayerAction {
	if g.State.CurrentActionIndex+1 > int64(len(g.State.PlayerActions)) {
		return nil
	}
	action := g.State.PlayerActions[g.State.CurrentActionIndex+1]
	return action
}

func (g *Gameplay) current() *zb.PlayerAction {
	action := g.State.PlayerActions[g.State.CurrentActionIndex]
	return action
}

func (g *Gameplay) activePlayer() *zb.PlayerState {
	return g.State.PlayerStates[g.State.CurrentPlayerIndex]
}

func (g *Gameplay) changePlayerTurn() {
	g.State.CurrentPlayerIndex = (g.State.CurrentPlayerIndex + 1) % int32(len(g.State.PlayerStates))
}

func (g *Gameplay) captureErrorAndStop(err error) stateFn {
	g.err = err
	return nil
}

func (g *Gameplay) isEnded() bool {
	for _, player := range g.State.PlayerStates {
		if player.Hp <= 0 {
			return true
		}
	}
	return false
}

func (g *Gameplay) PrintState() {
	state := g.State
	fmt.Printf("============StateInfo=============\n")
	fmt.Printf("Current Player Index: %v\n", state.CurrentPlayerIndex)

	for i, player := range g.State.PlayerStates {
		if g.State.CurrentPlayerIndex == int32(i) {
			fmt.Printf("Player%d: %s 🧟\n", i+1, player.Id)
		} else {
			fmt.Printf("Player%d: %s\n", i+1, player.Id)
		}
		fmt.Printf("\thp: %v\n", player.Hp)
		fmt.Printf("\tmana: %v\n", player.Mana)
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
	fmt.Printf("==================================\n")
}

func gameStart(g *Gameplay) stateFn {
	fmt.Printf("state: gameStart\n")
	if g.isEnded() {
		return nil
	}

	// determine the next action
	g.PrintState()
	next := g.next()
	if next == nil {
		return nil
	}

	switch next.ActionType {
	case zb.PlayerActionType_DrawCard:
		return actionDrawCard
	case zb.PlayerActionType_CardPlay:
		return actionCardPlay
	default:
		return nil
	}
}

func actionDrawCard(g *Gameplay) stateFn {
	fmt.Printf("state: %v\n", zb.PlayerActionType_DrawCard)
	if g.isEnded() {
		return nil
	}
	current := g.current()
	if current == nil {
		return nil
	}

	// check player turn
	if err := g.checkCurrentPlayer(current); err != nil {
		return g.captureErrorAndStop(err)
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
	g.PrintState()
	next := g.next()
	if next == nil {
		return nil
	}

	switch next.ActionType {
	case zb.PlayerActionType_EndTurn:
		return actionEndTurn
	case zb.PlayerActionType_DrawCard:
		return actionDrawCard
	case zb.PlayerActionType_CardPlay:
		return actionCardPlay
	case zb.PlayerActionType_CardAttack:
		return actionCardAttack
	default:
		return nil
	}
}

func actionCardPlay(g *Gameplay) stateFn {
	fmt.Printf("state: %v\n", zb.PlayerActionType_CardPlay)
	if g.isEnded() {
		return nil
	}
	current := g.current()
	if current == nil {
		return nil
	}

	// check player turn
	if err := g.checkCurrentPlayer(current); err != nil {
		return g.captureErrorAndStop(err)
	}

	// draw card
	// TODO: handle card limit on board
	if len(g.activePlayer().CardsInHand) > 0 {
		card := g.activePlayer().CardsInHand[0]
		g.activePlayer().CardsOnBoard = append(g.activePlayer().CardsOnBoard, card)
		g.activePlayer().CardsInHand = g.activePlayer().CardsInHand[1:]
	}

	// determine the next action
	g.PrintState()
	next := g.next()
	if next == nil {
		return nil
	}

	switch next.ActionType {
	case zb.PlayerActionType_EndTurn:
		return actionEndTurn
	case zb.PlayerActionType_DrawCard:
		return actionDrawCard
	case zb.PlayerActionType_CardPlay:
		return actionCardPlay
	case zb.PlayerActionType_CardAttack:
		return actionCardAttack
	default:
		return nil
	}
}

func actionCardAttack(g *Gameplay) stateFn {
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
		return g.captureErrorAndStop(err)
	}

	// TODO: card attack

	// determine the next action
	g.PrintState()
	next := g.next()
	if next == nil {
		return nil
	}

	switch next.ActionType {
	case zb.PlayerActionType_EndTurn:
		return actionEndTurn
	case zb.PlayerActionType_DrawCard:
		return actionDrawCard
	case zb.PlayerActionType_CardPlay:
		return actionCardPlay
	case zb.PlayerActionType_CardAttack:
		return actionCardAttack
	default:
		return nil
	}
}

func actionEndTurn(g *Gameplay) stateFn {
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
		return g.captureErrorAndStop(err)
	}
	// change player turn
	g.changePlayerTurn()

	// determine the next action
	g.PrintState()
	next := g.next()
	if next == nil {
		return nil
	}

	switch next.ActionType {
	case zb.PlayerActionType_EndTurn:
		return actionEndTurn
	case zb.PlayerActionType_DrawCard:
		return actionDrawCard
	case zb.PlayerActionType_CardPlay:
		return actionCardPlay
	case zb.PlayerActionType_CardAttack:
		return actionCardAttack
	default:
		return nil
	}
}

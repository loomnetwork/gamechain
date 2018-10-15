package battleground

import (
	"fmt"
	"math/rand"

	"github.com/loomnetwork/go-loom"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"

	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/pkg/errors"
)

const (
	mulliganCards = 3
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
	State          *zb.GameState
	stateFn        stateFn
	err            error
	customGameMode *CustomGameMode
	history        []*zb.HistoryData
}

type stateFn func(*Gameplay) stateFn

// NewGamePlay initializes GamePlay with default game state and run to the  latest state
func NewGamePlay(ctx contract.Context,
	id int64,
	version string,
	players []*zb.PlayerState,
	seed int64,
	customGameAddress *loom.Address,
) (*Gameplay, error) {
	var customGameMode *CustomGameMode
	if customGameAddress != nil {
		ctx.Logger().Info(fmt.Sprintf("Playing a custom game mode -%v", customGameAddress.String()))
		customGameMode = NewCustomGameMode(*customGameAddress)
	}

	state := &zb.GameState{
		Id:                 id,
		CurrentActionIndex: -1, // use -1 to avoid confict with default value
		PlayerStates:       players,
		CurrentPlayerIndex: -1, // use -1 to avoid confict with default value
		Randomseed:         seed,
		Version:            version,
	}
	g := &Gameplay{
		State:          state,
		customGameMode: customGameMode,
	}

	if err := g.createGame(); err != nil {
		return nil, err
	}

	if g.customGameMode != nil {
		err := g.customGameMode.UpdateInitialPlayerGameState(ctx, g.State)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("Error in custom game mode -%v", err))
			return nil, err
		}
	}
	if err := g.run(); err != nil {
		return nil, err
	}
	return g, nil
}

// GamePlayFrom initializes and run game to the latest state
func GamePlayFrom(state *zb.GameState) (*Gameplay, error) {
	g := &Gameplay{State: state}
	if err := g.run(); err != nil {
		return nil, err
	}
	return g, nil
}

func (g *Gameplay) createGame() error {
	// init players
	for i := 0; i < len(g.State.PlayerStates); i++ {
		g.State.PlayerStates[i].Hp = 20
		g.State.PlayerStates[i].Mana = 1
	}
	// coin toss for the first player
	r := rand.New(rand.NewSource(g.State.Randomseed))
	n := r.Int31n(int32(len(g.State.PlayerStates)))
	g.State.CurrentPlayerIndex = n

	// init hands
	for i := 0; i < len(g.State.PlayerStates); i++ {
		deck := g.State.PlayerStates[i].Deck
		g.State.PlayerStates[i].CardsInDeck = shuffleCardInDeck(deck, g.State.Randomseed)
		// draw cards 3 card for mulligan
		g.State.PlayerStates[i].CardsInHand = g.State.PlayerStates[i].CardsInDeck[:mulliganCards]
		g.State.PlayerStates[i].CardsInDeck = g.State.PlayerStates[i].CardsInDeck[mulliganCards:]
	}

	// add history data
	ps := make([]*zb.Player, len(g.State.PlayerStates))
	for i := range g.State.PlayerStates {
		ps[i] = &zb.Player{
			Id:   g.State.PlayerStates[i].Id,
			Deck: g.State.PlayerStates[i].Deck,
		}
	}
	// record history data
	g.history = append(g.history, &zb.HistoryData{
		Data: &zb.HistoryData_CreateGame{
			CreateGame: &zb.HistoryCreateGame{
				GameId:     g.State.Id,
				Players:    ps,
				Randomseed: g.State.Randomseed,
				Version:    g.State.Version,
			},
		},
	})
	return nil
}

// AddAction adds the given action and reruns the game state
func (g *Gameplay) AddAction(action *zb.PlayerAction) error {
	if err := g.checkCurrentPlayer(action); err != nil {
		return err
	}
	g.State.PlayerActions = append(g.State.PlayerActions, action)
	// resume the Gameplay
	return g.resume()
}

func (g *Gameplay) AddBundleAction(actions ...*zb.PlayerAction) error {
	for _, action := range actions {
		g.State.PlayerActions = append(g.State.PlayerActions, action)
	}
	// resume the Gameplay
	return g.resume()
}

func (g *Gameplay) checkCurrentPlayer(action *zb.PlayerAction) error {
	// skip checking for mulligan action
	if action.ActionType == zb.PlayerActionType_Mulligan {
		return nil
	}
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
	case zb.PlayerActionType_DrawCard:
		state = actionDrawCard
	case zb.PlayerActionType_CardPlay:
		state = actionCardPlay
	case zb.PlayerActionType_CardAttack:
		state = actionCardAttack
	case zb.PlayerActionType_CardAbilityUsed:
		state = actionCardAbilityUsed
	case zb.PlayerActionType_OverlordSkillUsed:
		state = actionOverloadSkillUsed
	case zb.PlayerActionType_EndTurn:
		state = actionEndTurn
	case zb.PlayerActionType_Mulligan:
		state = actionMulligan
	case zb.PlayerActionType_LeaveMatch:
		state = actionLeaveMatch
	case zb.PlayerActionType_RankBuff:
		state = actionRankBuff
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
	fmt.Printf("Is ended: %v, Winner: %s\n", state.IsEnded, state.Winner)
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
		fmt.Printf("\tcard in play (%d): %v\n", len(player.CardsInPlay), player.CardsInPlay)
		fmt.Printf("\tcard in deck (%d): %v\n", len(player.CardsInDeck), player.CardsInDeck)
	}

	fmt.Printf("History : count %v\n", len(g.history))
	for i, block := range g.history {
		fmt.Printf("\t[%d] %v\n", i, block)
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
	case zb.PlayerActionType_Mulligan:
		return actionMulligan
	case zb.PlayerActionType_CardPlay:
		return actionCardPlay
	case zb.PlayerActionType_CardAttack:
		return actionCardAttack
	case zb.PlayerActionType_CardAbilityUsed:
		return actionCardAbilityUsed
	case zb.PlayerActionType_OverlordSkillUsed:
		return actionOverloadSkillUsed
	case zb.PlayerActionType_EndTurn:
		return actionEndTurn
	case zb.PlayerActionType_LeaveMatch:
		return actionLeaveMatch
	case zb.PlayerActionType_RankBuff:
		return actionRankBuff
	default:
		return nil
	}
}

func actionMulligan(g *Gameplay) stateFn {
	fmt.Printf("state: %v\n", zb.PlayerActionType_Mulligan)
	if g.isEnded() {
		return nil
	}
	current := g.current()
	if current == nil {
		return nil
	}

	mulligan := current.GetMulligan()
	if mulligan == nil {
		return g.captureErrorAndStop(fmt.Errorf("expect mulligan action"))
	}
	var player *zb.PlayerState
	for i := 0; i < len(g.State.PlayerStates); i++ {
		if g.State.PlayerStates[i].Id == current.PlayerId {
			player = g.State.PlayerStates[i]
		}
	}
	if player == nil {
		return g.captureErrorAndStop(fmt.Errorf("player not found"))
	}

	// Check if all the mulliganed cards and number of card that can be mulligan
	if len(mulligan.MulliganedCards) > mulliganCards {
		return g.captureErrorAndStop(fmt.Errorf("number of mulligan card is exceed the maximum: %d", mulliganCards))
	}
	for _, card := range mulligan.MulliganedCards {
		_, found := containCardInCardList(card, player.CardsInHand)
		if !found {
			return g.captureErrorAndStop(fmt.Errorf("invalid mulligan card"))
		}
	}

	// keep only the cards in in mulligan
	keepCards := make([]*zb.CardInstance, 0)
	for _, mcard := range mulligan.MulliganedCards {
		card, found := containCardInCardList(mcard, player.CardsInHand)
		if found {
			keepCards = append(keepCards, card)
		}
	}

	// if the card in hand not match with the keep card, draw new cards
	rerollCards := make([]*zb.CardInstance, 0)
	if len(keepCards) == 0 {
		rerollCards = append(rerollCards, player.CardsInHand...)
	} else {
		for _, card := range player.CardsInHand {
			_, found := containCardInCardList(card, keepCards)
			if !found {
				rerollCards = append(rerollCards, card)
			}
		}
	}

	// set card in hands
	player.CardsInHand = keepCards
	// place cards back to deck
	player.CardsInDeck = append(player.CardsInDeck, rerollCards...)
	// draw card to replace the reroll cards
	for range rerollCards {
		player.CardsInHand = append(player.CardsInHand, player.CardsInDeck[0])
		// TODO: return to deck with random order
		player.CardsInDeck = player.CardsInDeck[1:]
	}

	// determine the next action
	g.PrintState()
	next := g.next()
	if next == nil {
		return nil
	}

	switch next.ActionType {
	case zb.PlayerActionType_Mulligan:
		return actionMulligan
	case zb.PlayerActionType_CardPlay:
		return actionCardPlay
	case zb.PlayerActionType_CardAttack:
		return actionCardAttack
	case zb.PlayerActionType_CardAbilityUsed:
		return actionCardAbilityUsed
	case zb.PlayerActionType_OverlordSkillUsed:
		return actionOverloadSkillUsed
	case zb.PlayerActionType_EndTurn:
		return actionEndTurn
	case zb.PlayerActionType_LeaveMatch:
		return actionLeaveMatch
	case zb.PlayerActionType_RankBuff:
		return actionRankBuff
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
	if len(g.activePlayer().CardsInDeck) < 1 {
		return g.captureErrorAndStop(errors.New("Can't draw card. No more cards in deck"))
	}

	card := g.activePlayer().CardsInDeck[0]
	g.activePlayer().CardsInHand = append(g.activePlayer().CardsInHand, card)
	// remove card from CardsInDeck
	g.activePlayer().CardsInDeck = g.activePlayer().CardsInDeck[1:]

	// record history data
	g.history = append(g.history, &zb.HistoryData{
		Data: &zb.HistoryData_FullInstance{
			FullInstance: &zb.HistoryFullInstance{
				InstanceId: card.InstanceId,
				Attack:     card.Attack,
				Defense:    card.Defence,
			},
		},
	})

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
	case zb.PlayerActionType_CardAbilityUsed:
		return actionCardAbilityUsed
	case zb.PlayerActionType_OverlordSkillUsed:
		return actionOverloadSkillUsed
	case zb.PlayerActionType_LeaveMatch:
		return actionLeaveMatch
	case zb.PlayerActionType_RankBuff:
		return actionRankBuff
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
	// TODO: handle card limit
	if len(g.activePlayer().CardsInHand) < 1 {
		return g.captureErrorAndStop(errors.New("Can't play card. No cards in hand"))
	}
	card := g.activePlayer().CardsInHand[0]
	g.activePlayer().CardsInPlay = append(g.activePlayer().CardsInPlay, card)
	g.activePlayer().CardsInHand = g.activePlayer().CardsInHand[1:]

	// record history data
	g.history = append(g.history, &zb.HistoryData{
		Data: &zb.HistoryData_FullInstance{
			FullInstance: &zb.HistoryFullInstance{
				InstanceId: card.InstanceId,
				Attack:     card.Attack,
				Defense:    card.Defence,
			},
		},
	})

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
	case zb.PlayerActionType_CardAbilityUsed:
		return actionCardAbilityUsed
	case zb.PlayerActionType_OverlordSkillUsed:
		return actionOverloadSkillUsed
	case zb.PlayerActionType_LeaveMatch:
		return actionLeaveMatch
	case zb.PlayerActionType_RankBuff:
		return actionRankBuff
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

	// record history data
	g.history = append(g.history, &zb.HistoryData{
		Data: &zb.HistoryData_ChangeInstance{
			ChangeInstance: &zb.HistoryInstance{
				InstanceId: 1, // TODO change to the actual card id
				Value:      2,
			},
		},
	})

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
	case zb.PlayerActionType_CardAbilityUsed:
		return actionCardAbilityUsed
	case zb.PlayerActionType_OverlordSkillUsed:
		return actionOverloadSkillUsed
	case zb.PlayerActionType_LeaveMatch:
		return actionLeaveMatch
	case zb.PlayerActionType_RankBuff:
		return actionRankBuff
	default:
		return nil
	}
}

func actionCardAbilityUsed(g *Gameplay) stateFn {
	fmt.Printf("state: %v\n", zb.PlayerActionType_CardAbilityUsed)
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

	// TODO: card ability

	// TODO: record history data

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
	case zb.PlayerActionType_CardAbilityUsed:
		return actionCardAbilityUsed
	case zb.PlayerActionType_OverlordSkillUsed:
		return actionOverloadSkillUsed
	case zb.PlayerActionType_LeaveMatch:
		return actionLeaveMatch
	case zb.PlayerActionType_RankBuff:
		return actionRankBuff
	default:
		return nil
	}
}

func actionOverloadSkillUsed(g *Gameplay) stateFn {
	fmt.Printf("state: %v\n", zb.PlayerActionType_OverlordSkillUsed)
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

	// TODO: overload skill

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
	case zb.PlayerActionType_CardAbilityUsed:
		return actionCardAbilityUsed
	case zb.PlayerActionType_OverlordSkillUsed:
		return actionOverloadSkillUsed
	case zb.PlayerActionType_LeaveMatch:
		return actionLeaveMatch
	case zb.PlayerActionType_RankBuff:
		return actionRankBuff
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
	case zb.PlayerActionType_CardAbilityUsed:
		return actionCardAbilityUsed
	case zb.PlayerActionType_OverlordSkillUsed:
		return actionOverloadSkillUsed
	case zb.PlayerActionType_LeaveMatch:
		return actionLeaveMatch
	case zb.PlayerActionType_RankBuff:
		return actionRankBuff
	default:
		return nil
	}
}

func actionLeaveMatch(g *Gameplay) stateFn {
	fmt.Printf("state: %v\n", zb.PlayerActionType_LeaveMatch)
	if g.isEnded() {
		return nil
	}
	// current action
	current := g.current()
	if current == nil {
		return nil
	}

	// update the winner of the game
	var winner string
	for _, player := range g.State.PlayerStates {
		if player.Id != current.PlayerId {
			winner = player.Id
		}
	}
	g.State.Winner = winner
	g.State.IsEnded = true

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
	case zb.PlayerActionType_CardAbilityUsed:
		return actionCardAbilityUsed
	case zb.PlayerActionType_OverlordSkillUsed:
		return actionOverloadSkillUsed
	case zb.PlayerActionType_RankBuff:
		return actionRankBuff
	default:
		return nil
	}
}

func actionRankBuff(g *Gameplay) stateFn {
	fmt.Printf("state: %v\n", zb.PlayerActionType_RankBuff)
	if g.isEnded() {
		return nil
	}
	// current action
	current := g.current()
	if current == nil {
		return nil
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
	case zb.PlayerActionType_CardAbilityUsed:
		return actionCardAbilityUsed
	case zb.PlayerActionType_OverlordSkillUsed:
		return actionOverloadSkillUsed
	case zb.PlayerActionType_RankBuff:
		return actionRankBuff
	default:
		return nil
	}
}

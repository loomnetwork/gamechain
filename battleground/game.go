package battleground

import (
	"bytes"
	"fmt"
	"math/rand"

	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	"github.com/pkg/errors"
)

const (
	defaultTurnTime      = 120
	defaultMulliganCards = 3
	maxMulliganCards     = 10
	maxCardsInPlay       = 6
	maxCardsInHand       = 10
	maxGooVials          = 10
)

var (
	errInvalidPlayer         = errors.New("invalid player")
	errCurrentActionNotfound = errors.New("current action not found")
	errInvalidAction         = errors.New("invalid action")
	errNotEnoughPlayer       = errors.New("not enough players")
	errAlreadyTossCoin       = errors.New("already tossed coin")
	errNoCurrentPlayer       = errors.New("no current player")
	errCardNotFoundInHand    = errors.New("card not found in hand")
	errLimitExceeded         = errors.New("max card limit exceeded")
	errNoCardsInHand         = errors.New("Can't play card. No cards in hand")
	errInsufficientGoo       = errors.New("insufficient goo")
)

type Gameplay struct {
	State                  *zb.GameState
	stateFn                stateFn
	cardLibrary            *zb.CardList
	err                    error
	customGameMode         *CustomGameMode
	history                []*zb.HistoryData
	ctx                    *contract.Context
	ClientSideRuleOverride bool         //disables all checks to ensure the client can work before server is fully implemented
	logger                 *loom.Logger // optional logger
}

type stateFn func(*Gameplay) stateFn

// NewGamePlay initializes GamePlay with default game state and run to the  latest state
func NewGamePlay(ctx contract.Context,
	id int64,
	version string,
	players []*zb.PlayerState,
	seed int64,
	customGameAddress *loom.Address,
	clientSideRuleOverride bool,
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
		CreatedAt:          ctx.Now().Unix(),
	}
	g := &Gameplay{
		State:                  state,
		customGameMode:         customGameMode,
		ctx:                    &ctx,
		ClientSideRuleOverride: clientSideRuleOverride,
		logger:                 ctx.Logger(),
	}

	var err error
	g.cardLibrary, err = getCardLibrary(ctx, version)
	if err != nil {
		return nil, err
	}

	err = populateDeckCards(ctx, g.cardLibrary, players)
	if err != nil {
		return nil, err
	}

	if err = g.createGame(ctx); err != nil {
		return nil, err
	}

	if err = g.run(); err != nil {
		return nil, err
	}
	return g, nil
}

// GamePlayFrom initializes and run game to the latest state
func GamePlayFrom(state *zb.GameState, override bool) (*Gameplay, error) {
	g := &Gameplay{State: state, ClientSideRuleOverride: override}
	if err := g.run(); err != nil {
		return nil, err
	}
	return g, nil
}

func (g *Gameplay) createGame(ctx contract.Context) error {
	// init players
	for i := 0; i < len(g.State.PlayerStates); i++ {
		g.State.PlayerStates[i].Defense = 20
		g.State.PlayerStates[i].CurrentGoo = 0
		g.State.PlayerStates[i].GooVials = 0
		g.State.PlayerStates[i].TurnTime = defaultTurnTime
		g.State.PlayerStates[i].InitialCardsInHandCount = defaultMulliganCards
		g.State.PlayerStates[i].MaxCardsInPlay = maxCardsInPlay
		g.State.PlayerStates[i].MaxCardsInHand = maxCardsInHand
		g.State.PlayerStates[i].MaxGooVials = maxGooVials
	}
	// coin toss for the first player
	r := rand.New(rand.NewSource(g.State.Randomseed))
	n := r.Int31n(int32(len(g.State.PlayerStates)))
	g.State.CurrentPlayerIndex = n

	if g.customGameMode != nil {
		err := g.customGameMode.CallHookBeforeMatchStart(ctx, g)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("Error in custom game mode -%v", err))
			return err
		}
	}

	// init hands
	for i := 0; i < len(g.State.PlayerStates); i++ {
		g.State.PlayerStates[i].CardsInDeck = shuffleCardInDeck(g.State.PlayerStates[i].CardsInDeck, g.State.Randomseed, i)
		// draw cards 3 card for mulligan
		g.State.PlayerStates[i].CardsInHand = g.State.PlayerStates[i].CardsInDeck[:g.State.PlayerStates[i].InitialCardsInHandCount]
		g.State.PlayerStates[i].CardsInDeck = g.State.PlayerStates[i].CardsInDeck[g.State.PlayerStates[i].InitialCardsInHandCount:]
	}

	if g.customGameMode != nil {
		err := g.customGameMode.CallHookAfterInitialDraw(ctx, g)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("Error in custom game mode -%v", err))
			return err
		}
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
	// skip checking for allowed actions
	if action.ActionType == zb.PlayerActionType_Mulligan || action.ActionType == zb.PlayerActionType_LeaveMatch {
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
	g.debugf("Gameplay stopped at action index %d, err=%v\n", g.State.CurrentActionIndex, g.err)
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

	g.debugf("Gameplay resumed at action index %d\n", g.State.CurrentActionIndex)
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
	if g.State.CurrentActionIndex < 0 {
		return nil
	}
	if g.State.CurrentActionIndex+1 > int64(len(g.State.PlayerActions)) {
		return nil
	}
	action := g.State.PlayerActions[g.State.CurrentActionIndex+1]
	return action
}

func (g *Gameplay) current() *zb.PlayerAction {
	if g.State.CurrentActionIndex < 0 {
		return nil
	}
	if g.State.CurrentActionIndex+1 > int64(len(g.State.PlayerActions)) {
		return nil
	}
	action := g.State.PlayerActions[g.State.CurrentActionIndex]
	return action
}

func (g *Gameplay) activePlayer() *zb.PlayerState {
	return g.State.PlayerStates[g.State.CurrentPlayerIndex]
}

func (g *Gameplay) activePlayerOpponent() *zb.PlayerState {
	for i, p := range g.State.PlayerStates {
		if int32(i) != g.State.CurrentPlayerIndex {
			return p
		}
	}

	return nil
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
		if player.Defense <= 0 {
			return true
		}
	}
	return false
}

func (g *Gameplay) validateGameState() error {
	for _, player := range g.State.PlayerStates {
		if player.MaxCardsInPlay < 1 || player.MaxCardsInPlay > maxCardsInPlay {
			return fmt.Errorf(
				"MaxCardsInPlay must be in range [%d;%d], current value %d",
				1,
				maxCardsInPlay,
				player.MaxCardsInPlay,
			)
		}

		if player.MaxCardsInHand < 1 || player.MaxCardsInHand > maxCardsInHand {
			return fmt.Errorf(
				"MaxCardsInHand must be in range [%d;%d], current value %d",
				1,
				maxCardsInHand,
				player.MaxCardsInHand,
			)
		}

		if player.GooVials < 1 || player.GooVials > maxGooVials {
			return fmt.Errorf(
				"GooVials must be in range [%d;%d], current value %d",
				1,
				maxGooVials,
				player.MaxGooVials,
			)
		}

		if player.InitialCardsInHandCount > maxMulliganCards {
			return fmt.Errorf(
				"InitialCardsInHandCount (%d) can't be larger than %d",
				player.InitialCardsInHandCount,
				maxMulliganCards,
			)
		}

		if player.InitialCardsInHandCount < 0 {
			return fmt.Errorf(
				"InitialCardsInHandCount (%d) can't be less than %d",
				player.InitialCardsInHandCount,
				0,
			)
		}

		if player.InitialCardsInHandCount > player.MaxCardsInHand {
			return fmt.Errorf(
				"InitialCardsInHandCount (%d) can't be larger than MaxCardsInHand (%d)",
				player.InitialCardsInHandCount,
				player.MaxCardsInHand,
			)
		}

		if player.TurnTime < 0 {
			return fmt.Errorf(
				"TurnTime must be larger than %d, current value %d",
				0,
				player.TurnTime,
			)
		}
	}

	return nil
}

func (g *Gameplay) SetLogger(logger *loom.Logger) {
	g.logger = logger
}

func (g *Gameplay) debugf(msg string, values ...interface{}) {
	if g.logger == nil {
		return
	}
	g.logger.Debug(fmt.Sprintf(msg, values...))
}

func (g *Gameplay) PrintState() {
	if g.logger == nil {
		return
	}
	state := g.State
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "============StateInfo=============\n")
	fmt.Fprintf(buf, "Is ended: %v, Winner: %s\n", state.IsEnded, state.Winner)
	fmt.Fprintf(buf, "Current Player Index: %v\n", state.CurrentPlayerIndex)

	for i, player := range g.State.PlayerStates {
		if g.State.CurrentPlayerIndex == int32(i) {
			fmt.Fprintf(buf, "Player%d: %s 🧟\n", i+1, player.Id)
		} else {
			fmt.Fprintf(buf, "Player%d: %s\n", i+1, player.Id)
		}
		fmt.Fprintf(buf, "\tdefense: %v\n", player.Defense)
		fmt.Fprintf(buf, "\tcurrent goo: %v\n", player.CurrentGoo)
		fmt.Fprintf(buf, "\tgoo vials: %v\n", player.GooVials)
		fmt.Fprintf(buf, "\thas drawn card: %v\n", player.HasDrawnCard)
		fmt.Fprintf(buf, "\tcard in hand (%d): %v\n", len(player.CardsInHand), player.CardsInHand)
		fmt.Fprintf(buf, "\tcard in play (%d): %v\n", len(player.CardsInPlay), player.CardsInPlay)
		fmt.Fprintf(buf, "\tcard in deck (%d): %v\n", len(player.CardsInDeck), player.CardsInDeck)
		fmt.Fprintf(buf, "\tcard in graveyard (%d): %v\n", len(player.CardsInGraveyard), player.CardsInGraveyard)
		fmt.Fprintf(buf, "\n") // extra line
	}

	fmt.Fprintf(buf, "History : count %v\n", len(g.history))
	for i, block := range g.history {
		fmt.Fprintf(buf, "\t[%d] %v\n", i, block)
	}

	fmt.Fprintf(buf, "Actions: count %v\n", len(state.PlayerActions))
	for i, action := range state.PlayerActions {
		if int64(i) == state.CurrentActionIndex {
			fmt.Fprintf(buf, "   -->> [%d] %v\n", i, action)
		} else {
			fmt.Fprintf(buf, "\t[%d] %v\n", i, action)
		}
	}
	fmt.Fprintf(buf, "Current Action Index: %v\n", state.CurrentActionIndex)
	fmt.Fprintf(buf, "==================================\n")
	g.debugf(buf.String())
}

func gameStart(g *Gameplay) stateFn {
	g.debugf("state: gameStart\n")
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
	g.debugf("state: %v\n", zb.PlayerActionType_Mulligan)
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
	if len(mulligan.MulliganedCards) > int(player.InitialCardsInHandCount) {
		return g.captureErrorAndStop(fmt.Errorf("number of mulligan card is exceed the maximum: %d", player.InitialCardsInHandCount))
	}
	for _, card := range mulligan.MulliganedCards {
		_, _, found := findCardInCardList(card, player.CardsInHand)
		if !found {
			return g.captureErrorAndStop(fmt.Errorf("invalid mulligan card"))
		}
	}

	// keep only the cards in in mulligan
	keepCards := make([]*zb.CardInstance, 0)
	for _, mcard := range mulligan.MulliganedCards {
		_, card, found := findCardInCardList(mcard, player.CardsInHand)
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
			_, _, found := findCardInCardList(card, keepCards)
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
	g.debugf("state: %v\n", zb.PlayerActionType_DrawCard)
	if g.isEnded() {
		return nil
	}
	current := g.current()
	if current == nil {
		return nil
	}

	// check player turn
	// TODO: for now, skip player check so we dont break client logic at match start due to mulligans
	// if err := g.checkCurrentPlayer(current); err != nil {
	// 	return g.captureErrorAndStop(err)
	// }

	// check if player has already drawn a card after starting new turn
	// if g.activePlayer().HasDrawnCard {
	// 	g.err = errInvalidAction
	// 	return nil
	// }

	// draw card
	if len(g.activePlayer().CardsInDeck) < 1 {
		return g.captureErrorAndStop(errors.New("Can't draw card. No more cards in deck"))
	}

	// handle card limit in hand
	if len(g.activePlayer().CardsInHand)+1 > int(g.activePlayer().MaxCardsInHand) {
		// TODO: assgin g.err
		return nil
	}

	// TODO: for now we just trust the client and draw the card it tells us
	// card := g.activePlayer().CardsInDeck[0]
	// if card.InstanceId != current.GetDrawCard().CardInstance.InstanceId {
	// 	return g.captureErrorAndStop(errors.New("Client drew a card but server could not verify it"))
	// }

	var cardIndexInDeck int
	var card *zb.CardInstance
	for i, cardInDeck := range g.activePlayer().CardsInDeck {
		if cardInDeck.InstanceId == current.GetDrawCard().CardInstance.InstanceId {
			card = cardInDeck
			cardIndexInDeck = i
			g.activePlayer().CardsInHand = append(g.activePlayer().CardsInHand, cardInDeck)
			break
		}
	}

	if card == nil {
		return g.captureErrorAndStop(errors.New("Can't draw card. Card not found in deck"))
	}

	// remove card from CardsInDeck
	// g.activePlayer().CardsInDeck = g.activePlayer().CardsInDeck[1:]
	g.activePlayer().CardsInDeck = append(g.activePlayer().CardsInDeck[:cardIndexInDeck], g.activePlayer().CardsInDeck[cardIndexInDeck+1:]...)

	// card drawn, don't allow another draw until next turn
	g.activePlayer().HasDrawnCard = true

	// record history data
	g.history = append(g.history, &zb.HistoryData{
		Data: &zb.HistoryData_FullInstance{
			FullInstance: &zb.HistoryFullInstance{
				InstanceId: card.InstanceId,
				Attack:     card.Attack,
				Defense:    card.Defense,
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
	g.debugf("state: %v\n", zb.PlayerActionType_CardPlay)
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

	cardPlay := current.GetCardPlay()

	if !g.ClientSideRuleOverride {
		card := cardPlay.Card

		// check card limit on board
		if len(g.activePlayer().CardsInPlay)+1 > int(g.activePlayer().MaxCardsInPlay) {
			if g.ClientSideRuleOverride {
				g.debugf("ClientSideRuleOverride-" + errLimitExceeded.Error())
			} else {
				return g.captureErrorAndStop(errLimitExceeded)
			}
		}

		activeCardsInHand := g.activePlayer().CardsInHand
		// TODO: handle card limit
		if len(activeCardsInHand) == 0 {
			if g.ClientSideRuleOverride {
				g.debugf("ClientSideRuleOverride-" + errNoCardsInHand.Error())
			} else {
				return g.captureErrorAndStop(errNoCardsInHand)
			}
		}

		// get card instance from cardsInHand list
		cardIndex, card, found := findCardInCardListInstanceID(cardPlay.Card, activeCardsInHand)
		if !found {
			if g.ClientSideRuleOverride {
				g.debugf("ClientSideRuleOverride-" + errCardNotFoundInHand.Error())
				next := g.next()
				if next == nil {
					return nil
				}
			} else {
				return g.captureErrorAndStop(errCardNotFoundInHand)
			}
		}

		// put card on board
		g.activePlayer().CardsInPlay = append(g.activePlayer().CardsInPlay, card)
		// remove card from hand
		activeCardsInHand = append(activeCardsInHand[:cardIndex], activeCardsInHand[cardIndex+1:]...)
		g.activePlayer().CardsInHand = activeCardsInHand

		// record history data
		g.history = append(g.history, &zb.HistoryData{
			Data: &zb.HistoryData_FullInstance{
				FullInstance: &zb.HistoryFullInstance{
					InstanceId: card.InstanceId,
					Attack:     card.Attack,
					Defense:    card.Defense,
				},
			},
		})
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
	case zb.PlayerActionType_LeaveMatch:
		return actionLeaveMatch
	case zb.PlayerActionType_RankBuff:
		return actionRankBuff
	default:
		return nil
	}
}

func actionCardAttack(g *Gameplay) stateFn {
	g.debugf("state: %v\n", zb.PlayerActionType_CardAttack)
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

	var attacker *zb.CardInstance
	var target *zb.CardInstance
	var attackerIndex int
	var targetIndex int

	switch current.GetCardAttack().AffectObjectType {
	case zb.AffectObjectType_CHARACTER:
		if len(g.activePlayer().CardsInPlay) <= 0 {
			if g.ClientSideRuleOverride {
				g.debugf("No cards on board to attack with")
				g.PrintState()
				next := g.next()
				if next == nil {
					return nil
				}
			} else {
				return g.captureErrorAndStop(errors.New("No cards on board to attack with"))
			}
		}
		if len(g.activePlayerOpponent().CardsInPlay) <= 0 {
			if g.ClientSideRuleOverride {
				g.debugf("No cards on board to attack with")
				g.PrintState()
				next := g.next()
				if next == nil {
					return nil
				}
			} else {
				return g.captureErrorAndStop(errors.New("No cards on board to attack"))
			}
		}
		for i, card := range g.activePlayer().CardsInPlay {
			if card.InstanceId == current.GetCardAttack().Attacker.InstanceId {
				attacker = card
				attackerIndex = i
				break
			}
			if g.ClientSideRuleOverride {
				g.debugf("zb.AffectObjectType_CHARACTER-Attacker not found\n")
				g.PrintState()
				next := g.next()
				if next == nil {
					return nil
				}
			} else {
				return g.captureErrorAndStop(errors.New("Attacker not found"))
			}
		}

		for i, card := range g.activePlayerOpponent().CardsInPlay {
			if card.InstanceId == current.GetCardAttack().Target.InstanceId {
				target = card
				targetIndex = i
				break
			}
			return g.captureErrorAndStop(errors.New("Target not found"))
		}

		attacker.Defense -= target.Attack
		target.Defense -= attacker.Attack

		if attacker.Defense <= 0 {
			g.activePlayer().CardsInPlay = append(g.activePlayer().CardsInPlay[:attackerIndex], g.activePlayer().CardsInPlay[attackerIndex+1:]...)
			g.activePlayer().CardsInGraveyard = append(g.activePlayer().CardsInGraveyard, attacker)
		}
		if target.Defense <= 0 {
			g.activePlayerOpponent().CardsInPlay = append(g.activePlayerOpponent().CardsInPlay[:targetIndex], g.activePlayerOpponent().CardsInPlay[targetIndex+1:]...)
			g.activePlayerOpponent().CardsInGraveyard = append(g.activePlayerOpponent().CardsInGraveyard, target)
		}

	case zb.AffectObjectType_PLAYER:
		if len(g.activePlayer().CardsInPlay) <= 0 {
			if g.ClientSideRuleOverride {
				g.debugf("No cards on board to attack with")
				g.PrintState()
				next := g.next()
				if next == nil {
					return nil
				}
			} else {
				return g.captureErrorAndStop(errors.New("No cards on board to attack with"))
			}
		}
		for i, card := range g.activePlayer().CardsInPlay {
			if card.InstanceId == current.GetCardAttack().Attacker.InstanceId {
				attacker = card
				attackerIndex = i
				break
			}
			if g.ClientSideRuleOverride {
				g.debugf("zb.AffectObjectType_PLAYER:-Attacker not found\n")
				g.PrintState()
				next := g.next()
				if next == nil {
					return nil
				}
			} else {
				return g.captureErrorAndStop(errors.New("Attacker not found"))
			}
		}

		g.activePlayerOpponent().Defense -= attacker.Attack

		if g.activePlayerOpponent().Defense <= 0 {
			g.State.Winner = g.activePlayer().Id
			g.State.IsEnded = true
		}
	}

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
	g.debugf("state: %v\n", zb.PlayerActionType_CardAbilityUsed)
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
	g.debugf("state: %v\n", zb.PlayerActionType_OverlordSkillUsed)
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
	g.debugf("state: %v\n", zb.PlayerActionType_EndTurn)
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

	// allow the new player to draw card on new turn
	g.activePlayer().HasDrawnCard = false

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
	g.debugf("state: %v\n", zb.PlayerActionType_LeaveMatch)
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
	g.debugf("state: %v\n", zb.PlayerActionType_RankBuff)
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

package battleground

import (
	"bytes"
	"fmt"
	"github.com/loomnetwork/gamechain/types/zb/zb_data"
	"github.com/loomnetwork/gamechain/types/zb/zb_enums"
	"math/rand"
	"sort"

	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/gamechain/types/zb"
	loom "github.com/loomnetwork/go-loom"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	"github.com/pkg/errors"
)

const (
	defaultTurnTime        = 120
	defaultMulliganCards   = 3
	defaultOverlordDefense = 50
	maxMulliganCards       = 10
	maxCardsInPlay         = 6
	maxCardsInHand         = 10
	maxGooVials            = 10
)

var (
	errInvalidPlayer         = errors.New("invalid player")
	errCurrentActionNotfound = errors.New("current action not found")
	errInvalidAction         = errors.New("invalid action")
	errNotEnoughPlayer       = errors.New("not enough players")
	errAlreadyTossCoin       = errors.New("already tossed coin")
	errNoCurrentPlayer       = errors.New("no current player")
	errLimitExceeded         = errors.New("max card limit exceeded")
	errNoCardsInHand         = errors.New("Can't play card. No cards in hand")
	errInsufficientGoo       = errors.New("insufficient goo")
	errCheatsRequired        = errors.New("cheats are required for this action")
)

type Gameplay struct {
	State               *zb_data.GameState
	stateFn             stateFn
	cardLibrary         *zb_data.CardList
	err                 error
	customGameMode      *CustomGameMode
	history             []*zb_data.HistoryData
	ctx                 *contract.Context
	useBackendGameLogic bool // when false, disables all checks to ensure the client can work before server is fully implemented
	actionOutcomes      []*zb_data.PlayerActionOutcome
	playersDebugCheats  []*zb_data.DebugCheatsConfiguration
	logger              *loom.Logger // optional logger
}

type stateFn func(*Gameplay) stateFn

// NewGamePlay initializes GamePlay with default game state and run to the  latest state
func NewGamePlay(ctx contract.Context,
	id int64,
	version string,
	players []*zb_data.PlayerState,
	seed int64,
	customGameAddress *loom.Address,
	useBackendGameLogic bool,
	playersDebugCheats []*zb_data.DebugCheatsConfiguration,
) (*Gameplay, error) {
	var customGameMode *CustomGameMode
	if customGameAddress != nil {
		ctx.Logger().Info(fmt.Sprintf("Playing a custom game mode -%v", customGameAddress.String()))
		customGameMode = NewCustomGameMode(*customGameAddress)
	}

	// So we won't have to do nil checks everywhere along the way
	if playersDebugCheats == nil {
		playersDebugCheats = []*zb_data.DebugCheatsConfiguration{{}, {}}
	}

	// Ensure that same random seed will result in the same player order,
	// no matter which player joined the pool earlier
	type playerDataTuple struct {
		playerState       *zb_data.PlayerState
		playerDebugCheats *zb_data.DebugCheatsConfiguration
	}

	playersData := make([]*playerDataTuple, len(players), len(players))
	for i, player := range players {
		playersData[i] = &playerDataTuple{}
		playersData[i].playerState = player
		playersData[i].playerDebugCheats = playersDebugCheats[i]
	}

	sort.SliceStable(playersData, func(i, j int) bool {
		return playersData[i].playerState.Id < playersData[j].playerState.Id
	})

	for i, playerData := range playersData {
		playerData.playerState.Index = int32(i)
		players[i] = playerData.playerState
		playersDebugCheats[i] = playerData.playerDebugCheats
	}

	state := &zb_data.GameState{
		Id:                 id,
		CurrentActionIndex: -1, // use -1 to avoid confict with default value
		PlayerStates:       players,
		CurrentPlayerIndex: -1, // use -1 to avoid confict with default value
		RandomSeed:         seed,
		Version:            version,
		CreatedAt:          ctx.Now().Unix(),
	}
	g := &Gameplay{
		State:               state,
		customGameMode:      customGameMode,
		ctx:                 &ctx,
		useBackendGameLogic: useBackendGameLogic,
		logger:              ctx.Logger(),
		playersDebugCheats:  playersDebugCheats,
	}

	var err error
	g.cardLibrary, err = getCardLibrary(ctx, version)
	if err != nil {
		return nil, err
	}

	err = populateDeckCards(g.cardLibrary, players, useBackendGameLogic)
	if err != nil {
		return nil, err
	}

	if err = g.createGame(ctx); err != nil {
		return nil, err
	}

	if err := saveInitialGameState(ctx, g.State); err != nil {
		return nil, err
	}

	if err = g.run(); err != nil {
		return nil, err
	}
	return g, nil
}

// GamePlayFrom initializes and run game to the latest state
func GamePlayFrom(state *zb_data.GameState, useBackendGameLogic bool, playersDebugCheats []*zb_data.DebugCheatsConfiguration) (*Gameplay, error) {
	g := &Gameplay{
		State:               state,
		useBackendGameLogic: useBackendGameLogic,
		playersDebugCheats:  playersDebugCheats,
	}
	if err := g.run(); err != nil {
		return nil, err
	}
	return g, nil
}

func (g *Gameplay) createGame(ctx contract.Context) error {
	gamechainState, err := loadState(ctx)
	if err != nil {
		return err
	}

	defaultDefense := defaultOverlordDefense
	if gamechainState.DefaultPlayerDefense > 0 {
		defaultDefense = int(gamechainState.DefaultPlayerDefense)
	}

	// init players
	for i := 0; i < len(g.State.PlayerStates); i++ {
		g.State.PlayerStates[i].Defense = int32(defaultDefense)
		g.State.PlayerStates[i].CurrentGoo = 0
		g.State.PlayerStates[i].GooVials = 0
		g.State.PlayerStates[i].TurnTime = defaultTurnTime
		g.State.PlayerStates[i].InitialCardsInHandCount = defaultMulliganCards
		g.State.PlayerStates[i].MaxCardsInPlay = maxCardsInPlay
		g.State.PlayerStates[i].MaxCardsInHand = maxCardsInHand
		g.State.PlayerStates[i].MaxGooVials = maxGooVials
	}
	// coin toss for the first player
	r := rand.New(rand.NewSource(g.State.RandomSeed))
	n := r.Int31n(int32(len(g.State.PlayerStates)))
	g.State.CurrentPlayerIndex = n

	// force first player cheat
loop:
	for i := 0; i < len(g.State.PlayerStates); i++ {
		for j := 0; j < len(g.State.PlayerStates); j++ {
			if g.playersDebugCheats[j].Enabled && g.playersDebugCheats[j].ForceFirstTurnUserId != "" && g.playersDebugCheats[j].ForceFirstTurnUserId == g.State.PlayerStates[i].Id {
				g.State.CurrentPlayerIndex = int32(i)
				break loop
			}
		}
	}

	// custom mode pre-match hook
	if g.customGameMode != nil {
		err := g.customGameMode.CallHookBeforeMatchStart(ctx, g)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("Error in custom game mode -%v", err))
			return err
		}
	}

	// init cards
	for i := 0; i < len(g.State.PlayerStates); i++ {
		playerState := g.State.PlayerStates[i]
		if !(g.playersDebugCheats[i].Enabled && g.playersDebugCheats[i].DisableDeckShuffle) {
			playerState.CardsInDeck = shuffleCardInDeck(playerState.CardsInDeck, g.State.RandomSeed, i)
		}

		// draw cards 3 card for mulligan
		// HOTFIX: TODO: Check this again
		if len(playerState.CardsInDeck) > int(playerState.InitialCardsInHandCount) {
			playerState.CardsInHand = playerState.CardsInDeck[:playerState.InitialCardsInHandCount]
			playerState.CardsInDeck = playerState.CardsInDeck[playerState.InitialCardsInHandCount:]
			for i := 0; i < len(playerState.CardsInHand); i++ {
				playerState.CardsInHand[i].Zone = zb_enums.Zone_HAND
			}
			for i := 0; i < len(playerState.CardsInDeck); i++ {
				playerState.CardsInDeck[i].Zone = zb_enums.Zone_DECK
			}
		}
	}

	// init instance IDs
	// 0 and 1 are reserved for overlords
	// ID 0 is the overlord of the player that has the first turn
	// ID 1 is the overlord of the other player that has the first turn
	// Card ID's start with the player that has the first turn
	assignInstanceIds := func(playerState *zb_data.PlayerState, currentInstanceId *int32) {
		for _, card := range playerState.CardsInPlay {
			card.InstanceId = &zb_data.InstanceId{Id: *currentInstanceId}
			*currentInstanceId++
		}

		for _, card := range playerState.CardsInHand {
			card.InstanceId = &zb_data.InstanceId{Id: *currentInstanceId}
			*currentInstanceId++
		}

		for _, card := range playerState.CardsInDeck {
			card.InstanceId = &zb_data.InstanceId{Id: *currentInstanceId}
			*currentInstanceId++
		}

		for _, card := range playerState.CardsInGraveyard {
			card.InstanceId = &zb_data.InstanceId{Id: *currentInstanceId}
			*currentInstanceId++
		}
	}
	var instanceId int32 = 2
	if g.State.CurrentPlayerIndex == 0 {
		g.State.PlayerStates[0].InstanceId = &zb_data.InstanceId{Id: 0}
		g.State.PlayerStates[1].InstanceId = &zb_data.InstanceId{Id: 1}
		assignInstanceIds(g.State.PlayerStates[0], &instanceId)
		assignInstanceIds(g.State.PlayerStates[1], &instanceId)
	} else {
		g.State.PlayerStates[0].InstanceId = &zb_data.InstanceId{Id: 1}
		g.State.PlayerStates[1].InstanceId = &zb_data.InstanceId{Id: 0}
		assignInstanceIds(g.State.PlayerStates[1], &instanceId)
		assignInstanceIds(g.State.PlayerStates[0], &instanceId)
	}

	g.State.NextInstanceId = instanceId

	if g.customGameMode != nil {
		err := g.customGameMode.CallHookAfterInitialDraw(ctx, g)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("Error in custom game mode -%v", err))
			return err
		}
	}

	// first player draws a card immediately
	if err := g.drawCard(g.activePlayer(), 1); err != nil {
		return err
	}

	// give initial 1 vial and 1 goo
	addGooVialAndFillAll(g.activePlayer())
	//addGooVialAndFillAll(g.activePlayerOpponent())

	// add history data
	ps := make([]*zb.Player, len(g.State.PlayerStates))
	for i := range g.State.PlayerStates {
		ps[i] = &zb.Player{
			Id:   g.State.PlayerStates[i].Id,
			Deck: g.State.PlayerStates[i].Deck,
		}
	}
	// record history data
	g.history = append(g.history, &zb_data.HistoryData{
		Data: &zb_data.HistoryData_CreateGame{
			CreateGame: &zb.HistoryCreateGame{
				GameId:     g.State.Id,
				Players:    ps,
				RandomSeed: g.State.RandomSeed,
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
	if action.ActionType == zb.PlayerActionType_Mulligan ||
		action.ActionType == zb.PlayerActionType_LeaveMatch ||
		action.ActionType == zb.PlayerActionType_CheatDestroyCardsOnBoard {
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
	case zb.PlayerActionType_CheatDestroyCardsOnBoard:
		state = actionCheatDestroyCardsOnBoard
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

func (g *Gameplay) activePlayer() *zb_data.PlayerState {
	return g.State.PlayerStates[g.State.CurrentPlayerIndex]
}

func (g *Gameplay) activePlayerDebugCheats() *zb_data.DebugCheatsConfiguration {
	return g.playersDebugCheats[g.State.CurrentPlayerIndex]
}

func (g *Gameplay) activePlayerOpponent() *zb_data.PlayerState {
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

// gives the player a new goo vial and fills up all their vials
func addGooVialAndFillAll(playerState *zb_data.PlayerState) {
	if playerState.GooVials < playerState.MaxGooVials {
		playerState.GooVials++
	}
	playerState.CurrentGoo = playerState.GooVials
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
	g.logger.Info(fmt.Sprintf(msg, values...))
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
	fmt.Fprintf(buf, "Ability Outcomes:\n")
	for i, outcome := range g.actionOutcomes {
		fmt.Fprintf(buf, "\t[%d] %v\n", i, outcome)
	}
	fmt.Fprintf(buf, "==================================\n")
}

func (g *Gameplay) DebugState() {
	state := g.State
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "============StateInfo=============\n")
	fmt.Fprintf(buf, "Is ended: %v, Winner: %s\n", state.IsEnded, state.Winner)
	fmt.Fprintf(buf, "Current Player Index: %v\n", state.CurrentPlayerIndex)

	formatAbility := func(abilities []*zb_data.CardAbilityInstance) string {
		b := new(bytes.Buffer)
		for _, a := range abilities {
			b.WriteString(fmt.Sprintf("[%v, active=%v]\n", a.AbilityType, a.IsActive))
		}
		return b.String()
	}

	for i, player := range g.State.PlayerStates {
		if g.State.CurrentPlayerIndex == int32(i) {
			fmt.Fprintf(buf, "Player%d: %s 🧟\n", i+1, player.Id)
		} else {
			fmt.Fprintf(buf, "Player%d: %s\n", i+1, player.Id)
		}
		fmt.Fprintf(buf, "\tstats:\n")
		fmt.Fprintf(buf, "\t\tdefense: %v\n", player.Defense)
		fmt.Fprintf(buf, "\t\tcurrent goo: %v\n", player.CurrentGoo)
		fmt.Fprintf(buf, "\t\tgoo vials: %v\n", player.GooVials)
		fmt.Fprintf(buf, "\t\thas drawn card: %v\n", player.HasDrawnCard)
		fmt.Fprintf(buf, "\tmulligan (%d):\n", len(player.MulliganCards))
		for _, card := range player.MulliganCards {
			fmt.Fprintf(buf, "\t\tName:%s\n", card.Prototype.Name)
		}
		fmt.Fprintf(buf, "\tcard in hand (%d):\n", len(player.CardsInHand))
		for _, card := range player.CardsInHand {
			fmt.Fprintf(buf, "\t\tId:%-2d Name:%-14s Dmg:%2d Def:%2d Goo:%2d, Zone:%0v, OwnerIndex:%d %s\n", card.InstanceId.Id, card.Prototype.Name, card.Instance.Damage, card.Instance.Defense, card.Instance.Cost, card.Zone, card.OwnerIndex, formatAbility(card.AbilitiesInstances))
		}
		fmt.Fprintf(buf, "\tcard in play (%d):\n", len(player.CardsInPlay))
		for _, card := range player.CardsInPlay {
			fmt.Fprintf(buf, "\t\tId:%-2d Name:%-14s Dmg:%2d Def:%2d Goo:%2d, Zone:%0v, OwnerIndex:%d %s\n", card.InstanceId.Id, card.Prototype.Name, card.Instance.Damage, card.Instance.Defense, card.Instance.Cost, card.Zone, card.OwnerIndex, formatAbility(card.AbilitiesInstances))
		}
		fmt.Fprintf(buf, "\tcard in deck (%d):\n", len(player.CardsInDeck))
		for _, card := range player.CardsInDeck {
			fmt.Fprintf(buf, "\t\tId:%-2d Name:%-14s Dmg:%2d Def:%2d Goo:%2d, Zone:%0v, OwnerIndex:%d %s\n", card.InstanceId.Id, card.Prototype.Name, card.Instance.Damage, card.Instance.Defense, card.Instance.Cost, card.Zone, card.OwnerIndex, formatAbility(card.AbilitiesInstances))
		}
		fmt.Fprintf(buf, "\tcard in graveyard (%d):\n", len(player.CardsInGraveyard))
		for _, card := range player.CardsInGraveyard {
			fmt.Fprintf(buf, "\t\tId:%-2d Name:%-14s Dmg:%2d Def:%2d Goo:%2d, Zone:%0v, OwnerIndex:%d %s\n", card.InstanceId.Id, card.Prototype.Name, card.Instance.Damage, card.Instance.Defense, card.Instance.Cost, card.Zone, card.OwnerIndex, formatAbility(card.AbilitiesInstances))
		}
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

	fmt.Fprintf(buf, "Ability Outcomes:\n")
	for i, outcome := range g.actionOutcomes {
		fmt.Fprintf(buf, "\t[%d] %v\n", i, outcome)
	}

	fmt.Fprintf(buf, "==================================\n")
	fmt.Println(buf.String())
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
	case zb.PlayerActionType_CheatDestroyCardsOnBoard:
		return actionCheatDestroyCardsOnBoard
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

	if g.useBackendGameLogic {
		mulligan := current.GetMulligan()
		if mulligan == nil {
			return g.captureErrorAndStop(fmt.Errorf("expect mulligan action"))
		}
		var player *zb_data.PlayerState
		var playerIndex int
		for i := 0; i < len(g.State.PlayerStates); i++ {
			if g.State.PlayerStates[i].Id == current.PlayerId {
				player = g.State.PlayerStates[i]
				playerIndex = i
			}
		}
		if player == nil {
			return g.captureErrorAndStop(fmt.Errorf("player not found"))
		}

		if player.TurnNumber > 0 {
			return g.captureErrorAndStop(fmt.Errorf("Mulligan not allowed after game has started"))
		}

		// Check if all the mulliganed cards and number of card that can be mulligan
		if len(mulligan.MulliganedCards) > int(player.InitialCardsInHandCount) {
			return g.captureErrorAndStop(fmt.Errorf("number of mulligan card is exceed the maximum: %d", player.InitialCardsInHandCount))
		}
		mulliganCards := make([]*zb_data.CardInstance, 0)
		for _, card := range mulligan.MulliganedCards {
			handCards := player.CardsInHand[:player.InitialCardsInHandCount]
			_, mulliganCard, found := findCardInCardListByInstanceId(card, handCards)
			if !found {
				return g.captureErrorAndStop(fmt.Errorf("invalid mulligan card"))
			}
			mulliganCards = append(mulliganCards, mulliganCard)
		}

		// draw card to replace the reroll cards
		for i := 0; i < len(mulliganCards); i++ {
			// move card from hand to deck
			cardInstance := NewCardInstance(mulliganCards[i], g)
			if err := cardInstance.Mulligan(); err != nil {
				return g.captureErrorAndStop(err)
			}
		}

		// re-shuffle cards in deck if player mulligan more than one card
		if len(mulliganCards) > 0 {
			shuffleCardInDeck(player.CardsInDeck, g.State.RandomSeed, playerIndex)
		}
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
	case zb.PlayerActionType_CheatDestroyCardsOnBoard:
		return actionCheatDestroyCardsOnBoard
	default:
		return nil
	}
}

func (g *Gameplay) drawCard(player *zb_data.PlayerState, count int) error {
	if g.useBackendGameLogic {
		// check if player has already drawn a card after starting new turn
		if player.HasDrawnCard {
			g.err = errInvalidAction
			return nil
		}

		for i := 0; i < count; i++ {
			// draw card
			if len(player.CardsInDeck) < 1 {
				break
			}

			// handle card limit in hand
			if len(player.CardsInHand)+1 > int(player.MaxCardsInHand) {
				// TODO: assgin g.err
				return nil
			}

			card := player.CardsInDeck[0]
			cardInstance := NewCardInstance(card, g)
			cardInstance.MoveZone(zb_enums.Zone_DECK, zb_enums.Zone_HAND)
		}
	} else {
		// do nothing, client currently doesn't care about this at all
	}

	// card drawn, don't allow another draw until next turn
	player.HasDrawnCard = true

	return nil
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

	if g.useBackendGameLogic {
		card := cardPlay.Card

		// check card limit on board
		if len(g.activePlayer().CardsInPlay)+1 > int(g.activePlayer().MaxCardsInPlay) {
			return g.captureErrorAndStop(errLimitExceeded)
		}

		activeCardsInHand := g.activePlayer().CardsInHand
		// TODO: handle card limit
		if len(activeCardsInHand) == 0 {
			return g.captureErrorAndStop(errNoCardsInHand)
		}

		// get card instance from cardsInHand list
		_, cardInstance, found := findCardInCardListByInstanceId(cardPlay.Card, activeCardsInHand)
		if !found {
			err := fmt.Errorf(
				"card (instance id: %d) not found in hand",
				cardPlay.Card.Id,
			)
			return g.captureErrorAndStop(err)
		}

		// check goo cost
		if !(g.activePlayerDebugCheats().Enabled && g.activePlayerDebugCheats().IgnoreGooRequirements) {
			if cardInstance.Instance.Cost > g.activePlayer().CurrentGoo {
				err := fmt.Errorf("Not enough goo to play card with instanceId %d", cardPlay.Card.Id)
				return g.captureErrorAndStop(err)
			}

			// change player goo
			// TODO: abilities that change goo vials, overflow etc
			g.activePlayer().CurrentGoo -= cardInstance.Instance.Cost
		}

		instance := NewCardInstance(cardInstance, g)
		if err := instance.Play(); err != nil {
			return g.captureErrorAndStop(err)
		}

		// record history data
		g.history = append(g.history, &zb_data.HistoryData{
			Data: &zb_data.HistoryData_FullInstance{
				FullInstance: &zb.HistoryFullInstance{
					InstanceId: card,
					Damage:     cardInstance.Instance.Damage,
					Defense:    cardInstance.Instance.Defense,
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
	case zb.PlayerActionType_CheatDestroyCardsOnBoard:
		return actionCheatDestroyCardsOnBoard
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

	if g.useBackendGameLogic {
		if len(g.activePlayer().CardsInPlay) <= 0 {
			return g.captureErrorAndStop(errors.New("No cards on board to attack with"))
		}
		cardAttack := current.GetCardAttack()
		if cardAttack == nil {
			return g.captureErrorAndStop(errors.New("No card attack speficied"))
		}

		var attacker *zb_data.CardInstance
		for _, card := range g.activePlayer().CardsInPlay {
			if proto.Equal(card.InstanceId, cardAttack.Attacker) {
				attacker = card
				break
			}
		}

		if attacker == nil {
			return g.captureErrorAndStop(errors.New("Attacker not found"))
		}

		targetInstanceID := cardAttack.Target.InstanceId.Id
		// instance id 0 and 1 are reserved for overlord
		if targetInstanceID == 0 || targetInstanceID == 1 {
			if g.activePlayer().InstanceId.Id == targetInstanceID {
				return g.captureErrorAndStop(errors.New("Can't attack own overlord"))
			}
			attackerInstance := NewCardInstance(attacker, g)
			attackerInstance.AttackOverlord(g.activePlayerOpponent(), g.activePlayer())
		} else {
			// attack card
			if len(g.activePlayerOpponent().CardsInPlay) <= 0 {
				return g.captureErrorAndStop(errors.New("No cards on board to attack"))

			}
			var target *zb_data.CardInstance
			for _, card := range g.activePlayerOpponent().CardsInPlay {
				if proto.Equal(card.InstanceId, current.GetCardAttack().Target.InstanceId) {
					target = card
					break
				}
			}
			if target == nil {
				return g.captureErrorAndStop(errors.New("Target not found"))
			}

			g.debugf(
				"card {instanceId: %d, name: %s} attacking card {instanceId: %d, name: %s}",
				attacker.InstanceId,
				attacker.Prototype.Name,
				target.InstanceId,
				target.Prototype.Name,
			)

			attackerInstance := NewCardInstance(attacker, g)
			targetInstance := NewCardInstance(target, g)
			err := attackerInstance.Attack(targetInstance)
			if err != nil {
				return g.captureErrorAndStop(err)
			}
		}
	}

	// record history data
	g.history = append(g.history, &zb_data.HistoryData{
		Data: &zb_data.HistoryData_ChangeInstance{
			ChangeInstance: &zb.HistoryInstance{
				InstanceId: &zb_data.InstanceId{Id: 1}, // TODO change to the actual card id
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
	case zb.PlayerActionType_CheatDestroyCardsOnBoard:
		return actionCheatDestroyCardsOnBoard
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

	if g.useBackendGameLogic {
		// TODO: Fix me
		cardAbilityUsed := current.GetCardAbilityUsed()
		if cardAbilityUsed == nil {
			return g.captureErrorAndStop(fmt.Errorf("no card ability used specified"))
		}
		card := cardAbilityUsed.Card
		if card == nil {
			return g.captureErrorAndStop(fmt.Errorf("no card in card ability used"))
		}

		activeCards := g.activePlayer().CardsInPlay

		// Because the game client sends cardabilityused before sending cardplay, we need to do this
		// but once the game client is fixed, this line needs to be removed
		activeCards = append(activeCards, g.activePlayer().CardsInHand...)

		// get card instance from cardsInPlay list
		_, cardInstance, found := findCardInCardListByInstanceId(card, activeCards)
		if !found {
			err := fmt.Errorf(
				"card (instance id: %d) not found in play",
				card.Id,
			)
			return g.captureErrorAndStop(err)
		}

		cardAbilityUsedInstance := NewCardInstance(cardInstance, g)

		targets := []*CardInstance{}

		// the target can be opponent's cards
		activeCards = append(activeCards, g.activePlayerOpponent().CardsInPlay...)

		for _, target := range cardAbilityUsed.Targets {
			_, cardInstance, found := findCardInCardListByInstanceId(target.InstanceId, activeCards)
			if !found {
				err := fmt.Errorf(
					"card (instance id: %d) not found in play",
					target.InstanceId,
				)
				return g.captureErrorAndStop(err)
			}
			targetCardInstance := NewCardInstance(cardInstance, g)
			targets = append(targets, targetCardInstance)
		}

		if err := cardAbilityUsedInstance.UseAbility(targets); err != nil {
			return g.captureErrorAndStop(err)
		}

	}

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
	case zb.PlayerActionType_CheatDestroyCardsOnBoard:
		return actionCheatDestroyCardsOnBoard
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

	// TODO: overlord skill

	// determine the next action
	g.PrintState()
	next := g.next()
	if next == nil {
		return nil
	}

	switch next.ActionType {
	case zb.PlayerActionType_EndTurn:
		return actionEndTurn
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
	case zb.PlayerActionType_CheatDestroyCardsOnBoard:
		return actionCheatDestroyCardsOnBoard
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

	g.activePlayer().TurnNumber++

	previousPlayerTurnNumber := g.activePlayer().TurnNumber

	// change player turn
	g.changePlayerTurn()

	// add GooVial to active player
	addGooVialAndFillAll(g.activePlayer())

	// allow the new player to draw card on new turn
	g.activePlayer().HasDrawnCard = false

	// Draw the card. If this is the first move of the second player, they get 2 cards
	var cardsToDraw int
	if previousPlayerTurnNumber == 1 && g.activePlayer().TurnNumber == 0 {
		cardsToDraw = 2
	} else {
		cardsToDraw = 1
	}

	if err := g.drawCard(g.activePlayer(), cardsToDraw); err != nil {
		return g.captureErrorAndStop(err)
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
	case zb.PlayerActionType_CheatDestroyCardsOnBoard:
		return actionCheatDestroyCardsOnBoard
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
	case zb.PlayerActionType_CheatDestroyCardsOnBoard:
		return actionCheatDestroyCardsOnBoard
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
	case zb.PlayerActionType_CheatDestroyCardsOnBoard:
		return actionCheatDestroyCardsOnBoard
	default:
		return nil
	}
}

func actionCheatDestroyCardsOnBoard(g *Gameplay) stateFn {
	g.debugf("state: %v\n", zb.PlayerActionType_CheatDestroyCardsOnBoard)
	if g.isEnded() {
		return nil
	}
	// current action
	current := g.current()
	if current == nil {
		return nil
	}

	destroyedCards := current.GetCheatDestroyCardsOnBoard().DestroyedCards
	for _, destroyedCard := range destroyedCards {
		destroyedCardFound := false
		for playerStateIndex, playerState := range g.State.PlayerStates {
			if !g.playersDebugCheats[playerStateIndex].Enabled {
				return g.captureErrorAndStop(errCheatsRequired)
			}

			temp := playerState.CardsInPlay[:0]
			for _, card := range playerState.CardsInPlay {
				if card.InstanceId.Id == destroyedCard.Id {
					destroyedCardFound = true
				} else {
					temp = append(temp, card)
				}
			}
			playerState.CardsInPlay = temp
		}

		if !destroyedCardFound && g.useBackendGameLogic {
			return g.captureErrorAndStop(fmt.Errorf("card with instance id %d not found", destroyedCard.Id))
		}
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
	case zb.PlayerActionType_CheatDestroyCardsOnBoard:
		return actionCheatDestroyCardsOnBoard
	default:
		return nil
	}
}

func calculateOverlordLevel(overlordLevelingData *zb.OverlordLevelingData, overlord *zb.Overlord) int32 {
	var level = int32(overlord.Level)
	for overlord.Experience >= getRequiredExperienceForLevel(overlordLevelingData, level + 1) && level < overlordLevelingData.MaxLevel {
		level++
	}

	return level
}

func getRequiredExperienceForLevel(overlordLevelingData *zb.OverlordLevelingData, level int32) int64 {
	if level <= 1 {
		return 0
	}
	
	var fixed = overlordLevelingData.Fixed
	var experienceStep = overlordLevelingData.ExperienceStep
	var requiredExperience = int64(fixed) + int64(experienceStep)*(int64(level - 1))
	return requiredExperience
}

func getLevelReward(overlordLevelingData *zb.OverlordLevelingData, level int32) *zb.LevelReward {
	for i := 0; i < len(overlordLevelingData.Rewards); i++ {
		if overlordLevelingData.Rewards[i].Level == level {
			return overlordLevelingData.Rewards[i]
		}
	}

	return nil
}
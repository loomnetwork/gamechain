package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/gogo/protobuf/jsonpb"
	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/plugin"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	"github.com/loomnetwork/zombie_battleground/battleground"
	"github.com/loomnetwork/zombie_battleground/types/zb"
	log "github.com/sirupsen/logrus"
)

var pubKeyHexString = "e4008e26428a9bca87465e8de3a8d0e9c37a56ca619d3d6202b0567528786618"

func main() {
	if len(os.Args) == 0 {
		log.Error("GamePlay JSON file not provided")
		os.Exit(1)
	}

	// read game replay json
	fname := os.Args[1]
	f, err := os.Open(fname)
	if err != nil {
		log.Error("error opening json file: ", err)
		os.Exit(1)
	}

	var gameReplay zb.GameReplay
	err = jsonpb.Unmarshal(f, &gameReplay)
	if err != nil {
		log.Error("error unmarshalling json: ", err)
		os.Exit(1)
	}

	// set up fake context
	zbContract := &battleground.ZombieBattleground{}
	log.Info("Setting up fake context")
	fakeCtx := setupFakeContext()

	// initialise game chain
	log.Info("Initialising gamechain")
	initFile, err := os.Open("init.json")
	if err != nil {
		log.WithError(err).Error("error opening init.json")
		os.Exit(1)
	}

	var initRequest zb.InitRequest
	err = jsonpb.Unmarshal(initFile, &initRequest)
	if err != nil {
		log.WithError(err).Error("error unmarshalling init.json")
		os.Exit(1)
	}

	err = zbContract.Init(*fakeCtx, &initRequest)
	if err != nil {
		log.WithError(err).Error("error calling Init transaction")
		return
	}

	// log the game play being replayed
	replayedGameReplay := zb.GameReplay{
		ReplayVersion: gameReplay.ReplayVersion,
		RandomSeed:    gameReplay.RandomSeed,
	}

	// initialise game state
	log.Info("Initialising states")
	err = initialiseStates(*fakeCtx, zbContract, &gameReplay, &replayedGameReplay)
	if err != nil {
		log.WithError(err).Error("error initialising state")
		os.Exit(1)
	}

	// start replaying the actions and validate states after each transition
	log.Info("Starting replay and validate")
	err = replayAndValidate(*fakeCtx, zbContract, &gameReplay, &replayedGameReplay)
	if err != nil {
		log.WithError(err).Error("error while validating gameplay")
		//os.Exit(1)
	}

	fnameTrimmed := strings.TrimSuffix(fname, ".json")
	fnameReplayed := fnameTrimmed + "_replayed.json"
	outFile, err := os.Create(fnameReplayed)
	if err != nil {
		log.WithError(err).Errorf("error creating file %s", fnameReplayed)
	}

	err = new(jsonpb.Marshaler).Marshal(outFile, &replayedGameReplay)
	if err != nil {
		log.WithError(err).Error("error writing output to file")
	}
	log.Infof("Gameplay validation completed, output written to %s", fnameReplayed)
}

func setupFakeContext() *contract.Context {
	pubKey, _ := hex.DecodeString(pubKeyHexString)

	addr := &loom.Address{
		Local: loom.LocalAddressFromPublicKey(pubKey),
	}

	ctx := contract.WrapPluginContext(
		plugin.CreateFakeContext(*addr, *addr),
	)
	return &ctx
}

func initialiseStates(ctx contract.Context, zbContract *battleground.ZombieBattleground, gameReplay, replayedGameReplay *zb.GameReplay) error {
	actionList := gameReplay.Events
	initialState := actionList[0]

	// set up user accounts
	log.Info("Initialising user accounts and setting up match")
	playerStates := initialState.Match.PlayerStates
	var err error
	var newMatch *zb.Match
	for _, ps := range playerStates {
		err = zbContract.CreateAccount(ctx, &zb.UpsertAccountRequest{
			UserId:  ps.Id,
			Version: gameReplay.ReplayVersion,
		})
		if err != nil {
			return err
		}

		err = zbContract.EditDeck(ctx, &zb.EditDeckRequest{
			UserId: ps.Id,
			Deck:   ps.Deck,
		})
		if err != nil {
			return err
		}

		findMatchResp, err := zbContract.FindMatch(ctx, &zb.FindMatchRequest{
			UserId:     ps.Id,
			DeckId:     ps.Deck.Id,
			Version:    gameReplay.ReplayVersion,
			RandomSeed: gameReplay.RandomSeed,
		})
		if err != nil {
			return err
		}

		// the second iteration of the loop should give us a useful match state
		newMatch = findMatchResp.Match
	}

	// initialise the game state
	log.Info("Initialising game state")
	/*
		err = zbContract.SetMatch(ctx, &zb.SetMatchRequest{
			Match: initialState.Match,
		})
		if err != nil {
			return err
		}

		err = zbContract.SetGameState(ctx, &zb.SetGameStateRequest{
			GameState: initialState.GameState,
		})
		if err != nil {
			return err
		}
	*/
	getGSResp, err := zbContract.GetGameState(ctx, &zb.GetGameStateRequest{
		MatchId: newMatch.Id,
	})
	if err != nil {
		return err
	}

	playerEvent := &zb.PlayerActionEvent{
		Match:     newMatch,
		GameState: getGSResp.GameState,
	}

	replayedGameReplay.Events = append(replayedGameReplay.Events, playerEvent)

	return nil
}

func replayAndValidate(ctx contract.Context, zbContract *battleground.ZombieBattleground, gameReplay, replayedGameReplay *zb.GameReplay) error {
	actionList := gameReplay.Events
	replayActionList := actionList[1:]
	for _, replayAction := range replayActionList {
		actionReq := zb.PlayerActionRequest{
			MatchId:      1, //replayAction.Match.Id,
			PlayerAction: replayAction.PlayerAction,
		}
		log.Info("replaying action: ", actionReq)
		actionResp, err := zbContract.SendPlayerAction(ctx, &actionReq)
		if err != nil {
			return fmt.Errorf("error sending action %v: %v", actionReq.PlayerAction, err)
		}

		playerEvent := &zb.PlayerActionEvent{
			PlayerActionType: actionReq.PlayerAction.ActionType,
			UserId:           actionReq.PlayerAction.PlayerId,
			PlayerAction:     actionReq.PlayerAction,
			Match:            actionResp.Match,
			GameState:        actionResp.GameState,
		}
		replayedGameReplay.Events = append(replayedGameReplay.Events, playerEvent)

		newGameState := actionResp.GameState

		logGameState := replayAction.GameState

		log.Info("Comparing game states")
		err = compareGameStates(newGameState, logGameState)
		if err != nil {
			log.Error("game states do not match: ", err)
		}
	}
	return nil
}

func compareGameStates(newGameState, logGameState *zb.GameState) error {
	if newGameState.CurrentPlayerIndex != logGameState.CurrentPlayerIndex {
		log.Error("currentPlayerIndex doesn't match")
	}

	if newGameState.CurrentActionIndex != logGameState.CurrentActionIndex {
		log.Error("currentActionIndex doesn't match")
	}

	newPlayerStates := newGameState.PlayerStates
	logPlayerStates := logGameState.PlayerStates

	log.Info("Comparing player states")
	err := comparePlayerStates(newPlayerStates, logPlayerStates)
	if err != nil {
		log.Error("player states do not match: ", err)
	}
	return nil
}

func comparePlayerStates(newPlayerStates, logPlayerStates []*zb.PlayerState) error {
	for _, newPlayerState := range newPlayerStates {
		for _, logPlayerState := range logPlayerStates {
			if newPlayerState.Id == logPlayerState.Id {

				fmt.Println("comparing player state for: ", newPlayerState.Id)
				// TODO: deep compare using some library??
				// hp
				if newPlayerState.Hp != logPlayerState.Hp {
					return fmt.Errorf("hp doesn't match")
				}

				// mana
				if newPlayerState.Mana != logPlayerState.Mana {
					return fmt.Errorf("mana doesn't match")
				}

				// current action
				if newPlayerState.CurrentAction != logPlayerState.CurrentAction {
					return fmt.Errorf("current action doesn't match")
				}

				// overlord instance

				// cards in hand
				if len(newPlayerState.CardsInHand) != len(logPlayerState.CardsInHand) {
					return fmt.Errorf("card in hand don't match")
				}

				// cards in deck
				if len(newPlayerState.CardsInDeck) != len(logPlayerState.CardsInDeck) {
					return fmt.Errorf("card in deck don't match")
				}
				// deck
				fmt.Println("----------")
			}
		}
	}
	return nil
}

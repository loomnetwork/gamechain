package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

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
	f, err := os.Open(os.Args[1])
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

	j, _ := json.Marshal(gameReplay)
	fmt.Println(string(j))

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

	// initialise game state
	log.Info("Initialising states")
	err = initialiseStates(*fakeCtx, zbContract, &gameReplay)
	if err != nil {
		log.Error("error initialising state: ", err)
		os.Exit(1)
	}

	// start replaying the actions and validate states after each transition
	log.Info("Starting replay and validate")
	err = replayAndValidate(*fakeCtx, zbContract, &gameReplay)
	if err != nil {
		fmt.Println("error while validating gameplay: ", err)
		os.Exit(1)
	}
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

func initialiseStates(ctx contract.Context, zbContract *battleground.ZombieBattleground, gameReplay *zb.GameReplay) error {
	actionList := gameReplay.Events
	initialState := actionList[0]
	// set up user accounts
	log.Info("Initialising user accounts")
	playerStates := initialState.Match.PlayerStates
	var err error
	for _, ps := range playerStates {
		err = zbContract.CreateAccount(ctx, &zb.UpsertAccountRequest{
			UserId:  ps.Id,
			Version: "v1",
		})
		if err != nil {
			return err
		}
	}

	// initialise the game state
	log.Info("Initialising game state")
	err = zbContract.SetMatch(ctx, &zb.SetMatchRequest{
		Match: initialState.Match,
	})
	if err != nil {
		return err
	}

	initialGameState := initialState.GameState
	gs := &zb.GameState{
		Id:                 initialGameState.Id,
		IsEnded:            initialGameState.IsEnded,
		CurrentPlayerIndex: initialGameState.CurrentPlayerIndex,
		PlayerStates:       initialGameState.PlayerStates,
		CurrentActionIndex: initialGameState.CurrentActionIndex,
		Randomseed:         initialGameState.Randomseed,
		PlayerActions:      initialGameState.PlayerActions,
	}

	err = zbContract.SetGameState(ctx, &zb.SetGameStateRequest{
		GameState: gs,
	})
	if err != nil {
		return err
	}
	return nil
}

func replayAndValidate(ctx contract.Context, zbContract *battleground.ZombieBattleground, gameReplay *zb.GameReplay) error {
	actionList := gameReplay.Events
	replayActionList := actionList[1:]
	for _, replayAction := range replayActionList {
		actionReq := zb.PlayerActionRequest{
			MatchId:      replayAction.Match.Id,
			PlayerAction: replayAction.PlayerAction,
		}
		log.Info("replaying action: ", actionReq)
		actionResp, err := zbContract.SendPlayerAction(ctx, &actionReq)
		if err != nil {
			return err
		}
		newGameState := actionResp.GameState
		newPlayerStates := newGameState.PlayerStates

		logGameState := replayAction.GameState
		logPlayerStates := logGameState.PlayerStates

		log.Info("comparing states")
		err = comparePlayerStates(newPlayerStates, logPlayerStates)
		if err != nil {
			log.Error("player states do not match: ", err)
		}

	}
	return nil
}

func comparePlayerStates(newPlayerStates, logPlayerStates []*zb.PlayerState) error {
	for _, newPlayerState := range newPlayerStates {
		for _, logPlayerState := range logPlayerStates {
			if newPlayerState.Id == logPlayerState.Id {

				fmt.Println("comparing state for user ", newPlayerState.Id)
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

			}
		}
	}
	return nil
}

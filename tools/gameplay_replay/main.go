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
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Error("error opening json file: ", err)
		return
	}

	var gameReplay zb.GameReplay
	err = jsonpb.Unmarshal(f, &gameReplay)
	if err != nil {
		log.Error("error unmarshalling json: ", err)
		return
	}

	j, _ := json.Marshal(gameReplay)
	fmt.Println(string(j))

	zbContract := &battleground.ZombieBattleground{}
	log.Info("setting up fake context")
	fakeCtx := setupFakeContext()

	initFile, err := os.Open("init.json")
	if err != nil {
		log.WithError(err).Error("error opening init.json")
		return
	}

	log.Info("Initialising gamechain")
	var initRequest zb.InitRequest
	err = jsonpb.Unmarshal(initFile, &initRequest)
	if err != nil {
		log.WithError(err).Error("error unmarshalling init.json")
		return
	}

	err = zbContract.Init(*fakeCtx, &initRequest)
	if err != nil {
		log.WithError(err).Error("error calling Init transaction")
		return
	}

	actionList := gameReplay.Events

	log.Info("Initialising states")
	initialState := actionList[0]
	err = initialiseStates(*fakeCtx, zbContract, initialState)
	if err != nil {
		log.Error("error initialising state: ", err)
		return
	}

	log.Info("Starting replay and validate")
	log.Info(len(actionList))
	err = replayAndValidate(*fakeCtx, zbContract, actionList[1:])
	if err != nil {
		fmt.Println("error while validating gameplay: ", err)
		return
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

func initialiseStates(ctx contract.Context, zbContract *battleground.ZombieBattleground, initialState *zb.PlayerActionEvent) error {
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
	return nil
}

func replayAndValidate(ctx contract.Context, zbContract *battleground.ZombieBattleground, replayActionList []*zb.PlayerActionEvent) error {
	for _, replayAction := range replayActionList {
		actionReq := zb.PlayerActionRequest{
			MatchId:      replayAction.Match.Id,
			PlayerAction: replayAction.PlayerAction,
		}
		log.Info("replaying action: ", actionReq)
		actionResp, err := zbContract.SendPlayerAction(ctx, &actionReq)
		if err != nil {
			log.Error("error: ", err)
		}
		newGameState := actionResp.GameState
		newPlayerStates := newGameState.PlayerStates

		logPlayerStates := replayAction.Match.PlayerStates
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
				// TODO: compare using some library??
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

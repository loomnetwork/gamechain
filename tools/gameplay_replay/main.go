package main

import (
	"database/sql"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/loomnetwork/gamechain/battleground"
	"github.com/loomnetwork/gamechain/types/zb"
	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/plugin"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	log "github.com/sirupsen/logrus"
)

var (
	genesisFile = flag.String("genesis", "tools/gameplay_replay/init.json", "path to genesis file")
)

var (
	pubKeyHexString = "e4008e26428a9bca87465e8de3a8d0e9c37a56ca619d3d6202b0567528786618"
	db              *sql.DB
	readFromDB      bool
)

func main() {
	flag.Parse()

	readFromDB, _ = strconv.ParseBool(os.Getenv("READ_FROM_DB"))

	var gameReplay zb.GameReplay
	var fname string

	if readFromDB {
		var err error
		db, err = connectToDb()
		if err != nil {
			log.Println(err)
		}
		defer db.Close()
		if len(os.Args) == 0 {
			log.Fatal("Need match id argument")
		}
		row := db.QueryRow("SELECT * FROM replays WHERE match_id = ?", os.Args[1])

		var id int
		var replayJSON string
		err = row.Scan(&id, &replayJSON)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(replayJSON)
		err = jsonpb.UnmarshalString(replayJSON, &gameReplay)
		if err != nil {
			log.Error("error unmarshalling json: ", err)
			os.Exit(1)
		}

		fname = fmt.Sprintf("match%d.json", id)
	} else {
		if len(os.Args) == 1 {
			log.Error("GamePlay JSON file not provided")
			os.Exit(1)
		}

		// read game replay json
		fname = os.Args[1]
		// path := filepath.Join(basepath, "../../replays", fname)
		f, err := os.Open(fname)
		if err != nil {
			log.Error("error opening json file: ", err)
			os.Exit(1)
		}
		defer f.Close()

		err = jsonpb.Unmarshal(f, &gameReplay)
		if err != nil {
			log.Error("error unmarshalling json: ", err)
			os.Exit(1)
		}
	}

	// set up fake context
	zbContract := &battleground.ZombieBattleground{}
	log.Info("Setting up fake context")
	fakeCtx := setupFakeContext()

	// initialise game chain
	log.Info("Initialising gamechain using init.json")
	// initFilePath := filepath.Join(basepath, "init.json")
	initFile, err := os.Open(*genesisFile)
	if err != nil {
		log.WithError(err).Error("error opening init.json")
		os.Exit(1)
	}
	defer initFile.Close()

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
		Actions: gameReplay.Actions,
		Blocks:  gameReplay.Blocks,
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
		os.Exit(1)
	}

	// fnameTrimmed := strings.TrimSuffix(fname, ".json")
	// fnameReplayed := fnameTrimmed + "_replayed.json"
	// pathReplayed := filepath.Join(basepath, "../../replays", fnameReplayed)
	// outFile, err := os.Create(pathReplayed)
	// if err != nil {
	// 	log.WithError(err).Errorf("error creating file %s", pathReplayed)
	// }

	// err = new(jsonpb.Marshaler).Marshal(outFile, &replayedGameReplay)
	// if err != nil {
	// 	log.WithError(err).Error("error writing output to file")
	// }
	// log.Infof("Gameplay validation completed, output written to %s", pathReplayed)
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
	if len(gameReplay.Blocks) == 0 {
		return fmt.Errorf("no history data")
	}
	initialblock := gameReplay.Blocks[0]
	game := initialblock.GetCreateGame()
	if game == nil {
		return fmt.Errorf("expected createGame in the first history block")
	}

	// set up user accounts
	log.Info("Initialising user accounts and setting up match")
	var err error
	var newMatch *zb.Match
	for _, ps := range game.Players {
		err = zbContract.CreateAccount(ctx, &zb.UpsertAccountRequest{
			UserId:  ps.Id,
			Version: game.Version,
		})
		if err != nil {
			return errors.Wrapf(err, "error creating user account")
		}

		err = zbContract.EditDeck(ctx, &zb.EditDeckRequest{
			UserId:  ps.Id,
			Deck:    ps.Deck,
			Version: game.Version,
		})
		if err != nil {
			return err
		}

		findMatchResp, err := zbContract.FindMatch(ctx, &zb.FindMatchRequest{
			UserId:     ps.Id,
			DeckId:     ps.Deck.Id,
			Version:    game.Version,
			RandomSeed: game.Randomseed,
		})
		if err != nil {
			return err
		}

		// the second iteration of the loop should give us a useful match state
		newMatch = findMatchResp.Match
	}

	// initialise the game state
	log.Info("Initialising game state")
	_, err = zbContract.GetGameState(ctx, &zb.GetGameStateRequest{
		MatchId: newMatch.Id,
	})
	if err != nil {
		return err
	}
	return nil
}

func replayAndValidate(ctx contract.Context, zbContract *battleground.ZombieBattleground, gameReplay, replayedGameReplay *zb.GameReplay) error {
	req := zb.BundlePlayerActionRequest{
		MatchId:       1,
		PlayerActions: gameReplay.Actions,
	}
	resp, err := zbContract.SendBundlePlayerAction(ctx, &req)
	if err != nil {
		return err
	}

	historyBlock := replayedGameReplay.Blocks[1:]

	// validate history from recorded game replay and contract call
	log.Info("Comparing history blocks")
	if len(resp.History) != len(historyBlock) {
		return fmt.Errorf("expected history len: %d, get %d", len(historyBlock), len(resp.History))
	}
	for i := 0; i < len(resp.History); i++ {
		if !reflect.DeepEqual(resp.History[i], historyBlock[i]) {
			log.Errorf("different history blocks: %d", i)
			return fmt.Errorf("different history blocks: %d", i)
		}
	}

	log.Info("Comparing complete - no difference")
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

func connectToDb() (*sql.DB, error) {
	dbURL := os.Getenv("DATABASE_URL")
	var dbName string
	if dbURL == "" {
		dbUserName := os.Getenv("DATABASE_USERNAME")
		dbName = os.Getenv("DATABASE_NAME")
		dbPass := os.Getenv("DATABASE_PASS")
		dbHost := os.Getenv("DATABASE_HOST")
		dbPort := os.Getenv("DATABASE_PORT")
		if len(dbHost) == 0 {
			dbHost = "127.0.0.1"
		}
		if len(dbUserName) == 0 {
			dbUserName = "root"
		}
		if len(dbName) == 0 {
			dbName = "zb_replays"
		}
		if len(dbPort) == 0 {
			dbPort = "3306"
		}
		dbURL = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true", dbUserName, dbPass, dbHost, dbPort, dbName)
	}
	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		return nil, err
	}
	return db, nil
}

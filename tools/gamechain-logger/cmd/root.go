package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	raven "github.com/getsentry/raven-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	"github.com/loomnetwork/gamechain/types/zb"
	loom "github.com/loomnetwork/go-loom/plugin/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	rpcclient "github.com/tendermint/tendermint/rpc/lib/client"
)

var rootCmd = &cobra.Command{
	Use:          "gamechain-logger",
	Short:        "Loom Gamechain logger",
	Long:         `A logger that captures events from Gamechain and creates game metadata`,
	Example:      `  gamechain-logger`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return run()
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().String("db-url", "", "MySQL Connection URL")
	rootCmd.PersistentFlags().String("db-host", "127.0.0.1", "MySQL host")
	rootCmd.PersistentFlags().String("db-port", "3306", "MySQL port")
	rootCmd.PersistentFlags().String("db-name", "loomauth", "MySQL database name")
	rootCmd.PersistentFlags().String("db-user", "root", "MySQL database user")
	rootCmd.PersistentFlags().String("db-password", "", "MySQL database password")
	rootCmd.PersistentFlags().String("replay-dir", "replay", "replay directory")
	rootCmd.PersistentFlags().String("ws-url", "ws://localhost:9999/queryws", "WebSocket Connection URL")
	rootCmd.PersistentFlags().String("ev-url", "http://localhost:9999", "Event Indexer RPC Host")
	rootCmd.PersistentFlags().String("contract-name", "zombiebattleground:1.0.0", "Contract Name")
	rootCmd.PersistentFlags().Int("reconnect-interval", 1000, "Reconnect interval in MS")
	rootCmd.PersistentFlags().Int("block-interval", 20, "Amount of blocks to fetch")
	rootCmd.PersistentFlags().String("sentry-dsn", "", "sentry DSN, blank locally cause we dont want to send errors locally")
	rootCmd.PersistentFlags().String("sentry-environment", "", "sentry environment, leave it blank for localhost")

	viper.BindPFlag("db-url", rootCmd.PersistentFlags().Lookup("db-url"))
	viper.BindPFlag("db-host", rootCmd.PersistentFlags().Lookup("db-host"))
	viper.BindPFlag("db-port", rootCmd.PersistentFlags().Lookup("db-port"))
	viper.BindPFlag("db-name", rootCmd.PersistentFlags().Lookup("db-name"))
	viper.BindPFlag("db-user", rootCmd.PersistentFlags().Lookup("db-user"))
	viper.BindPFlag("db-password", rootCmd.PersistentFlags().Lookup("db-password"))
	viper.BindPFlag("replay-dir", rootCmd.PersistentFlags().Lookup("replay-dir"))
	viper.BindPFlag("ws-url", rootCmd.PersistentFlags().Lookup("ws-url"))
	viper.BindPFlag("ev-url", rootCmd.PersistentFlags().Lookup("ev-url"))
	viper.BindPFlag("contract-name", rootCmd.PersistentFlags().Lookup("contract-name"))
	viper.BindPFlag("reconnect-interval", rootCmd.PersistentFlags().Lookup("reconnect-interval"))
	viper.BindPFlag("block-interval", rootCmd.PersistentFlags().Lookup("block-interval"))
	viper.BindPFlag("sentry-dsn", rootCmd.PersistentFlags().Lookup("sentry-dsn"))
	viper.BindPFlag("sentry-environment", rootCmd.PersistentFlags().Lookup("sentry-environment"))
}

func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	sentryDsn := viper.GetString("sentry-dsn")
	sentryEnvironment := viper.GetString("sentry-environment")

	raven.SetEnvironment(sentryEnvironment)
	raven.SetDSN(sentryDsn)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		raven.CaptureErrorAndWait(err, map[string]string{})
		log.Error(err)
		os.Exit(1)
	}
}

func run() error {
	var (
		dbURL             = viper.GetString("db-url")
		dbHost            = viper.GetString("db-host")
		dbPort            = viper.GetString("db-port")
		dbName            = viper.GetString("db-name")
		dbUser            = viper.GetString("db-user")
		dbPassword        = viper.GetString("db-password")
		wsURL             = viper.GetString("ws-url")
		evURL             = viper.GetString("ev-url")
		contractName      = viper.GetString("contract-name")
		reconnectInterval = viper.GetInt("reconnect-interval")
		blockInterval     = viper.GetInt("block-interval")
	)

	var parsedURL *url.URL
	var URLType string
	var err error

	dbConnStr := dbURL
	if dbURL == "" {
		dbConnStr = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)
	}
	log.Printf("connecting to database host %s", dbHost)

	db, err := connectDb(dbConnStr)
	if err != nil {
		return errors.Wrapf(err, "fail to connect to database")
	}
	log.Printf("connected to database host %s", dbHost)
	defer db.Close()

	if evURL != "" {
		URLType = "ev"
		parsedURL, err = url.Parse(evURL)
		if err != nil {
			return errors.Wrapf(err, "Error parsing url %s", wsURL)
		}
	} else if wsURL != "" {
		URLType = "ws"
		parsedURL, err = url.Parse(wsURL)
		if err != nil {
			return errors.Wrapf(err, "Error parsing url %s", wsURL)
		}
	}
	if parsedURL.String() == "" {
		return errors.New("Eventstore or WebSocket Connection URL (--ev-url or --ws-url) is required")
	}

	// control channels
	doneC := make(chan struct{})
	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigC)

	reconnectIntervalDur := time.Duration(int64(reconnectInterval)) * time.Millisecond
	r := NewRunner(parsedURL.String(), URLType, db, 10, reconnectIntervalDur, blockInterval, contractName)
	go r.Start()
	go func() {
		select {
		case <-sigC:
			log.Println("stopping logger...")
			r.Stop()
			close(doneC)
		}
	}()

	<-doneC

	return nil
}

func connectGamechain(wsURL string) (*websocket.Conn, error) {
	subscribeCommand := struct {
		Method  string            `json:"method"`
		JSONRPC string            `json:"jsonrpc"`
		Params  map[string]string `json:"params"`
		ID      string            `json:"id"`
	}{"subevents", "2.0", make(map[string]string), "dummy"}
	subscribeMsg, err := json.Marshal(subscribeCommand)
	if err != nil {
		return nil, errors.Wrapf(err, "Cannot marshal command to json")
	}

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "Fail to connect to %s", wsURL)
	}
	if err := conn.WriteMessage(websocket.TextMessage, subscribeMsg); err != nil {
		return nil, err
	}
	return conn, nil
}

func queryEventStore(evURL string, fromBlock uint64, interval uint64, contract string) (*loom.ContractEventsResult, error) {
	log.Println("Querying Events from Height: ", fromBlock)

	rpcClient := rpcclient.NewJSONRPCClient(evURL)
	params := map[string]interface{}{
		"fromBlock": fromBlock,
		"toBlock":   fromBlock + interval,
		"contract":  contract,
	}
	result := &loom.ContractEventsResult{}
	_, err := rpcClient.Call("contractevents", params, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func connectDb(dbURL string) (*gorm.DB, error) {
	db, err := gorm.Open("mysql", dbURL)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func writeReplayFile(topic string, event zb.PlayerActionEvent) ([]byte, error) {
	dir := viper.GetString("replay-dir")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if e := os.MkdirAll(dir, os.ModePerm); e != nil {
			return nil, e
		}
	}

	filename := fmt.Sprintf("%s.json", topic)
	path := filepath.Join(dir, filename)

	fmt.Println("Writing to file: ", path)

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if event.Block == nil {
		return nil, nil
	}

	var replay zb.GameReplay
	if fi, _ := f.Stat(); fi.Size() > 0 {
		if err := jsonpb.Unmarshal(f, &replay); err != nil {
			log.Println(err)
			return nil, err
		}
	}

	if event.PlayerAction != nil {
		replay.Blocks = append(replay.Blocks, event.Block.List...)
		replay.Actions = append(replay.Actions, event.PlayerAction)
	} else {
		replay.Blocks = append(replay.Blocks, event.Block.List...)
	}

	m := jsonpb.Marshaler{}
	result, err := m.MarshalToString(&replay)
	if err != nil {
		return nil, err
	}

	if err := ioutil.WriteFile(path, []byte(result), os.ModePerm); err != nil {
		return nil, err
	}

	return []byte(result), nil
}

package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:          "gamechain-logger [url]",
	Short:        "Loom Gamechain logger",
	Long:         `A logger that captures events from Gamechain and creates game metadata`,
	Example:      `  gamechain-logger ws://localhost:9999/queryws replays`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			cmd.Usage()
			return fmt.Errorf("Gamechain websocket URL endpoint required")
		}
		return run(args[0])
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
	viper.BindPFlag("db-url", rootCmd.PersistentFlags().Lookup("db-url"))
	viper.BindPFlag("db-host", rootCmd.PersistentFlags().Lookup("db-host"))
	viper.BindPFlag("db-port", rootCmd.PersistentFlags().Lookup("db-port"))
	viper.BindPFlag("db-name", rootCmd.PersistentFlags().Lookup("db-name"))
	viper.BindPFlag("db-user", rootCmd.PersistentFlags().Lookup("db-user"))
	viper.BindPFlag("db-password", rootCmd.PersistentFlags().Lookup("db-password"))
	viper.BindPFlag("replay-dir", rootCmd.PersistentFlags().Lookup("replay-dir"))
}

func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(wsURL string) error {
	var (
		dbURL      = viper.GetString("DATABASE_URL")
		dbHost     = viper.GetString("DATABASE_HOST")
		dbPort     = viper.GetString("DATABASE_PORT")
		dbName     = viper.GetString("DATABASE_NAME")
		dbUser     = viper.GetString("DATABASE_USERNAME")
		dbPassword = viper.GetString("DATABASE_PASS")
	)

	parsedURL, err := url.Parse(wsURL)
	if err != nil {
		return errors.Wrapf(err, "Error parsing url %s", wsURL)
	}

	// db should be optional?
	dbConnStr := dbURL
	if dbURL == "" {
		dbConnStr = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true", dbUser, url.QueryEscape(dbPassword), dbHost, dbPort, dbName)
	}
	log.Printf("connecting to database host %s", dbHost)

	db, err := connectDb(dbConnStr)
	if err != nil {
		return errors.Wrapf(err, "fail to connect to database")
	}
	log.Printf("connected to database host %s", dbHost)
	defer db.Close()

	// control channels
	doneC := make(chan struct{})
	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigC)

	r := NewRunner(parsedURL.String(), db, 10)
	// Start is not blocking
	r.Start()

	go func() {
		select {
		case err := <-r.Error():
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v", err)
			}
			r.Stop()
			close(doneC)
		case <-sigC:
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
		replay.Blocks = event.Block.List
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

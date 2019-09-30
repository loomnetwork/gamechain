package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
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
	"github.com/jinzhu/gorm"
	"github.com/loomnetwork/gamechain/types/zb/zb_data"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	rootCmd.PersistentFlags().String("chain-id", "default", "Chain Id")
	rootCmd.PersistentFlags().String("read-uri", "http://localhost:46658/query", "URI for quering app state")
	rootCmd.PersistentFlags().String("contract-name", "zombiebattleground:1.0.0", "Contract Name")
	rootCmd.PersistentFlags().Int("reconnect-interval", 10, "Reconnect interval in seconds")
	rootCmd.PersistentFlags().Int("poll-interval", 10, "Poll interval in seconds")
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
	viper.BindPFlag("chain-id", rootCmd.PersistentFlags().Lookup("chain-id"))
	viper.BindPFlag("read-uri", rootCmd.PersistentFlags().Lookup("read-uri"))
	viper.BindPFlag("contract-name", rootCmd.PersistentFlags().Lookup("contract-name"))
	viper.BindPFlag("poll-interval", rootCmd.PersistentFlags().Lookup("poll-interval"))
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
		log.Println(err)
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
		chainID           = viper.GetString("chain-id")
		readURI           = viper.GetString("read-uri")
		contractName      = viper.GetString("contract-name")
		pollInterval      = viper.GetInt("poll-interval")
		reconnectInterval = viper.GetInt("reconnect-interval")
		blockInterval     = viper.GetInt("block-interval")
	)

	var parsedURL *url.URL
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

	parsedURL, err = url.Parse(readURI)
	if err != nil {
		return errors.Wrapf(err, "Error parsing url %s", readURI)
	}

	// control channels
	doneC := make(chan struct{})
	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigC)

	config := &Config{
		ReadURI:           parsedURL.String(),
		ChainID:           chainID,
		ReconnectInterval: time.Duration(int64(reconnectInterval)) * time.Second,
		PollInterval:      time.Duration(int64(pollInterval)) * time.Second,
		ContractName:      contractName,
		BlockInterval:     blockInterval,
	}

	r := NewRunner(db, config)
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

func connectDb(dbURL string) (*gorm.DB, error) {
	db, err := gorm.Open("mysql", dbURL)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func writeReplayFile(topic string, event zb_data.PlayerActionEvent) ([]byte, error) {
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

	var replay zb_data.GameReplay
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

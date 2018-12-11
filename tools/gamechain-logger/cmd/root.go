package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/Jeffail/gabs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	"github.com/loomnetwork/gamechain/models"
	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:          "gamechain-logger [url]",
	Short:        "Loom Gamechain logger",
	Long:         `A logger that captures events from Gamechain and creates game metadata`,
	Example:      `  gamechain-logger ws://localhost:9999/queryws`,
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
	rootCmd.PersistentFlags().String("db-name", "loom", "MySQL database name")
	rootCmd.PersistentFlags().String("db-user", "root", "MySQL database user")
	rootCmd.PersistentFlags().String("db-password", "", "MySQL database password")
	viper.BindPFlag("db-url", rootCmd.PersistentFlags().Lookup("db-url"))
	viper.BindPFlag("db-host", rootCmd.PersistentFlags().Lookup("db-host"))
	viper.BindPFlag("db-port", rootCmd.PersistentFlags().Lookup("db-port"))
	viper.BindPFlag("db-name", rootCmd.PersistentFlags().Lookup("db-name"))
	viper.BindPFlag("db-user", rootCmd.PersistentFlags().Lookup("db-user"))
	viper.BindPFlag("db-password", rootCmd.PersistentFlags().Lookup("db-password"))
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
		dbURL      = viper.GetString("db-url")
		dbHost     = viper.GetString("db-host")
		dbPort     = viper.GetString("db-port")
		dbName     = viper.GetString("db-name")
		dbUser     = viper.GetString("db-user")
		dbPassword = viper.GetString("db-password")
	)

	parsedURL, err := url.Parse(wsURL)
	if err != nil {
		return errors.Wrapf(err, "Error parsing url %s", wsURL)
	}

	conn, err := connectGamechain(parsedURL.String())
	if err != nil {
		return err
	}
	defer conn.Close()
	log.Printf("connected to gamechain url %s", wsURL)

	// db should be optional?
	dbConnStr := dbURL
	if dbURL != "" {
		dbConnStr = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true", dbUser, url.QueryEscape(dbPassword), dbHost, dbPort, dbName)
	}
	log.Printf("connected to database host %s", dbHost)

	db, err := connectDb(dbConnStr)
	if err != nil {
		return errors.Wrapf(err, "fail to connect to database")
	}
	defer db.Close()
	err = db.AutoMigrate(&models.Match{}, &models.Replay{}, &models.Deck{}, &models.DeckCard{}).Error
	if err != nil {
		return err
	}

	// TODO: need to run with recovery and need to capture SIG TERM
	wsLoop(conn, db)

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

func writeReplayFile(topic string, body []byte) ([]byte, error) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	path := filepath.Join(basepath, "../../replays/")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}

	filename := fmt.Sprintf("%s.json", topic)
	path = filepath.Join(path, filename)

	fmt.Println("Writing to file: ", path)

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var event zb.PlayerActionEvent
	err = proto.Unmarshal(body, &event)
	if err != nil {
		return nil, err
	}

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

func wsLoop(conn *websocket.Conn, db *gorm.DB) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}

		msgJSON, _ := gabs.ParseJSON(message)
		result := msgJSON.Path("result")

		results, _ := result.Children()
		if len(results) != 0 {
			pluginName := result.Path("plugin_name").Data().(string)
			if strings.HasPrefix(pluginName, "zombiebattleground") {

				height := int(result.Path("block_height").Data().(float64))
				log.Printf("height: %d", height)
				topics, _ := result.Path("topics").Children()
				fmt.Println("Getting event with topics: ", topics)
				topic := strings.Trim(strings.Replace(topics[0].String(), ":", "", -1), "\"")
				var secondaryTopic string
				if len(topics) > 1 {
					secondaryTopic = strings.Trim(strings.Replace(topics[1].String(), ":", "", -1), "\"")
				}
				encodedBody, ok := result.Path("encoded_body").Data().(string)
				if !ok {
					log.Println("Error getting encoded_body from message")
				}
				body, err := base64.StdEncoding.DecodeString(encodedBody)
				if err != nil {
					log.Println("Error decoding encoded_body")
				}

				switch {
				case secondaryTopic == "zombiebattlegroundfindmatch" || secondaryTopic == "zombiebattlegroundacceptmatch":
					var event zb.PlayerActionEvent
					err = proto.Unmarshal(body, &event)
					if err != nil {
						log.Println(err)
					}

					if secondaryTopic == "zombiebattlegroundfindmatch" {
						match := models.Match{}
						err = db.Where(&models.Match{ID: event.Match.Id}).First(&match).Error
						if err == nil {
							continue
						}
						match = models.Match{
							ID:              event.Match.Id,
							Player1ID:       event.Match.PlayerStates[0].Id,
							Player2ID:       event.Match.PlayerStates[1].Id,
							Player1Accepted: event.Match.PlayerStates[0].MatchAccepted,
							Player2Accepted: event.Match.PlayerStates[1].MatchAccepted,
							Player1DeckID:   event.Match.PlayerStates[0].Deck.Id,
							Player2DeckID:   event.Match.PlayerStates[1].Deck.Id,
							Status:          event.Match.Status.String(),
							Version:         event.Match.Version,
							RandomSeed:      event.Match.RandomSeed,
						}
						err = db.Create(&match).Error
						if err != nil {
							log.Println("Error creating match: ", err)
						}
					} else {
						match := models.Match{}

						err = db.Where(&models.Match{ID: event.Match.Id}).First(&match).Error
						if err != nil {
							log.Println("Error getting match from DB: ", err)
							continue
						}

						match.Player1Accepted = event.Match.PlayerStates[0].MatchAccepted
						match.Player2Accepted = event.Match.PlayerStates[1].MatchAccepted
						match.Status = event.Match.Status.String()

						err = db.Save(&match).Error
						if err != nil {
							log.Println("Error updating match: ", err)
						}
					}
				case secondaryTopic == "endgame":
					var event zb.PlayerActionEvent
					err = proto.Unmarshal(body, &event)
					if err != nil {
						fmt.Println(err)
					}

					match := models.Match{}

					err = db.Where(&models.Match{ID: event.Match.Id}).First(&match).Error
					if err != nil {
						log.Println("Error getting match from DB: ", err)
						continue
					}

					match.WinnerID = event.Block.List[0].GetEndGame().WinnerId
					match.Status = event.Match.Status.String()

					err = db.Save(&match).Error
					if err != nil {
						log.Println("Error updating match: ", err)
					}
				case topic == "zombiebattlegroundcreatedeck":
					var event zb.CreateDeckEvent
					err := proto.Unmarshal(body, &event)
					if err != nil {
						log.Println("Error unmarshaling event: ", err)
						continue
					}

					log.Printf("Saving deck with deck ID %d, userid %s, name %s to DB", event.Deck.Id, event.UserId, event.Deck.Name)

					cards := []models.DeckCard{}
					for _, card := range event.Deck.Cards {
						cards = append(cards, models.DeckCard{
							CardName: card.CardName,
							Amount:   card.Amount,
						})
					}

					fmt.Printf("DECK MSG: %+v", event)

					deck := models.Deck{
						UserID:           event.UserId,
						DeckID:           event.Deck.Id,
						Name:             event.Deck.Name,
						HeroID:           event.Deck.HeroId,
						Cards:            cards,
						PrimarySkillID:   0,
						SecondarySkillID: 0,
						Version:          event.Version,
						SenderAddress:    event.SenderAddress,
					}

					err = db.Create(&deck).Error
					if err != nil {
						log.Println("Error saving deck: ", err)
					}
					log.Printf("Saved deck with deck ID %d, userid %s, name %s to DB", event.Deck.Id, event.UserId, event.Deck.Name)
				case topic == "zombiebattlegroundeditdeck":
					var event zb.EditDeckEvent
					err := proto.Unmarshal(body, &event)
					if err != nil {
						log.Println("Error unmarshaling event: ", err)
						continue
					}

					log.Printf("Editing deck with deck ID %d, userid %s, name %s", event.Deck.Id, event.UserId, event.Deck.Name)

					deck := models.Deck{}

					err = db.Where(&models.Deck{UserID: event.UserId, DeckID: event.Deck.Id}).First(&deck).Error
					if err != nil {
						log.Println("Error getting deck from DB: ", err)
						continue
					}

					cards := []models.DeckCard{}
					for _, card := range event.Deck.Cards {
						cards = append(cards, models.DeckCard{
							CardName: card.CardName,
							Amount:   card.Amount,
						})
					}

					db.Model(&deck).Association("Cards").Replace(cards)

					deck.Name = event.Deck.Name
					deck.HeroID = event.Deck.HeroId
					deck.PrimarySkillID = 0
					deck.SecondarySkillID = 0
					deck.Version = event.Version
					deck.SenderAddress = event.SenderAddress

					err = db.Save(&deck).Error
					if err != nil {
						log.Println("Error updating deck: ", err)
					}
					log.Printf("Saved deck with deck ID %d, userid %s, name %s", event.Deck.Id, event.UserId, event.Deck.Name)
				case topic == "zombiebattlegrounddeletedeck":
					var event zb.DeleteDeckEvent
					err := proto.Unmarshal(body, &event)
					if err != nil {
						log.Println("Error unmarshaling event: ", err)
						continue
					}

					log.Printf("Deleting deck with deck ID %d, userid %s from DB", event.DeckId, event.UserId)

					db.Where(&models.Deck{UserID: event.UserId, DeckID: event.DeckId}).Delete(models.Deck{})

					log.Printf("Deleted deck with deck ID %d, userid %s from DB", event.DeckId, event.UserId)
				case strings.HasPrefix(topic, "match"):
					replay, err := writeReplayFile(topic, body)
					if err != nil {
						log.Println("Error writing replay file: ", err)
					}

					matchID, err := strconv.ParseInt(topic[5:], 10, 64)
					if err != nil {
						log.Println(err)
					}
					log.Printf("Saving replay with match ID %d to DB", matchID)

					dbReplay := models.Replay{}

					err = db.Where(&models.Replay{MatchID: matchID}).First(&dbReplay).Error
					if err == nil {
						db.First(&dbReplay)
						dbReplay.ReplayJSON = replay
						db.Save(&dbReplay)
					} else if gorm.IsRecordNotFoundError(err) {
						// insert
						dbReplay.MatchID = matchID
						dbReplay.ReplayJSON = replay
						db.Create(&dbReplay)
					} else {
						log.Println("Error getting replay: ", err)
					}
				default:
					continue
				}
			}
		}
	}
}

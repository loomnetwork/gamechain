package main

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
	"time"

	"github.com/Jeffail/gabs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	"github.com/loomnetwork/gamechain/types/zb"
)

var (
	wsURL string
	db    *gorm.DB
)

type Match struct {
	ID              int64 `gorm:"PRIMARY_KEY,auto_increment:false"`
	Player1ID       string
	Player2ID       string
	Player1Accepted bool
	Player2Accepted bool
	Status          string
	Version         string
	RandomSeed      int64
	Replay          Replay
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Replay struct {
	gorm.Model
	MatchID    int64
	ReplayJSON json.RawMessage
}

func main() {
	wsURL = os.Getenv("wsURL")
	if len(wsURL) == 0 {
		wsURL = "ws://localhost:9999/queryws"
	}
	log.Printf("wsURL - %s", wsURL)

	var err error
	db, err = connectToDb()
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	wsLoop()
}

func wsLoop() {
	subscribeCommand := struct {
		Method  string            `json:"method"`
		JSONRPC string            `json:"jsonrpc"`
		Params  map[string]string `json:"params"`
		ID      string            `json:"id"`
	}{"subevents", "2.0", make(map[string]string), "dummy"}
	subscribeMsg, err := json.Marshal(subscribeCommand)
	if err != nil {
		log.Fatal("Cannot marshal command to json")
	}

	parsedURL, err := url.Parse(wsURL)
	if err != nil {
		log.Fatal("Error parsing url: ", err)
	}

	u := url.URL{Scheme: parsedURL.Scheme, Host: parsedURL.Host, Path: parsedURL.Path}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
	if err := c.WriteMessage(websocket.TextMessage, subscribeMsg); err != nil {
		log.Fatal("Error writing command:", err)
	}

	for {
		_, message, err := c.ReadMessage()
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
				var extraTopic string
				if len(topics) > 1 {
					extraTopic = strings.Trim(strings.Replace(topics[1].String(), ":", "", -1), "\"")
				}
				encodedBody := result.Path("encoded_body").Data().(string)
				body, _ := base64.StdEncoding.DecodeString(encodedBody)

				if extraTopic == "zombiebattlegroundfindmatch" || extraTopic == "zombiebattlegroundacceptmatch" {
					var event zb.PlayerActionEvent
					err = proto.Unmarshal(body, &event)
					if err != nil {
						fmt.Println(err)
					}

					if extraTopic == "zombiebattlegroundfindmatch" {
						match := Match{
							ID:              event.Match.Id,
							Player1ID:       event.Match.PlayerStates[0].Id,
							Player2ID:       event.Match.PlayerStates[1].Id,
							Player1Accepted: event.Match.PlayerStates[0].MatchAccepted,
							Player2Accepted: event.Match.PlayerStates[1].MatchAccepted,
							Status:          event.Match.Status.String(),
							Version:         event.Match.Version,
							RandomSeed:      event.Match.RandomSeed,
						}
						db.Create(&match)
					} else {
						match := Match{}

						err = db.Where(&Match{ID: event.Match.Id}).First(&match).Error
						if err != nil {
							log.Println("Error getting match from DB: ", err)
							continue
						}

						match.Player1ID = event.Match.PlayerStates[0].Id
						match.Player2ID = event.Match.PlayerStates[1].Id
						match.Player1Accepted = event.Match.PlayerStates[0].MatchAccepted
						match.Player2Accepted = event.Match.PlayerStates[1].MatchAccepted
						match.Status = event.Match.Status.String()
						match.Version = event.Match.Version
						match.RandomSeed = event.Match.RandomSeed

						db.Save(&match)
					}

					continue
				}

				if !strings.HasPrefix(topic, "match") {
					continue
				}

				replay, err := writeReplayFile(topic, body)
				if err != nil {
					log.Println("Error writing replay file: ", err)
				}

				matchID, err := strconv.ParseInt(topic[5:], 10, 64)
				if err != nil {
					log.Println(err)
				}
				log.Printf("Saving replay with match ID %d to DB", matchID)

				dbReplay := Replay{}

				err = db.Where(&Replay{MatchID: matchID}).First(&dbReplay).Error
				if err == nil {
					db.First(&dbReplay)
					dbReplay.ReplayJSON = replay
					db.Save(&dbReplay)
				} else {
					// insert
					dbReplay.MatchID = matchID
					dbReplay.ReplayJSON = replay
					db.Create(&dbReplay)
				}
			}
		}
	}
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

func connectToDb() (*gorm.DB, error) {
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
	db, err := gorm.Open("mysql", dbURL)
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&Match{}, &Replay{})
	// _, err = db.Exec("CREATE TABLE IF NOT EXISTS zb_replays (match_id INT, replay_json MEDIUMBLOB, PRIMARY KEY (match_id))")
	// if err != nil {
	// 	return nil, err
	// }
	// _, err = db.Exec("CREATE TABLE IF NOT EXISTS zb_matches (match_id INT, player1_id VARCHAR(255), player2_id VARCHAR(255), player1_accepted BOOL DEFAULT false, player2_accepted BOOL DEFAULT false, status INT, version VARCHAR(32), randomseed INT, created_at TIMESTAMP NOT NULL DEFAULT NOW(), updated_at TIMESTAMP NOT NULL DEFAULT NOW() ON UPDATE now(), PRIMARY KEY (match_id))")
	// if err != nil {
	// 	return nil, err
	// }
	return db, nil
}

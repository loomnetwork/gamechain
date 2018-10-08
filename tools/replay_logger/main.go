package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	_ "github.com/go-sql-driver/mysql"

	"encoding/base64"
	"strings"

	"github.com/Jeffail/gabs"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gorilla/websocket"
	"github.com/loomnetwork/gamechain/types/zb"
)

var (
	wsURL    string
	saveToDB bool
	db       *sql.DB
)

func main() {
	wsURL = os.Getenv("wsURL")
	if len(wsURL) == 0 {
		wsURL = "ws://localhost:9999/queryws"
	}
	log.Printf("wsURL - %s", wsURL)
	saveToDB, _ = strconv.ParseBool(os.Getenv("SAVE_TO_DB"))

	if saveToDB {
		var err error
		db, err = connectToDb()
		if err != nil {
			log.Println(err)
		}
		defer db.Close()
	}

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
				topic := strings.Trim(strings.Replace(topics[0].String(), ":", "", -1), "\"")
				encodedBody := result.Path("encoded_body").Data().(string)
				body, _ := base64.StdEncoding.DecodeString(encodedBody)

				fmt.Println("Getting event with topic: ", topic)

				if !strings.HasPrefix(topic, "match") {
					continue
				}

				replay, err := writeReplayFile(topic, body)
				if err != nil {
					log.Println("Error writing replay file: ", err)
				}

				if saveToDB {
					matchID, err := strconv.ParseInt(topic[5:], 10, 64)
					if err != nil {
						log.Println(err)
					}
					log.Printf("Saving replay with match ID %d to DB", matchID)
					_, err = db.Exec(`INSERT INTO replays set match_id=?, replay_json=? ON DUPLICATE KEY UPDATE replay_json = ?`, matchID, replay, replay)
					if err != nil {
						log.Println("Error saving replay to DB: ", err)
					}
				}

			}
		}
	}
}

func writeReplayFile(topic string, body []byte) ([]byte, error) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	filename := fmt.Sprintf("replays/%s.json", topic)
	path := filepath.Join(basepath, "../../", filename)

	fmt.Println("Writing to file: ", path)

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	var event zb.PlayerActionEvent
	if err := jsonpb.UnmarshalString(string(body), &event); err != nil {
		return nil, err
	}

	var replay zb.GameReplay
	if fi, _ := f.Stat(); fi.Size() > 0 {
		if err := jsonpb.Unmarshal(f, &replay); err != nil {
			log.Println(err)
			return nil, err
		}
	} else {
		replay.Events = []*zb.PlayerActionEvent{}
		bodyJSON, _ := gabs.ParseJSON(body)
		seed, err := strconv.ParseInt(bodyJSON.Path("gameState.randomseed").Data().(string), 10, 64)
		if err != nil {
			return nil, err
		}
		replay.RandomSeed = seed
		version := bodyJSON.Path("match.version").Data()
		if version == nil {
			version = "v1" //TODO: make sure we always have a version
		}
		replay.ReplayVersion = version.(string)
	}

	replay.Events = append(replay.Events, &event)

	f.Close()

	m := jsonpb.Marshaler{}
	result, err := m.MarshalToString(&replay)
	if err != nil {
		return nil, err
	}

	if err := ioutil.WriteFile(path, []byte(result), 0644); err != nil {
		return nil, err
	}

	return []byte(result), nil
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
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS replays (match_id INT, replay_json MEDIUMBLOB, PRIMARY KEY (match_id))")
	if err != nil {
		return nil, err
	}
	return db, nil
}

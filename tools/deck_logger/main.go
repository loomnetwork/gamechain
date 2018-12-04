package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/Jeffail/gabs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
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
				topic := strings.Trim(strings.Replace(topics[0].String(), ":", "", -1), "\"")
				encodedBody := result.Path("encoded_body").Data().(string)
				body, _ := base64.StdEncoding.DecodeString(encodedBody)

				fmt.Println("Getting event with topic: ", topic)

				if "zombiebattlegroundcreatedeck" != topic {
					continue
				}

				fmt.Println("CREATE DECK TOPIC RECEIVED")
				fmt.Printf("body: %+v", string(body))

				var event zb.CreateDeckEvent
				err := proto.Unmarshal(body, &event)
				if err != nil {
					log.Println("Error unmarshaling event: ", err)
					continue
				}

				fmt.Printf("event: %+v", event)

				log.Printf("Saving deck with deck ID %d to DB", event.Deck.Id)
				log.Printf("Saving deck with deck userid %s to DB", event.UserId)
				log.Printf("Saving deck with deck name %s to DB", event.Deck.Name)

				m := jsonpb.Marshaler{}
				deckJSON, err := m.MarshalToString(event.Deck)
				if err != nil {
					log.Println("Error marshaling deck to json: ", err)
					continue
				}

				_, err = db.Exec(`INSERT INTO zb_decks set user=?, deck_id=?, deck_name=?, deck_json=?, version=?, sender_address=?`, event.UserId, event.Deck.Id, event.Deck.Name, deckJSON, event.Version, event.SenderAddress)
				if err != nil {
					log.Println("Error saving replay to DB: ", err)
				}

			}
		}
	}
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
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS zb_decks (user VARCHAR(255), deck_id INT, deck_name VARCHAR(255), deck_json MEDIUMBLOB, version VARCHAR(32), sender_address VARCHAR(255), created_at TIMESTAMP NOT NULL DEFAULT NOW(), PRIMARY KEY (user, deck_id))")
	if err != nil {
		return nil, err
	}
	return db, nil
}

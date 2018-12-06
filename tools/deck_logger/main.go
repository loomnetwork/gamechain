package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/Jeffail/gabs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	"github.com/loomnetwork/gamechain/types/zb"
)

var (
	wsURL string
	db    *gorm.DB
)

type Deck struct {
	gorm.Model
	UserID           string `gorm:"UNIQUE_INDEX:idx_userid_deckid"`
	DeckID           int64  `gorm:"UNIQUE_INDEX:idx_userid_deckid"`
	Name             string
	HeroID           int64
	Cards            []DeckCard
	PrimarySkillID   int
	SecondarySkillID int
	Version          string
	SenderAddress    string
}

type DeckCard struct {
	gorm.Model
	DeckID   uint
	CardName string
	Amount   int64
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
				topic := strings.Trim(strings.Replace(topics[0].String(), ":", "", -1), "\"")
				encodedBody := result.Path("encoded_body").Data().(string)
				body, _ := base64.StdEncoding.DecodeString(encodedBody)

				fmt.Println("Getting event with topic: ", topic)

				switch topic {
				case "zombiebattlegroundcreatedeck":
					var event zb.CreateDeckEvent
					err := proto.Unmarshal(body, &event)
					if err != nil {
						log.Println("Error unmarshaling event: ", err)
						continue
					}

					log.Printf("Saving deck with deck ID %d, userid %s, name %s to DB", event.Deck.Id, event.UserId, event.Deck.Name)

					cards := []DeckCard{}
					for _, card := range event.Deck.Cards {
						cards = append(cards, DeckCard{
							CardName: card.CardName,
							Amount:   card.Amount,
						})
					}

					fmt.Printf("DECK MSG: %+v", event)

					deck := Deck{
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

					db.NewRecord(deck)
					db.Create(&deck)
					log.Printf("Saved deck with deck ID %d, userid %s, name %s to DB", event.Deck.Id, event.UserId, event.Deck.Name)
					// _, err = db.Exec(`INSERT INTO zb_decks set user=?, deck_id=?, deck_name=?, deck_json=?, version=?, sender_address=?`, event.UserId, event.Deck.Id, event.Deck.Name, deckJSON, event.Version, event.SenderAddress)
					// if err != nil {
					// 	log.Println("Error saving deck to DB: ", err)
					// }
				case "zombiebattlegroundeditdeck":
					var event zb.EditDeckEvent
					err := proto.Unmarshal(body, &event)
					if err != nil {
						log.Println("Error unmarshaling event: ", err)
						continue
					}

					log.Printf("Editing deck with deck ID %d, userid %s, name %s", event.Deck.Id, event.UserId, event.Deck.Name)

					deck := Deck{}

					err = db.Where(&Deck{UserID: event.UserId, DeckID: event.Deck.Id}).First(&deck).Error
					if err != nil {
						log.Println("Error getting deck from DB: ", err)
						continue
					}

					cards := []DeckCard{}
					for _, card := range event.Deck.Cards {
						cards = append(cards, DeckCard{
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

					db.Save(&deck)
					log.Printf("Saved deck with deck ID %d, userid %s, name %s", event.Deck.Id, event.UserId, event.Deck.Name)
					// _, err = db.Exec(`UPDATE zb_decks SET deck_name=?, deck_json=?, version=?, sender_address=? WHERE user=? AND deck_id=?`, event.Deck.Name, deckJSON, event.Version, event.SenderAddress, event.UserId, event.Deck.Id)
					// if err != nil {
					// 	log.Println("Error saving deck to DB: ", err)
					// }
				case "zombiebattlegrounddeletedeck":
					var event zb.DeleteDeckEvent
					err := proto.Unmarshal(body, &event)
					if err != nil {
						log.Println("Error unmarshaling event: ", err)
						continue
					}

					log.Printf("Deleting deck with deck ID %d, userid %s from DB", event.DeckId, event.UserId)

					db.Where(&Deck{UserID: event.UserId, DeckID: event.DeckId}).Delete(Deck{})

					log.Printf("Deleted deck with deck ID %d, userid %s from DB", event.DeckId, event.UserId)
					// _, err = db.Exec(`DELETE FROM zb_decks WHERE user=? AND deck_id=?`, event.UserId, event.DeckId)
					// if err != nil {
					// 	log.Println("Error saving deck to DB: ", err)
					// }
				default:
					continue
				}

			}
		}
	}
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
	db.AutoMigrate(&Deck{}, &DeckCard{})
	// _, err = db.Exec("CREATE TABLE IF NOT EXISTS zb_decks (user VARCHAR(255), deck_id INT, deck_name VARCHAR(255), deck_json MEDIUMBLOB, version VARCHAR(32), sender_address VARCHAR(255), created_at TIMESTAMP NOT NULL DEFAULT NOW(), updated_at TIMESTAMP NOT NULL DEFAULT NOW() ON UPDATE now(), PRIMARY KEY (user, deck_id))")
	// if err != nil {
	// 	return nil, err
	// }
	return db, nil
}

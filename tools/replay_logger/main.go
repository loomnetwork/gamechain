package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"

	"encoding/base64"
	"strings"

	"github.com/Jeffail/gabs"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gorilla/websocket"
	"github.com/loomnetwork/zombie_battleground/types/zb"
)

var (
	wsURL string
)

func main() {
	wsURL = os.Getenv("wsURL")
	if len(wsURL) == 0 {
		wsURL = "ws://localhost:9999/queryws"
	}
	log.Printf("wsURL - %s", wsURL)
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
		log.Printf("result: %s", result.String())
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

				if err := writeReplayFile(topic, body); err != nil {
					log.Println("Error writing replay file: ", err)
				}
			}
		}

		log.Printf("recv: %s", message)
	}
}

func writeReplayFile(topic string, body []byte) error {
	pwd, _ := os.Getwd()
	filename := fmt.Sprintf("replays/%s.json", topic)
	path := filepath.Join(pwd, filename)

	fmt.Println("Writing to file: ", path)

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	var replay zb.GameReplay
	var event zb.PlayerActionEvent
	if fi, _ := f.Stat(); fi.Size() > 0 {
		if err := jsonpb.Unmarshal(f, &replay); err != nil {
			log.Println(err)
			return err
		}
	} else {
		replay.Events = []*zb.PlayerActionEvent{}
	}

	if err := jsonpb.UnmarshalString(string(body), &event); err != nil {
		return err
	}
	replay.Events = append(replay.Events, &event)

	f.Close()

	m := jsonpb.Marshaler{}
	result, err := m.MarshalToString(&replay)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(path, []byte(result), 0644); err != nil {
		return err
	}

	return nil
}

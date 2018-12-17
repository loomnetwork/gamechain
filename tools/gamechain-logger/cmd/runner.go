package cmd

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/jinzhu/gorm"
	"github.com/loomnetwork/gamechain/battleground"
	"github.com/loomnetwork/go-loom/client"
	"github.com/loomnetwork/go-loom/plugin/types"
)

type Runner struct {
	db     *gorm.DB
	eventC chan *types.EventData
	stopC  chan struct{}
	errC   chan error
	wsURL  string
}

func NewRunner(wsURL string, db *gorm.DB, n int) *Runner {
	return &Runner{
		wsURL:  wsURL,
		db:     db,
		stopC:  make(chan struct{}),
		errC:   make(chan error),
		eventC: make(chan *types.EventData, n),
	}
}

func (r *Runner) Start() {
	go r.watchTopic()
	go r.processEvent()
}

func (r *Runner) Stop() {
	close(r.stopC)
}

func (r *Runner) Error() chan error {
	return r.errC
}

func (r *Runner) watchTopic() {
	log.Printf("connecting to %s", r.wsURL)
	conn, err := connectGamechain(r.wsURL)
	if err != nil {
		select {
		case r.errC <- err:
			return
		}
	}
	defer conn.Close()

	log.Printf("connected to %s", r.wsURL)
	log.Printf("watching events from %s", r.wsURL)
	var unmarshaler jsonpb.Unmarshaler
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("error reading from websocket:", err)
			return
		}

		var resp client.RPCResponse
		if err := json.Unmarshal(message, &resp); err != nil {
			log.Println("error parsing jsonrpc response", err)
			continue
		}

		var eventData types.EventData
		if err = unmarshaler.Unmarshal(bytes.NewBuffer(resp.Result), &eventData); err != nil {
			log.Println("error parsing event data", err)
			continue
		}

		// only zombiebattleground smart contract
		if !strings.HasPrefix(eventData.PluginName, "zombiebattleground") {
			continue
		}

		select {
		case r.eventC <- &eventData:
		case <-r.stopC:
			return
		}
	}
}

func (r *Runner) processEvent() {
	for {
		select {
		case eventData := <-r.eventC:
			for _, topic := range eventData.Topics {
				var topicHandler TopicHandler
				switch topic {
				case battleground.TopicFindMatchEvent:
					topicHandler = FindMatchHandler
				case battleground.TopicAcceptMatchEvent:
					topicHandler = AcceptMatchHandler
				case battleground.TopicCreateDeckEvent:
					topicHandler = CreateDeckHandler
				case battleground.TopicEditDeckEvent:
					topicHandler = EditDeckHandler
				case battleground.TopicDeleteDeckEvent:
					topicHandler = DeleteDeckHandler
				default:
					if strings.HasPrefix(topic, "match:") {
						topicHandler = MatchHandler
					}
				}

				if topicHandler != nil {
					err := topicHandler(eventData, r.db)
					if err != nil {
						log.Println("error calling topic handler:", err)
					}
				}
			}
		case <-r.stopC:
			return
		}
	}
}

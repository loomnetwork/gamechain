package cmd

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	"github.com/loomnetwork/gamechain/battleground"
	"github.com/loomnetwork/go-loom/client"
	"github.com/loomnetwork/go-loom/plugin/types"
	"github.com/pkg/errors"
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
	go r.processEvent()
	for {
		// delay before connecting again
		time.Sleep(500 * time.Millisecond)
		log.Printf("connecting to %s", r.wsURL)
		conn, err := connectGamechain(r.wsURL)
		if err != nil {
			log.Println(err)
			continue
		}
		defer conn.Close()
		r.watchTopic(conn)
	}
}

func (r *Runner) Stop() {
	close(r.stopC)
}

func (r *Runner) Error() chan error {
	return r.errC
}

func (r *Runner) watchTopic(conn *websocket.Conn) error {
	log.Printf("connected to %s", r.wsURL)
	log.Printf("watching events from %s", r.wsURL)
	var unmarshaler jsonpb.Unmarshaler
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			return errors.Wrapf(err, "error reading from websocket")
		}

		var resp client.RPCResponse
		if err := json.Unmarshal(message, &resp); err != nil {
			return errors.Wrapf(err, "error parsing jsonrpc response")
		}

		var eventData types.EventData
		if err = unmarshaler.Unmarshal(bytes.NewBuffer(resp.Result), &eventData); err != nil {
			return errors.Wrapf(err, "error parsing event data")
		}

		// only zombiebattleground smart contract
		if !strings.HasPrefix(eventData.PluginName, "zombiebattleground") {
			continue
		}

		select {
		case r.eventC <- &eventData:
		case <-r.stopC:
			return nil
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

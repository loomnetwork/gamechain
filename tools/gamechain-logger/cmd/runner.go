package cmd

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"
	"time"

	raven "github.com/getsentry/raven-go"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/jinzhu/gorm"
	"github.com/loomnetwork/gamechain/battleground"
	"github.com/loomnetwork/go-loom/client"
	"github.com/loomnetwork/go-loom/plugin/types"
	"github.com/loomnetwork/loomauth/models"
	"github.com/pkg/errors"
)

type Runner struct {
	db                *gorm.DB
	eventC            chan *types.EventData
	stopC             chan struct{}
	errC              chan error
	URL               string
	URLType           string
	reconnectInterval time.Duration
	blockInterval     int
	contractName      string
}

func NewRunner(URL string, URLType string, db *gorm.DB, n int, reconnectInterval time.Duration, blockInterval int, contractName string) *Runner {
	return &Runner{
		URL:               URL,
		URLType:           URLType,
		db:                db,
		stopC:             make(chan struct{}),
		errC:              make(chan error),
		eventC:            make(chan *types.EventData, n),
		reconnectInterval: reconnectInterval,
		blockInterval:     blockInterval,
		contractName:      contractName,
	}
}

// Start runs the loop to watch topic. It's a blocking call.
func (r *Runner) Start() {
	go r.processEvent()
	for {
		err := r.watchTopic()
		if err == nil {
			break
		}
		log.Printf("error: %v", err)
		raven.CaptureErrorAndWait(err, map[string]string{})
		// delay before connecting again
		time.Sleep(r.reconnectInterval)
	}
}

func (r *Runner) Stop() {
	close(r.stopC)
}

func (r *Runner) Error() chan error {
	return r.errC
}

func (r *Runner) watchTopic() error {
	if r.URLType == "ev" {
		ticker := time.NewTicker(r.reconnectInterval)
		for {
			select {
			case <-ticker.C:
				height := models.ZbHeightCheck{}
				err := r.db.Where(&models.ZbHeightCheck{Key: 1}).First(&height).Error
				if err != nil && !gorm.IsRecordNotFoundError(err) {
					return err
				}
				fromBlock := height.LastBlockHeight + 1
				toBlock := fromBlock + uint64(r.blockInterval) - 1
				result, err := queryEventStore(r.URL, fromBlock, toBlock, r.contractName)
				if err != nil {
					return err
				}

				var newBlockHeight uint64

				for _, ev := range result.Events {
					r.eventC <- ev
				}
				newBlockHeight = result.ToBlock

				if newBlockHeight > 0 {
					err = UpdateBlockHeight(r.db, newBlockHeight)
					if err != nil {
						return err
					}
				}
			case <-r.stopC:
				ticker.Stop()
				return nil
			}

		}
	} else {
		// Websocket
		log.Printf("connecting to chain %s", r.URL)
		conn, err := connectGamechain(r.URL)
		if err != nil {
			return err
		}
		defer conn.Close()

		log.Printf("connected to %s", r.URL)
		log.Printf("watching events from %s", r.URL)
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
						err = errors.Wrapf(err, "error calling topic handler")
						log.Println(err)
						log.Printf("event: %+v", eventData)
						raven.CaptureErrorAndWait(err, map[string]string{})
					}
				}
			}
		case <-r.stopC:
			return
		}
	}
}

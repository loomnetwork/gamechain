PKG = github.com/loomnetwork/gamechain
PKG_BATTLEGROUND = $(PKG)/battleground

GIT_SHA = `git rev-parse --verify HEAD`
BUILD_DATE = `date -Iseconds`

PROTOC = protoc --plugin=./protoc-gen-gogo -I. -I$(GOPATH)/src -I/usr/local/include
PLUGIN_DIR = $(GOPATH)/src/github.com/loomnetwork/go-loom
GOGO_PROTOBUF_DIR = $(GOPATH)/src/github.com/gogo/protobuf
LOOMCHAIN_DIR = $(GOPATH)/src/github.com/loomnetwork/loomchain
HASHICORP_DIR = $(GOPATH)/src/github.com/hashicorp/go-plugin
PROMETHEUS_PROCFS_DIR = $(GOPATH)/src/github.com/prometheus/procfs

GOFLAGS_BASE = -X $(PKG_BATTLEGROUND).BuildDate=$(BUILD_DATE) -X $(PKG_BATTLEGROUND).BuildGitSha=$(GIT_SHA) -X $(PKG_BATTLEGROUND).BuildNumber=$(BUILD_NUMBER)
GOFLAGS = -ldflags "$(GOFLAGS_BASE)"

all: build-ext cli

build: contracts/zombiebattleground.so.1.0.0

build-ext: contracts/zombiebattleground.1.0.0

cli: bin/zb-cli

oracle: bin/gcoracle

tools: bin/zb-enum-gen bin/zb-console-game

gamechain-logger: proto bin/gamechain-logger

gamechain-replay: proto bin/gamechain-replay

gamechain-debugger: bin/zb-cli bin/gamechain-debugger

bin/zb-cli:
	go build -o $@ $(PKG)/cli

bin/zb-enum-gen:
	go build -o $@ tools/cmd/templates/main.go

bin/zb-console-game:
	go build -o $@ tools/cmd/console_game/main.go

bin/gamechain-logger:
	go build -o $@ $(PKG)/tools/gamechain-logger

bin/gamechain-replay:
	go build -o $@ $(PKG)/tools/gamechain-replay

bin/gamechain-debugger:
	packr2 build -o $@ $(PKG)/tools/gamechain-debugger

bin/gcoracle:
	go build -o $@ $(PKG)/tools/gcoracle

contracts/zombiebattleground.so.1.0.0: proto
	go build $(GOFLAGS) -buildmode=plugin -o $@ $(PKG)/plugin

contracts/zombiebattleground.1.0.0: proto
	go build $(GOFLAGS) -o $@ $(PKG)/plugin

protoc-gen-gogo:
	go build github.com/gogo/protobuf/protoc-gen-gogo

%.pb.go: %.proto protoc-gen-gogo
	if [ -e "protoc-gen-gogo.exe" ]; then mv protoc-gen-gogo.exe protoc-gen-gogo; fi
	$(PROTOC) --gogo_out=$(GOPATH)/src $(PKG)/$<

%.cs: %.proto protoc-gen-gogo
	if [ -e "protoc-gen-gogo.exe" ]; then mv protoc-gen-gogo.exe protoc-gen-gogo; fi
	cp $< $<-cs.bak
	grep -vw 'import "github.com/gogo/protobuf/gogoproto/gogo.proto";' $<-cs.bak | sed -e 's/\[[^][]*\]//g;s/option (gogoproto\.[a-z_]*) = false;//g' > $<-cs && rm $<-cs.bak
	$(PROTOC) --csharp_out=./$(dir $<) $(PKG)/$<-cs
	rm $<-cs
	find ./$(dir $<)*.cs -type f -exec sed -i.bak 's/global::Google.Protobuf/global::Loom.Google.Protobuf/g' {} \;
	find ./$(dir $<)*.cs.bak -type f -exec rm {} \;

proto: types/zb/zb_data/zb_data.pb.go \
    types/zb/zb_enums/zb_enums.pb.go \
    types/zb/zb_calls/zb_calls.pb.go \
    types/zb/zb_custombase/zb_custombase.pb.go \
    types/zb/zb_data/zb_data.cs \
    types/zb/zb_enums/zb_enums.cs \
    types/zb/zb_calls/zb_calls.cs \
    types/zb/zb_custombase/zb_custombase.cs \
    types/oracle/oracle.pb.go \
    types/nullable/nullable_pb/nullable.pb.go \
    types/nullable/nullable_pb/nullable.cs

$(PLUGIN_DIR):
	git clone -q git@github.com:loomnetwork/go-loom.git $@

$(LOOMCHAIN_DIR):
	git clone -q git@github.com:loomnetwork/loomchain.git $@

deps: $(PLUGIN_DIR) $(LOOMCHAIN_DIR)
	go get \
		github.com/golang/dep/cmd/dep \
		github.com/spf13/cobra \
		github.com/spf13/viper \
		github.com/pkg/errors \
		github.com/stretchr/testify/assert\
		github.com/hashicorp/go-plugin \
		github.com/google/uuid \
		github.com/grpc-ecosystem/go-grpc-prometheus \
		github.com/prometheus/client_golang/prometheus \
		github.com/loomnetwork/e2e \
		github.com/iancoleman/strcase \
		github.com/jroimartin/gocui \
		github.com/Jeffail/gabs \
		github.com/gorilla/websocket \
		github.com/go-sql-driver/mysql \
		gopkg.in/yaml.v2 \
		github.com/sirupsen/logrus \
		gopkg.in/check.v1 \
		github.com/kr/logfmt \
		github.com/phonkee/go-pubsub \
		github.com/jinzhu/gorm \
		github.com/mattn/go-sqlite3 \
		github.com/dgrijalva/jwt-go \
		github.com/getsentry/raven-go \
		github.com/gobuffalo/packr/v2 \
		github.com/gobuffalo/packr/v2/... \
		github.com/gorilla/mux 
		
	go install github.com/golang/dep/cmd/dep
	# Need loomchain to run e2e test
	rm -rf $(PROMETHEUS_PROCFS_DIR)
	cd $(LOOMCHAIN_DIR) && make deps && make && cp loom $(GOPATH)/bin


abigen:
	go build github.com/ethereum/go-ethereum/cmd/abigen
	mkdir tmp_build || true
	# Need to run truffle compile and compile over latest ABI for a zombie battleground solidity mode
	cat ./ethcontract/zbgame_mode.json | jq '.abi' > ./tmp_build/eth_game_mode_contract.abi
	./abigen --abi ./tmp_build/eth_game_mode_contract.abi --pkg ethcontract --type ZGCustomGameMode --out ethcontract/zb_gamemode.go

oracle-abigen:
	go build github.com/ethereum/go-ethereum/cmd/abigen
	./abigen --abi oracle/abi/ZBGCardABI.json --pkg ethcontract --type ZBGCard --out oracle/ethcontract/zbgcard.go

test:
	#TODO fix go vet in tests
	go test -timeout=30m -vet=off -v ./... -tags evm

clean:
	go clean
	rm -f \
		protoc-gen-gogo \
		types/zb/zb.pb.go \
		types/zb/Zb.cs \
		types/zb/zb_data/zb_data.pb.go \
		types/zb/zb_enums/zb_enums.pb.go \
		types/zb/zb_calls/zb_calls.pb.go \
		types/zb/zb_custombase/zb_custombase.pb.go \
		types/zb/zb_data/ZbData.cs \
		types/zb/zb_calls/ZbCalls.cs \
		types/zb/zb_enums/ZbEnums.cs \
		types/zb/zb_custombase/ZbCustombase.cs \
		types/oracle/oracle.pb.go \
		types/nullable/nullable_pb/nullable.pb.go \
		types/nullable/nullable_pb/Nullable.cs \
		types/nullable/nullable_test_pb/nullable_test.pb.go \
		contracts/zombiebattleground.so.1.0.0 \
		contracts/zombiebattleground.1.0.0 \
		bin/zb-cli \
		bin/zb-enum-gen \
		bin/gamechain-logger \
		bin/gamechain-replay \
		bin/gamechain-debugger


.PHONY: all clean test deps proto cli zb_console_game tools bin/zb-enum-gen bin/gamechain-logger abigen bin/gcoracle oracle-abigen

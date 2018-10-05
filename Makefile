PKG = github.com/loomnetwork/zombie_battleground
GIT_SHA = `git rev-parse --verify HEAD`
PROTOC = protoc --plugin=./protoc-gen-gogo -I. -Ivendor -I$(GOPATH)/src -I/usr/local/include
PLUGIN_DIR = $(GOPATH)/src/github.com/loomnetwork/go-loom
GOGO_PROTOBUF_DIR = $(GOPATH)/src/github.com/gogo/protobuf
LOOMCHAIN_DIR = $(GOPATH)/src/github.com/loomnetwork/loomchain

all: build cli

build: contracts/zombiebattleground.1.0.0

cli: bin/zb-cli

enum_gen: bin/zb-enum-gen

replay_logger: bin/replay-logger

bin/zb-cli:
	go build -o $@ $(PKG)/cli

bin/zb-enum-gen:
	go build -o $@ $(PKG)/tools/cmd

bin/replay-logger:
	go build -o $@ $(PKG)/tools/replay_logger

contracts/zombiebattleground.1.0.0: proto
	go build -o $@ $(PKG)/plugin

protoc-gen-gogo:
	go build github.com/gogo/protobuf/protoc-gen-gogo

%.pb.go: %.proto protoc-gen-gogo
	if [ -e "protoc-gen-gogo.exe" ]; then mv protoc-gen-gogo.exe protoc-gen-gogo; fi
	$(PROTOC) --gogo_out=$(GOPATH)/src $(PKG)/$<

%.cs: %.proto protoc-gen-gogo
	if [ -e "protoc-gen-gogo.exe" ]; then mv protoc-gen-gogo.exe protoc-gen-gogo; fi
	$(PROTOC) --csharp_out=./types/zb $(PKG)/$<
	sed -i.bak 's/global::Google.Protobuf/global::Loom.Google.Protobuf/g' ./types/zb/Zb.cs && rm ./types/zb/Zb.cs.bak

proto: types/zb/zb.pb.go types/zb/zb.cs

$(PLUGIN_DIR):
	git clone -q git@github.com:loomnetwork/go-loom.git $@

$(LOOMCHAIN_DIR):
	git clone -q git@github.com:loomnetwork/loomchain.git $@

deps: $(PLUGIN_DIR) $(LOOMCHAIN_DIR)
	go get \
		github.com/golang/dep/cmd/dep \
		github.com/gogo/protobuf/jsonpb \
		github.com/gogo/protobuf/proto \
		github.com/spf13/cobra \
		github.com/pkg/errors \
		github.com/stretchr/testify/assert\
		github.com/hashicorp/go-plugin \
		github.com/google/uuid \
		github.com/grpc-ecosystem/go-grpc-prometheus \
		github.com/prometheus/client_golang/prometheus \
		github.com/loomnetwork/e2e \
		github.com/iancoleman/strcase \
		github.com/Jeffail/gabs \
		github.com/gorilla/websocket \
		github.com/sirupsen/logrus
	go install github.com/golang/dep/cmd/dep
	cd $(LOOMCHAIN_DIR) && make deps && make && cp loom $(GOPATH)/bin
	cd $(GOGO_PROTOBUF_DIR) && git checkout 1ef32a8b9fc3f8ec940126907cedb5998f6318e4

test:
	go test -v ./...

clean:
	go clean
	rm -f \
		protoc-gen-gogo \
		types/zb/zb.pb.go \
		types/zb/Zb.cs \
		contracts/zombiebattleground.1.0.0 \
		bin/zb-cli \
		bin/zb-enum-gen \
		bin/replay-logger

.PHONY: all clean test deps proto cli

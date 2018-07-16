PKG = github.com/loomnetwork/zombie_battleground
GIT_SHA = `git rev-parse --verify HEAD`
PROTOC = protoc --plugin=./protoc-gen-gogo -I. -Ivendor -I$(GOPATH)/src -I/usr/local/include
PLUGIN_DIR = $(GOPATH)/src/github.com/loomnetwork/go-loom

all: build cli

build: contracts/zombiebattleground.1.0.0 contracts/zombiebattleground-card.1.0.0

cli: bin/zb-cli

bin/zb-cli:
	go build -o $@ $(PKG)/cli

contracts/zombiebattleground.1.0.0: proto
	go build -o $@ $(PKG)/plugin

contracts/zombiebattleground-card.1.0.0: proto
	go build -o $@ $(PKG)/card/plugin

protoc-gen-gogo:
	go build github.com/gogo/protobuf/protoc-gen-gogo

%.pb.go: %.proto protoc-gen-gogo
	if [ -e "protoc-gen-gogo.exe" ]; then mv protoc-gen-gogo.exe protoc-gen-gogo; fi
	$(PROTOC) --gogo_out=$(GOPATH)/src $(PKG)/$<

%.cs: %.proto protoc-gen-gogo
	if [ -e "protoc-gen-gogo.exe" ]; then mv protoc-gen-gogo.exe protoc-gen-gogo; fi
	$(PROTOC) --csharp_out=./types/zb $(PKG)/$<

proto: types/zb/zb.pb.go types/zb/zb.cs 

deps: $(PLUGIN_DIR)
	cd $(PLUGIN_DIR) && git pull
	go get \
		github.com/loomnetwork/go-loom \
		github.com/gogo/protobuf/jsonpb \
		github.com/gogo/protobuf/proto \
		github.com/spf13/cobra \
		github.com/pkg/errors \
    	github.com/hashicorp/go-plugin \
		github.com/google/uuid

clean:
	go clean
	rm -f \
		protoc-gen-gogo \
		types/zb/zb.pb.go \
		types/zb/Zb.cs \
		contracts/zombiebattleground.so.1.0.0 \
		contracts/zombiebattleground-card.so.1.0.0 \
		bin/zb-cli

.PHONY: all clean test deps proto cli

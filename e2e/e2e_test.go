package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/loomnetwork/e2e/common"
	"github.com/stretchr/testify/assert"
)

func setupInternalPlugin(dir string) error {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	binary, err := exec.LookPath("go")
	if err != nil {
		return err
	}

	cmd := exec.Cmd{
		Path: binary,
		Args: []string{binary, "build", "-buildmode", "plugin", "-o", path.Join(dir, "zombiebattleground.so.1.0.0"), "github.com/loomnetwork/gamechain/plugin"},
	}
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func rune2e(t *testing.T, name string, testFile string, validators int, accounts int, genFile string) {
	singlenode, _ := strconv.ParseBool(os.Getenv("SINGLENODE"))
	// skip multi-node tests?
	if singlenode && validators > 1 {
		return
	}

	config, err := common.NewConfig(name, testFile, genFile, validators, accounts)
	if err != nil {
		t.Fatal(err)
	}
	binary, err := exec.LookPath("go")
	if err != nil {
		t.Fatal(err)
	}
	// required binary
	cmd := exec.Cmd{
		Dir:  config.BaseDir,
		Path: binary,
		Args: []string{binary, "build", "-o", "zb-cli", "github.com/loomnetwork/gamechain/cli"},
	}
	if err := cmd.Run(); err != nil {
		t.Fatal(fmt.Errorf("fail to execute command: %s\n%v", strings.Join(cmd.Args, " "), err))
	}

	if err := common.DoRun(*config); err != nil {
		t.Fatal(err)
	}

	// pause before running the next test
	time.Sleep(3000 * time.Millisecond)

	// clean up test data if successful
	os.RemoveAll(config.BaseDir)

}

func TestE2EAccountDeck(t *testing.T) {
	tests := []struct {
		name       string
		testFile   string
		validators int
		accounts   int
		genFile    string
	}{
		{"zb-account-1", "test_account.toml", 1, 10, "e2e.genesis.json"},
		{"zb-account-4", "test_account.toml", 4, 10, "e2e.genesis.json"},
		{"zb-deck-1", "test_deck.toml", 1, 10, "e2e.genesis.json"},
		{"zb-deck-4", "test_deck.toml", 4, 10, "e2e.genesis.json"},
		{"zb-overlord-1", "test_overlord.toml", 1, 10, "e2e.genesis.json"},
		{"zb-overlord-4", "test_overlord.toml", 4, 10, "e2e.genesis.json"},
	}

	// required to have loom binary
	common.LoomPath = "loom"
	common.ContractDir = "./contracts"
	// required internal contract to resolve port conflicts
	err := setupInternalPlugin(common.ContractDir)
	assert.Nil(t, err)

	for _, test := range tests {
		rune2e(t, test.name, test.testFile, test.validators, test.accounts, test.genFile)
	}
}

func TestE2EMatchMaking(t *testing.T) {
	tests := []struct {
		name       string
		testFile   string
		validators int
		accounts   int
		genFile    string
	}{
		{"zb-findmatch-1", "test_findmatch.toml", 1, 10, "e2e.genesis.json"},
		{"zb-findmatch-4", "test_findmatch.toml", 4, 10, "e2e.genesis.json"},
	}

	// required to have loom binary
	common.LoomPath = "loom"
	common.ContractDir = "./contracts"
	// required internal contract to resolve port conflicts
	err := setupInternalPlugin(common.ContractDir)
	assert.Nil(t, err)

	for _, test := range tests {
		rune2e(t, test.name, test.testFile, test.validators, test.accounts, test.genFile)
	}
}

func TestE2EGameplay(t *testing.T) {
	tests := []struct {
		name       string
		testFile   string
		validators int
		accounts   int
		genFile    string
	}{
		{"zb-gameplay-1", "test_gameplay.toml", 1, 10, "e2e.genesis.json"},
		//{"zb-gameplay-4", "test_gameplay.toml", 4, 10, "e2e.genesis.json"},
	}

	// required to have loom binary
	common.LoomPath = "loom"
	common.ContractDir = "./contracts"
	// required internal contract to resolve port conflicts
	err := setupInternalPlugin(common.ContractDir)
	assert.Nil(t, err)

	for _, test := range tests {
		rune2e(t, test.name, test.testFile, test.validators, test.accounts, test.genFile)
	}
}

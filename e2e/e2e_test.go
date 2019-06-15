package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/loomnetwork/e2e/common"
	assert "github.com/stretchr/testify/require"
)

type testdata struct {
	name       string
	testFile   string
	validators int
	accounts   int
	genFile    string
	replacementTokens map[string]string
}

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

func prepareTestFile(t *testing.T, testFile string, replacementTokens map[string]string) string {
	data, err := ioutil.ReadFile(testFile)
	assert.Nil(t, err)
	str := string(data)

	if len(replacementTokens) > 0 {
		for tokenKey, tokenValue := range replacementTokens {
			tokenKey = "{{" + tokenKey + "}}"

			str = strings.Replace(str, tokenKey, tokenValue, -1)
		}
	}

	_, testFileName := filepath.Split(testFile)

	tempFile, err := ioutil.TempFile("", testFileName)
	assert.Nil(t, err)
	_, err = tempFile.WriteString(str)
	assert.Nil(t, err)
	defer tempFile.Close()

	return tempFile.Name()
}

func runE2ETests(t *testing.T, tests []testdata) {
	// required to have loom binary
	common.LoomPath = "loom"
	common.ContractDir = "./contracts"
	// required internal contract to resolve port conflicts
	err := setupInternalPlugin(common.ContractDir)
	assert.Nil(t, err)

	for _, test := range tests {
		tempTestFilePath := prepareTestFile(t, test.testFile, test.replacementTokens)

		//noinspection ALL
		defer os.Remove(tempTestFilePath)
		rune2e(t, test.name, tempTestFilePath, test.validators, test.accounts, test.genFile)
	}
}

func TestE2EAccount(t *testing.T) {
	tests := []testdata {
		{"zb-account-1", "test_account.toml", 1, 10, "../test_data/simple-genesis.json", nil},
		{"zb-account-4", "test_account.toml", 4, 10, "../test_data/simple-genesis.json", nil},
	}

	runE2ETests(t, tests)
}

func TestE2EDeck(t *testing.T) {
	tests := []testdata {
		{"zb-deck-1", "test_deck.toml", 1, 10, "../test_data/simple-genesis.json", nil},
		{"zb-deck-4", "test_deck.toml", 4, 10, "../test_data/simple-genesis.json", nil},
	}

	runE2ETests(t, tests)
}

func TestE2EOverlord(t *testing.T) {
	tests := []testdata {
		{"zb-overlord-1", "test_overlord.toml", 1, 10, "../test_data/simple-genesis.json", nil},
		{"zb-overlord-4", "test_overlord.toml", 4, 10, "../test_data/simple-genesis.json", nil},
	}

	runE2ETests(t, tests)
}

func TestE2EMatchMaking(t *testing.T) {
	tests := []testdata {
		{"zb-findmatch-1", "test_findmatch.toml", 1, 10, "../test_data/simple-genesis.json", nil},
		{"zb-findmatch-4", "test_findmatch.toml", 4, 10, "../test_data/simple-genesis.json", nil},
	}

	runE2ETests(t, tests)
}

func TestE2EGameplay(t *testing.T) {
	replacementTokens := map[string]string{
		"BackendLogicEnabled": "false",
	}
	tests := []testdata{
		{"zb-gameplay-1", "test_gameplay.toml", 1, 10, "../test_data/simple-genesis.json", replacementTokens},
		{"zb-gameplay-1", "test_gameplay.toml", 4, 10, "../test_data/simple-genesis.json", replacementTokens},
	}

	runE2ETests(t, tests)
}

func TestE2EGameplayBackendLogic(t *testing.T) {
	replacementTokens := map[string]string{
		"BackendLogicEnabled": "true",
	}
	tests := []testdata{
		{"zb-gameplay-1", "test_gameplay.toml", 1, 10, "../test_data/simple-genesis.json", replacementTokens},
		{"zb-gameplay-1", "test_gameplay.toml", 4, 10, "../test_data/simple-genesis.json", replacementTokens},
	}

	runE2ETests(t, tests)
}
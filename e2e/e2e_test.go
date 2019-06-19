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

var (
	isZbCliBuilt = false
	isGamechainBuilt = false
)

func getLoomchainDir() string {
	goPath := os.Getenv("GOPATH")
	return path.Join(goPath, "src", "github.com/loomnetwork/loomchain")
}

func buildGamechain(t *testing.T) {
	if isGamechainBuilt {
		return
	}

	// build loomchain with gamechain compiled in
	if err := runCommand(getLoomchainDir(), "make", "gamechain"); err != nil {
		t.Fatal(err)
	}

	isGamechainBuilt = true
}

func buildZbCli( t *testing.T) {
	if isZbCliBuilt {
		return
	}

	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	if err := runCommand(filepath.Dir(currentDir), "make", "cli", "-B"); err != nil {
		t.Fatal(err)
	}

	isZbCliBuilt = true
}

func runCommand(workingDir, binary string, args ...string) error {
	binary, err := exec.LookPath(binary)
	if err != nil {
		return err
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	//noinspection ALL
	defer os.Chdir(currentDir)
	err = os.Chdir(workingDir)
	if err != nil {
		return err
	}

	args = append([]string{binary}, args...)

	cmd := exec.Cmd{
		Path: binary,
		Args: args,
	}
	combinedOutput, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("fail to execute command: %s\n%v\n%s", strings.Join(cmd.Args, " "), err, string(combinedOutput))
	}

	return nil
}

func runE2E(t *testing.T, name string, testFile string, validators int, accounts int, genFile string) {
	singleNode, _ := strconv.ParseBool(os.Getenv("SINGLENODE"))

	// skip multi-node tests?
	if singleNode && validators > 1 {
		return
	}

	config, err := common.NewConfig(name, testFile, genFile, validators, accounts)
	if err != nil {
		t.Fatal(err)
	}

	buildZbCli(t)

	// copy zb-cli to basedir
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	zbCliBytes, err := ioutil.ReadFile(path.Join(currentDir, "..", "bin", "zb-cli"))
	assert.Nil(t, err)
	err = ioutil.WriteFile(path.Join(config.BaseDir, "zb-cli"), zbCliBytes, os.ModePerm)
	assert.Nil(t, err)

	if err := common.DoRun(*config); err != nil {
		t.Fatal(err)
	}

	// pause before running the next test
	time.Sleep(3000 * time.Millisecond)

	// clean up test data if successful
	//noinspection ALL
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
	//noinspection ALL
	defer tempFile.Close()

	return tempFile.Name()
}

func runE2ETests(t *testing.T, tests []testdata) {
	// required to have loom binary
	common.LoomPath = path.Join(getLoomchainDir(), "gamechain")
	common.ContractDir = "./contracts"

	buildGamechain(t)

	// create ContractDir to make e2e happy
	if err := os.MkdirAll(common.ContractDir, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	for _, test := range tests {
		tempTestFilePath := prepareTestFile(t, test.testFile, test.replacementTokens)

		//noinspection ALL
		defer os.Remove(tempTestFilePath)
		runE2E(t, test.name, tempTestFilePath, test.validators, test.accounts, test.genFile)
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
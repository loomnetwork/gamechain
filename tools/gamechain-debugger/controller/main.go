package controller

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gorilla/mux"
	"github.com/loomnetwork/gamechain/tools/gamechain-debugger/view"
)

type MainController struct {
	gamechainDebugger   *view.View
	clientStateDebugger *view.View
	cliDebugger         *view.View
	privateKeyFilePath  string
	cliFilePath         string
}

func NewMainController(cliFilePath string, privateKeyFilePath string) *MainController {
	mc := MainController{
		gamechainDebugger:   view.NewView("base", "gamechain_debugger.html"),
		clientStateDebugger: view.NewView("base", "client_state_debugger.html"),
		cliDebugger:         view.NewView("base", "cli_debugger.html"),
		cliFilePath:         cliFilePath,
		privateKeyFilePath:  privateKeyFilePath,
	}
	return &mc
}

func WriteFile(filename string, fileContent []byte) error {
	err := ioutil.WriteFile(filename, fileContent, 0777)
	if err != nil {
		return err
	}
	return nil
}

func (MainController *MainController) RunMatchCmd(matchId string, cmdName string) ([]byte, error) {
	cmd := exec.Command(MainController.cliFilePath, cmdName, "-k", MainController.privateKeyFilePath, "-m", matchId, "-O", "json")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println(cmd.Stderr)
		return out.Bytes(), err
	}

	return out.Bytes(), nil
}

func ZipFiles(filename string, files []string) error {

	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	for _, file := range files {
		if err = AddFileToZip(zipWriter, file); err != nil {
			return err
		}
	}
	return nil
}

func AddFileToZip(zipWriter *zip.Writer, filename string) error {

	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Name = filename
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}

func (MainController *MainController) GamechainDebugger(w http.ResponseWriter, r *http.Request) {
	MainController.gamechainDebugger.Render(w, nil)
}
func (MainController *MainController) ClientStateDebugger(w http.ResponseWriter, r *http.Request) {
	MainController.clientStateDebugger.Render(w, nil)
}
func (MainController *MainController) CliDebugger(w http.ResponseWriter, r *http.Request) {
	MainController.cliDebugger.Render(w, nil)
}
func (MainController *MainController) RunCli(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["args"]
	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'args' is missing")
		return
	}

	keyArray := strings.Split(keys[0], " ")
	args := []string{}
	args = append(args, MainController.cliFilePath)
	args = append(args, keyArray...)

	cmd := exec.Command(strings.Join(args, " "))
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println(cmd.Stderr)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out.Bytes())

}
func (MainController *MainController) SaveState(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	MatchId := vars["MatchId"]

	gameState, err := MainController.RunMatchCmd(MatchId, "get_game_state")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(err.Error()))
		return
	}
	initState, err := MainController.RunMatchCmd(MatchId, "get_initial_game_state")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(err.Error()))
		return
	}
	match, err := MainController.RunMatchCmd(MatchId, "get_match")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	gameStateFilename := "./tmp/game_state.json"
	initStateFilename := "./tmp/init_game_state.json"
	matchFilename := "./tmp/match.json"

	os.Mkdir("./tmp", 0777)
	if err := WriteFile(initStateFilename, initState); err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	if err := WriteFile(gameStateFilename, gameState); err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	if err := WriteFile(matchFilename, match); err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	if err := ZipFiles("match_data.zip", []string{gameStateFilename, initStateFilename, matchFilename}); err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	os.RemoveAll("./tmp/")

	dbFile, err := ioutil.ReadFile("match_data.zip")
	if err != nil {
		fmt.Fprintf(w, "%s", err)
	}
	b := bytes.NewBuffer(dbFile)
	filename := "match_data_" + MatchId + ".zip"
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	if _, err := b.WriteTo(w); err != nil {
		fmt.Fprintf(w, "%s", err)
	}

}
func (MainController *MainController) GetState(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	MatchId := vars["MatchId"]
	output, err := MainController.RunMatchCmd(MatchId, "get_game_state")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}

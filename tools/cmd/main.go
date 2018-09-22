package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"
)

type Ability struct {
	Name  string
	Types []string
}

func (seq *Ability) Next() (s string) {
	s = seq.Types[0]
	if len(seq.Types) > 1 {
		seq.Types = seq.Types[1:]
	}
	return
}

//FileStruct abilities dataset from json
type FileStruct struct {
	Abilities []Ability
}

var fns = template.FuncMap{
	"last": func(x int, a []string) bool {
		fmt.Printf("x -%d len-%v\n", x, len(a))
		return x == len(a)-1
	},
}

func outputTemplate(outputfile string, ab *FileStruct, templateBaseName, templatefile string) {
	filename := outputfile

	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	t, err := template.New(templateBaseName).Funcs(fns).ParseFiles(templatefile)
	if err != nil {
		panic(err)
	}
	fmt.Printf("outputfile file -%s\n", filename)
	err = t.Execute(f, ab)
	if err != nil {
		panic(err)
	}

	f.Sync()
}
func main() {
	ab := &FileStruct{}

	byteValue, err := ioutil.ReadFile("tools/cmd/templates/inputdata.json")
	if err != nil {
		panic(err)
	}

	json.Unmarshal(byteValue, &ab)

	outputTemplate("Enumerators.cs", ab, "csharpabilities.parse", "tools/cmd/templates/csharpabilities.parse")
}

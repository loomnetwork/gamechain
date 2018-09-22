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

func main() {
	ab := &FileStruct{}

	byteValue, err := ioutil.ReadFile("tools/cmd/templates/inputdata.json")
	if err != nil {
		panic(err)
	}

	json.Unmarshal(byteValue, &ab)
	fmt.Printf("ab-%v\n", ab.Abilities[0].Types[0])

	t, err := template.New("csharpabilities.parse").Funcs(fns).ParseFiles("tools/cmd/templates/csharpabilities.parse")
	if err != nil {
		panic(err)
	}
	err = t.Execute(os.Stdout, ab)
	if err != nil {
		panic(err)
	}
}

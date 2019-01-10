package generator

import (
	"flag"
	"fmt"
	"go/types"
	"golang.org/x/tools/go/loader"
	"log"
	"os"
	"strings"
)

type Generator struct {
	targetPackagePath string
	targetPackage string
	protoPackage  string
	outputPath    string
}

func NewGenerator(targetPackagePath string, targetPackage string, protoPackage string, outputPath string) *Generator {
	generator := &Generator{
		targetPackagePath: targetPackagePath,
		targetPackage: targetPackage,
		protoPackage: protoPackage,
		outputPath: outputPath,
	}

	return generator
}

func (generator *Generator) Generate() {
	var conf loader.Config

	argsWithProg := os.Args
	fmt.Println(argsWithProg)
	for _, p := range flag.Args() {
		conf.ImportWithTests(p)
	}

	conf.ImportWithTests(argsWithProg[1])
	prog, err := conf.Load()
	if err != nil {
		log.Fatal(err)
	}

	var interestPackages []*types.Package
	for p1, p := range prog.AllPackages {
		if strings.Contains(p1.Name(), "pbgraphserialization") {
			interestPackages = append(interestPackages, p1)
		}
		fmt.Println(p1.String() + ": " + p.String())
	}

	fmt.Println(interestPackages)
}

package generator

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io/ioutil"
)

var generateCmdArgs struct {
	targetPackagePath string
	targetPackageName string
	protoPackageName  string
	outputPath        string
}

var generateCmd = &cobra.Command{
	Use:   "pbgraphserialization-gen",
	Short: "Protobuf graph serialization code generator tool",
	RunE: func(cmd *cobra.Command, args []string) error {
		generator, err := NewGenerator(generateCmdArgs.targetPackagePath, generateCmdArgs.targetPackageName, generateCmdArgs.protoPackageName)
		if err != nil {
			return err
		}

		if len(generator.ProgramLoadErrors) > 0 {
			fmt.Printf("Found %d errors while loading code, ignored\n", len(generator.ProgramLoadErrors))
		}

		err = generator.AddEnabledTypesFromCode()
		if err != nil {
			return err
		}

		code, err := generator.Generate()
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(generateCmdArgs.outputPath, []byte(code), 0)
		if err != nil {
			return errors.Wrap(err, "error while writing output file")
		}

		fmt.Printf("Written generated code to '%s'\n", generateCmdArgs.outputPath)

		return nil
	},
}

func Execute() error {
	generateCmd.PersistentFlags().StringVarP(&generateCmdArgs.targetPackagePath, "targetPackagePath", "", "", "Path to the target package root")
	generateCmd.PersistentFlags().StringVarP(&generateCmdArgs.targetPackageName, "targetPackageName", "", "", "Target package name to generate serialization code for")
	generateCmd.PersistentFlags().StringVarP(&generateCmdArgs.protoPackageName, "protoPackageName", "", "", "Package name of Protobuf-generated code corresponding to the target package")
	generateCmd.PersistentFlags().StringVarP(&generateCmdArgs.outputPath, "outputPath", "", "", "File path to the generated file")

	generateCmd.MarkPersistentFlagRequired("targetPackagePath")
	generateCmd.MarkPersistentFlagRequired("targetPackageName")
	generateCmd.MarkPersistentFlagRequired("protoPackageName")
	generateCmd.MarkPersistentFlagRequired("outputPath")

	return generateCmd.Execute()
}

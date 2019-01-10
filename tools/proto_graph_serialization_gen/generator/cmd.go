package generator

import (
	"github.com/spf13/cobra"
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
		generator := NewGenerator(generateCmdArgs.targetPackagePath, generateCmdArgs.targetPackageName, generateCmdArgs.protoPackageName, generateCmdArgs.outputPath)
		return generator.Generate()
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

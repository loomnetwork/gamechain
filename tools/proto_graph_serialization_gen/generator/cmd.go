package generator

import (
	"github.com/spf13/cobra"
)

var generateCmdArgs struct {
	targetPackagePath string
	targetPackage string
	protoPackage  string
	outputPath    string
}

var generateCmd = &cobra.Command{
	Use:   "proto-graph-pbgraphserialization-gen",
	Short: "Protobuf graph pbgraphserialization code generator tool",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func Execute() error {
	generateCmd.PersistentFlags().StringVarP(&generateCmdArgs.targetPackagePath, "targetPackagePath", "", "", "Path to the target package root")
	generateCmd.PersistentFlags().StringVarP(&generateCmdArgs.targetPackage, "targetPackage", "", "", "Target package name to generate pbgraphserialization code for")
	generateCmd.PersistentFlags().StringVarP(&generateCmdArgs.protoPackage, "protoPackage", "", "", "Package name of Protobuf-generated code corresponding to the target package")
	generateCmd.PersistentFlags().StringVarP(&generateCmdArgs.outputPath, "outputPath", "", "", "File path of the generated file")

	generateCmd.MarkPersistentFlagRequired("targetPackagePath")
	generateCmd.MarkPersistentFlagRequired("targetPackage")
	generateCmd.MarkPersistentFlagRequired("protoPackage")
	generateCmd.MarkPersistentFlagRequired("outputPath")

	return generateCmd.Execute()
}

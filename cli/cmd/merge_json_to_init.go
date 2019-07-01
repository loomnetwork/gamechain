package cmd

import (
	"fmt"
	"github.com/loomnetwork/gamechain/tools/battleground_utility"
	"github.com/loomnetwork/gamechain/types/zb/zb_data"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
)

var mergeJsonToInitCmdArgs struct {
	initJsonTemplateFile string

	defaultDecksFile string
	defaultCollectionFile string
	cardLibraryFile string
	overlordsFile string
	aiDecksFile string
	overlordLevelingFile string

	outputFile string
}

var mergeJsonToInitCmd = &cobra.Command{
	Use:   "merge_json_to_init",
	Short: "merges init data from separate JSON files into a single one",
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		var initData zb_data.InitData
		var defaultDecksData zb_data.DefaultDecksDataContainer
		var defaultCollectionData zb_data.DefaultCollectionDataContainer
		var cardLibraryData zb_data.CardLibraryDataContainer
		var overlordsPrototypesData zb_data.OverlordPrototypesDataContainer
		var aiDecksData zb_data.AIDecksDataContainer
		var overlordLevelingData zb_data.OverlordLevelingDataContainer

		// Read pieces
		err = battleground_utility.ReadJsonFileToProtoMessage(mergeJsonToInitCmdArgs.initJsonTemplateFile, &initData)
		if err != nil {
			return errors.Wrap(err, "error parsing " + mergeJsonToInitCmdArgs.initJsonTemplateFile)
		}

		err = battleground_utility.ReadJsonFileToProtoMessage(mergeJsonToInitCmdArgs.defaultDecksFile, &defaultDecksData)
		if err != nil {
			return errors.Wrap(err, "error parsing " + mergeJsonToInitCmdArgs.defaultDecksFile)
		}

		err = battleground_utility.ReadJsonFileToProtoMessage(mergeJsonToInitCmdArgs.defaultCollectionFile, &defaultCollectionData)
		if err != nil {
			return errors.Wrap(err, "error parsing " + mergeJsonToInitCmdArgs.defaultCollectionFile)
		}

		err = battleground_utility.ReadJsonFileToProtoMessage(mergeJsonToInitCmdArgs.cardLibraryFile, &cardLibraryData)
		if err != nil {
			return errors.Wrap(err, "error parsing " + mergeJsonToInitCmdArgs.cardLibraryFile)
		}

		err = battleground_utility.ReadJsonFileToProtoMessage(mergeJsonToInitCmdArgs.overlordsFile, &overlordsPrototypesData)
		if err != nil {
			return errors.Wrap(err, "error parsing " + mergeJsonToInitCmdArgs.overlordsFile)
		}

		err = battleground_utility.ReadJsonFileToProtoMessage(mergeJsonToInitCmdArgs.aiDecksFile, &aiDecksData)
		if err != nil {
			return errors.Wrap(err, "error parsing " + mergeJsonToInitCmdArgs.aiDecksFile)
		}

		err = battleground_utility.ReadJsonFileToProtoMessage(mergeJsonToInitCmdArgs.overlordLevelingFile, &overlordLevelingData)
		if err != nil {
			return errors.Wrap(err, "error parsing " + mergeJsonToInitCmdArgs.overlordLevelingFile)
		}

		// Merge
		initData.DefaultDecks = defaultDecksData.DefaultDecks
		initData.DefaultCollection = defaultCollectionData.DefaultCollection
		initData.Cards = cardLibraryData.Cards
		initData.Overlords = overlordsPrototypesData.Overlords
		initData.AiDecks = aiDecksData.AiDecks
		initData.OverlordLeveling = overlordLevelingData.OverlordLeveling

		// Write merged file
		mergedFile, err := os.Create(mergeJsonToInitCmdArgs.outputFile)
		if err != nil {
			return errors.Wrap(err, "error creating output file")
		}

		defer func() {
			if err := mergedFile.Close(); err != nil {
				panic(err)
			}
		}()

		err = battleground_utility.PrintProtoMessageAsJson(mergedFile, &initData)
		if err != nil {
			return errors.Wrap(err, "error writing output file")
		}

		fmt.Printf("merged file written to %s\n", mergeJsonToInitCmdArgs.outputFile)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(mergeJsonToInitCmd)

	mergeJsonToInitCmd.Flags().StringVarP(&mergeJsonToInitCmdArgs.initJsonTemplateFile, "initJsonTemplate", "", "init_template.json", "init.json template file without data")

	mergeJsonToInitCmd.Flags().StringVarP(&mergeJsonToInitCmdArgs.defaultDecksFile, "defaultDecks", "", "default_decks.json", "default decks JSON data file")
	mergeJsonToInitCmd.Flags().StringVarP(&mergeJsonToInitCmdArgs.defaultCollectionFile, "defaultCollection", "", "default_collection.json", "default collection JSON data file")
	mergeJsonToInitCmd.Flags().StringVarP(&mergeJsonToInitCmdArgs.cardLibraryFile, "cardLibrary", "", "card_library.json", "card library JSON data file")
	mergeJsonToInitCmd.Flags().StringVarP(&mergeJsonToInitCmdArgs.overlordsFile, "overlords", "", "overlords.json", "overlord prototypes JSON data file")
	mergeJsonToInitCmd.Flags().StringVarP(&mergeJsonToInitCmdArgs.aiDecksFile, "aiDecks", "", "ai_decks.json", "AI decks JSON data file")
	mergeJsonToInitCmd.Flags().StringVarP(&mergeJsonToInitCmdArgs.overlordLevelingFile, "overlordLeveling", "", "overlord_leveling.json", "Overlord leveling JSON data file")

	mergeJsonToInitCmd.Flags().StringVarP(&mergeJsonToInitCmdArgs.outputFile, "outputFile", "o", "update_init_merged.json", "path to the merged output file")
}

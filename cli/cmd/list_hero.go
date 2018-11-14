package cmd

import (
	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var listHeroForUserCmdArgs struct {
	userID string
}

var listHeroForUserCmd = &cobra.Command{
	Use:   "list_hero",
	Short: "list hero for user",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}

		req := zb.ListHeroesRequest{
			UserId: listHeroForUserCmdArgs.userID,
		}
		result := zb.ListHeroesResponse{}

		_, err := commonTxObjs.contract.StaticCall("ListHeroes", &req, callerAddr, &result)
		if err != nil {
			return err
		}

		printProtoMessageAsJsonToStdout(&result)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listHeroForUserCmd)

	listHeroForUserCmd.Flags().StringVarP(&listHeroForUserCmdArgs.userID, "userId", "u", "loom", "UserId of account")
}

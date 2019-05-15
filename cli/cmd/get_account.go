package cmd

import (
	"fmt"
	"strings"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/loomnetwork/gamechain/types/zb"
	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var getAccCmdArgs struct {
	userID string
}

var getAccountCmd = &cobra.Command{
	Use:   "get_account",
	Short: "gets account data for zombiebattleground",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}

		req := &zb_calls.GetAccountRequest{
			UserId: getAccCmdArgs.userID,
		}
		var result zb_data.Account

		_, err := commonTxObjs.contract.StaticCall("GetAccount", req, callerAddr, &result)
		if err != nil {
			return fmt.Errorf("error encountered while calling GetAccount: %s", err.Error())
		}

		switch strings.ToLower(rootCmdArgs.outputFormat) {
		case "json":
			output, err := new(jsonpb.Marshaler).MarshalToString(&result)
			if err != nil {
				return err
			}
			fmt.Println(string(output))
		default:
			fmt.Printf("User: %s\n", result.UserId)
			fmt.Printf("Image: %s\n", result.Image)
			fmt.Printf("Game Membership Tier: %d\n", result.GameMembershipTier)
			fmt.Printf("Elo Score: %d\n", result.EloScore)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getAccountCmd)

	getAccountCmd.Flags().StringVarP(&getAccCmdArgs.userID, "userId", "u", "loom", "UserId of account")
}

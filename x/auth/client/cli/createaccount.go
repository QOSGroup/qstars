package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/QOSGroup/qstars/stub"
	"github.com/QOSGroup/qstars/wire"
)

// GetAccountCmd returns a query account that will display the state of the
// account at a given address.
//, decoder auth.AccountDecoder
func CreateAccountCmd(cdc *wire.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "createaccount",
		Short: "create an account",
		//Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			output := stub.AccountCreateStr()

			fmt.Println(output)
			return nil
		},
	}
}

// Copyright 2018 The QOS Authors

package auth

import (
	"errors"
	"fmt"
	"github.com/QOSGroup/qstars/slim"

	"github.com/QOSGroup/qstars/wire"
	"github.com/spf13/cobra"

)

// CreateAccountCmd returns a new account
func CreateAccountCmd(cdc *wire.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "createaccount password",
		Short: "create an account",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("需要参数password")
			}
			password := args[0]

			acc := slim.AccountCreate(password)

			output, err := wire.MarshalJSONIndent(cdc, acc)
			if err != nil {
				return err
			}

			fmt.Println(string(output))

			return nil
		},
	}
}

// GetAccountCmd returns a query account that will display the state of the
// account at a given address.
func GetAccountCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account [address]",
		Short: "Query account balance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// find the key to look up the account
			addr := args[0]
			acc, err := QueryAccount(cdc, addr)
			if err != nil {
				return err
			}

			output, err := wire.MarshalJSONIndent(cdc, acc)
			if err != nil {
				return err
			}

			fmt.Println(string(output))
			return nil
		},
	}

	return cmd
}

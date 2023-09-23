package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func main() {
	s := NewServer()

	rootCmd := &cobra.Command{}

	createWalletCmd := &cobra.Command{
		Use: "create-wallet",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString("name")
			passphrase, _ := cmd.Flags().GetString("passphrase")
			if len(name) == 0 || len(passphrase) == 0 {
				fmt.Println("name and passphrase are required")
				return
			}

			fmt.Println(name, passphrase)

			//err := s.CreateWallet(name, passphrase)
			//if err != nil {
			//	fmt.Println(err)
			//	return
			//}

			fmt.Println("wallet created")
		},
	}

	mintAssetCmd := &cobra.Command{
		Use: "mint-asset",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString("nameAsset")
			amount, _ := cmd.Flags().GetUint64("amount")
			s.MintAsset(name, amount)
		},
	}

	createWalletCmd.PersistentFlags().String("name", "", "wallet name")
	createWalletCmd.PersistentFlags().String("passphrase", "", "wallet passphrase")
	rootCmd.AddCommand(createWalletCmd)

	mintAssetCmd.PersistentFlags().String("nameAsset", "", "asset name")
	mintAssetCmd.PersistentFlags().Uint64("amount", 0, "asset amount")
	rootCmd.AddCommand(mintAssetCmd)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}

}

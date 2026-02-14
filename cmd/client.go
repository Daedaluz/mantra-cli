package cmd

import (
	"github.com/spf13/cobra"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "client commands",
}

func init() {
	rootCmd.AddCommand(clientCmd)
	clientCmd.PersistentFlags().StringP("domain", "d", envStr("DOMAIN", ""), "domain to manage [$DOMAIN]")
	clientCmd.PersistentFlags().StringP("client-id", "t", envStr("CLIENT_ID", ""), "client id [$CLIENT_ID]")
	clientCmd.PersistentFlags().StringP("client-secret", "p", envStr("CLIENT_SECRET", ""), "client secret [$CLIENT_SECRET]")
}

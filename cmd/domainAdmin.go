package cmd

import (
	"github.com/spf13/cobra"
)

// domainAdminCmd represents the domainAdmin command
var domainAdminCmd = &cobra.Command{
	Use:   "domainAdmin",
	Short: "domain admin commands",
}

func init() {
	rootCmd.AddCommand(domainAdminCmd)
	domainAdminCmd.PersistentFlags().StringP("domain", "d", envStr("DOMAIN", ""), "domain to manage [$DOMAIN]")
	domainAdminCmd.PersistentFlags().StringP("client-id", "t", envStr("CLIENT_ID", ""), "client id [$CLIENT_ID]")
	domainAdminCmd.PersistentFlags().StringP("client-secret", "p", envStr("CLIENT_SECRET", ""), "client secret [$CLIENT_SECRET]")
}

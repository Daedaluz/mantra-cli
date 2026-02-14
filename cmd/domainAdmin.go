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

	domainDefault := ""
	clientIDDefault := ""
	clientSecretDefault := ""
	if activeCtx != nil {
		domainDefault = ctxStr(activeCtx.Domain, domainDefault)
		clientIDDefault = ctxStr(activeCtx.ClientID, clientIDDefault)
		clientSecretDefault = ctxStr(activeCtx.ClientSecret, clientSecretDefault)
	}

	domainAdminCmd.PersistentFlags().StringP("domain", "d", envStr("DOMAIN", domainDefault), "domain to manage [$DOMAIN]")
	domainAdminCmd.PersistentFlags().StringP("client-id", "t", envStr("CLIENT_ID", clientIDDefault), "client id [$CLIENT_ID]")
	domainAdminCmd.PersistentFlags().StringP("client-secret", "p", envStr("CLIENT_SECRET", clientSecretDefault), "client secret [$CLIENT_SECRET]")
}

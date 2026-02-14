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

	domainDefault := ""
	clientIDDefault := ""
	clientSecretDefault := ""
	if activeCtx != nil {
		domainDefault = ctxStr(activeCtx.Domain, domainDefault)
		clientIDDefault = ctxStr(activeCtx.ClientID, clientIDDefault)
		clientSecretDefault = ctxStr(activeCtx.ClientSecret, clientSecretDefault)
	}

	clientCmd.PersistentFlags().StringP("domain", "d", envStr("DOMAIN", domainDefault), "domain to manage [$DOMAIN]")
	clientCmd.PersistentFlags().StringP("client-id", "t", envStr("CLIENT_ID", clientIDDefault), "client id [$CLIENT_ID]")
	clientCmd.PersistentFlags().StringP("client-secret", "p", envStr("CLIENT_SECRET", clientSecretDefault), "client secret [$CLIENT_SECRET]")
}

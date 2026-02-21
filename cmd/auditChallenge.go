package cmd

import (
	"fmt"

	"github.com/daedaluz/mantra-cli/internal"
	admin "github.com/daedaluz/mantra-cli/lib/grpc/service"
	"github.com/spf13/cobra"
)

var auditChallengeCmd = &cobra.Command{
	Use:   "auditChallenge <challenge-id>",
	Short: "Audit a specific challenge",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		domain, _ := cmd.Flags().GetString("domain")

		ctx, api := internal.GetAdminAPIClient(cmd)
		e, err := api.AuditChallenge(ctx, &admin.AuditChallengeRequest{
			Domain:      domain,
			ChallengeId: args[0],
		})
		if err != nil {
			fmt.Printf("Error auditing challenge: %v\n", err)
			return
		}
		printAuditEntry(e)
	},
}

func init() {
	adminCmd.AddCommand(auditChallengeCmd)
	auditChallengeCmd.Flags().StringP("domain", "d", envStr("DOMAIN", ""), "domain to audit [$DOMAIN]")
}

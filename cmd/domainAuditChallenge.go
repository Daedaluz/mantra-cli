package cmd

import (
	"fmt"

	"github.com/daedaluz/mantra-cli/internal"
	admin "github.com/daedaluz/mantra-cli/lib/grpc/service"
	"github.com/spf13/cobra"
)

var domainAuditChallengeCmd = &cobra.Command{
	Use:   "auditChallenge <challenge-id>",
	Short: "Audit a specific challenge in this domain",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx, api := internal.GetDomainAdminAPIClient(cmd)
		e, err := api.AuditChallenge(ctx, &admin.DomainAuditChallengeRequest{
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
	domainAdminCmd.AddCommand(domainAuditChallengeCmd)
}

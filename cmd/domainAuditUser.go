package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/daedaluz/mantra-cli/internal"
	admin "github.com/daedaluz/mantra-cli/lib/grpc/service"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var domainAuditUserCmd = &cobra.Command{
	Use:   "auditUser <user-id>",
	Short: "Audit authentication history for a user in this domain",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		from, _ := cmd.Flags().GetString("from")
		to, _ := cmd.Flags().GetString("to")
		verbose, _ := cmd.Flags().GetBool("verbose")

		req := &admin.DomainAuditUserRequest{
			UserId: args[0],
		}
		if from != "" {
			t, err := time.Parse(time.RFC3339, from)
			if err != nil {
				fmt.Printf("Invalid --from time: %v\n", err)
				return
			}
			req.From = timestamppb.New(t)
		}
		if to != "" {
			t, err := time.Parse(time.RFC3339, to)
			if err != nil {
				fmt.Printf("Invalid --to time: %v\n", err)
				return
			}
			req.To = timestamppb.New(t)
		}

		ctx, api := internal.GetDomainAdminAPIClient(cmd)
		resp, err := api.AuditUser(ctx, req)
		if err != nil {
			fmt.Printf("Error auditing user: %v\n", err)
			return
		}
		if len(resp.Entries) == 0 {
			fmt.Println("No audit entries found.")
			return
		}
		if verbose {
			for _, e := range resp.Entries {
				printAuditEntry(e)
				fmt.Println()
			}
			return
		}
		out := tabwriter.NewWriter(os.Stdout, 10, 5, 1, ' ', 0)
		_, _ = fmt.Fprintf(out, "CHALLENGE\tTITLE\tSTATUS\tVERIFIED\tCREATED\tSIGNED\n")
		_, _ = fmt.Fprintf(out, "---------\t-----\t------\t--------\t-------\t------\n")
		for _, e := range resp.Entries {
			created := ""
			if e.CreatedAt != nil {
				created = e.CreatedAt.AsTime().Format("2006-01-02 15:04:05")
			}
			signed := ""
			if e.SignedAt != nil {
				signed = e.SignedAt.AsTime().Format("2006-01-02 15:04:05")
			}
			_, _ = fmt.Fprintf(out, "%s\t%s\t%s\t%v\t%s\t%s\n",
				e.ChallengeId, e.Title, e.Status.String(), e.Verified, created, signed,
			)
		}
		_ = out.Flush()
	},
}

func init() {
	domainAdminCmd.AddCommand(domainAuditUserCmd)
	domainAuditUserCmd.Flags().String("from", "", "filter entries from this time (RFC3339)")
	domainAuditUserCmd.Flags().String("to", "", "filter entries to this time (RFC3339)")
	domainAuditUserCmd.Flags().BoolP("verbose", "v", false, "show detailed entry information")
}

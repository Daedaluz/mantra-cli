package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/daedaluz/mantra-cli/internal"
	admin "github.com/daedaluz/mantra-cli/lib/grpc/service"
	"github.com/spf13/cobra"
)

// listClientsCmd represents the listClients command
var listClientsCmd = &cobra.Command{
	Use:   "listClients",
	Short: "List domain clients with client ID and client secret",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, api := internal.GetDomainAdminAPIClient(cmd)
		adminOnly, _ := cmd.Flags().GetBool("admin-only")

		resp, err := api.ListClients(ctx, &admin.ListClientsRequest{
			AdminOnly: adminOnly,
		})
		if err != nil {
			fmt.Printf("Error listing clients: %v\n", err)
			return
		}

		if len(resp.Clients) == 0 {
			fmt.Println("No clients found.")
			return
		}

		out := tabwriter.NewWriter(os.Stdout, 10, 5, 1, ' ', 0)
		_, _ = fmt.Fprintf(out, "CLIENT ID\tCLIENT SECRET\tNAME\tADMIN\tCREATED\n")
		_, _ = fmt.Fprintf(out, "---------\t-------------\t----\t-----\t-------\n")
		for _, c := range resp.Clients {
			created := ""
			if c.CreatedAt != nil {
				created = c.CreatedAt.AsTime().Format("2006-01-02 15:04:05")
			}
			name := c.Name
			if name == "" {
				name = "-"
			}
			_, _ = fmt.Fprintf(out, "%s\t%s\t%s\t%v\t%s\n",
				c.ClientId,
				c.ClientSecret,
				name,
				c.Admin,
				created,
			)
		}
		_ = out.Flush()
	},
}

func init() {
	domainAdminCmd.AddCommand(listClientsCmd)
	listClientsCmd.Flags().Bool("admin-only", false, "List only admin clients")
}

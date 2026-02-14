/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"text/tabwriter"

	"github.com/daedaluz/mantra-cli/cmd/mantra-cli/internal"
	admin "github.com/daedaluz/mantra-cli/lib/grpc/service"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/status"
)

// listDomainsCmd represents the listDomains command
var listDomainsCmd = &cobra.Command{
	Use:   "listDomains [domain...]",
	Short: "List all domains, optionally filtered by domain names",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, api := internal.GetAdminAPIClient(cmd)
		resp, err := api.ListDomains(ctx, &admin.ListDomainsRequest{
			Domains: args,
		})
		if err != nil {
			// Print detailed error message
			if status, ok := status.FromError(err); ok {
				fmt.Printf("Error listing domain: %s (code: %s)\n", status.Message(), status.Code())
			} else {
				fmt.Printf("Error listing domain: %v\n", err)
			}
			return
		}
		out := tabwriter.NewWriter(os.Stdout, 10, 5, 1, ' ', 0)
		_, _ = fmt.Fprintf(out, "DOMAIN\tNAME\tDESCRIPTION\n")
		_, _ = fmt.Fprintf(out, "------\t----\t-----------\n")
		for _, domain := range resp.Domains {
			fmt.Fprintf(out, "%s\t%s\t%s\n", domain.Domain, domain.Name, domain.Description)
		}
		_ = out.Flush()
	},
}

func init() {
	adminCmd.AddCommand(listDomainsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listDomainsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listDomainsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

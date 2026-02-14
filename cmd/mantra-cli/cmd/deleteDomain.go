/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/daedaluz/mantra-cli/cmd/mantra-cli/internal"
	admin "github.com/daedaluz/mantra-cli/lib/grpc/service"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/status"
)

// deleteDomainCmd represents the deleteDomain command
var deleteDomainCmd = &cobra.Command{
	Use:   "deleteDomain",
	Short: "Delete a domain",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx, api := internal.GetAdminAPIClient(cmd)
		_, err := api.DeleteDomain(ctx, &admin.DeleteDomainRequest{
			Domain: args[0],
		})
		if err != nil {
			// Print detailed error message
			if status, ok := status.FromError(err); ok {
				fmt.Printf("Error deleting domain: %s (code: %s)\n", status.Message(), status.Code())
			} else {
				fmt.Printf("Error deleting domain: %v\n", err)
			}
			return
		}
		fmt.Print("Domain deleted successfully\n")
	},
}

func init() {
	adminCmd.AddCommand(deleteDomainCmd)
}

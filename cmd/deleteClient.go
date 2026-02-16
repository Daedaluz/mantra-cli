package cmd

import (
	"fmt"

	"github.com/daedaluz/mantra-cli/internal"
	admin "github.com/daedaluz/mantra-cli/lib/grpc/service"
	"github.com/spf13/cobra"
)

var deleteClientCmd = &cobra.Command{
	Use:   "deleteClient [client-id]",
	Short: "Delete a client",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx, api := internal.GetDomainAdminAPIClient(cmd)
		_, err := api.DeleteClient(ctx, &admin.DeleteClientRequest{
			ClientId: args[0],
		})
		if err != nil {
			fmt.Printf("Error deleting client: %v\n", err)
			return
		}
		fmt.Println("Client deleted successfully")
	},
}

func init() {
	domainAdminCmd.AddCommand(deleteClientCmd)
}

package cmd

import (
	"fmt"

	"github.com/daedaluz/mantra-cli/internal"
	admin "github.com/daedaluz/mantra-cli/lib/grpc/service"
	"github.com/spf13/cobra"
)

var deleteUserCmd = &cobra.Command{
	Use:   "deleteUser [user-id]",
	Short: "Delete a user",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx, api := internal.GetDomainAdminAPIClient(cmd)
		_, err := api.DeleteUser(ctx, &admin.DeleteUserRequest{
			UserId: args[0],
		})
		if err != nil {
			fmt.Printf("Error deleting user: %v\n", err)
			return
		}
		fmt.Println("User deleted successfully")
	},
}

func init() {
	domainAdminCmd.AddCommand(deleteUserCmd)
}

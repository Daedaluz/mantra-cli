package cmd

import (
	"fmt"

	"github.com/daedaluz/mantra-cli/internal"
	admin "github.com/daedaluz/mantra-cli/lib/grpc/service"
	"github.com/spf13/cobra"
)

var createClientCmd = &cobra.Command{
	Use:   "createClient",
	Short: "Create a new client",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, api := internal.GetDomainAdminAPIClient(cmd)
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		icon, _ := cmd.Flags().GetString("icon")
		isAdmin, _ := cmd.Flags().GetBool("admin")

		resp, err := api.CreateClient(ctx, &admin.CreateClientRequest{
			Name:        name,
			Description: description,
			Icon:        icon,
			Admin:       isAdmin,
		})
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			return
		}

		fmt.Printf("Client created successfully\n")
		fmt.Printf("  Client ID:     %s\n", resp.ClientId)
		fmt.Printf("  Client Secret: %s\n", resp.ClientSecret)
		fmt.Printf("  Name:          %s\n", resp.Name)
		fmt.Printf("  Description:   %s\n", resp.Description)
		fmt.Printf("  Icon:          %s\n", resp.Icon)
		fmt.Printf("  Admin:         %v\n", resp.Admin)
	},
}

func init() {
	domainAdminCmd.AddCommand(createClientCmd)
	createClientCmd.Flags().StringP("name", "n", "", "client name")
	createClientCmd.Flags().StringP("description", "", "", "client description")
	createClientCmd.Flags().StringP("icon", "", "", "client icon")
	createClientCmd.Flags().Bool("admin", false, "grant admin privileges")
}

package cmd

import (
	"fmt"

	"github.com/daedaluz/mantra-cli/internal"
	admin "github.com/daedaluz/mantra-cli/lib/grpc/service"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/status"
)

// createDomainCmd represents the createDomain command
var createDomainCmd = &cobra.Command{
	Use:   "createDomain",
	Short: "Create a new domain",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		ctx, api := internal.GetAdminAPIClient(cmd)
		// if number of args is two, set description to name
		if len(args) == 2 {
			args = append(args, args[1])
		}
		regPath := "/reg"
		signPath := "/auth"
		if activeCtx != nil {
			if activeCtx.RegisterPath != "" {
				regPath = activeCtx.RegisterPath
			}
			if activeCtx.AuthPath != "" {
				signPath = activeCtx.AuthPath
			}
		}

		resp, err := api.CreateDomain(ctx, &admin.CreateDomainRequest{
			Domain:           args[0],
			Name:             args[1],
			Description:      args[2],
			RegistrationPath: regPath,
			SignPath:         signPath,
		})
		if err != nil {
			// Print detailed error message
			if status, ok := status.FromError(err); ok {
				fmt.Printf("Error creating domain: %s (code: %s)\n", status.Message(), status.Code())
			} else {
				fmt.Printf("Error creating domain: %v\n", err)
			}
			return
		}
		fmt.Printf("Domain created successfully:\n%s\n", resp)
	},
}

func init() {
	adminCmd.AddCommand(createDomainCmd)
}

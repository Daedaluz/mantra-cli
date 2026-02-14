package cmd

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/daedaluz/mantra-cli/internal"
	admin "github.com/daedaluz/mantra-cli/lib/grpc/service"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/emptypb"
)

// listUsersCmd represents the listUsers command
var listUsersCmd = &cobra.Command{
	Use:   "listUsers [user-id]",
	Short: "List users, or show details for a specific user",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, api := internal.GetDomainAdminAPIClient(cmd)

		if len(args) > 0 {
			// Show details for a specific user
			resp, err := api.GetUser(ctx, &admin.GetUserRequest{
				UserId: args[0],
			})
			if err != nil {
				fmt.Printf("Error getting user: %v\n", err)
				return
			}
			fmt.Printf("Domain:     %s\n", resp.Domain)
			fmt.Printf("User ID:    %s\n", resp.UserId)
			if resp.CreatedAt != nil {
				fmt.Printf("Created At: %s\n", resp.CreatedAt.AsTime().String())
			}
			if len(resp.Keys) > 0 {
				fmt.Printf("\nKeys:\n")
				out := tabwriter.NewWriter(os.Stdout, 10, 5, 1, ' ', 0)
				_, _ = fmt.Fprintf(out, "  NAME\tDISABLED\tATTACHMENT\tTRANSPORTS\tCREATED\tLAST USED\n")
				_, _ = fmt.Fprintf(out, "  ----\t--------\t----------\t----------\t-------\t---------\n")
				for _, key := range resp.Keys {
					created := ""
					if key.CreatedAt != nil {
						created = key.CreatedAt.AsTime().Format("2006-01-02 15:04:05")
					}
					lastUsed := ""
					if key.LastUsed != nil {
						lastUsed = key.LastUsed.AsTime().Format("2006-01-02 15:04:05")
					}
					name := key.Name
					if name == "" {
						name = base64.RawURLEncoding.EncodeToString(key.KeyId)
					}
					_, _ = fmt.Fprintf(out, "  %s\t%v\t%s\t%s\t%s\t%s\n",
						name,
						key.Disabled,
						key.Attachment,
						strings.Join(key.Transports, ","),
						created,
						lastUsed,
					)
				}
				_ = out.Flush()
			}
			return
		}

		// List all users
		resp, err := api.ListUsers(ctx, &emptypb.Empty{})
		if err != nil {
			fmt.Printf("Error listing users: %v\n", err)
			return
		}
		for _, user := range resp.Users {
			fmt.Println(user)
		}
	},
}

func init() {
	domainAdminCmd.AddCommand(listUsersCmd)
}

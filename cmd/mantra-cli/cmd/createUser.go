package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/daedaluz/mantra-cli/cmd/mantra-cli/internal"
	"github.com/daedaluz/mantra-cli/lib/grpc/common"
	admin "github.com/daedaluz/mantra-cli/lib/grpc/service"
	"github.com/mdp/qrterminal"
	"github.com/spf13/cobra"
)

func mustParseInt64(s string) int64 {
	var i int64
	if s == "" {
		return 300
	}
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("failed to parse int64 from %q: %v", s,
			err))
	}
	return i
}

// clearScreen clears the terminal screen
func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func getValueFromFlag(cmd *cobra.Command, flag string) string {
	value, err := cmd.Flags().GetString(flag)
	if err != nil {
		return ""
	}
	return value
}

// createUserCmd represents the createUser command
var createUserCmd = &cobra.Command{
	Use:   "createUser",
	Short: "Create a new user",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, api := internal.GetDomainAdminAPIClient(cmd)
		geohash := getValueFromFlag(cmd, "location")
		locationRequired, _ := cmd.Flags().GetBool("location-required")

		var location *common.Location
		if geohash != "" {
			location = &common.Location{
				Geohash:  geohash,
				Required: locationRequired,
			}
		}

		challenge, err := api.CreateUser(ctx, &admin.CreateUserRequest{
			UserId:    getValueFromFlag(cmd, "user-id"),
			Name:      getValueFromFlag(cmd, "name"),
			Timeout:   mustParseInt64(getValueFromFlag(cmd, "timeout")),
			ReturnUrl: "",
			Location:  location,
		})
		if err != nil {
			fmt.Printf("Error creating user: %v\n", err)
			return
		}
		start := time.Now()
		for {
			animated, err := challenge.CalculateAnimatedToken(start)
			if err != nil {
				fmt.Printf("Error calculating token: %v\n", err)
				return
			}
			static, err := challenge.CalculateStaticToken()
			if err != nil {
				fmt.Printf("Error calculating token: %v\n", err)
				return
			}
			animatedLink := fmt.Sprintf("https://%s/reg?t=%s", cmd.Flag("domain").Value.String(), animated)
			staticLink := fmt.Sprintf("https://%s/reg?t=%s", cmd.Flag("domain").Value.String(), static)
			status, err := api.CollectChallenge(ctx, &admin.CollectRequest{
				ChallengeId:     challenge.ChallengeId,
				ChallengeSecret: challenge.ChallengeSecret,
			})
			if err != nil {
				fmt.Printf("Error collecting challenge: %v\n", err)
				return
			}
			switch status.Status {
			case common.Status_Cancelled:
				fmt.Println("Challenge cancelled")
				return
			case common.Status_Expired:
				fmt.Println("Challenge expired")
				return
			case common.Status_Completed:
				clearScreen()
				fmt.Println("Challenge completed successfully")
				x, err := json.MarshalIndent(status, "", "  ")
				if err != nil {
					fmt.Printf("Error marshalling status: %v\n", err)
					return
				}
				fmt.Println(string(x))
				return
			}
			clearScreen()
			fmt.Printf("To complete the registration, scan or follow the link below:\n\n")
			qrterminal.GenerateHalfBlock(animatedLink, qrterminal.L, os.Stdout)
			fmt.Printf("Status: %s\n", status.Status.String())
			fmt.Println("Animated link:", animatedLink)
			fmt.Println("Static link:  ", staticLink)
			time.Sleep(1 * time.Second)
		}
	},
}

func init() {
	domainAdminCmd.AddCommand(createUserCmd)
	createUserCmd.Flags().StringP("user-id", "u", "", "Optional user ID (if not provided, a random one will be generated)")
	createUserCmd.Flags().StringP("name", "n", "", "Name of the key to register")
	createUserCmd.Flags().Int64P("timeout", "o", 300, "Timeout in seconds for the registration challenge")
	createUserCmd.Flags().StringP("location", "l", "", "Geohash of the registration location (triggers location request from authenticator)")
	createUserCmd.Flags().Bool("location-required", false, "Require location and verify it before accepting the registration")
}

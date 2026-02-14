package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/daedaluz/mantra-cli/internal"
	"github.com/daedaluz/mantra-cli/lib/grpc/common"
	admin "github.com/daedaluz/mantra-cli/lib/grpc/service"
	"github.com/mdp/qrterminal"
	"github.com/spf13/cobra"
)

// addKeyCmd represents the addKey command
var addKeyCmd = &cobra.Command{
	Use:   "addKey",
	Short: "Register an additional key for an existing user",
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

		challenge, err := api.AddKey(ctx, &admin.AddKeyRequest{
			UserId:    getValueFromFlag(cmd, "user-id"),
			Name:      getValueFromFlag(cmd, "name"),
			Timeout:   mustParseInt64(getValueFromFlag(cmd, "timeout")),
			ReturnUrl: "",
			Location:  location,
		})
		if err != nil {
			fmt.Printf("Error adding key: %v\n", err)
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
			regPath := "/reg"
			if activeCtx != nil && activeCtx.RegisterPath != "" {
				regPath = activeCtx.RegisterPath
			}
			animatedLink := fmt.Sprintf("https://%s%s?t=%s", cmd.Flag("domain").Value.String(), regPath, animated)
			staticLink := fmt.Sprintf("https://%s%s?t=%s", cmd.Flag("domain").Value.String(), regPath, static)
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
			fmt.Printf("To register the new key, scan or follow the link below:\n\n")
			qrterminal.GenerateHalfBlock(animatedLink, qrterminal.L, os.Stdout)
			fmt.Printf("Status: %s\n", status.Status.String())
			fmt.Println("Animated link:", animatedLink)
			fmt.Println("Static link:  ", staticLink)
			time.Sleep(1 * time.Second)
		}
	},
}

func init() {
	domainAdminCmd.AddCommand(addKeyCmd)
	addKeyCmd.Flags().StringP("user-id", "u", "", "User ID to add a key for")
	addKeyCmd.Flags().StringP("name", "n", "", "Name of the key to register")
	addKeyCmd.Flags().Int64P("timeout", "o", 300, "Timeout in seconds for the registration challenge")
	addKeyCmd.Flags().StringP("location", "l", "", "Geohash of the registration location (triggers location request from authenticator)")
	addKeyCmd.Flags().Bool("location-required", false, "Require location and verify it before accepting the registration")
	_ = addKeyCmd.MarkFlagRequired("user-id")
}

package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/daedaluz/mantra-cli/internal"
	"github.com/daedaluz/mantra-cli/lib/grpc/client"
	"github.com/daedaluz/mantra-cli/lib/grpc/common"
	"github.com/mdp/qrterminal"
	"github.com/spf13/cobra"
)

// authenticateCmd represents the authenticate command
var authenticateCmd = &cobra.Command{
	Use:   "authenticate",
	Short: "Create an authentication (signature) challenge",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, api := internal.GetClientApiClient(cmd)
		geohash := getValueFromFlag(cmd, "location")
		locationRequired, _ := cmd.Flags().GetBool("location-required")

		var location *common.Location
		if geohash != "" {
			location = &common.Location{
				Geohash:  geohash,
				Required: locationRequired,
			}
		}

		challenge, err := api.Sign(ctx, &client.SignRequest{
			User:             getValueFromFlag(cmd, "user-id"),
			Timeout:          mustParseInt64(getValueFromFlag(cmd, "timeout")),
			UserVerification: getValueFromFlag(cmd, "user-verification"),
			Mediation:        getValueFromFlag(cmd, "mediation"),
			Message:          getValueFromFlag(cmd, "message"),
			Data:             []byte{},
			ReturnUrl:        "",
			Location:         location,
		})
		if err != nil {
			fmt.Printf("Error creating challenge: %v\n", err)
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
			animatedLink := fmt.Sprintf("https://%s/auth?t=%s", cmd.Flag("domain").Value.String(), animated)
			staticLink := fmt.Sprintf("https://%s/auth?t=%s", cmd.Flag("domain").Value.String(), static)
			status, err := api.Collect(ctx, &client.CollectRequest{
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
			case common.Status_Rejected:
				clearScreen()
				fmt.Printf("Status: %s\n", status.Status.String())
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
			fmt.Printf("To authenticate, scan or follow the link below:\n\n")
			qrterminal.GenerateHalfBlock(animatedLink, qrterminal.L, os.Stdout)
			fmt.Printf("Status: %s\n", status.Status.String())
			fmt.Println("Animated link:", animatedLink)
			fmt.Println("Static link:  ", staticLink)
			time.Sleep(1 * time.Second)
		}
	},
}

func init() {
	domainAdminCmd.AddCommand(authenticateCmd)
	authenticateCmd.Flags().StringP("user-id", "u", "", "User ID to authenticate (optional, leave empty for any user)")
	authenticateCmd.Flags().StringP("user-verification", "v", "preferred", "User verification requirement (required, discouraged, preferred)")
	authenticateCmd.Flags().StringP("mediation", "m", "optional", "Mediation requirement (silent, optional, required)")
	authenticateCmd.Flags().StringP("message", "M", "Please authenticate", "Message to display on the authenticator")
	authenticateCmd.Flags().Int64P("timeout", "o", 300, "Timeout in seconds")
	authenticateCmd.Flags().StringP("location", "l", "", "Geohash of the challenge location (triggers location request from authenticator)")
	authenticateCmd.Flags().Bool("location-required", false, "Require location and verify it before accepting the signature")
}

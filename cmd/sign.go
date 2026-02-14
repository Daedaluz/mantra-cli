/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
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

// signCmd represents the sign command
var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "Create a signature challenge",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, api := internal.GetClientApiClient(cmd)
		challenge, err := api.Sign(ctx, &client.SignRequest{
			User:             getValueFromFlag(cmd, "user-id"),
			Timeout:          mustParseInt64(getValueFromFlag(cmd, "timeout")),
			UserVerification: getValueFromFlag(cmd, "user-verification"),
			Mediation:        getValueFromFlag(cmd, "mediation"),
			Message:          getValueFromFlag(cmd, "message"),
			Data:             []byte{},
			ReturnUrl:        "",
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
			authPath := "/auth"
			if activeCtx != nil && activeCtx.AuthPath != "" {
				authPath = activeCtx.AuthPath
			}
			animatedLink := fmt.Sprintf("https://%s%s?t=%s", cmd.Flag("domain").Value.String(), authPath, animated)
			staticLink := fmt.Sprintf("https://%s%s?t=%s", cmd.Flag("domain").Value.String(), authPath, static)
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
			fmt.Printf("To sign, scan or follow the link below:\n\n")
			qrterminal.GenerateHalfBlock(animatedLink, qrterminal.L, os.Stdout)
			fmt.Printf("Status: %s\n", status.Status.String())
			fmt.Println("Animated link:", animatedLink)
			fmt.Println("Static link:  ", staticLink)
			time.Sleep(1 * time.Second)
		}
	},
}

func init() {
	clientCmd.AddCommand(signCmd)
	signCmd.Flags().StringP("user-id", "u", "", "User ID to authenticate as (Optional)")
	signCmd.Flags().StringP("user-verification", "v", "preferred", "User verification requirement (required, discouraged, preferred, required)")
	signCmd.Flags().StringP("mediation", "m", "optional", "Mediation requirement (silent, optional, required)")
	signCmd.Flags().StringP("message", "M", "Please sign this challenge", "Message to display on the authenticator")
	signCmd.Flags().Int64P("timeout", "o", 300, "Timeout in seconds")
}

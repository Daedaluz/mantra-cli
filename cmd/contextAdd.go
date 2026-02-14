package cmd

import (
	"fmt"

	"github.com/daedaluz/mantra-cli/internal/config"
	"github.com/spf13/cobra"
)

var contextAddCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Add a client context",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		apiRef, _ := cmd.Flags().GetString("api")
		domain, _ := cmd.Flags().GetString("domain")
		clientID, _ := cmd.Flags().GetString("client-id")
		clientSecret, _ := cmd.Flags().GetString("client-secret")
		registerPath, _ := cmd.Flags().GetString("register-path")
		authPath, _ := cmd.Flags().GetString("auth-path")

		cfg := config.Load()

		if _, ok := cfg.APIs[apiRef]; !ok {
			fmt.Printf("Error: API %q does not exist. Add it first with 'mantra-cli api add'.\n", apiRef)
			return
		}

		cfg.Contexts[name] = &config.Context{
			API:          apiRef,
			Domain:       domain,
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RegisterPath: registerPath,
			AuthPath:     authPath,
		}

		// Auto-set as current if it's the first context
		if len(cfg.Contexts) == 1 {
			cfg.CurrentContext = name
		}

		if err := config.Save(cfg); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			return
		}
		if cfg.CurrentContext == name {
			fmt.Printf("Context %q added and set as current.\n", name)
		} else {
			fmt.Printf("Context %q added.\n", name)
		}
	},
}

func init() {
	contextCmd.AddCommand(contextAddCmd)
	contextAddCmd.Flags().String("api", "", "API to reference (required)")
	contextAddCmd.MarkFlagRequired("api")
	contextAddCmd.Flags().String("domain", "", "Domain name")
	contextAddCmd.Flags().String("client-id", "", "Client ID")
	contextAddCmd.Flags().String("client-secret", "", "Client secret")
	contextAddCmd.Flags().String("register-path", "/reg", "Registration URL path")
	contextAddCmd.Flags().String("auth-path", "/auth", "Authentication URL path")
}

package cmd

import (
	"fmt"

	"github.com/daedaluz/mantra-cli/internal/config"
	"github.com/spf13/cobra"
)

var contextUseCmd = &cobra.Command{
	Use:   "use <name>",
	Short: "Set the active context",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		cfg := config.Load()

		if _, ok := cfg.Contexts[name]; !ok {
			fmt.Printf("Error: context %q does not exist.\n", name)
			return
		}

		cfg.CurrentContext = name
		if err := config.Save(cfg); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			return
		}
		fmt.Printf("Switched to context %q.\n", name)
	},
}

func init() {
	contextCmd.AddCommand(contextUseCmd)
}

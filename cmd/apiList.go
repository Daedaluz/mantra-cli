package cmd

import (
	"fmt"
	"sort"
	"text/tabwriter"
	"os"

	"github.com/daedaluz/mantra-cli/internal/config"
	"github.com/spf13/cobra"
)

var apiListCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured APIs",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		if len(cfg.APIs) == 0 {
			fmt.Println("No APIs configured.")
			return
		}
		names := make([]string, 0, len(cfg.APIs))
		for name := range cfg.APIs {
			names = append(names, name)
		}
		sort.Strings(names)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tSERVER\tAPI KEY\tPLAINTEXT\tSKIP VERIFY")
		for _, name := range names {
			api := cfg.APIs[name]
			keyStatus := "<not set>"
			if api.APIKey != "" {
				keyStatus = "<set>"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%v\t%v\n", name, api.Server, keyStatus, api.Plaintext, api.SkipVerify)
		}
		w.Flush()
	},
}

func init() {
	apiCmd.AddCommand(apiListCmd)
}

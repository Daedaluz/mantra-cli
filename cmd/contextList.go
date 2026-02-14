package cmd

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/daedaluz/mantra-cli/internal/config"
	"github.com/spf13/cobra"
)

var contextListCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured contexts",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		if len(cfg.Contexts) == 0 {
			fmt.Println("No contexts configured.")
			return
		}
		names := make([]string, 0, len(cfg.Contexts))
		for name := range cfg.Contexts {
			names = append(names, name)
		}
		sort.Strings(names)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "  \tNAME\tAPI\tDOMAIN")
		for _, name := range names {
			ctx := cfg.Contexts[name]
			marker := " "
			if name == cfg.CurrentContext {
				marker = "*"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", marker, name, ctx.API, ctx.Domain)
		}
		w.Flush()
	},
}

func init() {
	contextCmd.AddCommand(contextListCmd)
}

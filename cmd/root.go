package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mantra-cli",
	Short: "command line interface for Mantra",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func envStr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envBool(key string) bool {
	v := strings.ToLower(os.Getenv(key))
	return v == "1" || v == "true" || v == "yes"
}

func init() {
	rootCmd.PersistentFlags().StringP("server", "s", envStr("SERVER", "mantra-api.inits.se:443"), "Hostname of the server to connect to [$SERVER]")
	rootCmd.PersistentFlags().Bool("plaintext", envBool("PLAINTEXT"), "Use plaintext gRPC (no TLS) [$PLAINTEXT]")
	rootCmd.PersistentFlags().Bool("skip-verify", envBool("SKIP_VERIFY"), "Skip server certificate verification (TLS only) [$SKIP_VERIFY]")
}

package cmd

import (
	"os"
	"strings"

	"github.com/daedaluz/mantra-cli/internal/config"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mantra-cli",
	Short: "command line interface for Mantra",
}

var loadedCfg = config.Load()
var activeCtx, activeAPI = loadedCfg.ResolveContext()

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

func ctxStr(ctxVal, fallback string) string {
	if ctxVal != "" {
		return ctxVal
	}
	return fallback
}

func init() {
	serverDefault := "mantra-api.inits.se:443"
	if activeAPI != nil {
		serverDefault = ctxStr(activeAPI.Server, serverDefault)
	}

	rootCmd.PersistentFlags().StringP("server", "s", envStr("SERVER", serverDefault), "Hostname of the server to connect to [$SERVER]")
	plaintextDefault := envBool("PLAINTEXT")
	skipVerifyDefault := envBool("SKIP_VERIFY")
	if activeAPI != nil {
		if !plaintextDefault {
			plaintextDefault = activeAPI.Plaintext
		}
		if !skipVerifyDefault {
			skipVerifyDefault = activeAPI.SkipVerify
		}
	}

	rootCmd.PersistentFlags().Bool("plaintext", plaintextDefault, "Use plaintext gRPC (no TLS) [$PLAINTEXT]")
	rootCmd.PersistentFlags().Bool("skip-verify", skipVerifyDefault, "Skip server certificate verification (TLS only) [$SKIP_VERIFY]")
}

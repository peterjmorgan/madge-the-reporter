package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "Madge",
	Short: "Report Phylum Issues to ticketing systems",
	Long:  `Report Phylum Issues to ticketing systems`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.PhylumSyringeGitlab.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Debug logging")
	rootCmd.PersistentFlags().BoolP("dry-run", "D", false, "Dry Run")
	// TODO: consider removing these. Mostly for testing for a specific use case. Perhaps moving to the environment is better
	//rootCmd.PersistentFlags().BoolP("mine-only", "m", false, "(Gitlab) Only projects owned by the user")
	//rootCmd.PersistentFlags().Int32P("ratelimit", "r", 100, "Rate Limit (X/reqs/sec) ")
	//rootCmd.PersistentFlags().StringP("proxyUrl", "p", "", "proxy (https://url:port)")
}

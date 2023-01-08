package cmd

import (
	"fmt"
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

	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.PhylumSyringeGitlab.yaml)")
	rootCmd.PersistentFlags().StringP("configFile", "c", "", "config file (default is $PWD/madge_config.yaml")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Debug logging")
	rootCmd.PersistentFlags().BoolP("dry-run", "D", false, "Dry Run")
	rootCmd.PersistentFlags().BoolP("version", "v", false, "Print version")

	versionFlag, _ := rootCmd.Flags().GetBool("version")
	if versionFlag {
		fmt.Printf("Version is %s", Version)
	}

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}

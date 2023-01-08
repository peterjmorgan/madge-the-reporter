package cmd

import (
	"context"
	"github.com/creativeprojects/go-selfupdate"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"runtime"
)

func init() {
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update",
	Run: func(cmd *cobra.Command, args []string) {
		//source, err := selfupdate.NewGitHubSource(selfupdate.GitHubConfig{})
		//if err != nil {
		//	log.Errorf("error occurred getting source: %w", err)
		//	return
		//}

		//TODO: get current version

		latest, found, err := selfupdate.DetectLatest(context.Background(), selfupdate.ParseSlug("peterjmorgan/madge-the-reporter"))
		if err != nil {
			log.Errorf("error occurred while detecting version: %w", err)
			return
		}
		if !found {
			log.Errorf("latest version for %s/%s could not be found from github repository", runtime.GOOS, runtime.GOARCH)
			return
		}

		if latest.LessOrEqual(version) {
			log.Printf("Current version (%s) is the latest", version)
			return
		}

		exe, err := os.Executable()
		if err != nil {
			log.Errorf("could not locate executable path")
			return
		}
		if err := selfupdate.UpdateTo(context.Background(), latest.AssetURL, latest.AssetName, exe); err != nil {
			log.Errorf("error occurred while updating binary: %w", err)
			return
		}
		log.Printf("Successfully updated to version %s", latest.Version())
	},
}

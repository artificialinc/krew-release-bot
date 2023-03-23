package main

import (
	"os"

	"github.com/rajatjindal/krew-release-bot/pkg/source/actions"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(actionCmd)
}

// actionCmd is the github action command
var actionCmd = &cobra.Command{
	Use:   "action",
	Short: "github action for updating plugin manifests in krew-index repo",
	Run: func(cmd *cobra.Command, args []string) {
		if os.Getenv("LOG_LEVEL") == "debug" {
			logrus.SetLevel(logrus.DebugLevel)
		}
		err := actions.RunAction()
		if err != nil {
			logrus.Fatal(err)
		}
	},
}

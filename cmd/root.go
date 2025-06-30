package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "agit",
	Short: "AGit helper for Gitea",
	Long:  `A command-line tool to create and manage pull requests using AGit workflow with Gitea`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(prCmd)
}
package cmd

import (
	"agit/internal/git"
	"agit/internal/gitea"
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List open pull requests",
	Long:    `List all open pull requests for the current repository`,
	RunE:    runList,
}

func runList(cmd *cobra.Command, args []string) error {
	repo, err := git.OpenRepository(".")
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}

	remoteURL, err := repo.GetRemoteURL("origin")
	if err != nil {
		return fmt.Errorf("failed to get remote URL: %w", err)
	}

	owner, repoName, baseURL, err := gitea.ParseRemoteURL(remoteURL)
	if err != nil {
		return fmt.Errorf("failed to parse remote URL: %w", err)
	}

	client, err := gitea.NewClient(baseURL)
	if err != nil {
		return fmt.Errorf("failed to create Gitea client: %w", err)
	}

	prs, err := client.ListPullRequests(owner, repoName)
	if err != nil {
		return fmt.Errorf("failed to list pull requests: %w", err)
	}

	if len(prs) == 0 {
		fmt.Println("No open pull requests found")
		return nil
	}

	fmt.Printf("Open Pull Requests for %s/%s:\n\n", owner, repoName)
	for _, pr := range prs {
		fmt.Printf("#%-4d %-50s %s\n", pr.Index, pr.Title, pr.Poster.UserName)
	}

	return nil
}
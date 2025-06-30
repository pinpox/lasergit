package cmd

import (
	"agit/internal/git"
	"agit/internal/gitea"
	"agit/internal/tui"
	"fmt"

	"github.com/spf13/cobra"
)

var prCmd = &cobra.Command{
	Use:     "pr",
	Aliases: []string{"p"},
	Short:   "Browse and checkout pull requests",
	Long:    `Browse open pull requests interactively. Use enter to checkout a PR or 'v' to view details.`,
	RunE:    runPR,
}

func runPR(cmd *cobra.Command, args []string) error {
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

	result, err := tui.ShowPRList(prs, owner, repoName)
	if err != nil {
		return fmt.Errorf("failed to show PR list: %w", err)
	}

	switch result.Action {
	case "checkout":
		if result.SelectedPR != nil {
			pr := result.SelectedPR
			branchName := fmt.Sprintf("agit-%d", pr.Index)
			
			fmt.Printf("🔄 Fetching PR #%d...\n", pr.Index)
			err = repo.FetchPullRequest("origin", int(pr.Index), branchName)
			if err != nil {
				return fmt.Errorf("failed to fetch PR: %w", err)
			}

			fmt.Printf("🔀 Checking out branch '%s'...\n", branchName)
			err = repo.CheckoutBranch(branchName)
			if err != nil {
				return fmt.Errorf("failed to checkout branch: %w", err)
			}

			fmt.Printf("✅ Successfully checked out PR #%d: %s\n", pr.Index, pr.Title)
			fmt.Printf("📍 You are now on branch '%s'\n", branchName)
		}
	case "view":
		if result.SelectedPR != nil {
			pr := result.SelectedPR
			fmt.Printf("\n✨ PR #%d: %s\n", pr.Index, pr.Title)
			if pr.Body != "" {
				fmt.Printf("\nDescription:\n%s\n", pr.Body)
			}
			fmt.Printf("\nAuthor: %s\n", pr.Poster.UserName)
			if pr.HTMLURL != "" {
				fmt.Printf("URL: %s\n", pr.HTMLURL)
			}
			if pr.Updated != nil {
				fmt.Printf("Updated: %s\n", pr.Updated.Format("2006-01-02 15:04"))
			}
		}
	case "refresh":
		return runPR(cmd, args)
	}

	return nil
}
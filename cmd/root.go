package cmd

import (
	"lasergit/internal/git"
	"lasergit/internal/gitea"
	"lasergit/internal/tui"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	rootRepoPath string
)

var rootCmd = &cobra.Command{
	Use:   "lasergit",
	Short: "AGit helper for Gitea",
	Long: `A unified interface to manage pull requests using AGit workflow with Gitea.

Browse, create, and checkout pull requests interactively:
• Navigate with ↑/↓ arrows  
• Press Enter to checkout a PR
• Press 'c' to create a new PR
• Press 'v' to view PR details
• Press 'r' to refresh the list
• Press 'q' or Esc to quit`,
	RunE: runRoot,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVar(&rootRepoPath, "repo", ".", "Path to git repository")
}

func runRoot(cmd *cobra.Command, args []string) error {
	return runPRLogic(rootRepoPath)
}

func runPRLogic(repoPath string) error {
	repo, err := git.OpenRepository(repoPath)
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

	currentBranch, err := repo.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	result, err := tui.ShowPRList(prs, owner, repoName, currentBranch)
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
	case "create":
		return handleCreatePR(repo)
	case "refresh":
		return runPRLogic(repoPath)
	}

	return nil
}

func handleCreatePR(repo *git.Repository) error {
	currentBranch, err := repo.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	targetBranch := "main" // Default target branch
	topicName := currentBranch

	result, err := tui.ShowCreatePRDialog(topicName, targetBranch)
	if err != nil {
		return fmt.Errorf("failed to get PR details: %w", err)
	}

	pushOptions := []string{
		fmt.Sprintf("topic=%s", result.Topic),
		fmt.Sprintf("title=%s", result.Title),
		fmt.Sprintf("description=%s", result.Description),
	}

	err = repo.PushAGit(result.Target, pushOptions)
	if err != nil {
		return fmt.Errorf("failed to push: %w", err)
	}

	fmt.Printf("✅ Successfully created PR for topic '%s' targeting '%s'\n", result.Topic, result.Target)
	return nil
}
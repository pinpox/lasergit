package cmd

import (
	"agit/internal/git"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var checkoutCmd = &cobra.Command{
	Use:   "checkout <pr-number>",
	Short: "Fetch and checkout a specific pull request",
	Long:  `Fetch a pull request and checkout it locally as agit-<pr-number>`,
	Args:  cobra.ExactArgs(1),
	RunE:  runCheckout,
}

func runCheckout(cmd *cobra.Command, args []string) error {
	prNumber, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid PR number: %s", args[0])
	}

	repo, err := git.OpenRepository(".")
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}

	branchName := fmt.Sprintf("agit-%d", prNumber)

	err = repo.FetchPullRequest("origin", prNumber, branchName)
	if err != nil {
		return fmt.Errorf("failed to fetch PR: %w", err)
	}

	err = repo.CheckoutBranch(branchName)
	if err != nil {
		return fmt.Errorf("failed to checkout branch: %w", err)
	}

	fmt.Printf("Successfully checked out PR #%d as branch '%s'\n", prNumber, branchName)
	return nil
}
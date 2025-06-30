package cmd

import (
	"agit/internal/git"
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	autoMode    bool
	topicName   string
	forceMode   bool
	prTitle     string
	prDesc      string
	targetBranch string
)

var createCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"c"},
	Short:   "Create a pull request using AGit",
	Long:    `Create a pull request by pushing to refs/for/<branch> with AGit workflow`,
	RunE:    runCreate,
}

func init() {
	createCmd.Flags().BoolVar(&autoMode, "auto", false, "Use latest commit message as PR title/description")
	createCmd.Flags().StringVar(&topicName, "topic", "", "Topic branch name")
	createCmd.Flags().BoolVar(&forceMode, "force", false, "Force push to update existing PR")
	createCmd.Flags().StringVar(&prTitle, "title", "", "PR title")
	createCmd.Flags().StringVar(&prDesc, "description", "", "PR description")
	createCmd.Flags().StringVar(&targetBranch, "branch", "main", "Target branch")
}

func runCreate(cmd *cobra.Command, args []string) error {
	repo, err := git.OpenRepository(".")
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}

	if topicName == "" {
		branch, err := repo.GetCurrentBranch()
		if err != nil {
			return fmt.Errorf("failed to get current branch: %w", err)
		}
		topicName = branch
	}

	var title, description string

	if autoMode {
		commit, err := repo.GetLatestCommit()
		if err != nil {
			return fmt.Errorf("failed to get latest commit: %w", err)
		}
		title = commit.Message
		description = commit.Message
	} else if prTitle != "" && prDesc != "" {
		title = prTitle
		description = prDesc
	} else {
		title, description, err = getInteractiveInput()
		if err != nil {
			return fmt.Errorf("failed to get interactive input: %w", err)
		}
	}

	pushOptions := []string{
		fmt.Sprintf("topic=%s", topicName),
		fmt.Sprintf("title=%s", title),
		fmt.Sprintf("description=%s", description),
	}

	if forceMode {
		pushOptions = append(pushOptions, "force-push=true")
	}

	err = repo.PushAGit(targetBranch, pushOptions)
	if err != nil {
		return fmt.Errorf("failed to push: %w", err)
	}

	fmt.Printf("Successfully created PR for topic '%s' targeting '%s'\n", topicName, targetBranch)
	return nil
}

func getInteractiveInput() (string, string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("PR Title: ")
	title, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}
	title = strings.TrimSpace(title)

	fmt.Print("PR Description (optional): ")
	description, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}
	description = strings.TrimSpace(description)

	return title, description, nil
}
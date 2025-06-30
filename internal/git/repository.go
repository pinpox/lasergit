package git

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/go-git/go-git/v5"
)

type Repository struct {
	repo *git.Repository
}

type Commit struct {
	Hash    string
	Message string
	Author  string
}

func OpenRepository(path string) (*Repository, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	return &Repository{repo: repo}, nil
}

func (r *Repository) GetCurrentBranch() (string, error) {
	head, err := r.repo.Head()
	if err != nil {
		return "", err
	}

	return head.Name().Short(), nil
}

func (r *Repository) GetRemoteURL(remoteName string) (string, error) {
	remote, err := r.repo.Remote(remoteName)
	if err != nil {
		return "", err
	}

	urls := remote.Config().URLs
	if len(urls) == 0 {
		return "", fmt.Errorf("no URLs found for remote %s", remoteName)
	}

	return urls[0], nil
}

func (r *Repository) GetLatestCommit() (*Commit, error) {
	head, err := r.repo.Head()
	if err != nil {
		return nil, err
	}

	commit, err := r.repo.CommitObject(head.Hash())
	if err != nil {
		return nil, err
	}

	return &Commit{
		Hash:    commit.Hash.String(),
		Message: strings.TrimSpace(commit.Message),
		Author:  commit.Author.Name,
	}, nil
}

func (r *Repository) PushAGit(targetBranch string, pushOptions []string) error {
	cmd := exec.Command("git", "push", "origin", fmt.Sprintf("HEAD:refs/for/%s", targetBranch))
	
	for _, option := range pushOptions {
		cmd.Args = append(cmd.Args, "-o", option)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git push failed: %s", string(output))
	}

	return nil
}

func (r *Repository) FetchPullRequest(remoteName string, prNumber int, branchName string) error {
	refSpec := fmt.Sprintf("pull/%d/head:%s", prNumber, branchName)
	
	cmd := exec.Command("git", "fetch", remoteName, refSpec)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git fetch failed: %s", string(output))
	}

	return nil
}

func (r *Repository) CheckoutBranch(branchName string) error {
	cmd := exec.Command("git", "checkout", branchName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git checkout failed: %s", string(output))
	}

	return nil
}
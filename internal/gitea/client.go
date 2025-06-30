package gitea

import (
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"

	"code.gitea.io/sdk/gitea"
)

type Client struct {
	client *gitea.Client
}

func NewClient(baseURL string) (*Client, error) {
	token := os.Getenv("GITEA_TOKEN")

	client, err := gitea.NewClient(baseURL, gitea.SetToken(token))
	if err != nil {
		return nil, err
	}

	return &Client{client: client}, nil
}

func (c *Client) ListPullRequests(owner, repo string) ([]*gitea.PullRequest, error) {
	prs, _, err := c.client.ListRepoPullRequests(owner, repo, gitea.ListPullRequestsOptions{
		State: gitea.StateOpen,
	})
	if err != nil {
		return nil, err
	}

	return prs, nil
}

func ParseRemoteURL(remoteURL string) (owner, repo, baseURL string, err error) {
	if strings.Contains(remoteURL, "@") && strings.Contains(remoteURL, ":") && !strings.HasPrefix(remoteURL, "http") {
		return parseSSHURL(remoteURL)
	} else if strings.HasPrefix(remoteURL, "https://") || strings.HasPrefix(remoteURL, "http://") {
		return parseHTTPSURL(remoteURL)
	}

	return "", "", "", fmt.Errorf("unsupported remote URL format: %s", remoteURL)
}

func parseSSHURL(remoteURL string) (owner, repo, baseURL string, err error) {
	re := regexp.MustCompile(`([^@]+)@([^:]+):([^/]+)/(.+)\.git$`)
	matches := re.FindStringSubmatch(remoteURL)

	if len(matches) != 5 {
		return "", "", "", fmt.Errorf("invalid SSH URL format: %s", remoteURL)
	}

	host := matches[2]
	owner = matches[3]
	repo = matches[4]
	baseURL = fmt.Sprintf("https://%s", host)

	return owner, repo, baseURL, nil
}

func parseHTTPSURL(remoteURL string) (owner, repo, baseURL string, err error) {
	u, err := url.Parse(remoteURL)
	if err != nil {
		return "", "", "", err
	}

	path := strings.TrimPrefix(u.Path, "/")
	path = strings.TrimSuffix(path, ".git")

	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		return "", "", "", fmt.Errorf("invalid repository path: %s", path)
	}

	owner = parts[0]
	repo = parts[1]
	baseURL = fmt.Sprintf("%s://%s", u.Scheme, u.Host)

	return owner, repo, baseURL, nil
}

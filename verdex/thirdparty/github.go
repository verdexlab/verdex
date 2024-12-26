package thirdparty

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v30/github"
	"golang.org/x/oauth2"
)

var clientTimeout = 60 * time.Second

func GitHubGetClients() (*github.Client, *http.Client) {
	httpClient := &http.Client{
		Timeout: clientTimeout,
	}

	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		httpClient = oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))
	}

	githubClient := github.NewClient(httpClient)
	return githubClient, httpClient
}

func GitHubListReleases(organization string, repository string, page int, perPage int, client *github.Client) ([]*github.RepositoryRelease, error) {
	opts := github.ListOptions{
		Page:    page,
		PerPage: perPage,
	}

	releases, ghRes, err := client.Repositories.ListReleases(context.Background(), organization, repository, &opts)
	if err != nil {
		return nil, gitHubParseError(ghRes, err, organization, repository)
	}

	return releases, nil
}

func GitHubGetLatestPrefixedRelease(organization string, repository string, prefix string, client *github.Client) (*github.RepositoryRelease, error) {
	page := 1
	perPage := 100

	releases, err := GitHubListReleases(organization, repository, page, perPage, client)
	if err != nil {
		return nil, err
	}

	for _, release := range releases {
		if strings.HasPrefix(*(release.Name), prefix) {
			return release, nil
		}
	}

	return nil, fmt.Errorf("no GitHub release found with prefix '%s'", prefix)
}

func gitHubParseError(ghRes *github.Response, err error, organization string, repository string) error {
	if ghRes != nil && ghRes.StatusCode == http.StatusNotFound {
		return fmt.Errorf("GitHub repository not found: %s/%s", organization, repository)
	} else if _, ok := err.(*github.RateLimitError); ok {
		return fmt.Errorf("hit GitHub ratelimit while downloading latest release, try to change GITHUB_TOKEN env variable")
	} else if ghRes != nil && (ghRes.StatusCode == http.StatusForbidden || ghRes.StatusCode == http.StatusUnauthorized) {
		return fmt.Errorf("GitHub authentication failed, try to unset GITHUB_TOKEN env variable")
	} else {
		return err
	}
}

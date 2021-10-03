package golangcilint

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/pkg/errors"
	"github.com/tomnomnom/linkheader"
)

const (
	releasesURL = "https://api.github.com/repos/golangci/golangci-lint/releases"
)

type GitHubRelease struct {
	Name        string `json:"name"`
	TagName     string `json:"tag_name"`
	URL         string `json:"url"`
	CreatedAt   string `json:"created_at"`
	PublishedAt string `json:"published_at"`
}

// GetGolangCILintVersions returns a list of versions of golangci-lint on GitHub.
func GetGolangCILintVersions() ([]string, error) {
	var versions []string

	for url := releasesURL; url != ""; {
		log.Printf("[golangci-lint] Fetching %v...\n", url)
		fetched, next, err := getGolangCILintVersions(url)
		if err != nil {
			return nil, errors.Wrapf(err, "error encountered while fetching %v", url)
		}
		versions = append(versions, fetched...)
		url = next
	}

	return versions, nil
}

func getGolangCILintVersions(url string) ([]string, string, error) {
	var versions []string

	// Fetch list of known releases from GitHub API.
	resp, err := http.Get(url)
	if err != nil {
		return nil, "", err
	}
	defer func() { _ = resp.Body.Close() }()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	// Handle invalid status code.
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, "", errors.Wrapf(err, "got status code %v", resp.StatusCode)
	}

	// Unmarshal and append version list.
	var releases []*GitHubRelease
	if err := json.Unmarshal(bytes, &releases); err != nil {
		return nil, "", errors.Wrapf(err, "json decode error")
	}
	for _, release := range releases {
		versions = append(versions, release.TagName)
	}

	// Get next page URL.
	var nextURL string
	for _, link := range linkheader.Parse(resp.Header.Get("link")) {
		if link.Rel == "next" {
			nextURL = link.URL
		}
	}

	return versions, nextURL, nil
}

package version_check

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bearer/bearer/cmd/bearer/build"
	"github.com/google/go-github/github"
)

func githubClient() *github.Client {
	var httpClient = &http.Client{Timeout: 10 * time.Second}
	return github.NewClient(httpClient)
}

func GithubBinaryVersionCheck(ctx context.Context, meta *VersionMeta) error {
	client := githubClient()
	release, _, err := client.Repositories.GetLatestRelease(ctx, "bearer", "bearer")

	if err == nil {
		version := strings.TrimPrefix(*release.Name, "v")
		if version != build.Version {
			meta.Binary.Latest = false
			meta.Binary.Message = fmt.Sprintf("You are running an outdated version of Bearer CLI, %s is now available. You can find update instructions at https://docs.bearer.com/reference/installation/#updating-bearer", *release.Name)
		} else {
			meta.Binary.Latest = true
		}
	}

	return err
}

func GithubLatestRules(ctx context.Context, meta *VersionMeta, languages []string) error {
	client := githubClient()
	release, _, err := client.Repositories.GetLatestRelease(ctx, "bearer", "bearer-rules")
	if err == nil {
		if release.TagName == nil {
			return errors.New("could not find valid release for rules from github")
		}
		meta.Rules.Version = release.TagName
		for _, asset := range release.Assets {
			for _, language := range languages {
				if asset.GetName() != language+".tar.gz" {
					continue
				}
				meta.Rules.Packages[language] = asset.GetBrowserDownloadURL()
			}
		}
	}

	return err
}

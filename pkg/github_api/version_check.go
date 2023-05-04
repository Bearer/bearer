package github_api

import (
	"context"
	"strings"

	"github.com/bearer/bearer/cmd/bearer/build"
	"github.com/bearer/bearer/pkg/util/output"
	"github.com/google/go-github/github"
	"github.com/rs/zerolog/log"
)

func VersionCheck(ctx context.Context, disableVersionCheck bool, Quiet bool) {
	if disableVersionCheck {
		log.Debug().Msgf("Version checking disabled. Skipping version check")
	} else {
		client := github.NewClient(nil)
		release, _, err := client.Repositories.GetLatestRelease(ctx, "bearer", "bearer")
		if err != nil {
			log.Debug().Msgf("couldn't retrieve latest release from GitHub %s", err)
		} else {
			version := strings.TrimPrefix(*release.Name, "v")
			if version != build.Version && build.Version != "dev" && !Quiet {
				output.StdErrLogger().Msgf("You are running an outdated version of Bearer CLI, %s is now available. You can find update instructions at https://docs.bearer.com/reference/installation/#updating-bearer", *release.Name)
			}
		}
	}
}

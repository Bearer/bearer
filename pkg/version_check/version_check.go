package version_check

import (
	"context"

	"github.com/bearer/bearer/cmd/bearer/build"
	"github.com/bearer/bearer/pkg/util/output"
	"github.com/rs/zerolog/log"
)

type VersionMeta struct {
	Rules  RuleVersionMeta
	Binary BinaryVersionMeta
}

type RuleVersionMeta struct {
	Version  *string
	Packages map[string]string
}

type BinaryVersionMeta struct {
	Latest  bool
	Message string
}

func GetVersionMeta(ctx context.Context, languages []string) (*VersionMeta, error) {
	meta, err := GetBearerVerionMeta(languages)
	if err != nil {
		log.Debug().Msgf("Bearer version API failed: %s", err)
		log.Debug().Msgf("Falling back to github version check")

		var meta = &VersionMeta{
			Rules: RuleVersionMeta{
				Packages: make(map[string]string),
			},
		}
		err := GithubBinaryVersionCheck(ctx, meta)

		if err != nil {
			return nil, err
		}

		if len(languages) != 0 {
			log.Debug().Msgf("Falling back to github rules downloads - this downloads the latest version of the rules which may not be compatible with old versions")
			err := GithubLatestRules(ctx, meta, languages)
			if err != nil {
				return nil, err
			}
		}

		return meta, nil
	}

	return meta, nil
}

func DisplayBinaryVersionWarning(meta *VersionMeta, Quiet bool) {
	if !meta.Binary.Latest {
		log.Debug().Msg("Binary version is outdated")
		if build.Version != "dev" && !Quiet {
			output.StdErrLog(meta.Binary.Message + "\n")
		} else {
			log.Debug().Msg(meta.Binary.Message)
		}
	} else {
		log.Debug().Msg("Binary version is up to date")
	}
}

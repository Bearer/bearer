package version_check

import (
	"context"

	"github.com/bearer/bearer/cmd/bearer/build"
	"github.com/bearer/bearer/pkg/flag"
	"github.com/bearer/bearer/pkg/util/output"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
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

func GetVersionMeta(ctx context.Context, languages []string) (meta *VersionMeta, err error) {
	meta, err = GetBearerVerionMeta(languages)
	if err != nil {
		log.Debug().Msgf("Bearer version API failed: %s", err)

		// set default data
		meta = &VersionMeta{
			Rules: RuleVersionMeta{
				Packages: make(map[string]string),
			},
			Binary: BinaryVersionMeta{
				Latest: true,
			},
		}

		if len(languages) != 0 {
			log.Debug().Msgf("Falling back to github rules downloads - this downloads the latest version of the rules which may not be compatible with old versions")
			err = GithubLatestRules(ctx, meta, languages)
			if err != nil {
				return
			}
		}

		if checkVersion() {
			log.Debug().Msgf("Falling back to github version check")
			err = GithubBinaryVersionCheck(ctx, meta)
			if err != nil {
				return
			}
		}
	}

	return
}

func DisplayBinaryVersionWarning(meta *VersionMeta, Quiet bool) {
	if !meta.Binary.Latest && checkVersion() {
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

func checkVersion() bool {
	return !viper.GetBool(flag.DisableVersionCheckFlag.ConfigName)
}

package version_check

import (
	"encoding/json"

	"github.com/spf13/viper"

	"github.com/bearer/bearer/api"
	"github.com/bearer/bearer/internal/flag"
)

func GetBearerVerionMeta(languages []string) (*VersionMeta, error) {
	var meta VersionMeta
	client := api.New(
		api.API{
			Host: viper.GetString(flag.HostFlag.ConfigName),
		},
	)
	data, err := client.Version(languages)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &meta)
	if err != nil {
		return nil, err
	}
	return &meta, nil
}

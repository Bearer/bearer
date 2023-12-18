package flag

import (
	"testing"

	flagtypes "github.com/bearer/bearer/internal/flag/types"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test_getStringSlice(t *testing.T) {
	type env struct {
		key   string
		value string
	}
	tests := []struct {
		name      string
		flag      *flagtypes.Flag
		flagValue interface{}
		env       env
		want      []string
	}{
		{
			name:      "happy path. Empty value",
			flag:      ScannerFlag,
			flagValue: "",
			want:      nil,
		},
		{
			name:      "happy path. String value",
			flag:      ScannerFlag,
			flagValue: "sast,secrets",
			want: []string{
				string(ScannerSAST),
				string(ScannerSecrets),
			},
		},
		{
			name: "happy path. Slice value",
			flag: ScannerFlag,
			flagValue: []string{
				"sast",
				"secrets",
			},
			want: []string{
				string(ScannerSAST),
				string(ScannerSecrets),
			},
		},
		{
			name: "happy path. Env value",
			flag: ScannerFlag,
			env: env{
				key:   "BEARER_SCANNER",
				value: "secrets,sast",
			},
			want: []string{
				string(ScannerSAST),
				string(ScannerSecrets),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.env.key == "" {
				viper.Set(tt.flag.ConfigName, tt.flagValue)
			} else {
				// err := viper.BindEnv(tt.flag.ConfigName, tt.env.key)
				// assert.NoError(t, err)

				t.Setenv(tt.env.key, tt.env.value)
			}

			sl := getStringSlice(tt.flag)
			assert.Equal(t, tt.want, sl)

			viper.Reset()
		})
	}
}

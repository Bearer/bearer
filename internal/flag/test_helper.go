package flag

import (
	"fmt"
	"testing"

	flagtypes "github.com/bearer/bearer/internal/flag/types"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type Env struct {
	key   string
	value string
}

type TestCase struct {
	name      string
	flag      *flagtypes.Flag
	flagValue interface{}
	env       Env
	want      []string
}

func RunFlagTest(testCase TestCase, t *testing.T) {
	t.Run(testCase.name, func(t *testing.T) {
		if testCase.env.key == "" {
			viper.Set(testCase.flag.ConfigName, testCase.flagValue)
		} else {
			err := BindViper(testCase.flag)
			if err != nil {
				assert.NoError(t, err)
			}

			t.Setenv(testCase.env.key, testCase.env.value)
		}

		fmt.Println(testCase.name)
		fmt.Println("envVar", viper.AllEnvVar())

		sl := getStringSlice(testCase.flag)
		assert.Equal(t, testCase.want, sl)

		viper.Reset()
	})
}

func RunFlagTests(tests []TestCase, t *testing.T) {
	for _, tt := range tests {
		RunFlagTest(tt, t)
	}
}

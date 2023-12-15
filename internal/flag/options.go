package flag

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/bearer/bearer/internal/types"
	"github.com/bearer/bearer/internal/util/set"
)

var ErrInvalidScannerReportCombination = errors.New("invalid scanner argument; privacy report requires sast scanner")

type Flag struct {
	// Name is for CLI flag and environment variable.
	// If this field is empty, it will be available only in config file.
	Name string

	// ConfigName is a key in config file. It is also used as a key of viper.
	ConfigName string

	// Shorthand is a shorthand letter.
	Shorthand string

	// Value is the default value. It must be filled to determine the flag type.
	Value interface{}

	// Usage explains how to use the flag.
	Usage string

	// DisableInConfig represents if flag should be present in config
	DisableInConfig bool

	// Do not show flag in the helper
	Hide bool

	// Deprecated represents if the flag is deprecated
	Deprecated bool

	// Additional environment variables to read the value from, in addition to the default
	EnvironmentVariables []string
}

type flagGroupBase struct {
	name  string
	flags []*Flag
}

type FlagGroup interface {
	Name() string
	Flags() []*Flag
	SetOptions(options *Options, args []string) error
}

type Flags []FlagGroup

// Options holds all the runtime configuration
type Options struct {
	ReportOptions
	RuleOptions
	ScanOptions
	RepositoryOptions
	GeneralOptions
	IgnoreAddOptions
	IgnoreShowOptions
	IgnoreMigrateOptions
	WorkerOptions
}

func addFlag(cmd *cobra.Command, flag *Flag) {
	if flag == nil || flag.Name == "" {
		return
	}

	flags := cmd.Flags()

	switch v := flag.Value.(type) {
	case int:
		flags.IntP(flag.Name, flag.Shorthand, v, flag.Usage)
	case string:
		flags.StringP(flag.Name, flag.Shorthand, v, flag.Usage)
	case []string:
		flags.StringSliceP(flag.Name, flag.Shorthand, v, flag.Usage)
	case bool:
		flags.BoolP(flag.Name, flag.Shorthand, v, flag.Usage)
	case time.Duration:
		flags.DurationP(flag.Name, flag.Shorthand, v, flag.Usage)
	}
}

func bind(cmd *cobra.Command, flag *Flag) error {
	if flag == nil {
		return nil
	} else if flag.Name == "" {
		// This flag is available only in bearer.yaml
		viper.SetDefault(flag.ConfigName, flag.Value)
		return nil
	}

	if err := viper.BindPFlag(flag.ConfigName, cmd.Flags().Lookup(flag.Name)); err != nil {
		return err
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("bearer")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	arguments := append(
		[]string{flag.ConfigName},
		flag.EnvironmentVariables...,
	)

	if err := viper.BindEnv(arguments...); err != nil {
		return err
	}

	return nil
}

func argsToMap(flag *Flag) map[string]bool {
	strSlice := getStringSlice(flag)

	result := make(map[string]bool)
	for _, str := range strSlice {
		result[str] = true
	}

	return result
}

func getString(flag *Flag) string {
	if flag == nil {
		return ""
	}

	return viper.GetString(flag.ConfigName)
}

func getStringSlice(flag *Flag) []string {
	if flag == nil {
		return nil
	}
	// viper always returns a string for ENV
	// https://github.com/spf13/viper/blob/419fd86e49ef061d0d33f4d1d56d5e2a480df5bb/viper.go#L545-L553
	// and uses strings.Field to separate values (whitespace only)
	// we need to separate env values with ','
	v := viper.GetStringSlice(flag.ConfigName)
	switch {
	case len(v) == 0: // no strings
		return nil
	case len(v) == 1 && strings.Contains(v[0], ","): // unseparated string
		v = strings.Split(v[0], ",")
	}
	return v
}

func getBool(flag *Flag) bool {
	if flag == nil {
		return false
	}
	return viper.GetBool(flag.ConfigName)
}

func getDuration(flag *Flag) time.Duration {
	if flag == nil {
		return 0
	}
	return viper.GetDuration(flag.ConfigName)
}

func getInteger(flag *Flag) int {
	if flag == nil {
		return -1
	}

	return viper.GetInt(flag.ConfigName)
}

func getSeverities(flag *Flag) set.Set[string] {
	result := set.New[string]()

	for _, value := range getStringSlice(flag) {
		if !slices.Contains(types.Severities, value) {
			return nil
		}

		result.Add(value)
	}

	return result
}

func (f *flagGroupBase) add(flag Flag) *Flag {
	f.flags = append(f.flags, &flag)
	return &flag
}

func (f *flagGroupBase) Name() string {
	return f.name
}

func (f *flagGroupBase) Flags() []*Flag {
	return f.flags
}

func (f Flags) AddFlags(cmd *cobra.Command) {
	for _, group := range f {
		for _, flag := range group.Flags() {
			addFlag(cmd, flag)
		}
	}

	cmd.Flags().SetNormalizeFunc(func(f *pflag.FlagSet, name string) pflag.NormalizedName {
		return pflag.NormalizedName(name)
	})
}

func (f Flags) Usages(cmd *cobra.Command) string {
	var usages string
	for _, group := range f {
		flags := pflag.NewFlagSet(cmd.Name(), pflag.ContinueOnError)
		lflags := cmd.LocalFlags()
		for _, flag := range group.Flags() {
			if flag == nil || flag.Name == "" || flag.Hide {
				continue
			}
			flags.AddFlag(lflags.Lookup(flag.Name))
		}
		if !flags.HasAvailableFlags() {
			continue
		}

		usages += fmt.Sprintf("%s Flags\n", group.Name())
		usages += flags.FlagUsages() + "\n"
	}
	return strings.TrimSpace(usages)
}

func (f Flags) Bind(cmd *cobra.Command) error {
	return f.bind(cmd, false)
}

func (f Flags) BindForConfigInit(cmd *cobra.Command) error {
	return f.bind(cmd, true)
}

func (f Flags) bind(cmd *cobra.Command, supportIgnoreConfig bool) error {
	for _, group := range f {
		for _, flag := range group.Flags() {
			if supportIgnoreConfig && flag.DisableInConfig {
				continue
			}

			if err := bind(cmd, flag); err != nil {
				return fmt.Errorf("flag groups: %w", err)
			}
		}
	}
	return nil
}

func (f Flags) ToOptions(args []string) (Options, error) {
	// 	var err error
	options := Options{}

	for _, group := range f {
		if err := group.SetOptions(&options, args); err != nil {
			return Options{}, fmt.Errorf("%s flags error: %w", group.Name(), err)
		}
	}

	if options.ReportOptions.Report == "privacy" && !slices.Contains(options.ScanOptions.Scanner, "sast") {
		return Options{}, ErrInvalidScannerReportCombination
	}

	return options, nil
}

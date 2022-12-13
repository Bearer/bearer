package flag

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang.org/x/xerrors"
)

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

	// Deprecated represents if the flag is deprecated
	Deprecated bool
}

type FlagGroup interface {
	Name() string
	Flags() []*Flag
}

type Flags struct {
	RepoFlagGroup    *RepoFlagGroup
	ReportFlagGroup  *ReportFlagGroup
	PolicyFlagGroup  *PolicyFlagGroup
	ProcessFlagGroup *ProcessFlagGroup
	ScanFlagGroup    *ScanFlagGroup
	WorkerFlagGroup  *WorkerFlagGroup
	GeneralFlagGroup *GeneralFlagGroup
}

// Options holds all the runtime configuration
type Options struct {
	RepoOptions
	ReportOptions
	PolicyOptions
	WorkerOptions
	ScanOptions
	GeneralOptions
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
		// This flag is available only in curio.yaml
		viper.SetDefault(flag.ConfigName, flag.Value)
		return nil
	}

	if err := viper.BindPFlag(flag.ConfigName, cmd.Flags().Lookup(flag.Name)); err != nil {
		return err
	}
	// We don't use viper.AutomaticEnv, so we need to add a prefix manually here.
	if err := viper.BindEnv(flag.ConfigName, strings.ToUpper("curio_"+strings.ReplaceAll(flag.Name, "-", "_"))); err != nil {
		return err
	}

	return nil
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

func getInt(flag *Flag) int {
	if flag == nil {
		return 0
	}
	return viper.GetInt(flag.ConfigName)
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

func (f *Flags) groups() []FlagGroup {
	var groups []FlagGroup
	// This order affects the usage message, so they are sorted by frequency of use.
	if f.ReportFlagGroup != nil {
		groups = append(groups, f.ReportFlagGroup)
	}
	if f.PolicyFlagGroup != nil {
		groups = append(groups, f.PolicyFlagGroup)
	}
	if f.ScanFlagGroup != nil {
		groups = append(groups, f.ScanFlagGroup)
	}
	if f.GeneralFlagGroup != nil {
		groups = append(groups, f.GeneralFlagGroup)
	}
	if f.WorkerFlagGroup != nil {
		groups = append(groups, f.WorkerFlagGroup)
	}
	if f.ProcessFlagGroup != nil {
		groups = append(groups, f.ProcessFlagGroup)
	}
	if f.RepoFlagGroup != nil {
		groups = append(groups, f.RepoFlagGroup)
	}

	return groups
}

func (f *Flags) AddFlags(cmd *cobra.Command) {
	for _, group := range f.groups() {
		for _, flag := range group.Flags() {
			addFlag(cmd, flag)
		}
	}

	cmd.Flags().SetNormalizeFunc(func(f *pflag.FlagSet, name string) pflag.NormalizedName {
		return pflag.NormalizedName(name)
	})
}

func (f *Flags) Usages(cmd *cobra.Command) string {
	var usages string
	for _, group := range f.groups() {

		flags := pflag.NewFlagSet(cmd.Name(), pflag.ContinueOnError)
		lflags := cmd.LocalFlags()
		for _, flag := range group.Flags() {
			if flag == nil || flag.Name == "" {
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

func (f *Flags) Bind(cmd *cobra.Command) error {
	return f.bind(cmd, false)
}

func (f *Flags) BindForConfigInit(cmd *cobra.Command) error {
	return f.bind(cmd, true)
}

func (f *Flags) bind(cmd *cobra.Command, supportIgnoreConfig bool) error {
	for _, group := range f.groups() {
		if group == nil {
			continue
		}
		for _, flag := range group.Flags() {
			if supportIgnoreConfig && flag.DisableInConfig {
				continue
			}

			if err := bind(cmd, flag); err != nil {
				return xerrors.Errorf("flag groups: %w", err)
			}
		}
	}
	return nil
}

// nolint: gocyclo
func (f *Flags) ToOptions(args []string) (Options, error) {
	var err error
	opts := Options{}

	if f.RepoFlagGroup != nil {
		opts.RepoOptions = f.RepoFlagGroup.ToOptions()
	}

	if f.ReportFlagGroup != nil {
		opts.ReportOptions = f.ReportFlagGroup.ToOptions()
	}

	if f.PolicyFlagGroup != nil {
		opts.PolicyOptions = f.PolicyFlagGroup.ToOptions(args)
	}

	if f.WorkerFlagGroup != nil {
		opts.WorkerOptions = f.WorkerFlagGroup.ToOptions()
	}

	if f.ScanFlagGroup != nil {
		opts.ScanOptions, err = f.ScanFlagGroup.ToOptions(args)
		if err != nil {
			return Options{}, xerrors.Errorf("scan flag error: %w", err)
		}
	}

	if f.GeneralFlagGroup != nil {
		opts.GeneralOptions = f.GeneralFlagGroup.ToOptions()
	}

	return opts, nil
}

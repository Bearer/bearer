package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/bearer/bearer/cmd/bearer/build"
	"github.com/bearer/bearer/internal/commands"
	"github.com/bearer/bearer/internal/flag"
	flagtypes "github.com/bearer/bearer/internal/flag/types"
	"github.com/bearer/bearer/internal/util/set"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type cmdOption struct {
	Name                 string
	Shorthand            string   `yaml:",omitempty"`
	DefaultValue         string   `yaml:"default_value,omitempty"`
	Usage                string   `yaml:",omitempty"`
	EnvironmentVariables []string `yaml:"environment_variables,omitempty"`
}

type cmdDoc struct {
	Name             string
	Synopsis         string      `yaml:",omitempty"`
	Description      string      `yaml:",omitempty"`
	Usage            string      `yaml:",omitempty"`
	Options          []cmdOption `yaml:",omitempty"`
	InheritedOptions []cmdOption `yaml:"inherited_options,omitempty"`
	Example          string      `yaml:",omitempty"`
	SeeAlso          []string    `yaml:"see_also,omitempty"`
	Aliases          []string    `yaml:"aliases,omitempty"`
}

var AllFlags = []*flagtypes.Flag{}
var envVars = viper.AllEnvVar()

func main() {
	AllFlags = append(AllFlags, flag.GeneralFlagGroup.Flags()...)
	AllFlags = append(AllFlags, flag.IgnoreAddFlagGroup.Flags()...)
	AllFlags = append(AllFlags, flag.IgnoreMigrateFlagGroup.Flags()...)
	AllFlags = append(AllFlags, flag.IgnoreShowFlagGroup.Flags()...)
	AllFlags = append(AllFlags, flag.ReportFlagGroup.Flags()...)
	AllFlags = append(AllFlags, flag.RepositoryFlagGroup.Flags()...)
	AllFlags = append(AllFlags, flag.RuleFlagGroup.Flags()...)
	AllFlags = append(AllFlags, flag.ScanFlagGroup.Flags()...)
	AllFlags = append(AllFlags, flag.WorkerFlagGroup.Flags()...)

	dir := "./docs/_data"
	if _, err := os.Stat(dir); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	cmd := commands.NewApp(build.Version, build.CommitSHA)
	err := writeDocs(
		cmd,
		dir,
	)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func writeDocs(cmd *cobra.Command, dir string) error {
	// Exit if there's an error
	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
			continue
		}
		if c.HasSubCommands() {
			for _, subCmd := range c.Commands() {
				if err := writeDocs(subCmd, dir); err != nil {
					return err
				}
			}
			continue
		}
		if err := writeDocs(c, dir); err != nil {
			return err
		}
	}

	// create a file to use
	basename := "bearer.yaml"

	if cmd.CommandPath() != "" {
		basename = fmt.Sprintf("%s.yaml", strings.Replace(cmd.CommandPath(), " ", "_", -1))
	}

	filename := filepath.Join(dir, basename)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// given the file
	err = GenYaml(cmd, f)

	return err
}

func GenYaml(cmd *cobra.Command, w io.Writer) error {
	return GenYamlCustom(cmd, w, func(s string) string { return s })
}

func GenYamlCustom(cmd *cobra.Command, w io.Writer, linkHandler func(string) string) error {
	cmd.InitDefaultHelpCmd()
	cmd.InitDefaultHelpFlag()

	yamlDoc := cmdDoc{}
	yamlDoc.Name = cmd.CommandPath()
	yamlDoc.Synopsis = forceMultiLine(cmd.Short)
	yamlDoc.Description = forceMultiLine(cmd.Long)

	if cmd.Runnable() {
		yamlDoc.Usage = cmd.UseLine()
	}

	if len(cmd.Example) > 0 {
		yamlDoc.Example = cmd.Example
	}

	flags := cmd.NonInheritedFlags()
	if flags.HasFlags() {
		yamlDoc.Options = genFlagResult(flags)
	}
	flags = cmd.InheritedFlags()
	if flags.HasFlags() {
		yamlDoc.InheritedOptions = genFlagResult(flags)
	}

	if hasSeeAlso(cmd) {
		result := []string{}
		if cmd.HasParent() {
			parent := cmd.Parent()
			result = append(result, parent.CommandPath()+" - "+parent.Short)
		}
		children := cmd.Commands()
		sort.Sort(byName(children))
		for _, child := range children {
			if !child.IsAvailableCommand() || child.IsAdditionalHelpTopicCommand() {
				continue
			}
			result = append(result, child.CommandPath()+" - "+child.Short)
		}
		yamlDoc.SeeAlso = result
	}

	yamlDoc.Aliases = cmd.Aliases

	final, err := yaml.Marshal(&yamlDoc)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if _, err := w.Write(final); err != nil {
		return err
	}
	return nil
}

func searchAppFlag(flag *pflag.Flag) *flagtypes.Flag {
	var foundFlag *flagtypes.Flag
	for _, f := range AllFlags {
		if f.Name == flag.Name {
			foundFlag = f
			break
		}
	}

	return foundFlag
}

func lookupEnvVariables(flag *pflag.Flag, otherFlag *flagtypes.Flag) []string {
	var vars []string

	if otherFlag == nil {
		vars = envVars[flag.Name]
	} else {
		vars = envVars[otherFlag.ConfigName]
	}

	return removeDup(vars)
}

func removeDup(vars []string) []string {
	unique := set.New[string]()
	unique.AddAll(vars)

	sorted := unique.Items()
	slices.Sort(sorted)

	return sorted
}

func genFlagResult(flags *pflag.FlagSet) []cmdOption {
	var result []cmdOption

	flags.VisitAll(func(flag *pflag.Flag) {
		appFlag := searchAppFlag(flag)
		environmentVariables := lookupEnvVariables(flag, appFlag)

		if appFlag != nil && appFlag.Hide {
			return
		}

		if !(len(flag.ShorthandDeprecated) > 0) && len(flag.Shorthand) > 0 {
			opt := cmdOption{
				flag.Name,
				flag.Shorthand,
				flag.DefValue,
				forceMultiLine(flag.Usage),
				environmentVariables,
			}
			result = append(result, opt)
		} else {
			opt := cmdOption{
				Name:                 flag.Name,
				DefaultValue:         forceMultiLine(flag.DefValue),
				Usage:                forceMultiLine(flag.Usage),
				EnvironmentVariables: environmentVariables,
			}
			result = append(result, opt)
		}
	})

	return result
}

func forceMultiLine(s string) string {
	if len(s) > 60 && !strings.Contains(s, "\n") {
		s = s + "\n"
	}
	return s
}

func hasSeeAlso(cmd *cobra.Command) bool {
	if cmd.HasParent() {
		return true
	}
	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
			continue
		}
		return true
	}
	return false
}

type byName []*cobra.Command

func (s byName) Len() int           { return len(s) }
func (s byName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s byName) Less(i, j int) bool { return s[i].Name() < s[j].Name() }

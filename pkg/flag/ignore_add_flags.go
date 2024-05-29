package flag

import flagtypes "github.com/bearer/bearer/pkg/flag/types"

type ignoreAddFlagGroup struct{ flagGroupBase }

var IgnoreAddFlagGroup = &ignoreAddFlagGroup{flagGroupBase{name: "Ignore Add"}}

var (
	AuthorFlag = IgnoreAddFlagGroup.add(flagtypes.Flag{
		Name:       "author",
		ConfigName: "ignore_add.author",
		Shorthand:  "a",
		Value:      FormatEmpty,
		Usage:      "Add author information to this ignored finding. (default output of \"git config user.name\")",
	})

	CommentFlag = IgnoreAddFlagGroup.add(flagtypes.Flag{
		Name:       "comment",
		ConfigName: "ignore_add.comment",
		Value:      FormatEmpty,
		Usage:      "Add a comment to this ignored finding.",
	})

	FalsePositiveFlag = IgnoreAddFlagGroup.add(flagtypes.Flag{
		Name:       "false-positive",
		ConfigName: "ignore_add.false-positive",
		Value:      false,
		Usage:      "Mark an this ignored finding as false positive.",
	})

	IgnoreAddForceFlag = IgnoreAddFlagGroup.add(flagtypes.Flag{
		Name:       "force",
		ConfigName: "ignore_add.force",
		Value:      false,
		Usage:      "Overwrite an existing ignored finding.",
	})
)

type IgnoreAddOptions struct {
	Author        string `mapstructure:"author" json:"author" yaml:"author"`
	Comment       string `mapstructure:"comment" json:"comment" yaml:"comment"`
	FalsePositive bool   `mapstructure:"false_positive" json:"false_positive" yaml:"false_positive"`
	Force         bool   `mapstructure:"ignore_add_force" json:"ignore_add_force" yaml:"ignore_add_force"`
}

func (ignoreAddFlagGroup) SetOptions(options *flagtypes.Options, args []string) error {
	options.IgnoreAddOptions = flagtypes.IgnoreAddOptions{
		Author:        getString(AuthorFlag),
		Comment:       getString(CommentFlag),
		FalsePositive: getBool(FalsePositiveFlag),
		Force:         getBool(IgnoreAddForceFlag),
	}

	return nil
}

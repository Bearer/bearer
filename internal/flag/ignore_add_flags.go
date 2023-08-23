package flag

var (
	AuthorFlag = Flag{
		Name:       "author",
		ConfigName: "ignore_add.author",
		Shorthand:  "a",
		Value:      FormatEmpty,
		Usage:      "Add author information to this ignored finding. (default output of \"git config user.name\")",
	}
	CommentFlag = Flag{
		Name:       "comment",
		ConfigName: "ignore_add.comment",
		Value:      FormatEmpty,
		Usage:      "Add a comment to this ignored finding.",
	}
	FalsePositiveFlag = Flag{
		Name:       "false-positive",
		ConfigName: "ignore_add.false-positive",
		Value:      false,
		Usage:      "Mark an this ignored finding as false positive.",
	}
	IgnoreAddForceFlag = Flag{
		Name:       "force",
		ConfigName: "ignore_add.force",
		Value:      false,
		Usage:      "Overwrite an existing ignored finding.",
	}
)

type IgnoreAddFlagGroup struct {
	AuthorFlag         *Flag
	CommentFlag        *Flag
	FalsePositiveFlag  *Flag
	IgnoreAddForceFlag *Flag
}

type IgnoreAddOptions struct {
	Author        string `mapstructure:"author" json:"author" yaml:"author"`
	Comment       string `mapstructure:"comment" json:"comment" yaml:"comment"`
	FalsePositive bool   `mapstructure:"false_positive" json:"false_positive" yaml:"false_positive"`
	Force         bool   `mapstructure:"ignore_add_force" json:"ignore_add_force" yaml:"ignore_add_force"`
}

func NewIgnoreAddFlagGroup() *IgnoreAddFlagGroup {
	return &IgnoreAddFlagGroup{
		AuthorFlag:         &AuthorFlag,
		CommentFlag:        &CommentFlag,
		FalsePositiveFlag:  &FalsePositiveFlag,
		IgnoreAddForceFlag: &IgnoreAddForceFlag,
	}
}

func (f *IgnoreAddFlagGroup) Name() string {
	return "IgnoreAdd"
}

func (f *IgnoreAddFlagGroup) Flags() []*Flag {
	return []*Flag{
		f.AuthorFlag,
		f.CommentFlag,
		f.FalsePositiveFlag,
		f.IgnoreAddForceFlag,
	}
}

func (f *IgnoreAddFlagGroup) ToOptions() IgnoreAddOptions {
	return IgnoreAddOptions{
		Author:        getString(f.AuthorFlag),
		Comment:       getString(f.CommentFlag),
		FalsePositive: getBool(f.FalsePositiveFlag),
		Force:         getBool(f.IgnoreAddForceFlag),
	}
}

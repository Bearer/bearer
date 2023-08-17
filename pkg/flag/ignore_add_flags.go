package flag

var (
	AuthorFlag = Flag{
		Name:       "author",
		ConfigName: "ignore_add.author",
		Shorthand:  "a",
		Value:      FormatEmpty,
		Usage:      "Specify report format (json, yaml, sarif, gitlab-sast, rdjson, html)",
	}
	CommentFlag = Flag{
		Name:       "comment",
		ConfigName: "ignore_add.comment",
		Value:      FormatEmpty,
		Usage:      "Add a comment to this ignored finding.",
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
	IgnoreAddForceFlag *Flag
}

type IgnoreAddOptions struct {
	Author  string `mapstructure:"author" json:"author" yaml:"author"`
	Comment string `mapstructure:"comment" json:"comment" yaml:"comment"`
	Force   bool   `mapstructure:"ignore_add_force" json:"ignore_add_force" yaml:"ignore_add_force"`
}

func NewIgnoreAddFlagGroup() *IgnoreAddFlagGroup {
	return &IgnoreAddFlagGroup{
		AuthorFlag:         &AuthorFlag,
		CommentFlag:        &CommentFlag,
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
		f.IgnoreAddForceFlag,
	}
}

func (f *IgnoreAddFlagGroup) ToOptions() IgnoreAddOptions {
	return IgnoreAddOptions{
		Author:  getString(f.AuthorFlag),
		Comment: getString(f.CommentFlag),
		Force:   getBool(f.IgnoreAddForceFlag),
	}
}

package flag

var (
	AuthorFlag = Flag{
		Name:       "author",
		ConfigName: "ignore.author",
		Shorthand:  "a",
		Value:      FormatEmpty,
		Usage:      "Specify report format (json, yaml, sarif, gitlab-sast, rdjson, html)",
	}
	CommentFlag = Flag{
		Name:       "comment",
		ConfigName: "ignore.comment",
		Value:      FormatEmpty,
		Usage:      "Add a comment to this ignored finding.",
	}
	IgnoreForceFlag = Flag{
		Name:       "force",
		ConfigName: "ignore.force",
		Value:      false,
		Usage:      "Overwrite an existing ignored finding.",
	}
)

type IgnoreFlagGroup struct {
	AuthorFlag      *Flag
	CommentFlag     *Flag
	IgnoreForceFlag *Flag
}

type IgnoreOptions struct {
	Author  string `mapstructure:"author" json:"author" yaml:"author"`
	Comment string `mapstructure:"comment" json:"comment" yaml:"comment"`
	Force   bool   `mapstructure:"force" json:"force" yaml:"force"`
}

func NewIgnoreFlagGroup() *IgnoreFlagGroup {
	return &IgnoreFlagGroup{
		AuthorFlag:      &AuthorFlag,
		CommentFlag:     &CommentFlag,
		IgnoreForceFlag: &IgnoreForceFlag,
	}
}

func (f *IgnoreFlagGroup) Name() string {
	return "Ignore"
}

func (f *IgnoreFlagGroup) Flags() []*Flag {
	return []*Flag{
		f.AuthorFlag,
		f.CommentFlag,
		f.IgnoreForceFlag,
	}
}

func (f *IgnoreFlagGroup) ToOptions() IgnoreOptions {
	return IgnoreOptions{
		Author:  getString(f.AuthorFlag),
		Comment: getString(f.CommentFlag),
		Force:   getBool(f.IgnoreForceFlag),
	}
}

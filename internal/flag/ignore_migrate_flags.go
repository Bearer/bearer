package flag

type ignoreMigrateFlagGroup struct{ flagGroupBase }

var IgnoreMigrateFlagGroup = &ignoreMigrateFlagGroup{flagGroupBase{name: "Ignore Migrate"}}

var (
	IgnoreMigrateForceFlag = IgnoreMigrateFlagGroup.add(Flag{
		Name:       "force",
		ConfigName: "ignore_migrate.force",
		Value:      false,
		Usage:      "Overwrite an existing ignored finding.",
	})
)

type IgnoreMigrateOptions struct {
	Force bool `mapstructure:"ignore_migrate_force" json:"ignore_migrate_force" yaml:"ignore_migrate_force"`
}

func (ignoreMigrateFlagGroup) SetOptions(options *Options, args []string) error {
	options.IgnoreMigrateOptions = IgnoreMigrateOptions{
		Force: getBool(IgnoreMigrateForceFlag),
	}

	return nil
}

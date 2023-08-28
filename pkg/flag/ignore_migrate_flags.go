package flag

var (
	IgnoreMigrateForceFlag = Flag{
		Name:       "force",
		ConfigName: "ignore_migrate.force",
		Value:      false,
		Usage:      "Overwrite an existing ignored finding.",
	}
)

type IgnoreMigrateFlagGroup struct {
	IgnoreMigrateForceFlag *Flag
}

type IgnoreMigrateOptions struct {
	Force bool `mapstructure:"ignore_migrate_force" json:"ignore_migrate_force" yaml:"ignore_migrate_force"`
}

func NewIgnoreMigrateFlagGroup() *IgnoreMigrateFlagGroup {
	return &IgnoreMigrateFlagGroup{
		IgnoreMigrateForceFlag: &IgnoreMigrateForceFlag,
	}
}

func (f *IgnoreMigrateFlagGroup) Name() string {
	return "IgnoreMigrate"
}

func (f *IgnoreMigrateFlagGroup) Flags() []*Flag {
	return []*Flag{
		f.IgnoreMigrateForceFlag,
	}
}

func (f *IgnoreMigrateFlagGroup) ToOptions() IgnoreMigrateOptions {
	return IgnoreMigrateOptions{
		Force: getBool(f.IgnoreMigrateForceFlag),
	}
}

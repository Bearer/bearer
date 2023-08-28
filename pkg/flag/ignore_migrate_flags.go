package flag

var (
	IgnoreMigrateForceFlag = Flag{
		Name:       "force",
		ConfigName: "ignore_migrate.force",
		Value:      false,
		Usage:      "Overwrite an existing ignored finding.",
	}
	IgnoreMigrateBearerIgnoreFileFlag = Flag{
		Name:            "config-file",
		ConfigName:      "config-file",
		Value:           "bearer.yml",
		Usage:           "Load configuration from the specified path.",
		DisableInConfig: true,
	}
	IgnoreMigrateConfigFileFlag = Flag{
		Name:            "config-file",
		ConfigName:      "config-file",
		Value:           "bearer.yml",
		Usage:           "Load configuration from the specified path.",
		DisableInConfig: true,
	}
)

type IgnoreMigrateFlagGroup struct {
	IgnoreMigrateForceFlag            *Flag
	IgnoreMigrateConfigFileFlag       *Flag
	IgnoreMigrateBearerIgnoreFileFlag *Flag
}

type IgnoreMigrateOptions struct {
	Force      bool   `mapstructure:"ignore_migrate_force" json:"ignore_migrate_force" yaml:"ignore_migrate_force"`
	ConfigFile string `mapstructure:"ignore_migrate_config_file" json:"ignore_migrate_config_file" yaml:"ignore_migrate_config_file"`
}

func NewIgnoreMigrateFlagGroup() *IgnoreMigrateFlagGroup {
	return &IgnoreMigrateFlagGroup{
		IgnoreMigrateForceFlag:            &IgnoreMigrateForceFlag,
		IgnoreMigrateBearerIgnoreFileFlag: &IgnoreMigrateBearerIgnoreFileFlag,
		IgnoreMigrateConfigFileFlag:       &IgnoreMigrateConfigFileFlag,
	}
}

func (f *IgnoreMigrateFlagGroup) Name() string {
	return "IgnoreMigrate"
}

func (f *IgnoreMigrateFlagGroup) Flags() []*Flag {
	return []*Flag{
		f.IgnoreMigrateForceFlag,
		f.IgnoreMigrateConfigFileFlag,
	}
}

func (f *IgnoreMigrateFlagGroup) ToOptions() IgnoreMigrateOptions {
	return IgnoreMigrateOptions{
		Force:      getBool(f.IgnoreMigrateForceFlag),
		ConfigFile: getString(f.IgnoreMigrateConfigFileFlag),
	}
}

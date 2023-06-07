package flag

var (
	PortFlag = Flag{
		Name:       "port",
		ConfigName: "process.port",
		Shorthand:  "p",
		Value:      "",
		Usage:      "Set the server's listening port.",
	}

	WorkerIDFlag = Flag{
		Name:       "worker-id",
		ConfigName: "process.worker-id",
		Value:      "",
		Usage:      "Set the worker's identifier.",
	}

	WorkerDebugProfileFlag = Flag{
		Name:       "debug-profile",
		ConfigName: "process.debug-profile",
		Value:      false,
		Usage:      "Generate profiling data for debugging",
	}
)

type ProcessFlagGroup struct {
	PortFlag               *Flag
	WorkerIDFlag           *Flag
	WorkerDebugProfileFlag *Flag
}

type ProcessOptions struct {
	WorkerID     string `mapstructure:"worker-id" json:"worker-id" yaml:"worker-id"`
	Port         string `mapstructure:"port" json:"port" yaml:"port"`
	DebugProfile bool   `mapstructure:"debug_profile" json:"debug_profile" yaml:"debug_profile"`
}

func NewProcessGroup() *ProcessFlagGroup {
	return &ProcessFlagGroup{
		PortFlag:               &PortFlag,
		WorkerIDFlag:           &WorkerIDFlag,
		WorkerDebugProfileFlag: &WorkerDebugProfileFlag,
	}
}

func (f *ProcessFlagGroup) Name() string {
	return "process"
}

func (f *ProcessFlagGroup) Flags() []*Flag {
	return []*Flag{
		f.PortFlag,
		f.WorkerIDFlag,
		f.WorkerDebugProfileFlag,
	}
}

func (f *ProcessFlagGroup) ToOptions() (ProcessOptions, error) {
	port := getString(f.PortFlag)
	workerID := getString(f.WorkerIDFlag)

	return ProcessOptions{
		Port:         port,
		WorkerID:     workerID,
		DebugProfile: getBool(f.WorkerDebugProfileFlag),
	}, nil
}

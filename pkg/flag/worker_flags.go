package flag

import flagtypes "github.com/bearer/bearer/pkg/flag/types"

type workerFlagGroup struct{ flagGroupBase }

var WorkerFlagGroup = &workerFlagGroup{flagGroupBase{name: "Worker"}}

var (
	ParentProcessIDFlag = WorkerFlagGroup.add(flagtypes.Flag{
		Name:       "parent-process-id",
		ConfigName: "worker.parent-process-id",
		Value:      -1,
	})
	PortFlag = WorkerFlagGroup.add(flagtypes.Flag{
		Name:       "port",
		ConfigName: "worker.port",
		Shorthand:  "p",
		Value:      "",
		Usage:      "Set the server's listening port.",
	})
	WorkerIDFlag = WorkerFlagGroup.add(flagtypes.Flag{
		Name:       "worker-id",
		ConfigName: "worker.id",
		Value:      "",
		Usage:      "Set the worker's identifier.",
	})
)

func init() {
	WorkerFlagGroup.add(*LogLevelFlag)
	WorkerFlagGroup.add(*QuietFlag)
	WorkerFlagGroup.add(*DebugProfileFlag)
}

type WorkerOptions struct {
	ParentProcessID int
	WorkerID        string `mapstructure:"worker-id" json:"worker-id" yaml:"worker-id"`
	Port            string `mapstructure:"port" json:"port" yaml:"port"`
}

func (workerFlagGroup) SetOptions(options *flagtypes.Options, args []string) error {
	options.WorkerOptions = flagtypes.WorkerOptions{
		ParentProcessID: getInteger(ParentProcessIDFlag),
		Port:            getString(PortFlag),
		WorkerID:        getString(WorkerIDFlag),
	}

	return nil
}

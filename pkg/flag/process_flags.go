package flag

var (
	PortFlag = Flag{
		Name:       "port",
		ConfigName: "port",
		Shorthand:  "p",
		Value:      "",
		Usage:      "server listening on what port",
	}
)

type ProcessFlagGroup struct {
	PortFlag *Flag
}

type ProcessOptions struct {
	Port string
}

func NewProcessGroup() *ProcessFlagGroup {
	return &ProcessFlagGroup{
		PortFlag: &PortFlag,
	}
}

func (f *ProcessFlagGroup) Name() string {
	return "process"
}

func (f *ProcessFlagGroup) Flags() []*Flag {
	return []*Flag{
		f.PortFlag,
	}
}

func (f *ProcessFlagGroup) ToOptions() (ProcessOptions, error) {
	port := getString(f.PortFlag)

	return ProcessOptions{
		Port: port,
	}, nil
}

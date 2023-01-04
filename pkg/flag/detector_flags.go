package flag

var (
	SkipDetectorFlag = Flag{
		Name:       "skip-detector",
		ConfigName: "detector.skip-detector",
		Value:      []string{},
		Usage:      "Specify the comma-separated ids of the detectors you would like to skip. Runs all other detectors.",
	}
	OnlyDetectorFlag = Flag{
		Name:       "only-detector",
		ConfigName: "detector.only-detector",
		Value:      []string{},
		Usage:      "Specify the comma-separated ids of the detectors you would like to run. Skips all other detectors.",
	}
)

type DetectorFlagGroup struct {
	SkipDetectorFlag *Flag
	OnlyDetectorFlag *Flag
}

type DetectorOptions struct {
	SkipDetector map[string]bool `mapstructure:"skip-detector" json:"skip-detector" yaml:"skip-detector"`
	OnlyDetector map[string]bool `mapstructure:"only-detector" json:"only-detector" yaml:"only-detector"`
}

func NewDetectorFlagGroup() *DetectorFlagGroup {
	return &DetectorFlagGroup{
		SkipDetectorFlag: &SkipDetectorFlag,
		OnlyDetectorFlag: &OnlyDetectorFlag,
	}
}

func (f *DetectorFlagGroup) Name() string {
	return "Detector"
}

func (f *DetectorFlagGroup) Flags() []*Flag {
	return []*Flag{
		f.SkipDetectorFlag,
		f.OnlyDetectorFlag,
	}
}

func (f *DetectorFlagGroup) ToOptions(args []string) DetectorOptions {
	return DetectorOptions{
		SkipDetector: argsToMap(f.SkipDetectorFlag),
		OnlyDetector: argsToMap(f.OnlyDetectorFlag),
	}
}

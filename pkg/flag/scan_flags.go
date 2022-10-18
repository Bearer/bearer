package flag

var (
	SkipDirsFlag = Flag{
		Name:       "skip-dirs",
		ConfigName: "scan.skip-dirs",
		Value:      []string{},
		Usage:      "specify the directories where the traversal is skipped",
	}
	SkipFilesFlag = Flag{
		Name:       "skip-files",
		ConfigName: "scan.skip-files",
		Value:      []string{},
		Usage:      "specify the file paths to skip traversal",
	}
)

type ScanFlagGroup struct {
	SkipDirs  *Flag
	SkipFiles *Flag
}

type ScanOptions struct {
	Target    string
	SkipDirs  []string
	SkipFiles []string
}

func NewScanFlagGroup() *ScanFlagGroup {
	return &ScanFlagGroup{
		SkipDirs:  &SkipDirsFlag,
		SkipFiles: &SkipFilesFlag,
	}
}

func (f *ScanFlagGroup) Name() string {
	return "Scan"
}

func (f *ScanFlagGroup) Flags() []*Flag {
	return []*Flag{f.SkipDirs, f.SkipFiles}
}

func (f *ScanFlagGroup) ToOptions(args []string) (ScanOptions, error) {
	var target string
	if len(args) == 1 {
		target = args[0]
	}

	return ScanOptions{
		Target:    target,
		SkipDirs:  getStringSlice(f.SkipDirs),
		SkipFiles: getStringSlice(f.SkipFiles),
	}, nil
}

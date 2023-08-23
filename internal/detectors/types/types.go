package types

import (
	"github.com/bearer/bearer/internal/report"
	"github.com/bearer/bearer/internal/util/file"
)

type DetectorConstructor func() Detector

type Detector interface {
	AcceptDir(dir *file.Path) (bool, error)
	ProcessFile(file *file.FileInfo, dir *file.Path, report report.Report) (bool, error)
}

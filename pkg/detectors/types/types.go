package types

import (
	"github.com/bearer/curio/pkg/report"
	"github.com/bearer/curio/pkg/util/file"
)

type DetectorConstructor func() Detector

type Detector interface {
	AcceptDir(dir *file.Path) (bool, error)
	ProcessFile(file *file.FileInfo, dir *file.Path, report report.Report) (bool, error)
}

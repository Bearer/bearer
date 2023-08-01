package files

import (
	"time"

	bbftypes "github.com/bearer/bearer/pkg/report/basebranchfindings/types"
)

type List struct {
	Files     []File
	BaseFiles []File
	Renames   map[string]string
	Chunks    map[string]bbftypes.Chunks
}

type File struct {
	Timeout  time.Duration
	FilePath string
}

type LineMapping struct {
	Base,
	Delta int
}

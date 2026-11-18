package files

import (
	"time"

	"github.com/bearer/bearer/pkg/git"
	"github.com/bearer/bearer/pkg/report/coverage"
)

type List struct {
	Files     []File
	BaseFiles []File
	Renames   map[string]string
	Chunks    map[string]git.Chunks
	Coverage  coverage.Discovery
}

type File struct {
	Timeout  time.Duration
	FilePath string
}

type LineMapping struct {
	Base,
	Delta int
}

package files

import (
	"time"

	"github.com/bearer/bearer/internal/git"
)

type List struct {
	Files     []File
	BaseFiles []File
	Renames   map[string]string
	Chunks    map[string]git.Chunks
}

type File struct {
	Timeout  time.Duration
	FilePath string
}

type LineMapping struct {
	Base,
	Delta int
}

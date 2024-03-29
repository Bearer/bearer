package cache

import (
	"os"
	"path/filepath"
)

// DefaultCacheDir returns/creates the cache-dir to be used for bearer operations
func DefaultDir() string {
	tmpDir, err := os.UserCacheDir()
	if err != nil {
		tmpDir = os.TempDir()
	}
	return filepath.Join(tmpDir, "bearer")
}

package memory

import (
	"fmt"
	"runtime"

	"github.com/bearer/curio/pkg/commands/process/settings"
)

type OutOfMemoryError struct {
	MaximumMemoryMb int
	FilesProcessed  int
}

func (err *OutOfMemoryError) Error() string {
	return fmt.Sprintf("maximum memory reached %dMB when processed %d number of files", err.MaximumMemoryMb, err.FilesProcessed)
}

type MemoryConstraint struct {
	settings       settings.TypeSettings
	processedFiles int
}

func NewConstraint(settings settings.TypeSettings) *MemoryConstraint {
	return &MemoryConstraint{
		settings:       settings,
		processedFiles: 0,
	}
}

func (constraint *MemoryConstraint) FileProcessed() {
	constraint.processedFiles++
}

// CheckOverflow returns OutOfMemory error upon overflow
func (constraint *MemoryConstraint) CheckOverflow() error {
	if constraint.processedFiles%constraint.settings.MemoryCheckEachFiles != 0 {
		return nil
	}

	if UsageMb() < constraint.settings.MaximumMemoryMb {
		return nil
	}

	runtime.GC()

	if UsageMb() < constraint.settings.MaximumMemoryMb {
		return nil
	}

	return &OutOfMemoryError{
		MaximumMemoryMb: int(constraint.settings.MaximumMemoryMb),
		FilesProcessed:  constraint.processedFiles,
	}
}

// returns memory used in Mbytes
func UsageMb() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return bToMb(m.Sys)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

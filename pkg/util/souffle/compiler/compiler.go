package compiler

import (
	"fmt"
	"os/exec"
)

func Compile(inputFilename, outputFilename string) error {
	cmd := exec.Command("souffle", "-g", outputFilename, inputFilename)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("souffle compilation error: %w", err)
	}

	return nil
}

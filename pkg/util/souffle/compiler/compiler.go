package compiler

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"regexp"
)

var logSplitPattern = regexp.MustCompile(`[\r\n]+`)

type debugLogWriter struct {
	allOutput bytes.Buffer
}

func (writer *debugLogWriter) Write(data []byte) (int, error) {
	writer.allOutput.Write(data)

	return len(data), nil
}

func Compile(inputFilename, outputFilename string) error {
	log.Printf("running souffle")
	cmd := exec.Command("souffle", "--verbose", "-g", outputFilename, inputFilename)

	logWriter := &debugLogWriter{}
	cmd.Stderr = logWriter
	cmd.Stdout = logWriter

	err := cmd.Run()

	for _, line := range logSplitPattern.Split(logWriter.allOutput.String(), -1) {
		log.Printf("[souffle] %s", line)
	}

	if err != nil {
		return fmt.Errorf("souffle compilation error: %w", err)
	}

	return nil
}

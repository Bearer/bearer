package tmpfile

import (
	"errors"
	"fmt"
	"os"

	"github.com/bearer/bearer/internal/util/output"
)

var ErrCreateFailed = errors.New("failed to create file")

func Create(ext string) string {
	outputFile, err := os.CreateTemp("", "*"+ext)
	if err != nil {
		output.Fatal(fmt.Sprintf("got create fail error %s %s", err, ErrCreateFailed))
	}
	outputFile.Close()

	return outputFile.Name()
}

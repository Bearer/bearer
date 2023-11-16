package git

import (
	"context"
	"strings"
)

func GetOriginURL(dir string) (string, error) {
	output, err := captureCommandBasic(
		context.TODO(),
		dir,
		"remote",
		"get-url",
		"origin",
	)

	return strings.TrimSpace(output), err
}

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

	if err != nil && strings.Contains(err.Error(), "No such remote 'origin'") {
		return "", nil
	}

	return strings.TrimSpace(output), err
}

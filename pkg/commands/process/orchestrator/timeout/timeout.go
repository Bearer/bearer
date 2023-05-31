package timeout

import (
	"io/fs"
	"time"

	"github.com/bearer/bearer/pkg/commands/process/settings"
)

func Assign(d fs.DirEntry, config settings.Config) time.Duration {
	var timeout time.Duration

	fileInfo, err := d.Info()
	if err != nil {
		return time.Duration(0)
	}

	timeout = config.Worker.TimeoutFileMinimum

	timeoutFileSize := time.Duration(fileInfo.Size() / int64(config.Worker.TimeoutFileBytesPerSecond) * int64(time.Second))
	if timeoutFileSize > timeout {
		if timeoutFileSize > config.Worker.TimeoutFileMaximum {
			timeout = config.Worker.TimeoutFileMaximum
		} else {
			timeout = timeoutFileSize
		}
	}

	return timeout
}

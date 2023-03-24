package s3

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"github.com/bearer/bearer/api"
	"github.com/google/uuid"
)

func SignForAPI(req *UploadRequestS3) (*api.RequestFileUpload, error) {
	fileUuid := uuid.NewString()

	reportFile, err := os.Open(req.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file for upload %e", err)
	}
	defer reportFile.Close()

	stats, err := reportFile.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat file %e", err)
	}

	hash := md5.New()
	_, err = io.Copy(hash, reportFile)
	if err != nil {
		return nil, fmt.Errorf("failed copying file content to hash %e", err)
	}

	checksumMD5 := hash.Sum(nil)

	return &api.RequestFileUpload{
		Checksum:        base64.StdEncoding.EncodeToString(checksumMD5[:]),
		ByteSize:        int(stats.Size()),
		UUID:            fileUuid,
		Prefix:          req.FilePrefix,
		ContentType:     req.ContentType,
		ContentEncoding: req.ContentEncoding,
	}, nil
}

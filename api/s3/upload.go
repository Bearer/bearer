package s3

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type UploadRequest struct {
	Client   *http.Client
	FilePath string
	FileSize int64
	URL      string
	Headers  map[string]string
}

func Upload(req UploadRequest) error {
	reportFile, err := os.Open(req.FilePath)
	if err != nil {
		return fmt.Errorf("failed to open file for uploading: %s", err)
	}
	defer reportFile.Close()

	request, err := http.NewRequest("PUT", req.URL, reportFile)
	request.ContentLength = req.FileSize
	if err != nil {
		return fmt.Errorf("failed to create upload request: %s", err)
	}
	defer request.Body.Close()

	for key, value := range req.Headers {
		request.Header.Add(key, value)
	}

	response, err := req.Client.Do(request)
	if err != nil {
		return fmt.Errorf("failed to upload file: %s", err)
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		responseBody, _ := io.ReadAll(response.Body)
		return fmt.Errorf("file upload returned error status: %d\n%s", response.StatusCode, string(responseBody))
	}

	return nil
}

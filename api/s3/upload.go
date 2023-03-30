package s3

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/bearer/bearer/api"
	"github.com/rs/zerolog/log"
)

type UploadRequest struct {
	Client   *http.Client
	FilePath string
	FileSize int64
	URL      string
	Headers  map[string]string
}

type UploadRequestS3 struct {
	Api             *api.API
	FilePath        string
	FilePrefix      string
	FileType        string
	ContentType     string
	ContentEncoding string
}

func GetSignedURL(req UploadRequest) error {
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

func UploadS3(req *UploadRequestS3) (fileUploadOffer *api.FileUploadOffer, err error) {
	requestFileUploadAction, err := SignForAPI(req)
	if err != nil {
		return nil, err
	}

	log.Debug().Msgf("Sending S3 upload request to Bearer API...")
	fileUploadOffer, err = req.Api.RequestFileUpload(*requestFileUploadAction, "")
	if err != nil {
		return nil, err
	}

	log.Debug().Msgf("Uploading file to Bearer S3...")
	err = GetSignedURL(UploadRequest{
		Client:   api.UploadClient,
		FilePath: req.FilePath,
		FileSize: int64(requestFileUploadAction.ByteSize),
		URL:      fileUploadOffer.DirectUpload.URL,
		Headers:  fileUploadOffer.DirectUpload.Headers,
	})

	if err != nil {
		return nil, err
	}

	return fileUploadOffer, nil
}

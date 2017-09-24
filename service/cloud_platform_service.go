package service

import (
	"net/http"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"mime/multipart"
	"path/filepath"
	"io"
)

type CloudPlatformService struct {
	baseUrl string
}

func NewCloudPlatformService(baseUrl string) CloudPlatformService {
	return CloudPlatformService{baseUrl: baseUrl}
}

func (service CloudPlatformService) SendCommand(commandName string, requestBody interface{}) (*http.Response, error) {
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(requestBody)

	url := fmt.Sprintf("%s/command/%s", service.baseUrl, commandName)

	return http.Post(url, "application/json; charset=utf-8", buffer)
}

func (service CloudPlatformService) UploadComponent(name string, label string, filePath string) (*http.Response, error) {
	extraParams := map[string]string{
		"name":  name,
		"label": label,
	}

	request, err := newFileUploadRequest(fmt.Sprintf("%s/component", service.baseUrl), extraParams, "component", filePath)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func newFileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}

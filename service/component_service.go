package service

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type ComponentService interface {
	upload(name string, label string) (*http.Response, error)
}

type ComponentServiceOp struct {
	baseUrl string
}

type UploadComponentParams struct {
	Name  string `url:"name,omitempty"`
	Label string `url:"label,omitempty"`
}

func NewComponentService(baseUrl string) *ComponentServiceOp {
	return &ComponentServiceOp{
		baseUrl,
	}
}

func (s ComponentServiceOp) Upload(name string, label string, filePath string) (*http.Response, error) {
	extraParams := map[string]string{
		"name":  name,
		"label": label,
	}

	request, err := newFileUploadRequest(fmt.Sprintf("%s/component", s.baseUrl), extraParams, "component", filePath)
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

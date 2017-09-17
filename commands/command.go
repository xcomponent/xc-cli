package commands

import (
	"net/http"
	"bytes"
	"encoding/json"
	"fmt"
)

const protocol = "https"

type Command interface {
	Execute(args []string) error
}

type Output struct {
	Content string `json:"content"`
}

type Service struct {
}

type ServiceOp struct {
	baseUrl string
}

func NewCommandService(baseUrl string) *ServiceOp {
	return &ServiceOp{baseUrl: baseUrl}
}

func (service ServiceOp) ExecuteCommand(commandName string, requestBody interface{}) (*http.Response, error) {
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(requestBody)

	url := fmt.Sprintf("%s/command/%s", service.baseUrl, commandName)

	return http.Post(url, "application/json; charset=utf-8", buffer)
}

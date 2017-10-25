package services

import (
	"github.com/daniellavoie/go-shim"
	"net/http"
	"encoding/json"
)

//go:generate counterfeiter -o servicesfake/fake_httpservice.go . HttpService
type HttpService interface {
	Get(url string) (resp *http.Response, err error)
	GetJson(url string, target interface{}) error
}

func NewHttpService(http goshim.Http) *HttpServiceImpl {
	return &HttpServiceImpl{http}
}

type HttpServiceImpl struct {
	http goshim.Http
}

func (httpService *HttpServiceImpl) Get(url string) (resp *http.Response, err error) {
	return httpService.http.Get(url)
}

func (httpService *HttpServiceImpl) GetJson(url string, target interface{}) error {
	response, err := httpService.http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return json.NewDecoder(response.Body).Decode(target)
}
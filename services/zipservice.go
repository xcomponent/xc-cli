package services

import (
	"github.com/daniellavoie/go-shim"
	"archive/zip"
)

type ZipService interface {
	OpenReader(name string) (*zip.ReadCloser, error)
}

func NewZipService(zip goshim.Zip) *ZipServiceImpl {
	return &ZipServiceImpl{zip}
}

type ZipServiceImpl struct {
	zip goshim.Zip
}

func (zipService *ZipServiceImpl) OpenReader(name string) (*zip.ReadCloser, error){
	return zipService.zip.OpenReader(name)
}

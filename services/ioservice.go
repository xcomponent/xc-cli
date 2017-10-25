package services

import (
	"io"
	"io/ioutil"
)

//go:generate counterfeiter -o servicesfake/fake_ioservice.go . IoService
type IoService interface {
	Copy(dst io.Writer, src io.Reader) (written int64, err error)
	TempDir(dir, prefix string) (name string, err error)
}

func NewIoService() *IoServiceImpl {
	return &IoServiceImpl{}
}

type IoServiceImpl struct{}

func (ioService *IoServiceImpl) Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	return io.Copy(dst, src)
}

func (ioService *IoServiceImpl) TempDir(dir, prefix string) (name string, err error) {
	return ioutil.TempDir(dir, prefix)
}

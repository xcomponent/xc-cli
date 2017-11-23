package services

import (
	"io"
	"io/ioutil"
	"os"
)

//go:generate counterfeiter -o servicesfake/fake_ioservice.go . IoService
type IoService interface {
	Copy(dst io.Writer, src io.Reader) (written int64, err error)
	ReadFile(filename string) ([]byte, error)
	TempDir(dir, prefix string) (name string, err error)
	WriteFile(filename string, data []byte, perm os.FileMode) error
}

func NewIoService() *IoServiceImpl {
	return &IoServiceImpl{}
}

type IoServiceImpl struct{}

func (ioService *IoServiceImpl) Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	return io.Copy(dst, src)
}

func (ioService *IoServiceImpl) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func (ioService *IoServiceImpl) TempDir(dir, prefix string) (name string, err error) {
	return ioutil.TempDir(dir, prefix)
}

func (ioService *IoServiceImpl) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return ioutil.WriteFile(filename, data, perm)
}


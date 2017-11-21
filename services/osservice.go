package services

import (
	"os"
)

//go:generate counterfeiter -o servicesfake/fake_osservice.go . OsService
type OsService interface {
	Create(name string) (*os.File, error)
	Exit(exitCode int)
	Getwd() (string, error)
	IsNotExist(err error) bool
	Mkdir(name string, perm os.FileMode) error
	MkdirAll(path string, perm os.FileMode) error
	Open(name string) (*os.File, error)
	OpenFile(name string, flag int, perm os.FileMode) (*os.File, error)
	Remove(name string) error
	RemoveAll(path string) error
	Stat(name string) (os.FileInfo, error)
}

func NewOsService() *OsServiceImpl {
	return &OsServiceImpl{}
}

type OsServiceImpl struct{}

func (osService *OsServiceImpl) Create(name string) (*os.File, error) {
	return os.Create(name)
}

func (osService *OsServiceImpl) Getwd() (string, error) {
	return os.Getwd()
}

func (osService *OsServiceImpl) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}

func (osService *OsServiceImpl) Exit(code int) {
	os.Exit(code)
}

func (OsService *OsServiceImpl) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}

func (osService *OsServiceImpl) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (osService *OsServiceImpl) Open(name string) (*os.File, error) {
	return os.Open(name)
}

func (osService *OsServiceImpl) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(name, flag, perm)
}

func (osService *OsServiceImpl) Remove(name string) error {
	return os.Remove(name)
}

func (osService *OsServiceImpl) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (osService *OsServiceImpl) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

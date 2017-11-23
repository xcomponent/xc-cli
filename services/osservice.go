package services

import (
	"os"
)

//go:generate counterfeiter -o servicesfake/fake_osservice.go . OsService
type OsService interface {
	Create(name string) (*os.File, error)
	Exit(exitCode int)
	Getwd() (string, error)
	GetPathSeperator() rune
	IsNotExist(err error) bool
	Mkdir(name string, perm os.FileMode) error
	MkdirAll(path string, perm os.FileMode) error
	Open(name string) (*os.File, error)
	OpenFile(name string, flag int, perm os.FileMode) (*os.File, error)
	Remove(name string) error
	RemoveAll(path string) error
	Rename(oldpath string, newpath string) error
	Stat(name string) (os.FileInfo, error)
}

func NewOsService() *OsServiceImpl {
	return &OsServiceImpl{}
}

type OsServiceImpl struct{}

func (*OsServiceImpl) Create(name string) (*os.File, error) {
	return os.Create(name)
}

func (*OsServiceImpl) Getwd() (string, error) {
	return os.Getwd()
}

func (*OsServiceImpl) GetPathSeperator() rune {
	return os.PathSeparator
}

func (*OsServiceImpl) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}

func (*OsServiceImpl) Exit(code int) {
	os.Exit(code)
}

func (*OsServiceImpl) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}

func (*OsServiceImpl) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (*OsServiceImpl) Open(name string) (*os.File, error) {
	return os.Open(name)
}

func (*OsServiceImpl) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(name, flag, perm)
}

func (*OsServiceImpl) Remove(name string) error {
	return os.Remove(name)
}

func (*OsServiceImpl) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (*OsServiceImpl) Rename(oldpath string, newpath string) error {
	return os.Rename(oldpath, newpath)
}

func (*OsServiceImpl) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

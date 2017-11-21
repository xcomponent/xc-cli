package services

import (
	"os/exec"
	"github.com/daniellavoie/go-shim"
)

type ExecService interface {
	Command(name string, arg ...string) *exec.Cmd
}

func NewExecService(exec goshim.Exec) *ExecServiceImpl {
	return &ExecServiceImpl{exec}
}

type ExecServiceImpl struct {
	exec goshim.Exec
}

func (execService *ExecServiceImpl) Command(name string, arg ...string) *exec.Cmd{
	return execService.exec.Command(name, arg...)
}
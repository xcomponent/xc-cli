package commands_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/xcomponent/xc-cli/commands"
	"github.com/daniellavoie/go-shim/zipshim/fake_zip"
	"archive/zip"
	"github.com/daniellavoie/go-shim/execshim/exec_fake"
	"os/exec"
	"github.com/xcomponent/xc-cli/services"
	"github.com/xcomponent/xc-cli/services/servicesfake"
	"github.com/daniellavoie/go-shim/httpshim"
	"encoding/json"
	"strings"
	"errors"
	"fmt"
	"net/http"
	"io"
	"os"
)

var _ = Describe("Install", func() {
	var installConfig = `{
  "xcStudioDistribs" : {
    "x86" : "https://test.com/xcomponent.msi",
    "amd64" : "https://test.com/xcomponent.msi"
  }
}`

	var installConfigUrl string
	var osName string
	var osArch string

	var err error
	var osFake = &servicesfake.FakeOsService{}
	var httpFake = &servicesfake.FakeHttpService{}
	var ioFake = &servicesfake.FakeIoService{}
	var zipFake *fake_zip.FakeZip
	var execFake *fake_exec.FakeExec

	BeforeEach(func() {
		installConfigUrl = "https://raw.githubusercontent.com/xcomponent/xc-cli/install-config-v1/install-config.json"
		osName = "windows"
		osArch = "amd64"

		var osService = services.NewOsService()
		osFake.StatStub = osService.Stat
		osFake.IsNotExistStub = osService.IsNotExist
		osFake.MkdirAllStub = osService.MkdirAll
		osFake.MkdirStub = osService.Mkdir
		osFake.OpenStub = osService.Open
		osFake.CreateStub = osService.Create

		var httpService = services.NewHttpService(&httpshim.HttpShim{})
		httpFake.GetStub = httpService.Get
		httpFake.GetJsonStub = func(url string, target interface{}) error {
			json.NewDecoder(strings.NewReader(installConfig)).Decode(target)

			return nil
		}

		var ioService = services.NewIoService()
		ioFake.CopyStub = ioService.Copy
		ioFake.TempDirStub = ioService.TempDir

		zipFake = new(fake_zip.FakeZip)
		zipFake.OpenReaderStub = zip.OpenReader

		execFake = new(fake_exec.FakeExec)
		execFake.CommandStub = func(name string, arg ...string) *exec.Cmd {
			if name == "powershell" || name == "msiexec" {
				return exec.Command("echo", "True")
			} else {
				return exec.Command(name, arg...)
			}
		}

		err = nil
	})

	JustBeforeEach(func() {
		var installCommand = commands.NewInstallCommand(
			osFake,
			ioFake,
			httpFake,
			services.NewExecService(execFake))

		err = installCommand.Install(installConfigUrl, osName, osArch)
	})

	Context("Not windows", func() {
		BeforeEach(func() {
			osName = "darwin"
		})

		It("should fail", func() {
			Expect(err.Error()).To(Equal("Install command is only supported on Windows"))
		})
	})

	Context("checkdot net command error", func() {
		BeforeEach(func() {
			execFake.CommandStub = func(name string, arg ...string) *exec.Cmd {
				return exec.Command("badcommand")
			}
		})

		It("should fail", func() {
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "exec: \"badcommand\": executable file not found in")).To(Equal(true))
		})
	})

	Context("donet not installed", func() {
		BeforeEach(func() {
			execFake.CommandStub = func(name string, arg ...string) *exec.Cmd {
				return exec.Command("echo", "False")
			}
		})

		It("should fail", func() {
			Expect(err.Error()).To(Equal("XComponent requires a version of DotNet higher than 4.5"))
		})
	})

	Context("install config fails to load", func() {
		var getJsonErr = errors.New("could not download config")

		BeforeEach(func() {
			httpFake.GetJsonStub = func(url string, target interface{}) error {
				return getJsonErr
			}
		})

		It("should fail", func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(getJsonErr.Error()))
		})
	})

	Context("unsupported arch", func() {
		BeforeEach(func() {
			osArch = "arm"
		})

		It("should fail", func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(fmt.Sprintf("xc install does not support %s arch", osArch)))
		})
	})

	Context("could not create temp dir for msi", func() {
		var tempDirErr = errors.New("temp dir error")

		BeforeEach(func() {
			ioFake.TempDirStub = func(dir, prefix string) (name string, err error) {
				return "", tempDirErr
			}
		})

		It("should fail", func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(tempDirErr.Error()))
		})
	})

	Context("could not write msi file", func() {
		var createErr = errors.New("create error")

		BeforeEach(func() {
			osFake.CreateStub = func(name string) (*os.File, error) {
				return nil, createErr
			}
		})

		It("should fail", func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(createErr.Error()))
		})
	})

	Context("could not download msi", func() {
		var downloadMsiErr = errors.New("download msi error")

		BeforeEach(func() {
			httpFake.GetStub = func(url string) (resp *http.Response, err error) {
				return nil, downloadMsiErr
			}
		})

		It("should fail", func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(downloadMsiErr.Error()))
		})
	})

	Context("could not write msi content", func() {
		var copyErr = errors.New("copy error")

		BeforeEach(func() {
			ioFake.CopyStub = func(dst io.Writer, src io.Reader) (written int64, err error) {
				return 0, copyErr
			}
		})

		It("should fail", func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(copyErr.Error()))
		})
	})

	Context("msi install fails", func() {
		BeforeEach(func() {
			execFake.CommandStub = func(name string, arg ...string) *exec.Cmd {
				if name == "powershell" {
					return exec.Command("echo", "True")
				} else {
					return exec.Command("bad-command")
				}
			}
		})

		It("should fail", func() {
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "exec: \"bad-command\": executable file not found in")).To(Equal(true))
		})
	})

	Context("successful install", func() {
		It("should not fail", func() {
			Expect(err).ToNot(HaveOccurred())
		})
	})
})

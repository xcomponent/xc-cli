package commands_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xcomponent/xc-cli/commands"
	"github.com/daniellavoie/go-shim/httpshim"
	"github.com/daniellavoie/go-shim/zipshim"
	"github.com/daniellavoie/go-shim/execshim"
	"github.com/xcomponent/xc-cli/services"
	"os"
)

var _ = Describe("Config", func() {
	It("should support 2 commands", func() {
		workDir, _ := os.Getwd()

		commands := GetCommands(
			workDir,
			services.NewOsService(),
			services.NewHttpService(new(httpshim.HttpShim)),
			services.NewIoService(),
			services.NewZipService(new(zipshim.ZipShim)),
			services.NewExecService(new(execshim.ExecShim)))

		Expect(commands).ShouldNot(BeNil())
		Expect(len(commands)).Should(Equal(4))
	})
})

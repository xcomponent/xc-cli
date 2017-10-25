package commands_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
    . "github.com/xcomponent/xc-cli/commands"
	"github.com/daniellavoie/go-shim/osshim"
	"github.com/daniellavoie/go-shim/httpshim"
	"github.com/daniellavoie/go-shim/ioshim"
)

var _ = Describe("Config", func() {
	It("should support 2 commands", func() {
		commands := GetCommands(new(osshim.OsShim), new(httpshim.HttpShim), new(ioshim.IoShim))

		Expect(commands).ShouldNot(BeNil())
		Expect(len(commands)).Should(Equal(2))
	})
})

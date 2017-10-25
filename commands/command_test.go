package commands_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
    . "github.com/xcomponent/xc-cli/commands"
)

var _ = Describe("Config", func() {
	Context("GetCommands", func() {
		It("should support 2 commands", func() {
			commands := GetCommands()

			Expect(commands).ShouldNot(BeNil())
			Expect(len(commands)).Should(Equal(2))
		})
	})
})

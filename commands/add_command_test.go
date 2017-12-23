package commands_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"github.com/xcomponent/xc-cli/commands"
	"os"
	"github.com/daniellavoie/go-shim/zipshim/fake_zip"
	"archive/zip"
	"github.com/xcomponent/xc-cli/services"
	"github.com/xcomponent/xc-cli/services/servicesfake"
	"github.com/daniellavoie/go-shim/httpshim"
	"path"
)

var _ = Describe("Add", func() {
	var projectFolder string
	var addTemplatesGithubOrg string
	var elementType string
	var elementName string
	var err error
	var osFake = &servicesfake.FakeOsService{}
	var httpFake = &servicesfake.FakeHttpService{}
	var ioFake = &servicesfake.FakeIoService{}
	var zipFake *fake_zip.FakeZip

	BeforeEach(func() {
		addTemplatesGithubOrg = "xcomponent-add-templates"
		elementType = "component"
		elementName = "test-component"

		var osService = services.NewOsService()
		osFake.StatStub = osService.Stat
		osFake.IsNotExistStub = osService.IsNotExist
		osFake.MkdirAllStub = osService.MkdirAll
		osFake.MkdirStub = osService.Mkdir
		osFake.OpenStub = osService.Open
		osFake.OpenFileStub = osService.OpenFile
		osFake.CreateStub = osService.Create
		osFake.RenameStub = osService.Rename

		var ioService = services.NewIoService()
		ioFake.CopyStub = ioService.Copy

		zipFake = new(fake_zip.FakeZip)
		zipFake.OpenReaderStub = zip.OpenReader

		var httpService = services.NewHttpService(&httpshim.HttpShim{})
		httpFake.GetStub = httpService.Get
		httpFake.GetJsonStub = httpService.GetJson

		err = nil
	})

	JustBeforeEach(func() {
		err = commands.NewAddCommand(
			osFake,
			httpFake,
			ioFake,
			services.NewZipService(zipFake),
		).Execute(projectFolder, elementType, elementName, addTemplatesGithubOrg)
	})

	Context("Add component", func() {
		BeforeEach(func() {
			var tempErr error
			projectFolder, tempErr = ioutil.TempDir("", "add-element-test")
			if tempErr != nil {
				panic(tempErr)
			}
		})

		It("should succeed", func() {
			Expect(err).ToNot(HaveOccurred())

			elementCreated(projectFolder, elementName)
		})
	})

	AfterEach(func() {
		os.RemoveAll(projectFolder)
	})
})

func elementCreated(projectFolder string, elementName string) {
	files, err := ioutil.ReadDir(path.Join(projectFolder, elementName))
	Expect(err).ShouldNot(HaveOccurred())
	Expect(len(files)).Should(Equal(1))
}

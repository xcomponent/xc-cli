package commands_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/xcomponent/xc-cli/commands"
	"github.com/xcomponent/xc-cli/services/servicesfake"
	"io/ioutil"
	"github.com/xcomponent/xc-cli/util"
	"os"
	"path/filepath"
	"path"
)

const fixtureDir = "fixtures/rename-project-test"
const fixtureExpectedDir = "fixtures/rename-project-test-expected"

var _ = Describe("Add", func() {
	Context("ExtractProjectName", func() {
		It("returns expected project name", func() {
			Expect("Toto").To(Equal(commands.ExtractProjectName("Toto_Model.xcml")))
		})
	})

	Context("Rename project", func() {
		var workDir string
		var err error
		var ioService *servicesfake.FakeIoService
		var osService *servicesfake.FakeOsService

		BeforeEach(func() {
			var setupErr error
			workDir, setupErr = ioutil.TempDir("", "")
			if setupErr != nil {
				panic(setupErr)
			}

			setupErr = copyFixtures(workDir)
			if setupErr != nil {
				panic(setupErr)
			}
			ioService = &servicesfake.FakeIoService{}
			osService = &servicesfake.FakeOsService{}

			ioService.ReadDirStub = ioutil.ReadDir
			ioService.ReadFileStub = ioutil.ReadFile
			ioService.WriteFileStub = ioutil.WriteFile

			osService.StatStub = os.Stat
			osService.RenameStub = os.Rename
		})

		JustBeforeEach(func() {
			command := commands.NewRenameProjectCommand(ioService, osService, util.NewFileUtils(osService, ioService))

			err = command.Execute(workDir, "Toto")
		})

		It("should rename project without error", func() {
			Expect(err).ToNot(HaveOccurred())

			Expect(checksProjectRenamedCorrectly(fixtureExpectedDir, workDir)).ToNot(HaveOccurred())
		})

		AfterEach(func() {
			cleanupErr := os.RemoveAll(workDir)
			if cleanupErr != nil {
				panic(cleanupErr)
			}
		})
	})
})

func copyFixtures(tempDir string) error {
	fixturesFolder, err := filepath.Abs(fixtureDir)
	if err != nil {
		return err
	}

	err = os.MkdirAll(path.Join(tempDir, "Configuration.NewProject", "Dev"), 0700)
	if err != nil {
		return err
	}

	err = util.CopyFile(path.Join(fixturesFolder, "NewProject_Model.xcml"),
		path.Join(tempDir, "NewProject_Model.xcml"),
		0700,
	)
	if err != nil {
		return err
	}

	err = util.CopyFile(path.Join(fixturesFolder, "Configuration.NewProject", "applications.xml"),
		path.Join(tempDir, "Configuration.NewProject", "applications.xml"),
		0700,
	)

	if err != nil {
		return err
	}

	return util.CopyFile(
		path.Join(fixturesFolder, "Configuration.NewProject", "Dev", "NewProject_Deployment_Configuration.xml"),
		path.Join(tempDir, "Configuration.NewProject", "Dev", "NewProject_Deployment_Configuration.xml"),
		0700,
	)
}

func checksProjectRenamedCorrectly(fixtureDir string, workDir string) error {
	err := compareFiles(fixtureDir, workDir, "Toto_Model.xcml")
	if err != nil {
		return err
	}
	err = compareFiles(fixtureDir, workDir, "Configuration.Toto", "applications.xml")
	if err != nil {
		return err
	}
	err = compareFiles(fixtureDir, workDir, "Configuration.Toto", "Dev", "Toto_Deployment_Configuration.xml")
	if err != nil {
		return err
	}

	return nil
}

func compareFiles(fixtureDir string, workDir string, elem ... string) error {
	expectedFileContent, err := ioutil.ReadFile(path.Join(fixtureDir, path.Join(elem...)))
	if err != nil {
		return err
	}

	fileContent, err := ioutil.ReadFile(path.Join(workDir, path.Join(elem...)))
	if err != nil {
		return err
	}

	Expect(string(expectedFileContent)).To(Equal(string(fileContent)))

	return nil
}

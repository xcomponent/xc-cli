package commands_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"path/filepath"
	"errors"
	"github.com/xcomponent/xc-cli/commands"
	"os"
	"fmt"
	"net/http"
	"io"
	"github.com/daniellavoie/go-shim/zipshim/fake_zip"
	"archive/zip"
	"github.com/xcomponent/xc-cli/services"
	"github.com/xcomponent/xc-cli/services/servicesfake"
	"github.com/daniellavoie/go-shim/httpshim"
)

var _ = Describe("Init", func() {
	var projectFolder string
	var githubOrg string
	var templateName string
	var projectName string
	var err error
	var osFake = &servicesfake.FakeOsService{}
	var httpFake = &servicesfake.FakeHttpService{}
	var ioFake = &servicesfake.FakeIoService{}
	var zipFake *fake_zip.FakeZip

	var makeDirErr error

	BeforeEach(func() {
		githubOrg = "xcomponent-templates"
		templateName = "default"
		projectName = "test-project"

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
		err = commands.NewInitCommand(
			httpFake,
			ioFake,
			osFake,
			services.NewZipService(zipFake),
		).Init(projectFolder, githubOrg, templateName, projectName)
	})

	Context("Init default project", func() {
		BeforeEach(func() {
			var tempErr error
			projectFolder, tempErr = ioutil.TempDir("", "init-test")
			if tempErr != nil {
				panic(tempErr)
			}
		})

		Context("default project name", func() {
			BeforeEach(func() {
				projectName = ""
			})

			It("should not fail", func() {
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("project folder does not already exists", func() {

			It("should download github template project", func() {
				projectFolderInitialized(projectFolder)
			})

			Context("can't create project folder", func() {
				BeforeEach(func() {
					makeDirErr = errors.New("Could not create project folder dir")
					osFake.MkdirAllStub = func(path string, perm os.FileMode) error {
						return makeDirErr
					}
				})

				It("project folder creation fails", func() {
					Expect(err).To(HaveOccurred())
					Expect(makeDirErr).To(Equal(err))
				})
			})
		})

		Context("project folder already exists", func() {
			BeforeEach(func() {
				projectFolder, err = ioutil.TempDir("", "init-test")
				Expect(err).ShouldNot(HaveOccurred())
			})

			Context("and can't be read", func() {
				var cannotReadErr error
				BeforeEach(func() {
					osFake.StatStub = func(name string) (os.FileInfo, error) {
						cannotReadErr = errors.New("Folder can't be read")

						return nil, cannotReadErr
					}
				})

				It("should return error", func() {
					Expect(err).Should(HaveOccurred())
					Expect(err).To(Equal(cannotReadErr))
				})
			})

			Context("is a file", func() {
				BeforeEach(func() {
					file, fileErr := os.Create(filepath.Join(projectFolder, "file"))
					Expect(fileErr).ToNot(HaveOccurred())

					projectFolder = file.Name()
				})

				It("should return error", func() {
					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(Equal(fmt.Sprintf("%s is not a directory", projectFolder)))
				})
			})

			Context("fails on empty check", func() {
				var openErr = errors.New("Failed to open file")

				BeforeEach(func() {
					osFake.OpenStub = func(name string) (*os.File, error) {
						return nil, openErr
					}
				})

				It("should return error", func() {
					Expect(err).Should(HaveOccurred())
					Expect(err).To(Equal(openErr))
				})
			})

			Context("and not empty", func() {
				BeforeEach(func() {
					_, fileErr := os.Create(filepath.Join(projectFolder, "file"))
					Expect(fileErr).ToNot(HaveOccurred())
				})

				It("should return error", func() {
					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(Equal(fmt.Sprintf("%s is not empty", projectFolder)))
				})
			})

			It("should download github template project", func() {
				projectFolderInitialized(projectFolder)
			})
		})

		Context("zip can't be written", func() {
			var createZipErr = errors.New("zip file cannot be written")

			BeforeEach(func() {
				osFake.CreateStub = func(name string) (*os.File, error) {
					return nil, createZipErr
				}
			})

			It("should return error", func() {
				Expect(err).Should(HaveOccurred())
				Expect(err).To(Equal(createZipErr))
			})
		})

		Context("download failure", func() {
			Context("http failure", func() {
				var httpErr = errors.New("could not download file")

				BeforeEach(func() {
					httpFake.GetStub = func(url string) (resp *http.Response, err error) {
						return nil, httpErr
					}
				})

				It("should return error", func() {
					Expect(err).Should(HaveOccurred())
					Expect(err).To(Equal(httpErr))
				})
			})

			Context("disk write failure", func() {
				var ioErr = errors.New("could not write file")

				BeforeEach(func() {
					ioFake.CopyStub = func(dst io.Writer, src io.Reader) (written int64, err error) {
						return 0, ioErr
					}
				})

				It("should return error", func() {
					Expect(err).Should(HaveOccurred())
					Expect(err).To(Equal(ioErr))
				})
			})
		})

		Context("unknow template", func() {
			BeforeEach(func() {
				templateName = "unknow-template"
			})

			It("should fail", func() {
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(Equal(fmt.Sprintf("Project template %s/%s could not be found", githubOrg, templateName)))
			})
		})

		Context("unzip failure", func() {
			var zipErr = errors.New("could not unzip file")

			BeforeEach(func() {
				zipFake.OpenReaderStub = func(name string) (*zip.ReadCloser, error) {
					return nil, zipErr
				}
			})

			It("should return error", func() {
				Expect(err).Should(HaveOccurred())
				Expect(err).To(Equal(zipErr))
			})
		})

		Context("create folder on unzip failure", func() {
			var mkdirErr = errors.New("could not create zip folder")
			var activateMkdirFailure = false

			BeforeEach(func() {
				zipFake.OpenReaderStub = func(name string) (*zip.ReadCloser, error) {
					activateMkdirFailure = true
					return zip.OpenReader(name)
				}

				osFake.MkdirAllStub = func(path string, perm os.FileMode) error {
					if activateMkdirFailure {
						return mkdirErr
					} else {
						return os.MkdirAll(path, perm)
					}
				}
			})

			It("should return error", func() {
				Expect(err).Should(HaveOccurred())
				Expect(err).To(Equal(mkdirErr))
			})
		})
	})

	AfterEach(func() {
		os.RemoveAll(projectFolder)
	})
})

func projectFolderInitialized(projectFolder string) {
	files, err := ioutil.ReadDir(projectFolder)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(len(files)).Should(Equal(3))
}

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
	"github.com/daniellavoie/go-shim/httpshim/fake_http"
	"github.com/daniellavoie/go-shim/ioshim/fake_io"
	"github.com/daniellavoie/go-shim/osshim/fake_os"
)

var _ = Describe("Init", func() {
	var projectFolder string
	var githubOrg = "xcomponent-templates"
	var templateName = "default"
	var err error
	var osshim *fake_os.FakeOs
	var httpFake *fake_http.FakeHttp
	var ioFake *fake_io.FakeIo

	var makeDirErr error

	BeforeEach(func() {
		osshim = new(fake_os.FakeOs)
		osshim.StatStub = os.Stat
		osshim.IsNotExistStub = os.IsNotExist
		osshim.MkdirAllStub = os.MkdirAll
		osshim.MkdirStub = os.Mkdir
		osshim.OpenStub = os.Open
		osshim.CreateStub = os.Create

		httpFake = new(fake_http.FakeHttp)
		httpFake.GetStub = http.Get

		ioFake = new(fake_io.FakeIo)
		ioFake.CopyStub = io.Copy

		err = nil
	})

	JustBeforeEach(func() {
		err = commands.Init(projectFolder, githubOrg, templateName, osshim, httpFake, ioFake);
	})

	Context("Init default project", func() {
		BeforeEach(func() {
			projectFolder = filepath.Join(os.TempDir(), "init-test")
		})

		Context("project folder does not already exists", func() {

			It("should download github template project", func() {
				projectFolderInitialized(projectFolder)
			})

			Context("can't create project folder", func() {
				BeforeEach(func() {
					makeDirErr = errors.New("Could not create project folder dir")
					osshim.MkdirAllStub = func(path string, perm os.FileMode) error {
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
					osshim.StatStub = func(name string) (os.FileInfo, error) {
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
					osshim.OpenStub = func(name string) (*os.File, error) {
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
				osshim.CreateStub = func(name string) (*os.File, error) {
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
	})

	AfterEach(func() {
		os.RemoveAll(projectFolder)
	})
})

func projectFolderInitialized(projectFolder string) {
	files, err := ioutil.ReadDir(projectFolder)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(len(files)).Should(Equal(2))
}

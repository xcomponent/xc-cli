package util_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/xcomponent/xc-cli/util"
	"github.com/xcomponent/xc-cli/services/servicesfake"
	"fmt"
	"io/ioutil"
	"path"
	"os"
	"errors"
	"io"
	"path/filepath"
)

var _ = Describe("FileUtils", func() {
	var projectName = "test-project"

	var fileUtil *util.FileUtilsImpl
	var fakeOs *servicesfake.FakeOsService
	var fakeIo *servicesfake.FakeIoService
	var err error

	BeforeEach(func() {
		fakeOs = &servicesfake.FakeOsService{}
		fakeOs.RenameStub = os.Rename
		fakeOs.GetPathSeperatorStub = func() rune {
			return os.PathSeparator
		}

		fakeIo = &servicesfake.FakeIoService{}
		fakeIo.ReadFileStub = ioutil.ReadFile
		fakeIo.WriteFileStub = ioutil.WriteFile

		fileUtil = util.NewFileUtils(fakeOs, fakeIo)
	})

	Context("ReplaceInFileName", func() {

		var oldFile string
		var expectedNewFile string
		var newFilename string

		BeforeEach(func() {
			tempDir, tempDirErr := ioutil.TempDir("", "xc-rename")
			if tempDirErr != nil {
				panic(tempDir)
			}

			oldFile = path.Join(tempDir, "NewProject_Model.xcml")
			file, createErr := os.Create(oldFile)
			if createErr != nil {
				panic(createErr)
			}
			defer file.Close()

			expectedNewFile = path.Join(tempDir, fmt.Sprintf("%s_Model.xcml", projectName))
		})

		JustBeforeEach(func() {
			newFilename, err = fileUtil.ReplaceInFileName(oldFile, "NewProject", projectName)
		})

		It("renames file without error ", func() {
			Expect(err).ToNot(HaveOccurred())
			Expect(newFilename).To(Equal(expectedNewFile))

			file, err := os.OpenFile(newFilename, os.O_RDONLY, 0666)
			Expect(err).ToNot(HaveOccurred())
			Expect(file.Name()).To(Equal(expectedNewFile))
		})

		Context("file name does not contain replacement", func() {
			BeforeEach(func() {
				oldFile = "test.txt"
				expectedNewFile = "test.txt"
			})

			It("should not fail", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(newFilename).To(Equal(expectedNewFile))
			})
		})

		Context("file rename fails", func() {
			var renameErr = errors.New("rename failed")

			BeforeEach(func() {
				fakeOs.RenameStub = func(oldpath string, newpath string) error {
					return renameErr
				}
			})

			It("should return error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(renameErr.Error()))
			})
		})

		Context("unix", func() {
			BeforeEach(func() {
				oldFile = "/tmp/NewProject_Model.xcml"
				expectedNewFile = fmt.Sprintf("/tmp/%s_Model.xcml", projectName)

				fakeOs.GetPathSeperatorStub = func() rune {
					return '/'
				}

				fakeOs.RenameStub = func(oldpath string, newpath string) error {
					return nil
				}
			})

			It("works with unix path seperator", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(newFilename).To(Equal(expectedNewFile))
			})
		})

		Context("windows", func() {
			BeforeEach(func() {
				oldFile = "C:\\tmp\\NewProject_Model.xcml"
				expectedNewFile = fmt.Sprintf("C:\\tmp\\%s_Model.xcml", projectName)

				fakeOs.GetPathSeperatorStub = func() rune {
					return '\\'
				}

				fakeOs.RenameStub = func(oldpath string, newpath string) error {
					return nil
				}
			})

			It("works with windows path seperator", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(newFilename).To(Equal(expectedNewFile))
			})
		})
	})

	Context("ReplaceInFile", func() {
		JustBeforeEach(func() {
			var fixturesFolder string
			var localErr error

			fixturesFolder, localErr = filepath.Abs("fixtures/replace-projectname-test")
			if localErr != nil {
				panic(localErr)
			}

			testFolder, localErr := ioutil.TempDir("", "xc-replace")
			if localErr != nil {
				panic(localErr)
			}

			copyDir(fixturesFolder, testFolder)

			err = fileUtil.RenameAndReplaceFiles(testFolder, "NewProject", projectName)
		})

		It("should not fail", func() {
			Expect(err).ToNot(HaveOccurred())

			// Check file content.
		})

		Context("read file error", func() {
			readFileErr := errors.New("read file error")

			BeforeEach(func() {
				fakeIo.ReadFileStub = func(filename string) ([]byte, error) {
					return nil, readFileErr
				}
			})

			It("should fail", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(readFileErr.Error()))
			})
		})

		Context("write file error", func() {
			writeFileErr := errors.New("write file error")

			BeforeEach(func() {
				fakeIo.WriteFileStub = func(filename string, data []byte, perm os.FileMode) error {
					return writeFileErr
				}
			})

			It("should fail", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(writeFileErr.Error()))
			})
		})

		Context("file rename fails", func() {
			renameErr := errors.New("rename failed")

			BeforeEach(func() {
				fakeOs.RenameStub = func(oldpath string, newpath string) error {
					if fakeOs.RenameCallCount() == 2 {
						return renameErr
					} else {
						return os.Rename(oldpath, newpath)
					}
				}
			})

			It("should fail", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(renameErr.Error()))
			})
		})
	})
})

func copyDir(source string, dest string) error {
	// get properties of source dir
	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	// create dest dir

	err = os.MkdirAll(dest, sourceinfo.Mode())
	if err != nil {
		return err
	}

	directory, err := os.Open(source)
	if err != nil {
		return err
	}

	objects, err := directory.Readdir(-1)
	if err != nil {
		return err
	}

	for _, obj := range objects {
		sourcefilepointer := source + "/" + obj.Name()

		destinationfilepointer := dest + "/" + obj.Name()

		if obj.IsDir() {
			err = copyDir(sourcefilepointer, destinationfilepointer)
		} else {
			err = copyFile(sourcefilepointer, destinationfilepointer)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func copyFile(source string, dest string) error {
	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}

	defer sourcefile.Close()

	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer destfile.Close()

	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		sourceinfo, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, sourceinfo.Mode())
		}
	}

	return nil
}

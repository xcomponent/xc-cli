package commands

import (
	. "os"
	"fmt"
	"errors"
	"archive/zip"
	"path/filepath"
	"strings"
	io "io"
	"github.com/daniellavoie/go-shim"
)

func Init(projectFolder string, githubOrg, templateName string, osshim goshim.Os, httpshim goshim.Http, io goshim.Io) error {
	err := prepareProjectFolder(projectFolder, osshim)
	if err != nil {
		return err
	}

	fileName, err := downloadTemplate(projectFolder, githubOrg, templateName, osshim, httpshim, io)
	defer clean(fileName, osshim)
	if err != nil {
		return err
	}

	return unzip(fileName, projectFolder, osshim, io)
}

func prepareProjectFolder(projectFolder string, osshim goshim.Os) error {
	fileInfo, err := osshim.Stat(projectFolder)
	if err != nil {
		if osshim.IsNotExist(err) {
			return osshim.MkdirAll(projectFolder, 0700)
		} else {
			return err
		}
	} else {
		if !fileInfo.IsDir() {
			return errors.New(fmt.Sprintf("%s is not a directory", projectFolder))
		}

		empty, err := isEmpty(projectFolder, osshim)
		if err != nil {
			return err
		} else if !empty {
			return errors.New(fmt.Sprintf("%s is not empty", projectFolder))
		}
	}

	return nil
}

func isEmpty(name string, osshim goshim.Os) (bool, error) {
	f, err := osshim.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

func downloadTemplate(projectFolder string, githubOrg string, templateName string, osshim goshim.Os, httpshim goshim.Http, io goshim.Io) (string, error) {
	fmt.Println("Downloading template")

	fileName := fmt.Sprintf("%s/%s.zip", projectFolder, templateName)

	out, err := osshim.Create(fileName)
	defer out.Close()
	if err != nil {
		return fileName, err
	}

	resp, err := httpshim.Get(fmt.Sprintf("https://github.com/%s/%s/archive/master.zip", githubOrg, templateName))
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return fileName, err
	}

	if resp.StatusCode == 404 {
		return fileName, errors.New(fmt.Sprintf("Project template %s/%s could not be found", githubOrg, templateName))
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fileName, err
	}

	return fileName, nil
}

func unzip(archive string, target string, osshim goshim.Os, io goshim.Io) error {
	fmt.Println(fmt.Sprintf("Extracting template from %s", archive))

	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

	if err := osshim.MkdirAll(target, 0755); err != nil {
		return err
	}

	for _, file := range reader.File {
		err = unzipFile(file, target, osshim, io)
		if err != nil {
			return err
		}
	}

	return nil
}

func unzipFile(file *zip.File, target string, osshim goshim.Os, io goshim.Io) error {
	fileName := file.Name[strings.Index(file.Name, "/")+1:len(file.Name)]

	path := filepath.Join(target, fileName)
	if file.FileInfo().IsDir() {
		osshim.MkdirAll(path, file.Mode())
		return nil
	}

	fileReader, err := file.Open()
	if err != nil {
		return err
	}
	defer fileReader.Close()

	targetFile, err := osshim.OpenFile(path, O_WRONLY|O_CREATE|O_TRUNC, file.Mode())
	if targetFile != nil {
		defer targetFile.Close()
	}
	if err != nil {
		return err
	}

	if _, err := io.Copy(targetFile, fileReader); err != nil {
		return err
	}

	return nil
}

func clean(fileName string, osshim goshim.Os) error {
	if fileName == "" {
		return nil
	}

	fmt.Println(fmt.Sprintf("Cleaning %s", fileName))

	return osshim.Remove(fileName)
}

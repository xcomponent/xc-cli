package commands

import (
	"fmt"
	"os"
	"net/http"
	"io"
	"errors"
	"archive/zip"
	"path/filepath"
	"strings"
)

func Init(projectFolder string, githubOrg, templateName string) error {
	err := prepareProjectFolder(projectFolder)
	if err != nil {
		return err
	}

	fileName, err := downloadTemplate(projectFolder, githubOrg, templateName)
	defer clean(fileName)
	if err != nil {
		return err
	}

	return unzip(fileName, projectFolder)
}

func prepareProjectFolder(projectFolder string) error {
	fileInfo, err := os.Stat(projectFolder)
	if err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(projectFolder, os.ModeDir)
		} else {
			return err
		}
	} else {
		if !fileInfo.IsDir() {
			return errors.New(fmt.Sprintf("%s is not a directory", projectFolder))
		}

		empty, err := isEmpty(projectFolder)
		if err != nil {
			return err
		} else if !empty {
			return errors.New(fmt.Sprintf("%s is not empty", projectFolder))
		}
	}

	return nil
}

func isEmpty(name string) (bool, error) {
	f, err := os.Open(name)
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

func downloadTemplate(projectFolder string, githubOrg string, templateName string) (string, error) {
	fmt.Println("Downloading template")

	fileName := fmt.Sprintf("%s/%s.zip", projectFolder, templateName)

	out, err := os.Create(fileName)
	defer out.Close()
	if err != nil {
		return fileName, err
	}

	resp, err := http.Get(fmt.Sprintf("https://github.com/%s/%s/archive/master.zip", githubOrg, templateName))
	defer resp.Body.Close()
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

func unzip(archive string, target string) error {
	fmt.Println(fmt.Sprintf("Extracting template from %s", archive))

	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	for _, file := range reader.File {
		err = unzipFile(file, target)
		if err != nil {
			return err
		}
	}

	return nil
}

func unzipFile(file *zip.File, target string) error {
	fileName := file.Name[strings.Index(file.Name, "/")+1:len(file.Name)]

	path := filepath.Join(target, fileName)
	if file.FileInfo().IsDir() {
		os.MkdirAll(path, file.Mode())
		return nil
	}

	fileReader, err := file.Open()
	if err != nil {
		return err
	}
	defer fileReader.Close()

	targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
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

func clean(fileName string) error {
	if fileName == "" {
		return nil
	}

	fmt.Println(fmt.Sprintf("Cleaning %s", fileName))

	return os.Remove(fileName)
}

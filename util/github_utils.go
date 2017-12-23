package util

import (
	"archive/zip"
	"errors"
	"fmt"
	"github.com/xcomponent/xc-cli/services"
	"path/filepath"
	"strings"
	"os"
)

type GitHubUtils struct {
	os   services.OsService
	http services.HttpService
	io   services.IoService
	zip  services.ZipService
}

func NewGitHubUtils(os services.OsService, http services.HttpService, io services.IoService,
	zip services.ZipService) *GitHubUtils {

	return &GitHubUtils{os, http, io, zip}
}

func (gitHubUtils *GitHubUtils) DownloadTemplate(templateDestination string, githubOrg string, templateName string,
	newName string) error {
	fileName, err := gitHubUtils.downloadZip(templateDestination, githubOrg, templateName)
	defer gitHubUtils.clean(fileName)

	if err != nil {
		return err
	}

	err = gitHubUtils.unzip(fileName, templateDestination)
	if err != nil {
		return err
	}

	return nil
}

func (gitHubUtils *GitHubUtils) clean(fileName string) error {
	return gitHubUtils.os.Remove(fileName)
}

func (gitHubUtils *GitHubUtils) downloadZip(projectFolder string, githubOrg string, templateName string) (string, error) {
	fileName := fmt.Sprintf("%s/%s.zip", projectFolder, templateName)

	out, err := gitHubUtils.os.Create(fileName)
	defer out.Close()
	if err != nil {
		return fileName, err
	}

	resp, err := gitHubUtils.http.Get(fmt.Sprintf("https://github.com/%s/%s/archive/master.zip", githubOrg, templateName))
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return fileName, err
	}

	if resp.StatusCode == 404 {
		return fileName, errors.New(fmt.Sprintf("Project template %s/%s could not be found", githubOrg, templateName))
	}

	_, err = gitHubUtils.io.Copy(out, resp.Body)
	if err != nil {
		return fileName, err
	}

	return fileName, nil
}

func (gitHubUtils *GitHubUtils) unzip(archive string, target string) error {
	fmt.Println(fmt.Sprintf("Extracting template from %s", archive))

	reader, err := gitHubUtils.zip.OpenReader(archive)
	if err != nil {
		return err
	}

	if err := gitHubUtils.os.MkdirAll(target, 0755); err != nil {
		return err
	}

	for _, file := range reader.File {
		err = gitHubUtils.unzipFile(file, target)
		if err != nil {
			return err
		}
	}

	return nil
}

func (gitHubUtils *GitHubUtils) unzipFile(file *zip.File, target string) error {
	fileName := file.Name[strings.Index(file.Name, "/")+1: len(file.Name)]

	path := filepath.Join(target, fileName)
	if file.FileInfo().IsDir() {
		gitHubUtils.os.MkdirAll(path, file.Mode())
		return nil
	}

	fileReader, err := file.Open()
	if err != nil {
		return err
	}
	defer fileReader.Close()

	targetFile, err := gitHubUtils.os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
	if targetFile != nil {
		defer targetFile.Close()
	}
	if err != nil {
		return err
	}

	if _, err := gitHubUtils.io.Copy(targetFile, fileReader); err != nil {
		return err
	}

	return nil
}

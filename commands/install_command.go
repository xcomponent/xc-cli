package commands

import (
	"fmt"
	"errors"
	"bytes"
	"strings"
	"github.com/xcomponent/xc-cli/services"
)

type InstallConfig struct {
	XcStudioDistribs map[string]string `json:"xcStudioDistribs"`
}

func NewInstallCommand(os services.OsService, io services.IoService, http services.HttpService,
	exec services.ExecService) *InstallCommand {
	return &InstallCommand{os, io, http, exec}
}

type InstallCommand struct {
	os   services.OsService
	io   services.IoService
	http services.HttpService
	exec services.ExecService
}

func (installCommand *InstallCommand) Install(installConfigUrl string, osName string, osArch string) error {

	err := installCommand.checkOs(osName)
	if err != nil {
		return err
	}

	err = installCommand.checkDotNet()
	if err != nil {
		return err
	}

	installConfig, err := installCommand.loadInstallConfig(installConfigUrl)
	if err != nil {
		return err
	}

	distribUrl := installConfig.XcStudioDistribs[osArch]
	if distribUrl == "" {
		return fmt.Errorf("xc install does not support %s arch", osArch)
	}

	msiPath, err := installCommand.downloadMsi(distribUrl)
	if msiPath != "" {
		defer fmt.Printf("Cleaning temporary directory %s.\n", msiPath)
		defer installCommand.os.RemoveAll(msiPath)
	}
	if err != nil {
		return err
	}

	err = installCommand.installXcStudio(msiPath)
	if err != nil {
		return err
	}

	fmt.Print("XC Studio installation completed.\n")

	return nil
}

func (installCommand *InstallCommand) checkDotNet() error {
	fmt.Print("Checking .Net version.\n")

	var outbuf bytes.Buffer

	cmd := installCommand.exec.Command("powershell", "-command", "\"Get-ChildItem 'hklm:SOFTWARE\\Microsoft\\NET Framework Setup\\NDP\\v4\\Full\\' | Get-ItemPropertyValue -Name Release | % { $_ -ge 394802 }\"")

	cmd.Stdout = &outbuf

	fmt.Printf("Executing command : %s\n", cmd.Args)

	if err := cmd.Run(); err != nil {
		return err
	}
	stdout := outbuf.String()

	if strings.HasPrefix(stdout, "False") {
		return errors.New("XComponent requires a version of DotNet higher than 4.5")
	}

	// powershell -command "Get-ChildItem 'hklm:SOFTWARE\Microsoft\NET Framework Setup\NDP\v4\Full\' | Get-ItemPropertyValue -Name Release | % { $_ -ge 394802 }"
	return nil
}

func (installCommand *InstallCommand) checkOs(osName string) error {
	if osName != "windows" {
		return errors.New("Install command is only supported on Windows")
	}

	return nil
}

func (installCommand *InstallCommand) downloadMsi(msiUrl string) (msiPath string, err error) {
	dir, err := installCommand.io.TempDir("", "xc-install")

	if err != nil {
		return "", err
	}

	path := fmt.Sprintf("%s\\xc-studio.msi", dir)

	fmt.Printf("Downloading %s to %s.\n", msiUrl, path)

	out, err := installCommand.os.Create(path)
	defer out.Close()
	if err != nil {
		return dir, err
	}

	resp, err := installCommand.http.Get(msiUrl)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return dir, err
	}

	_, err = installCommand.io.Copy(out, resp.Body)
	if err != nil {
		return dir, err
	}

	fmt.Printf("Download completed.\n")

	return path, nil
}

func (installCommand *InstallCommand) installXcStudio(msiPath string) error {
	cmd := installCommand.exec.Command("msiexec", "/i", msiPath, "/passive")

	fmt.Printf("Executing command : %s\n", cmd.Args)

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (installCommand *InstallCommand) loadInstallConfig(installConfigUrl string) (*InstallConfig, error) {
	installConfig := InstallConfig{}

	err := installCommand.http.GetJson(installConfigUrl, &installConfig)
	if err != nil {
		return nil, err
	}

	return &installConfig, err
}

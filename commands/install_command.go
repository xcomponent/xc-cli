package commands

import (
	"fmt"
	"runtime"
	"errors"
	"io/ioutil"
	"os"
	"net/http"
	"io"
	"encoding/json"
	"os/exec"
	"bytes"
	"strings"
)

type InstallConfig struct {
	XcStudioDistribs map[string]string `json:"xcStudioDistribs"`
}

func Install(installConfigUrl string, mockWindows bool, keepTempFiles bool) error {
	if mockWindows {
		err := checkOs()
		if err != nil {
			return err
		}

		err = checkDotNet()
		if err != nil {
			return err
		}
	}

	installConfig, err := loadInstallConfig(installConfigUrl)
	if err != nil {
		return err
	}

	distribUrl := installConfig.XcStudioDistribs[runtime.GOARCH]
	if distribUrl == "" {
		return errors.New(fmt.Sprintf("xc install does not support %s arch", runtime.GOARCH))
	}

	msiPath, err := downloadMsi(distribUrl)
	if !keepTempFiles && msiPath != "" {
		defer fmt.Printf("Cleaning temporary directory %s.\n", msiPath)
		defer os.RemoveAll(msiPath)
	}
	if err != nil {
		return err
	}

	err = installXcStudio(msiPath)
	if err != nil {
		return err
	}

	fmt.Print("XC Studio installation completed.\n")

	return nil
}

func checkDotNet() error {
	fmt.Print("Checking .Net version.\n")

	var outbuf bytes.Buffer

	cmd := exec.Command("powershell", "-command", "\"Get-ChildItem 'hklm:SOFTWARE\\Microsoft\\NET Framework Setup\\NDP\\v4\\Full\\' | Get-ItemPropertyValue -Name Release | % { $_ -ge 394802 }\"")

	cmd.Stdout = &outbuf

	fmt.Printf("Executing command : %s\n", cmd.Args)

	if err := cmd.Run(); err != nil {
		return err
	}
	stdout := outbuf.String()

	if strings.HasPrefix(stdout, "True") {
		return errors.New("XComponent requires a version of DotNet higher than 4.5")
	}

	// powershell -command "Get-ChildItem 'hklm:SOFTWARE\Microsoft\NET Framework Setup\NDP\v4\Full\' | Get-ItemPropertyValue -Name Release | % { $_ -ge 394802 }"
	return nil
}

func checkOs() error {
	if runtime.GOOS != "windows" {
		return errors.New("Install command is only supported on Windows")
	}

	return nil
}

func downloadMsi(msiUrl string) (msiPath string, err error) {
	dir, err := ioutil.TempDir("", "xc-install")

	if err != nil {
		return "", err
	}

	path := fmt.Sprintf("%s\\xc-studio.msi", dir)

	fmt.Printf("Downloading %s to %s.\n", msiUrl, path)

	out, err := os.Create(path)
	defer out.Close()
	if err != nil {
		return dir, err
	}

	resp, err := http.Get(msiUrl)
	defer resp.Body.Close()
	if err != nil {
		return dir, err
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return dir, err
	}

	fmt.Printf("Download completed.\n")

	return path, nil
}

func getJson(url string, target interface{}) error {
	client := &http.Client{}
	response, err := client.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return json.NewDecoder(response.Body).Decode(target)
}

func installXcStudio(msiPath string) error {
	cmd := exec.Command("msiexec", "/i", msiPath, "/passive")

	fmt.Printf("Executing command : %s\n", cmd.Args)

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func loadInstallConfig(installConfigUrl string) (*InstallConfig, error) {
	installConfig := InstallConfig{}

	err := getJson(installConfigUrl, &installConfig)
	if err != nil {
		return nil, err
	}

	return &installConfig, err
}

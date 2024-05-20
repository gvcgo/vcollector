package conda

import (
	"fmt"
	"os"
	"strings"

	"github.com/gvcgo/goutils/pkgs/gutils"
)

/*
Test miniconda.

https://repo.anaconda.com/miniconda/Miniconda3-latest-Windows-x86_64.exe
https://repo.anaconda.com/miniconda/Miniconda3-latest-MacOSX-x86_64.sh
https://repo.anaconda.com/miniconda/Miniconda3-latest-MacOSX-arm64.sh
https://repo.anaconda.com/miniconda/Miniconda3-latest-Linux-x86_64.sh
https://repo.anaconda.com/miniconda/Miniconda3-latest-Linux-aarch64.sh
*/
func IsCondaInstalled() bool {
	homeDir, _ := os.UserHomeDir()
	_, err := gutils.ExecuteSysCommand(true, homeDir, "conda", "--help")
	return err == nil
}

/*
subdirs:
https://conda.anaconda.org/conda-forge/
*/
var CondaPlatformList = []string{
	"linux-64",
	"linux-aarch64",
	"win-64",
	"win-arm64",
	"osx-64",
	"osx-arm64",
}

func ParseArch(platform string) (archStr string) {
	if strings.Contains(platform, "-64") {
		return "amd64"
	}
	if strings.Contains(platform, "-arm64") {
		return "arm64"
	}
	if strings.Contains(platform, "-aarch64") {
		return "arm64"
	}
	return
}

func ParseOS(platform string) (osStr string) {
	if strings.Contains(platform, "linux-") {
		return "linux"
	}
	if strings.Contains(platform, "win-") {
		return "windows"
	}
	if strings.Contains(platform, "osx-") {
		return "darwin"
	}
	return
}

/*
conda search --override-channels --channel conda-forge --skip-flexible-search --subdir osx-64 --full-name php
*/
var CondaSearchCommand = []string{
	"conda",
	"search",
	"--override-channels",
	"--channel",
	"conda-forge",
	"--skip-flexible-search",
}

func GetVersionForPlatform(platform, sdkName string) (vlist []string) {
	homeDir, _ := os.UserHomeDir()
	_cmd := append([]string{}, CondaSearchCommand...)
	_cmd = append(_cmd, "--subdir", platform, "--full-name", sdkName)
	r, err := gutils.ExecuteSysCommand(true, homeDir, _cmd...)
	if err == nil {
		vlist = ParseSearchResult(r.String())
	}
	return
}

func ParseSearchResult(content string) (vlist []string) {
	header := FindHeader(content)
	if header == "" {
		return
	}
	sList := strings.Split(content, header)
	if len(sList) == 2 {
		lines := strings.Split(sList[1], "\n")
		for _, line := range lines {
			vlist = append(vlist, FindVersion(strings.Split(line, " ")))
		}
	}
	return
}

func FindHeader(content string) (header string) {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "# Name") {
			return line
		}
	}
	return
}

func FindVersion(llist []string) string {
	newList := []string{}
	for _, item := range llist {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		newList = append(newList, item)
	}
	if len(newList) > 1 {
		return newList[1]
	}
	return ""
}

func SearchVersions(sdkName string) (result map[string][]string) {
	if !IsCondaInstalled() {
		fmt.Println("conda is not installed.")
		return
	}
	result = make(map[string][]string)
	for _, platform := range CondaPlatformList {
		vlist := GetVersionForPlatform(platform, sdkName)
		key := fmt.Sprintf("%s/%s", ParseOS(platform), ParseArch(platform))
		result[key] = vlist
	}
	return
}

func TestCondaSearch() {
	r := SearchVersions("php")
	fmt.Printf("%+v\n", r)
}

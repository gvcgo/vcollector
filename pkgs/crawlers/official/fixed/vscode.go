package fixed

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/gvcgo/vcollector/internal/iconf"
	"github.com/gvcgo/vcollector/internal/req"
	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
	"github.com/gvcgo/vcollector/pkgs/version"
)

func init() {
	crawler.RegisterCrawler(NewVSCode())
}

type CodePlatform struct {
	Os         string `json:"os"`
	PrettyName string `json:"prettyname"`
}

type CodeItem struct {
	Url      string       `josn:"url"`
	Sum      string       `json:"sha256hash"`
	Version  string       `json:"name"`
	Build    string       `json:"build"`
	Platform CodePlatform `json:"platform"`
}

type CodeProducts struct {
	Products []CodeItem `json:"products"`
}

/*
https://code.visualstudio.com/sha?build=stable
*/
type VSCode struct {
	DownloadUrl string
	SDKName     string
	Version     version.VersionList
	result      CodeProducts
}

func NewVSCode() (v *VSCode) {
	v = &VSCode{
		DownloadUrl: "https://code.visualstudio.com/sha?build=stable",
		SDKName:     "vscode",
		Version:     make(version.VersionList),
		result:      CodeProducts{Products: []CodeItem{}},
	}
	return
}

func (v *VSCode) GetSDKName() string {
	return v.SDKName
}

func (v *VSCode) filter(dUrl string) (r bool) {
	excludeList := []string{"_cli", ".tar.gz", "armhf", "armv7hl"}
	for _, excludeStr := range excludeList {
		if strings.Contains(dUrl, excludeStr) {
			return false
		}
	}
	if strings.HasSuffix(dUrl, ".exe") && !strings.Contains(dUrl, "User") {
		return true
	}
	if strings.HasSuffix(dUrl, ".tar.gz") {
		return true
	}
	if strings.HasSuffix(dUrl, ".zip") && strings.Contains(dUrl, "darwin") {
		return true
	}
	if strings.HasSuffix(dUrl, ".deb") || strings.HasSuffix(dUrl, ".rpm") {
		return true
	}
	return
}

func (v *VSCode) parseArch(vItem CodeItem) (archStr string) {
	ss := vItem.Platform.Os
	if strings.Contains(ss, "-x64") {
		return "amd64"
	}
	if strings.Contains(ss, "-arm64") {
		return "arm64"
	}
	if ss == "darwin" {
		return "amd64"
	}
	if strings.Contains(ss, "-universal") {
		return "any"
	}
	return
}

func (v *VSCode) parseOs(vItem CodeItem) (osStr string) {
	ss := vItem.Platform.Os
	if strings.Contains(ss, "win32") {
		return "windows"
	}
	if strings.Contains(ss, "darwin") {
		return "darwin"
	}
	if strings.Contains(ss, "linux") {
		return "linux"
	}
	return
}

func (v *VSCode) parseInstaller(vItem CodeItem) (installerStr string) {
	if strings.HasSuffix(vItem.Url, ".exe") {
		return version.Executable
	}
	if strings.HasSuffix(vItem.Url, ".tar.gz") || strings.HasSuffix(vItem.Url, ".zip") {
		return version.Unarchiver
	}
	if strings.HasSuffix(vItem.Url, ".deb") {
		return version.Dpkg
	}
	if strings.HasSuffix(vItem.Url, ".rpm") {
		return version.Rpm
	}
	return
}

func (v *VSCode) getResult() {
	req.GetJson(v.DownloadUrl, &v.result)
	for _, vItem := range v.result.Products {
		if v.filter(vItem.Url) {
			ver := version.Item{}
			ver.Url = vItem.Url
			ver.Arch = v.parseArch(vItem)
			ver.Os = v.parseOs(vItem)
			if ver.Arch == "" || ver.Os == "" {
				continue
			}
			ver.Sum = vItem.Sum
			if ver.Sum != "" {
				ver.SumType = "sha256"
			}
			ver.Installer = v.parseInstaller(vItem)
			vStr := vItem.Version
			if _, ok := v.Version[vStr]; !ok {
				v.Version[vStr] = version.Version{}
			}
			v.Version[vStr] = append(v.Version[vStr], ver)
		}
	}
}

func (v *VSCode) Start() {
	v.getResult()
}

func (v *VSCode) GetVersions() []byte {
	r, _ := v.Version.Marshal()
	return r
}

func (v *VSCode) HomePage() string {
	return "https://code.visualstudio.com/"
}

func (v *VSCode) GetInstallConf() (ic iconf.InstallerConfig) {
	return iconf.InstallerConfig{
		BinaryDirs: &iconf.DirItems{
			Windows: []iconf.DirPath{},
			MacOS:   []iconf.DirPath{},
			Linux:   []iconf.DirPath{},
		},
	}
}

func TestVSCode() {
	vv := NewVSCode()
	vv.Start()

	ff := "/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/vscode.json"
	content, _ := json.MarshalIndent(vv.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

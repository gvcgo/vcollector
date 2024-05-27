package mix

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gvcgo/vcollector/internal/conda"
	"github.com/gvcgo/vcollector/internal/gh"
	"github.com/gvcgo/vcollector/internal/iconf"
	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
	"github.com/gvcgo/vcollector/pkgs/version"
)

func init() {
	crawler.RegisterCrawler(NewErlang())
}

/*
For windows:
https://github.com/erlang/otp/releases
*/
type Erlang struct {
	SDKName  string
	RepoName string
	Version  version.VersionList
}

func NewErlang() (e *Erlang) {
	e = &Erlang{
		SDKName:  "erlang",
		RepoName: "erlang/otp",
		Version:  make(version.VersionList),
	}
	return
}

func (e *Erlang) GetSDKName() string {
	return e.SDKName
}

func (e *Erlang) GetVersions() []byte {
	r, _ := e.Version.Marshal()
	return r
}

func (e *Erlang) getFromConda() {
	vv := conda.SearchVersions(e.SDKName)
	for osArch, vList := range vv {
		sList := strings.Split(osArch, "/")
		if len(sList) == 2 {
			osStr, archStr := sList[0], sList[1]
			for _, vStr := range vList {
				item := version.Item{}
				item.Os = osStr
				item.Arch = archStr
				if _, ok := e.Version[vStr]; !ok {
					e.Version[vStr] = version.Version{}
				}
				item.Installer = version.Conda
				e.Version[vStr] = append(e.Version[vStr], item)
			}
		}
	}
}

func (e *Erlang) trimVersion(vStr string) string {
	vStr = strings.TrimPrefix(vStr, "v")
	return strings.TrimPrefix(vStr, "OTP-")
}

func (e *Erlang) getFromGithub() {
	eItemList := gh.GetReleaseItems(e.RepoName)
	for _, eItem := range eItemList {
		vStr := e.trimVersion(eItem.TagName)
		for _, aa := range eItem.Assets {
			if !strings.HasSuffix(aa.Name, ".exe") {
				continue
			}
			if strings.Contains(aa.Name, "win64") {
				item := version.Item{}
				item.Os = "windows"
				item.Arch = "amd64"
				item.Installer = version.Executable
				item.Url = aa.Url
				item.Size = aa.Size
				if _, ok := e.Version[vStr]; !ok {
					e.Version[vStr] = version.Version{}
				}
				e.Version[vStr] = append(e.Version[vStr], item)
			}
		}
	}
}

func (e *Erlang) Start() {
	e.getFromConda()
	e.getFromGithub()
}

func (e *Erlang) HomePage() string {
	return "https://www.erlang.org/"
}

func (e *Erlang) GetInstallConf() (ic iconf.InstallerConfig) {
	return iconf.InstallerConfig{
		FlagFiles: &iconf.FileItems{
			Windows: []string{"bin", "lib"},
			MacOS:   []string{"bin", "lib"},
			Linux:   []string{"bin", "lib"},
		},
		BinaryDirs: &iconf.DirItems{
			Windows: []iconf.DirPath{{"bin"}},
			MacOS:   []iconf.DirPath{{"bin"}},
			Linux:   []iconf.DirPath{{"bin"}},
		},
	}
}

func TestErlang() {
	nn := NewErlang()
	nn.Start()

	ff := fmt.Sprintf(
		"/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/%s.json",
		nn.SDKName,
	)
	content, _ := json.MarshalIndent(nn.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

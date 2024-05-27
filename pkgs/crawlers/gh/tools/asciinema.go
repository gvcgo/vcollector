package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gvcgo/vcollector/internal/gh"
	"github.com/gvcgo/vcollector/internal/iconf"
	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
	"github.com/gvcgo/vcollector/pkgs/crawlers/gh/searcher"
	"github.com/gvcgo/vcollector/pkgs/version"
)

func init() {
	crawler.RegisterCrawler(NewAsciinema())
}

type Asciinema struct {
	SDKName  string
	RepoName string
	searcher.GhSearcher
}

func NewAsciinema() (a *Asciinema) {
	a = &Asciinema{
		SDKName:  "acast",
		RepoName: "gvcgo/asciinema",
		GhSearcher: searcher.GhSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (a *Asciinema) GetSDKName() string {
	return a.SDKName
}

func (a *Asciinema) tagFilter(ri gh.ReleaseItem) bool {
	return searcher.GhVersionRegexp.FindString(ri.TagName) != ""
}

func (a *Asciinema) fileFilter(aa gh.Asset) bool {
	return !strings.Contains(aa.Url, "archive/refs/")
}

func (a *Asciinema) osParser(fName string) (osStr string) {
	if strings.Contains(fName, "darwin") {
		return "darwin"
	}
	if strings.Contains(fName, "linux") {
		return "linux"
	}
	if strings.Contains(fName, "windows") {
		return "windows"
	}
	return
}

func (a *Asciinema) archParser(fName string) (archStr string) {
	if strings.Contains(fName, "amd64") {
		return "amd64"
	}
	if strings.Contains(fName, "arm64") {
		return "arm64"
	}
	return
}

func (a *Asciinema) vParser(tagName string) (vStr string) {
	return strings.TrimPrefix(tagName, "v")
}

func (a *Asciinema) insParser(fName string) (insStr string) {
	return version.Unarchiver
}

func (a *Asciinema) Start() {
	a.GhSearcher.Search(
		a.RepoName,
		a.tagFilter,
		a.fileFilter,
		a.vParser,
		a.archParser,
		a.osParser,
		a.insParser,
		nil,
	)
}

func (a *Asciinema) GetVersions() []byte {
	r, _ := a.Version.Marshal()
	return r
}

func (a *Asciinema) HomePage() string {
	return "https://asciinema.org/"
}

func (a *Asciinema) GetInstallConf() (ic iconf.InstallerConfig) {
	return iconf.InstallerConfig{
		FlagFiles: &iconf.FileItems{
			Windows: []string{"acast.exe"},
			MacOS:   []string{"acast"},
			Linux:   []string{"acast"},
		},
		FlagDirExcepted: true,
		BinaryDirs: &iconf.DirItems{
			Windows: []iconf.DirPath{},
			MacOS:   []iconf.DirPath{},
			Linux:   []iconf.DirPath{},
		},
	}
}

func TestAsciinema() {
	nn := NewAsciinema()
	nn.Start()

	ff := fmt.Sprintf(
		"/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/%s.json",
		nn.SDKName,
	)
	content, _ := json.MarshalIndent(nn.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

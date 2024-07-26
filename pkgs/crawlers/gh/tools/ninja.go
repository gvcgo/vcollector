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
	crawler.RegisterCrawler(NewNinja())
}

type Ninja struct {
	SDKName  string
	RepoName string
	searcher.GhSearcher
}

func NewNinja() (n *Ninja) {
	n = &Ninja{
		SDKName:  "ninja",
		RepoName: "ninja-build/ninja",
		GhSearcher: searcher.GhSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (n *Ninja) GetSDKName() string {
	return n.SDKName
}

func (n *Ninja) tagFilter(ri gh.ReleaseItem) bool {
	return searcher.GhVersionRegexp.FindString(ri.TagName) != ""
}

func (n *Ninja) fileFilter(a gh.Asset) bool {
	return !strings.Contains(a.Url, "archive/refs/")
}

func (n *Ninja) osParser(fName string) (osStr string) {
	if strings.Contains(fName, "mac") {
		return "darwin"
	}
	if strings.Contains(fName, "linux") {
		return "linux"
	}
	if strings.Contains(fName, "-win") {
		return "windows"
	}
	return
}

func (n *Ninja) archParser(fName string) (archStr string) {
	if strings.Contains(fName, "arm64") || strings.Contains(fName, "aarch64") {
		return "arm64"
	}
	if strings.Contains(fName, "mac") {
		return "any"
	}
	return "amd64"
}

func (n *Ninja) vParser(tagName string) (vStr string) {
	return strings.TrimPrefix(tagName, "v")
}

func (n *Ninja) insParser(fName string) (insStr string) {
	return version.Unarchiver
}

func (n *Ninja) Start() {
	n.GhSearcher.Search(
		n.RepoName,
		n.tagFilter,
		n.fileFilter,
		n.vParser,
		n.archParser,
		n.osParser,
		n.insParser,
		nil,
	)
}

func (n *Ninja) GetVersions() []byte {
	r, _ := n.Version.Marshal()
	return r
}

func (n *Ninja) HomePage() string {
	return "https://ninja-build.org/"
}

func (n *Ninja) GetInstallConf() (ic iconf.InstallerConfig) {
	return iconf.InstallerConfig{
		FlagFiles: &iconf.FileItems{
			Windows: []string{"ninja.exe"},
			MacOS:   []string{"ninja"},
			Linux:   []string{"ninja"},
		},
		FlagDirExcepted: false,
		BinaryDirs: &iconf.DirItems{
			Windows: []iconf.DirPath{},
			MacOS:   []iconf.DirPath{},
			Linux:   []iconf.DirPath{},
		},
	}
}

func TestNinja() {
	nn := NewNinja()
	nn.Start()

	ff := fmt.Sprintf(
		"/home/moqsien/projects/go/src/gvcgo/vcollector/test/%s.json",
		nn.SDKName,
	)
	content, _ := json.MarshalIndent(nn.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

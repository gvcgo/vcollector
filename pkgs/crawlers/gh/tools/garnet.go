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
	crawler.RegisterCrawler(NewGarnet())
}

type Garnet struct {
	SDKName  string
	RepoName string
	searcher.GhSearcher
}

func NewGarnet() (g *Garnet) {
	g = &Garnet{
		SDKName:  "garnet",
		RepoName: "microsoft/garnet",
		GhSearcher: searcher.GhSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (g *Garnet) GetSDKName() string {
	return g.SDKName
}

func (g *Garnet) tagFilter(ri gh.ReleaseItem) bool {
	return searcher.GhVersionRegexp.FindString(ri.TagName) != ""
}

func (g *Garnet) fileFilter(a gh.Asset) bool {
	if strings.Contains(a.Url, "archive/refs/") {
		return false
	}
	if strings.HasSuffix(a.Name, ".nupkg") {
		return false
	}
	if strings.HasSuffix(a.Name, "portable.7z") {
		return false
	}
	return true
}

func (g *Garnet) osParser(fName string) (osStr string) {
	if strings.Contains(fName, "osx") {
		return "darwin"
	}
	if strings.Contains(fName, "linux") {
		return "linux"
	}
	if strings.Contains(fName, "win-") {
		return "windows"
	}
	return
}

func (g *Garnet) archParser(fName string) (archStr string) {
	if strings.Contains(fName, "x64") {
		return "amd64"
	}
	if strings.Contains(fName, "arm64") {
		return "arm64"
	}
	return
}

func (g *Garnet) vParser(tagName string) (vStr string) {
	return strings.TrimPrefix(tagName, "v")
}

func (g *Garnet) insParser(fName string) (insStr string) {
	return version.Unarchiver
}

func (g *Garnet) Start() {
	g.GhSearcher.Search(
		g.RepoName,
		g.tagFilter,
		g.fileFilter,
		g.vParser,
		g.archParser,
		g.osParser,
		g.insParser,
		nil,
	)
}

func (g *Garnet) GetVersions() []byte {
	rr, _ := g.Version.Marshal()
	return rr
}

func (g *Garnet) HomePage() string {
	return "https://github.com/microsoft/garnet"
}

func (g *Garnet) GetInstallConf() (ic iconf.InstallerConfig) {
	return iconf.InstallerConfig{
		FlagFiles: &iconf.FileItems{
			Windows: []string{"GarnetServer.exe"},
			MacOS:   []string{"GarnetServer"},
			Linux:   []string{"GarnetServer"},
		},
		FlagDirExcepted: true,
		BinaryDirs: &iconf.DirItems{
			Windows: []iconf.DirPath{},
			MacOS:   []iconf.DirPath{},
			Linux:   []iconf.DirPath{},
		},
	}
}

func TestGarnet() {
	nn := NewGarnet()
	nn.Start()

	ff := fmt.Sprintf(
		"/home/moqsien/projects/go/src/gvcgo/vcollector/test/%s.json",
		nn.SDKName,
	)
	content, _ := json.MarshalIndent(nn.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

package lans

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
	crawler.RegisterCrawler(NewGleam())
}

type Gleam struct {
	SDKName  string
	RepoName string
	searcher.GhSearcher
}

func NewGleam() (g *Gleam) {
	g = &Gleam{
		SDKName:  "gleam",
		RepoName: "gleam-lang/gleam",
		GhSearcher: searcher.GhSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (g *Gleam) GetSDKName() string {
	return g.SDKName
}

func (g *Gleam) tagFilter(ri gh.ReleaseItem) bool {
	return searcher.GhVersionRegexp.FindString(ri.TagName) != ""
}

func (g *Gleam) fileFilter(a gh.Asset) bool {
	if strings.Contains(a.Url, "archive/refs/") {
		return false
	}
	if strings.Contains(a.Name, "-browser") {
		return false
	}
	if strings.HasSuffix(a.Name, ".sha256") {
		return false
	}
	if strings.HasSuffix(a.Name, ".sha512") {
		return false
	}
	return true
}

func (g *Gleam) osParser(fName string) (osStr string) {
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

func (g *Gleam) archParser(fName string) (archStr string) {
	if strings.Contains(fName, "-x86_64") {
		return "amd64"
	}
	if strings.Contains(fName, "-aarch64") {
		return "arm64"
	}
	return
}

func (g *Gleam) vParser(tagName string) (vStr string) {
	return strings.TrimPrefix(tagName, "v")
}

func (g *Gleam) insParser(fName string) (insStr string) {
	return version.Unarchiver
}

func (g *Gleam) Start() {
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

func (g *Gleam) GetVersions() []byte {
	r, _ := g.Version.Marshal()
	return r
}

func (g *Gleam) HomePage() string {
	return "https://gleam.run/"
}

func (g *Gleam) GetInstallConf() (ic iconf.InstallerConfig) {
	return iconf.InstallerConfig{
		FlagFiles: &iconf.FileItems{
			Windows: []string{"gleam.exe"},
			MacOS:   []string{"gleam"},
			Linux:   []string{"gleam"},
		},
		FlagDirExcepted: true,
		BinaryDirs: &iconf.DirItems{
			Windows: []iconf.DirPath{},
			MacOS:   []iconf.DirPath{},
			Linux:   []iconf.DirPath{},
		},
	}
}

func TestGleam() {
	nn := NewGleam()
	nn.Start()

	ff := fmt.Sprintf(
		"/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/%s.json",
		nn.SDKName,
	)
	content, _ := json.MarshalIndent(nn.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

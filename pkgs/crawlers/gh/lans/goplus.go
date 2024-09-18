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
	crawler.RegisterCrawler(NewGoPlus())
}

type GoPlus struct {
	SDKName  string
	RepoName string
	searcher.GhSearcher
}

func NewGoPlus() (g *GoPlus) {
	g = &GoPlus{
		SDKName:  "goplus",
		RepoName: "goplus/gop",
		GhSearcher: searcher.GhSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (g *GoPlus) GetSDKName() string {
	return g.SDKName
}

func (g *GoPlus) tagFilter(ri gh.ReleaseItem) bool {
	return searcher.GhVersionRegexp.FindString(ri.TagName) != ""
}

func (g *GoPlus) fileFilter(a gh.Asset) bool {
	if strings.Contains(a.Url, "archive/refs/") {
		return false
	}
	if strings.HasSuffix(a.Name, ".txt") {
		return false
	}
	if strings.HasSuffix(a.Name, ".deb") {
		return false
	}
	if strings.HasSuffix(a.Name, ".rpm") {
		return false
	}
	return true
}

func (g *GoPlus) osParser(fName string) (osStr string) {
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

func (g *GoPlus) archParser(fName string) (archStr string) {
	if strings.Contains(fName, "-amd64") {
		return "amd64"
	}
	if strings.Contains(fName, "-arm64") {
		return "arm64"
	}
	return
}

func (g *GoPlus) vParser(tagName string) (vStr string) {
	return strings.TrimPrefix(tagName, "v")
}

func (g *GoPlus) insParser(fName string) (insStr string) {
	return version.Unarchiver
}

func (g *GoPlus) Start() {
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

func (g *GoPlus) GetVersions() []byte {
	r, _ := g.Version.Marshal()
	return r
}

func (g *GoPlus) HomePage() string {
	return "https://goplus.org/"
}

func (g *GoPlus) GetInstallConf() (ic iconf.InstallerConfig) {
	return iconf.InstallerConfig{
		FlagFiles: &iconf.FileItems{
			Windows: []string{"bin"},
			MacOS:   []string{"bin"},
			Linux:   []string{"bin"},
		},
		FlagDirExcepted: false,
		BinaryDirs: &iconf.DirItems{
			Windows: []iconf.DirPath{{"bin"}},
			MacOS:   []iconf.DirPath{{"bin"}},
			Linux:   []iconf.DirPath{{"bin"}},
		},
	}
}

func TestGoPlus() {
	bb := NewGoPlus()
	bb.Start()

	ff := fmt.Sprintf(
		"/home/moqsien/projects/go/src/gvcgo/vcollector/test/%s.json",
		bb.SDKName,
	)
	content, _ := json.MarshalIndent(bb.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

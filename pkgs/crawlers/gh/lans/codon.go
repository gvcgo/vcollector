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
	crawler.RegisterCrawler(NewCodon())
}

type Codon struct {
	SDKName  string
	RepoName string
	searcher.GhSearcher
}

func NewCodon() (c *Codon) {
	c = &Codon{
		SDKName:  "codon",
		RepoName: "exaloop/codon",
		GhSearcher: searcher.GhSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (c *Codon) GetSDKName() string {
	return c.SDKName
}

func (c *Codon) tagFilter(ri gh.ReleaseItem) bool {
	return searcher.GhVersionRegexp.FindString(ri.TagName) != ""
}

func (c *Codon) fileFilter(a gh.Asset) bool {
	return !strings.Contains(a.Url, "archive/refs/")
}

func (c *Codon) osParser(fName string) (osStr string) {
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

func (c *Codon) archParser(fName string) (archStr string) {
	if strings.Contains(fName, "-x86_64") {
		return "amd64"
	}
	if strings.Contains(fName, "arm64") {
		return "arm64"
	}
	return
}

func (c *Codon) vParser(tagName string) (vStr string) {
	return strings.TrimPrefix(searcher.GhVersionRegexp.FindString(tagName), "v")
}

func (c *Codon) insParser(fName string) (insStr string) {
	return version.Unarchiver
}

func (c *Codon) Start() {
	c.GhSearcher.Search(
		c.RepoName,
		c.tagFilter,
		c.fileFilter,
		c.vParser,
		c.archParser,
		c.osParser,
		c.insParser,
		nil,
	)
}

func (c *Codon) GetVersions() []byte {
	r, _ := c.Version.Marshal()
	return r
}

func (c *Codon) HomePage() string {
	return "https://github.com/exaloop/codon"
}

func (c *Codon) GetInstallConf() (ic iconf.InstallerConfig) {
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

func TestCodon() {
	cc := NewCodon()
	cc.Start()

	ff := fmt.Sprintf(
		"/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/%s.json",
		cc.SDKName,
	)
	content, _ := json.MarshalIndent(cc.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

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
	crawler.RegisterCrawler(NewCrystal())
}

type Crystal struct {
	SDKName  string
	RepoName string
	searcher.GhSearcher
}

func NewCrystal() (c *Crystal) {
	c = &Crystal{
		SDKName:  "crystal",
		RepoName: "crystal-lang/crystal",
		GhSearcher: searcher.GhSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (c *Crystal) GetSDKName() string {
	return c.SDKName
}

func (c *Crystal) tagFilter(ri gh.ReleaseItem) bool {
	return searcher.GVersionRegexp.FindString(ri.TagName) != ""
}

func (c *Crystal) fileFilter(a gh.Asset) bool {
	if strings.Contains(a.Url, "archive/refs/") {
		return false
	}
	if strings.HasSuffix(a.Name, ".pkg") {
		return false
	}
	if strings.HasSuffix(a.Name, "-linux-x86_64.tar.gz") {
		return false
	}
	if strings.HasSuffix(a.Name, ".exe") {
		return false
	}
	if strings.HasSuffix(a.Name, ".docs.tar.gz") {
		return false
	}
	return true
}

func (c *Crystal) osParser(fName string) (osStr string) {
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

func (c *Crystal) archParser(fName string) (archStr string) {
	if strings.Contains(fName, "x86_64") {
		return "amd64"
	}
	if strings.Contains(fName, "universal") {
		return "any"
	}
	if strings.Contains(fName, "aarch64") {
		return "arm64"
	}
	if strings.Contains(fName, "arm64") {
		return "arm64"
	}
	return
}

func (c *Crystal) vParser(tagName string) (vStr string) {
	return strings.TrimPrefix(searcher.GVersionRegexp.FindString(tagName), "v")
}

func (c *Crystal) insParser(fName string) (insStr string) {
	return version.Unarchiver
}

func (c *Crystal) Start() {
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

func (c *Crystal) GetVersions() []byte {
	r, _ := c.Version.Marshal()
	return r
}

func (c *Crystal) HomePage() string {
	return "https://crystal-lang.org/"
}

func (b *Crystal) GetInstallConf() (ic iconf.InstallerConfig) {
	return iconf.InstallerConfig{
		FlagFiles: &iconf.FileItems{
			Windows: []string{"lib", "crystal.exe"},
			MacOS:   []string{"bin", "src"},
			Linux:   []string{"bin", "lib"},
		},
		FlagDirExcepted: false,
		BinaryDirs: &iconf.DirItems{
			Windows: []iconf.DirPath{},
			MacOS:   []iconf.DirPath{{"bin"}},
			Linux:   []iconf.DirPath{{"bin"}},
		},
	}
}

func TestCrystal() {
	bb := NewCrystal()
	bb.Start()

	ff := fmt.Sprintf(
		"/home/moqsien/projects/go/src/gvcgo/vcollector/test/%s.json",
		bb.SDKName,
	)
	content, _ := json.MarshalIndent(bb.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

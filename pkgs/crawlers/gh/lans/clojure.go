package lans

import (
	"strings"

	"github.com/gvcgo/vcollector/internal/gh"
	"github.com/gvcgo/vcollector/internal/iconf"
	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
	"github.com/gvcgo/vcollector/pkgs/crawlers/gh/searcher"
	"github.com/gvcgo/vcollector/pkgs/version"
)

func init() {
	crawler.RegisterCrawler(NewClojure())
}

type Clojure struct {
	SDKName  string
	RepoName string
	searcher.GhSearcher
}

func NewClojure() (c *Clojure) {
	c = &Clojure{
		SDKName:  "clojure",
		RepoName: "clojure/brew-install",
		GhSearcher: searcher.GhSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (c *Clojure) GetSDKName() string {
	return c.SDKName
}

func (c *Clojure) GetVersions() []byte {
	r, _ := c.Version.Marshal()
	return r
}

func (c *Clojure) HomePage() string {
	return "https://clojure.org/"
}

func (c *Clojure) tagFilter(ri gh.ReleaseItem) bool {
	return searcher.GVersionRegexp.FindString(ri.TagName) != ""
}

func (c *Clojure) fileFilter(a gh.Asset) bool {
	return strings.HasPrefix(a.Name, "clojure") && strings.HasSuffix(a.Name, ".tar.gz")
}

func (c *Clojure) osParser(fName string) (osStr string) {
	return "any"
}

func (c *Clojure) archParser(fName string) (archStr string) {
	return "any"
}

func (c *Clojure) vParser(tagName string) (vStr string) {
	return strings.TrimPrefix(tagName, "v")
}

func (c *Clojure) insParser(fName string) (insStr string) {
	return version.Unarchiver
}

func (c *Clojure) Start() {
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

func (c *Clojure) GetInstallConf() (ic iconf.InstallerConfig) {
	return iconf.InstallerConfig{
		FlagFiles: &iconf.FileItems{
			Windows: []string{"clojure", "clj"},
			MacOS:   []string{"clojure", "clj"},
			Linux:   []string{"clojure", "clj"},
		},
		FlagDirExcepted: true,
		BinaryDirs: &iconf.DirItems{
			Windows: []iconf.DirPath{{"bin"}},
			MacOS:   []iconf.DirPath{{"bin"}},
			Linux:   []iconf.DirPath{{"bin"}},
		},
	}
}

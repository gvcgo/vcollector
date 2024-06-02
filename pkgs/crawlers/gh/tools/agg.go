package tools

import (
	"strings"

	"github.com/gvcgo/vcollector/internal/gh"
	"github.com/gvcgo/vcollector/internal/iconf"
	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
	"github.com/gvcgo/vcollector/pkgs/crawlers/gh/searcher"
	"github.com/gvcgo/vcollector/pkgs/version"
)

func init() {
	crawler.RegisterCrawler(NewAgg())
}

type Agg struct {
	SDKName  string
	RepoName string
	searcher.GhSearcher
}

func NewAgg() (a *Agg) {
	a = &Agg{
		SDKName:  "agg",
		RepoName: "asciinema/agg",
		GhSearcher: searcher.GhSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (a *Agg) GetSDKName() string {
	return a.SDKName
}

func (a *Agg) tagFilter(ri gh.ReleaseItem) bool {
	return searcher.GhVersionRegexp.FindString(ri.TagName) != ""
}

func (a *Agg) fileFilter(aa gh.Asset) bool {
	if strings.Contains(aa.Url, "archive/refs/") {
		return false
	}
	if strings.HasSuffix(aa.Name, "-linux-gnu") {
		return false
	}
	return true
}

func (a *Agg) osParser(fName string) (osStr string) {
	if strings.Contains(fName, "apple") {
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

func (a *Agg) archParser(fName string) (archStr string) {
	if strings.Contains(fName, "-x86_64") {
		return "amd64"
	}
	if strings.Contains(fName, "-aarch64") {
		return "arm64"
	}
	if strings.Contains(fName, "-arm64") {
		return "arm64"
	}
	return
}

func (a *Agg) vParser(tagName string) (vStr string) {
	return strings.TrimPrefix(tagName, "v")
}

func (a *Agg) insParser(fName string) (insStr string) {
	return version.Executable
}

func (a *Agg) Start() {
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

func (a *Agg) GetVersions() []byte {
	r, _ := a.Version.Marshal()
	return r
}

func (a *Agg) HomePage() string {
	return "https://github.com/asciinema/agg"
}

func (c *Agg) GetInstallConf() (ic iconf.InstallerConfig) {
	return iconf.InstallerConfig{
		FlagFiles: &iconf.FileItems{
			Windows: []string{"agg.exe"},
			MacOS:   []string{"agg"},
			Linux:   []string{"agg"},
		},
		FlagDirExcepted: true,
		BinaryDirs: &iconf.DirItems{
			Windows: []iconf.DirPath{},
			MacOS:   []iconf.DirPath{},
			Linux:   []iconf.DirPath{},
		},
		BinaryRename: &iconf.BinaryRename{
			NameFlag: "agg",
			RenameTo: "agg",
		},
	}
}

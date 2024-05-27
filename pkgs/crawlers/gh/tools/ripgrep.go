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
	crawler.RegisterCrawler(NewRipgrep())
}

type Ripgrep struct {
	SDKName  string
	RepoName string
	searcher.GhSearcher
}

func NewRipgrep() (g *Ripgrep) {
	g = &Ripgrep{
		SDKName:  "ripgrep",
		RepoName: "BurntSushi/ripgrep",
		GhSearcher: searcher.GhSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (r *Ripgrep) GetSDKName() string {
	return r.SDKName
}

func (r *Ripgrep) tagFilter(ri gh.ReleaseItem) bool {
	return searcher.GVersionRegexp.FindString(ri.TagName) != ""
}

func (r *Ripgrep) fileFilter(a gh.Asset) bool {
	if strings.Contains(a.Url, "archive/refs/") {
		return false
	}
	if strings.HasSuffix(a.Name, "-windows-msvc.zip") {
		return false
	}
	if strings.HasSuffix(a.Name, ".sha256") {
		return false
	}
	if strings.HasSuffix(a.Name, ".deb") {
		return false
	}
	return true
}

func (r *Ripgrep) osParser(fName string) (osStr string) {
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

func (r *Ripgrep) archParser(fName string) (archStr string) {
	if strings.Contains(fName, "-x86_64") {
		return "amd64"
	}
	if strings.Contains(fName, "-aarch64") {
		return "arm64"
	}
	return
}

func (r *Ripgrep) vParser(tagName string) (vStr string) {
	return strings.TrimPrefix(tagName, "v")
}

func (r *Ripgrep) insParser(fName string) (insStr string) {
	return version.Unarchiver
}

func (r *Ripgrep) Start() {
	r.GhSearcher.Search(
		r.RepoName,
		r.tagFilter,
		r.fileFilter,
		r.vParser,
		r.archParser,
		r.osParser,
		r.insParser,
		nil,
	)
}

func (r *Ripgrep) GetVersions() []byte {
	rr, _ := r.Version.Marshal()
	return rr
}

func (r *Ripgrep) HomePage() string {
	return "https://github.com/BurntSushi/ripgrep"
}

func (r *Ripgrep) GetInstallConf() (ic iconf.InstallerConfig) {
	return iconf.InstallerConfig{
		FlagFiles: &iconf.FileItems{
			Windows: []string{"rg.exe"},
			MacOS:   []string{"rg"},
			Linux:   []string{"rg"},
		},
		FlagDirExcepted: true,
		BinaryDirs: &iconf.DirItems{
			Windows: []iconf.DirPath{},
			MacOS:   []iconf.DirPath{},
			Linux:   []iconf.DirPath{},
		},
	}
}

func TestRipgrep() {
	nn := NewRipgrep()
	nn.Start()

	ff := fmt.Sprintf(
		"/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/%s.json",
		nn.SDKName,
	)
	content, _ := json.MarshalIndent(nn.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

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
	crawler.RegisterCrawler(NewRedka())
}

type Redka struct {
	SDKName  string
	RepoName string
	searcher.GhSearcher
}

func NewRedka() (r *Redka) {
	r = &Redka{
		SDKName:  "redka",
		RepoName: "nalgeon/redka",
		GhSearcher: searcher.GhSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (r *Redka) GetSDKName() string {
	return r.SDKName
}

func (r *Redka) tagFilter(ri gh.ReleaseItem) bool {
	return searcher.GhVersionRegexp.FindString(ri.TagName) != ""
}

func (r *Redka) fileFilter(a gh.Asset) bool {
	if strings.Contains(a.Url, "archive/refs/") {
		return false
	}
	if strings.HasSuffix(a.Name, ".txt") {
		return false
	}
	return true
}

func (r *Redka) osParser(fName string) (osStr string) {
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

func (r *Redka) archParser(fName string) (archStr string) {
	if strings.Contains(fName, "amd64") {
		return "amd64"
	}
	if strings.Contains(fName, "arm64") {
		return "arm64"
	}
	return
}

func (r *Redka) vParser(tagName string) (vStr string) {
	return strings.TrimPrefix(tagName, "v")
}

func (r *Redka) insParser(fName string) (insStr string) {
	return version.Unarchiver
}

func (r *Redka) Start() {
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

func (r *Redka) GetVersions() []byte {
	rr, _ := r.Version.Marshal()
	return rr
}

func (r *Redka) HomePage() string {
	return "https://github.com/nalgeon/redka"
}

func (r *Redka) GetInstallConf() (ic iconf.InstallerConfig) {
	return iconf.InstallerConfig{
		FlagFiles: &iconf.FileItems{
			Windows: []string{"redka.exe"},
			MacOS:   []string{"redka"},
			Linux:   []string{"redka"},
		},
		FlagDirExcepted: true,
		BinaryDirs: &iconf.DirItems{
			Windows: []iconf.DirPath{},
			MacOS:   []iconf.DirPath{},
			Linux:   []iconf.DirPath{},
		},
	}
}

func TestRedka() {
	nn := NewRedka()
	nn.Start()

	ff := fmt.Sprintf(
		"/home/moqsien/projects/go/src/gvcgo/vcollector/test/%s.json",
		nn.SDKName,
	)
	content, _ := json.MarshalIndent(nn.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

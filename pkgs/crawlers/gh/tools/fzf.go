package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gvcgo/vcollector/internal/gh"
	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
	"github.com/gvcgo/vcollector/pkgs/crawlers/gh/searcher"
	"github.com/gvcgo/vcollector/pkgs/version"
)

func init() {
	crawler.RegisterCrawler(NewFzf())
}

type Fzf struct {
	SDKName  string
	RepoName string
	searcher.GhSearcher
}

func NewFzf() (f *Fzf) {
	f = &Fzf{
		SDKName:  "fzf",
		RepoName: "junegunn/fzf",
		GhSearcher: searcher.GhSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (f *Fzf) GetSDKName() string {
	return f.SDKName
}

func (f *Fzf) tagFilter(ri gh.ReleaseItem) bool {
	return searcher.GVersionRegexp.FindString(ri.TagName) != ""
}

func (f *Fzf) fileFilter(a gh.Asset) bool {
	if strings.Contains(a.Url, "archive/refs/") {
		return false
	}
	if strings.HasSuffix(a.Name, ".txt") {
		return false
	}
	return true
}

func (f *Fzf) osParser(fName string) (osStr string) {
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

func (g *Fzf) archParser(fName string) (archStr string) {
	if strings.Contains(fName, "amd64") {
		return "amd64"
	}
	if strings.Contains(fName, "arm64") {
		return "arm64"
	}
	return
}

func (f *Fzf) vParser(tagName string) (vStr string) {
	return strings.TrimPrefix(tagName, "v")
}

func (f *Fzf) insParser(fName string) (insStr string) {
	return version.Unarchiver
}

func (f *Fzf) Start() {
	f.GhSearcher.Search(
		f.RepoName,
		f.tagFilter,
		f.fileFilter,
		f.vParser,
		f.archParser,
		f.osParser,
		f.insParser,
		nil,
	)
}

func (f *Fzf) GetVersions() []byte {
	r, _ := f.Version.Marshal()
	return r
}

func TestFzf() {
	nn := NewFzf()
	nn.Start()

	ff := fmt.Sprintf(
		"/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/%s.json",
		nn.SDKName,
	)
	content, _ := json.MarshalIndent(nn.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

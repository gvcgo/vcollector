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
	crawler.RegisterCrawler(NewLazygit())
}

type Lazygit struct {
	SDKName  string
	RepoName string
	searcher.GhSearcher
}

func NewLazygit() (l *Lazygit) {
	l = &Lazygit{
		SDKName:  "lazygit",
		RepoName: "jesseduffield/lazygit",
		GhSearcher: searcher.GhSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (l *Lazygit) GetSDKName() string {
	return l.SDKName
}

func (l *Lazygit) tagFilter(ri gh.ReleaseItem) bool {
	return searcher.GhVersionRegexp.FindString(ri.TagName) != ""
}

func (l *Lazygit) fileFilter(a gh.Asset) bool {
	if strings.Contains(a.Url, "archive/refs/") {
		return false
	}
	if strings.HasSuffix(a.Name, ".txt") {
		return false
	}
	return true
}

func (l *Lazygit) osParser(fName string) (osStr string) {
	fName = strings.ToLower(fName)
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

func (l *Lazygit) archParser(fName string) (archStr string) {
	if strings.Contains(fName, "x86_64") {
		return "amd64"
	}
	if strings.Contains(fName, "arm64") {
		return "arm64"
	}
	return
}

func (l *Lazygit) vParser(tagName string) (vStr string) {
	return strings.TrimPrefix(tagName, "v")
}

func (l *Lazygit) insParser(fName string) (insStr string) {
	return version.Unarchiver
}

func (l *Lazygit) Start() {
	l.GhSearcher.Search(
		l.RepoName,
		l.tagFilter,
		l.fileFilter,
		l.vParser,
		l.archParser,
		l.osParser,
		l.insParser,
		nil,
	)
}

func (l *Lazygit) GetVersions() []byte {
	r, _ := l.Version.Marshal()
	return r
}

func TestLazygit() {
	nn := NewLazygit()
	nn.Start()

	ff := fmt.Sprintf(
		"/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/%s.json",
		nn.SDKName,
	)
	content, _ := json.MarshalIndent(nn.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

package lsp

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
	crawler.RegisterCrawler(NewVAnalyzer())
}

type VAnalyzer struct {
	SDKName  string
	RepoName string
	searcher.GhSearcher
}

func NewVAnalyzer() (v *VAnalyzer) {
	v = &VAnalyzer{
		SDKName:  "v-analyzer",
		RepoName: "v-analyzer/v-analyzer",
		GhSearcher: searcher.GhSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (v *VAnalyzer) GetSDKName() string {
	return v.SDKName
}

func (v *VAnalyzer) tagFilter(ri gh.ReleaseItem) bool {
	return searcher.GVersionRegexp.FindString(ri.TagName) != ""
}

func (v *VAnalyzer) fileFilter(a gh.Asset) bool {
	return !strings.Contains(a.Url, "archive/refs/")
}

func (v *VAnalyzer) osParser(fName string) (osStr string) {
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

func (v *VAnalyzer) archParser(fName string) (archStr string) {
	if strings.Contains(fName, "-x86_64") {
		return "amd64"
	}
	if strings.Contains(fName, "-arm64") {
		return "arm64"
	}
	return
}

func (v *VAnalyzer) vParser(tagName string) (vStr string) {
	return strings.TrimPrefix(tagName, "v")
}

func (v *VAnalyzer) insParser(fName string) (insStr string) {
	return version.Unarchiver
}

func (v *VAnalyzer) Start() {
	v.GhSearcher.Search(
		v.RepoName,
		v.tagFilter,
		v.fileFilter,
		v.vParser,
		v.archParser,
		v.osParser,
		v.insParser,
		nil,
	)
}

func (v *VAnalyzer) GetVersions() []byte {
	r, _ := v.Version.Marshal()
	return r
}

func TestVAnalyzer() {
	nn := NewVAnalyzer()
	nn.Start()

	ff := fmt.Sprintf(
		"/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/%s.json",
		nn.SDKName,
	)
	content, _ := json.MarshalIndent(nn.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

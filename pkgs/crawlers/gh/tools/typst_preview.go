package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gvcgo/vcollector/internal/gh"
	"github.com/gvcgo/vcollector/pkgs/crawlers/gh/searcher"
	"github.com/gvcgo/vcollector/pkgs/version"
)

type TypstPreview struct {
	SDKName  string
	RepoName string
	searcher.GhSearcher
}

func NewTypstPreview() (t *TypstPreview) {
	t = &TypstPreview{
		SDKName:  "typst-preview",
		RepoName: "Enter-tainer/typst-preview",
		GhSearcher: searcher.GhSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (t *TypstPreview) GetSDKName() string {
	return t.SDKName
}

func (t *TypstPreview) tagFilter(ri gh.ReleaseItem) bool {
	return searcher.GhVersionRegexp.FindString(ri.TagName) != ""
}

func (t *TypstPreview) fileFilter(a gh.Asset) bool {
	if strings.Contains(a.Url, "archive/refs/") {
		return false
	}
	if strings.HasSuffix(a.Name, ".dwarf") {
		return false
	}
	if strings.HasSuffix(a.Name, ".debug") {
		return false
	}
	if strings.HasSuffix(a.Name, ".vsix") {
		return false
	}
	if strings.HasSuffix(a.Name, ".html") {
		return false
	}
	return true
}

func (t *TypstPreview) osParser(fName string) (osStr string) {
	if strings.Contains(fName, "darwin") {
		return "darwin"
	}
	if strings.Contains(fName, "linux") {
		return "linux"
	}
	if strings.Contains(fName, "-win32") {
		return "windows"
	}
	return
}

func (t *TypstPreview) archParser(fName string) (archStr string) {
	if strings.Contains(fName, "-x64") {
		return "amd64"
	}
	if strings.Contains(fName, "-arm64") {
		return "arm64"
	}
	return
}

func (t *TypstPreview) vParser(tagName string) (vStr string) {
	return strings.TrimPrefix(tagName, "v")
}

func (t *TypstPreview) insParser(fName string) (insStr string) {
	return version.Executable
}

func (t *TypstPreview) Start() {
	t.GhSearcher.Search(
		t.RepoName,
		t.tagFilter,
		t.fileFilter,
		t.vParser,
		t.archParser,
		t.osParser,
		t.insParser,
		nil,
	)
}

func TestTypstPreview() {
	nn := NewTypstPreview()
	nn.Start()

	ff := fmt.Sprintf(
		"/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/%s.json",
		nn.SDKName,
	)
	content, _ := json.MarshalIndent(nn.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

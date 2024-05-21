package lsp

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gvcgo/vcollector/internal/gh"
	"github.com/gvcgo/vcollector/pkgs/crawlers/gh/searcher"
	"github.com/gvcgo/vcollector/pkgs/version"
)

type Zls struct {
	SDKName  string
	RepoName string
	searcher.GhSearcher
}

func NewZls() (z *Zls) {
	z = &Zls{
		SDKName:  "zls",
		RepoName: "zigtools/zls",
		GhSearcher: searcher.GhSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (z *Zls) GetSDKName() string {
	return z.SDKName
}

func (z *Zls) tagFilter(ri gh.ReleaseItem) bool {
	return GVersionRegexp.FindString(ri.TagName) != ""
}

func (z *Zls) fileFilter(a gh.Asset) bool {
	if strings.Contains(a.Url, "archive/refs/") {
		return false
	}
	if strings.Contains(a.Name, "-wasm32") {
		return false
	}
	if strings.Contains(a.Name, "-x86-") {
		return false
	}
	return true
}

func (z *Zls) osParser(fName string) (osStr string) {
	if strings.Contains(fName, "macos") {
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

func (z *Zls) archParser(fName string) (archStr string) {
	if strings.Contains(fName, "x86_64") {
		return "amd64"
	}
	if strings.Contains(fName, "aarch64") {
		return "arm64"
	}
	return
}

func (z *Zls) vParser(tagName string) (vStr string) {
	return strings.TrimPrefix(tagName, "v")
}

func (z *Zls) insParser(fName string) (insStr string) {
	return version.Unarchiver
}

func (z *Zls) Start() {
	z.GhSearcher.Search(
		z.RepoName,
		z.tagFilter,
		z.fileFilter,
		z.vParser,
		z.archParser,
		z.osParser,
		z.insParser,
		nil,
	)
}

func TestZls() {
	nn := NewZls()
	nn.Start()

	ff := fmt.Sprintf(
		"/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/%s.json",
		nn.SDKName,
	)
	content, _ := json.MarshalIndent(nn.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

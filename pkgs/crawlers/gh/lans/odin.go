package lans

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gvcgo/vcollector/internal/gh"
	"github.com/gvcgo/vcollector/pkgs/crawlers/gh/searcher"
	"github.com/gvcgo/vcollector/pkgs/version"
)

type Odin struct {
	SDKName  string
	RepoName string
	searcher.GhSearcher
}

func NewOdin() (o *Odin) {
	o = &Odin{
		SDKName:  "odin",
		RepoName: "odin-lang/Odin",
		GhSearcher: searcher.GhSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (o *Odin) GetSDKName() string {
	return o.SDKName
}

func (o *Odin) tagFilter(ri gh.ReleaseItem) bool {
	if searcher.GhVersionRegexp.FindString(ri.TagName) != "" {
		return true
	}
	if strings.HasPrefix(ri.TagName, "dev-") {
		return true
	}
	return false
}

func (o *Odin) fileFilter(a gh.Asset) bool {
	return !strings.Contains(a.Url, "archive/refs/")
}

func (o *Odin) osParser(fName string) (osStr string) {
	if strings.Contains(fName, "macos") {
		return "darwin"
	}
	if strings.Contains(fName, "ubuntu") {
		return "linux"
	}
	if strings.Contains(fName, "windows") {
		return "windows"
	}
	return
}

func (o *Odin) archParser(fName string) (archStr string) {
	if strings.Contains(fName, "-amd64") {
		return "amd64"
	}
	if strings.Contains(fName, "-arm64") {
		return "arm64"
	}
	return
}

func (o *Odin) vParser(tagName string) (vStr string) {
	return strings.TrimPrefix(tagName, "v")
}

func (o *Odin) insParser(fName string) (insStr string) {
	return version.Unarchiver
}

func (o *Odin) Start() {
	o.GhSearcher.Search(
		o.RepoName,
		o.tagFilter,
		o.fileFilter,
		o.vParser,
		o.archParser,
		o.osParser,
		o.insParser,
		nil,
	)
}

func TestOdin() {
	nn := NewOdin()
	nn.Start()

	ff := fmt.Sprintf(
		"/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/%s.json",
		nn.SDKName,
	)
	content, _ := json.MarshalIndent(nn.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

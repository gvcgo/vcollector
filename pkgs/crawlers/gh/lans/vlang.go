package lans

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/gvcgo/vcollector/internal/gh"
	"github.com/gvcgo/vcollector/pkgs/crawlers/gh/searcher"
	"github.com/gvcgo/vcollector/pkgs/version"
)

var GVersionRegexp = regexp.MustCompile(`\d+(.\d+){2}`)

type Vlang struct {
	SDKName  string
	RepoName string
	searcher.GhSearcher
}

func NewVlang() (v *Vlang) {
	v = &Vlang{
		SDKName:  "v",
		RepoName: "vlang/v",
		GhSearcher: searcher.GhSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (v *Vlang) GetSDKName() string {
	return v.SDKName
}

func (v *Vlang) tagFilter(ri gh.ReleaseItem) bool {
	if GVersionRegexp.FindString(ri.TagName) != "" {
		return true
	}
	if strings.HasPrefix(ri.TagName, "weekly.") {
		return true
	}
	return false
}

func (v *Vlang) fileFilter(a gh.Asset) bool {
	return !strings.Contains(a.Url, "archive/refs/")
}

func (v *Vlang) osParser(fName string) (osStr string) {
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

func (v *Vlang) archParser(fName string) (archStr string) {
	if strings.Contains(fName, "x86_64") {
		return "amd64"
	}
	if strings.Contains(fName, "arm64") {
		return "arm64"
	}
	return "amd64"
}

func (v *Vlang) vParser(tagName string) (vStr string) {
	return strings.TrimPrefix(tagName, "v")
}

func (v *Vlang) insParser(fName string) (insStr string) {
	return version.Unarchiver
}

func (v *Vlang) Start() {
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

func TestVlang() {
	nn := NewVlang()
	nn.Start()

	ff := fmt.Sprintf(
		"/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/%s.json",
		nn.SDKName,
	)
	content, _ := json.MarshalIndent(nn.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

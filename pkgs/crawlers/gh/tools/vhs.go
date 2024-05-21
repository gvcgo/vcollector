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

type Vhs struct {
	SDKName  string
	RepoName string
	searcher.GhSearcher
}

func NewVhs() (v *Vhs) {
	v = &Vhs{
		SDKName:  "vhs",
		RepoName: "charmbracelet/vhs",
		GhSearcher: searcher.GhSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (v *Vhs) GetSDKName() string {
	return v.SDKName
}

func (v *Vhs) tagFilter(ri gh.ReleaseItem) bool {
	return searcher.GhVersionRegexp.FindString(ri.TagName) != ""
}

func (v *Vhs) fileFilter(a gh.Asset) bool {
	if strings.Contains(a.Url, "archive/refs/") {
		return false
	}
	excludeSuffixList := []string{
		".txt",
		".pem",
		".sig",
		".rpm",
		".sbom",
		".deb",
		".apk",
	}
	for _, ss := range excludeSuffixList {
		if strings.HasSuffix(a.Name, ss) {
			return false
		}
	}
	return true
}

func (v *Vhs) osParser(fName string) (osStr string) {
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

func (v *Vhs) archParser(fName string) (archStr string) {
	if strings.Contains(fName, "x86_64") {
		return "amd64"
	}
	if strings.Contains(fName, "arm64") {
		return "arm64"
	}
	return
}

func (v *Vhs) vParser(tagName string) (vStr string) {
	return strings.TrimPrefix(tagName, "v")
}

func (v *Vhs) insParser(fName string) (insStr string) {
	return version.Unarchiver
}

func (v *Vhs) Start() {
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

func TestVhs() {
	nn := NewVhs()
	nn.Start()

	ff := fmt.Sprintf(
		"/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/%s.json",
		nn.SDKName,
	)
	content, _ := json.MarshalIndent(nn.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

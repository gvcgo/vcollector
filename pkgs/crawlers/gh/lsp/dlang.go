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
	crawler.RegisterCrawler(NewDlangLsp())
}

type DlangLsp struct {
	SDKName  string
	RepoName string
	searcher.GhSearcher
}

func NewDlangLsp() (d *DlangLsp) {
	d = &DlangLsp{
		SDKName:  "served",
		RepoName: "Pure-D/serve-d",
		GhSearcher: searcher.GhSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (d *DlangLsp) GetSDKName() string {
	return d.SDKName
}

func (d *DlangLsp) tagFilter(ri gh.ReleaseItem) bool {
	return searcher.GhVersionRegexp.FindString(ri.TagName) != ""
}

func (d *DlangLsp) fileFilter(a gh.Asset) bool {
	if strings.Contains(a.Url, "archive/refs/") {
		return false
	}
	if strings.HasSuffix(a.Name, ".tar.xz") {
		return false
	}
	if strings.HasSuffix(a.Name, "-x86.zip") {
		return false
	}
	return true
}

func (d *DlangLsp) osParser(fName string) (osStr string) {
	if strings.Contains(fName, "-osx") {
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

func (d *DlangLsp) archParser(fName string) (archStr string) {
	if strings.Contains(fName, "-x86_64") {
		return "amd64"
	}
	if strings.Contains(fName, "-arm64") {
		return "arm64"
	}
	return
}

func (d *DlangLsp) vParser(tagName string) (vStr string) {
	return strings.TrimPrefix(tagName, "v")
}

func (d *DlangLsp) insParser(fName string) (insStr string) {
	return version.Unarchiver
}

func (d *DlangLsp) Start() {
	d.GhSearcher.Search(
		d.RepoName,
		d.tagFilter,
		d.fileFilter,
		d.vParser,
		d.archParser,
		d.osParser,
		d.insParser,
		nil,
	)
}

func (d *DlangLsp) GetVersions() []byte {
	r, _ := d.Version.Marshal()
	return r
}

func TestDlangLsp() {
	nn := NewDlangLsp()
	nn.Start()

	ff := fmt.Sprintf(
		"/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/%s.json",
		nn.SDKName,
	)
	content, _ := json.MarshalIndent(nn.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

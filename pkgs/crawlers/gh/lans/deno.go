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

type Deno struct {
	SDKName  string
	RepoName string
	searcher.GhSearcher
}

func NewDeno() (d *Deno) {
	d = &Deno{
		SDKName:  "deno",
		RepoName: "denoland/deno",
		GhSearcher: searcher.GhSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (d *Deno) GetSDKName() string {
	return d.SDKName
}

func (d *Deno) tagFilter(ri gh.ReleaseItem) bool {
	return GhVersionRegexp.FindString(ri.TagName) != ""
}

func (d *Deno) fileFilter(a gh.Asset) bool {
	if strings.Contains(a.Url, "archive/refs/") {
		return false
	}
	if strings.Contains(a.Name, "denort-") {
		return false
	}
	if strings.HasSuffix(a.Name, "src.tar.gz") {
		return false
	}
	if strings.HasSuffix(a.Name, ".deno.d.ts") {
		return false
	}
	return true
}

func (d *Deno) osParser(fName string) (osStr string) {
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

func (d *Deno) archParser(fName string) (archStr string) {
	if strings.Contains(fName, "-x86_64") {
		return "amd64"
	}
	if strings.Contains(fName, "-aarch64") {
		return "arm64"
	}
	return
}

func (d *Deno) vParser(tagName string) (vStr string) {
	return strings.TrimPrefix(tagName, "v")
}

func (d *Deno) insParser(fName string) (insStr string) {
	return version.Unarchiver
}

func (d *Deno) Start() {
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

func TestDeno() {
	bb := NewDeno()
	bb.Start()

	ff := fmt.Sprintf(
		"/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/%s.json",
		bb.SDKName,
	)
	content, _ := json.MarshalIndent(bb.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

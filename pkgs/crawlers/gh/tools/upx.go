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
	crawler.RegisterCrawler(NewUpx())
}

type Upx struct {
	SDKName  string
	RepoName string
	searcher.GhSearcher
}

func NewUpx() (u *Upx) {
	u = &Upx{
		SDKName:  "upx",
		RepoName: "upx/upx",
		GhSearcher: searcher.GhSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (u *Upx) GetSDKName() string {
	return u.SDKName
}

func (u *Upx) tagFilter(ri gh.ReleaseItem) bool {
	return searcher.VVersionRegexp.FindString(ri.TagName) != ""
}

func (g *Upx) fileFilter(a gh.Asset) bool {
	if strings.Contains(a.Url, "archive/refs/") {
		return false
	}
	if strings.Contains(a.Name, "-docs") {
		return false
	}
	if strings.Contains(a.Name, "-src") {
		return false
	}
	return true
}

func (g *Upx) osParser(fName string) (osStr string) {
	if strings.Contains(fName, "darwin") {
		return "darwin"
	}
	if strings.Contains(fName, "linux") {
		return "linux"
	}
	if strings.Contains(fName, "-win") {
		return "windows"
	}
	return
}

func (g *Upx) archParser(fName string) (archStr string) {
	if strings.Contains(fName, "-amd64") {
		return "amd64"
	}
	if strings.Contains(fName, "-arm64") {
		return "arm64"
	}
	if strings.Contains(fName, "-win64") {
		return "amd64"
	}
	return
}

func (g *Upx) vParser(tagName string) (vStr string) {
	return strings.TrimPrefix(tagName, "v")
}

func (g *Upx) insParser(fName string) (insStr string) {
	return version.Unarchiver
}

func (g *Upx) Start() {
	g.GhSearcher.Search(
		g.RepoName,
		g.tagFilter,
		g.fileFilter,
		g.vParser,
		g.archParser,
		g.osParser,
		g.insParser,
		nil,
	)
}

func (u *Upx) GetVersions() []byte {
	r, _ := u.Version.Marshal()
	return r
}

func (u *Upx) HomePage() string {
	return "https://upx.github.io/"
}

func TestUpx() {
	nn := NewUpx()
	nn.Start()

	ff := fmt.Sprintf(
		"/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/%s.json",
		nn.SDKName,
	)
	content, _ := json.MarshalIndent(nn.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

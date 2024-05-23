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
	crawler.RegisterCrawler(NewGsudo())
}

type Gsudo struct {
	SDKName  string
	RepoName string
	searcher.GhSearcher
}

func NewGsudo() (g *Gsudo) {
	g = &Gsudo{
		SDKName:  "gsudo",
		RepoName: "gerardog/gsudo",
		GhSearcher: searcher.GhSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (g *Gsudo) GetSDKName() string {
	return g.SDKName
}

func (g *Gsudo) tagFilter(ri gh.ReleaseItem) bool {
	return searcher.GhVersionRegexp.FindString(ri.TagName) != ""
}

func (g *Gsudo) fileFilter(a gh.Asset) bool {
	if strings.Contains(a.Url, "archive/refs/") {
		return false
	}
	return strings.HasSuffix(a.Name, "gsudo.portable.zip")
}

func (g *Gsudo) osParser(fName string) (osStr string) {
	return "windows"
}

func (g *Gsudo) archParser(fName string) (archStr string) {
	return "amd64"
}

func (g *Gsudo) vParser(tagName string) (vStr string) {
	return strings.TrimPrefix(tagName, "v")
}

func (g *Gsudo) insParser(fName string) (insStr string) {
	return version.Unarchiver
}

func (g *Gsudo) Start() {
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

func (g *Gsudo) GetVersions() []byte {
	r, _ := g.Version.Marshal()
	return r
}

func (g *Gsudo) HomePage() string {
	return "https://gerardog.github.io/gsudo/"
}

func TestGsudo() {
	nn := NewGsudo()
	nn.Start()

	ff := fmt.Sprintf(
		"/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/%s.json",
		nn.SDKName,
	)
	content, _ := json.MarshalIndent(nn.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

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
	crawler.RegisterCrawler(NewGitWin())
}

type GitWin struct {
	SDKName  string
	RepoName string
	searcher.GhSearcher
}

func NewGitWin() (g *GitWin) {
	g = &GitWin{
		SDKName:  "git-win",
		RepoName: "git-for-windows/git",
		GhSearcher: searcher.GhSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (g *GitWin) GetSDKName() string {
	return g.SDKName
}

func (g *GitWin) tagFilter(ri gh.ReleaseItem) bool {
	return searcher.GhVersionRegexp.FindString(ri.TagName) != ""
}

func (g *GitWin) fileFilter(a gh.Asset) bool {
	return strings.HasPrefix(a.Name, "PortableGit-")
}

func (g *GitWin) osParser(fName string) (osStr string) {
	return "windows"
}

func (g *GitWin) archParser(fName string) (archStr string) {
	if strings.Contains(fName, "-64-bit") {
		return "amd64"
	}
	return
}

func (g *GitWin) vParser(tagName string) (vStr string) {
	return strings.TrimPrefix(tagName, "v")
}

func (g *GitWin) insParser(fName string) (insStr string) {
	return version.Unarchiver
}

func (g *GitWin) Start() {
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

func (g *GitWin) GetVersions() []byte {
	r, _ := g.Version.Marshal()
	return r
}

func TestGitWin() {
	nn := NewGitWin()
	nn.Start()

	ff := fmt.Sprintf(
		"/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/%s.json",
		nn.SDKName,
	)
	content, _ := json.MarshalIndent(nn.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

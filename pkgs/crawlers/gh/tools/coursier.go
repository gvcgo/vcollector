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
	crawler.RegisterCrawler(NewCoursier())
}

/*
https://github.com/coursier/coursier
https://github.com/VirtusLab/coursier-m1
*/
type Coursier struct {
	SDKName   string
	RepoName1 string
	RepoName2 string
	searcher.GhSearcher
}

func NewCoursier() (c *Coursier) {
	c = &Coursier{
		SDKName:   "coursier",
		RepoName1: "coursier/coursier",
		RepoName2: "VirtusLab/coursier-m1",
		GhSearcher: searcher.GhSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (c *Coursier) GetSDKName() string {
	return c.SDKName
}

func (c *Coursier) tagFilter(ri gh.ReleaseItem) bool {
	return searcher.GhVersionRegexp.FindString(ri.TagName) != ""
}

func (c *Coursier) fileFilter(a gh.Asset) bool {
	if strings.HasPrefix(a.Name, "cs-") && strings.HasSuffix(a.Name, "-sdk.zip") {
		return true
	}
	return false
}

func (c *Coursier) osParser(fName string) (osStr string) {
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

func (c *Coursier) archParser(fName string) (archStr string) {
	if strings.Contains(fName, "-x86_64") {
		return "amd64"
	}
	if strings.Contains(fName, "-aarch64") {
		return "arm64"
	}
	return
}

func (c *Coursier) vParser(tagName string) (vStr string) {
	return strings.TrimPrefix(tagName, "v")
}

func (c *Coursier) insParser(fName string) (insStr string) {
	return version.Unarchiver
}

func (c *Coursier) Start() {
	repoList := []string{c.RepoName1, c.RepoName2}
	for _, repoName := range repoList {
		c.GhSearcher.Search(
			repoName,
			c.tagFilter,
			c.fileFilter,
			c.vParser,
			c.archParser,
			c.osParser,
			c.insParser,
			nil,
		)
	}
}

func (c *Coursier) GetVersions() []byte {
	r, _ := c.Version.Marshal()
	return r
}

func TestCoursier() {
	nn := NewCoursier()
	nn.Start()

	ff := fmt.Sprintf(
		"/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/%s.json",
		nn.SDKName,
	)
	content, _ := json.MarshalIndent(nn.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

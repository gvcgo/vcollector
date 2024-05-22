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
	crawler.RegisterCrawler(NewCMake())
}

type CMake struct {
	SDKName  string
	RepoName string
	searcher.GhSearcher
}

func NewCMake() (c *CMake) {
	c = &CMake{
		SDKName:  "cmake",
		RepoName: "Kitware/CMake",
		GhSearcher: searcher.GhSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (c *CMake) GetSDKName() string {
	return c.SDKName
}

func (c *CMake) tagFilter(ri gh.ReleaseItem) bool {
	return searcher.GhVersionRegexp.FindString(ri.TagName) != ""
}

func (c *CMake) fileFilter(a gh.Asset) bool {
	if strings.Contains(a.Url, "archive/refs/") {
		return false
	}
	if strings.HasSuffix(a.Name, ".json") {
		return false
	}
	if strings.HasSuffix(a.Name, ".sh") {
		return false
	}
	if strings.HasSuffix(a.Name, ".dmg") {
		return false
	}
	if strings.HasSuffix(a.Name, ".txt") {
		return false
	}
	if strings.HasSuffix(a.Name, ".asc") {
		return false
	}
	if strings.HasSuffix(a.Name, ".msi") {
		return false
	}
	if strings.Contains(a.Name, "macos10.10-") {
		return false
	}
	return true
}

func (c *CMake) osParser(fName string) (osStr string) {
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

func (c *CMake) archParser(fName string) (archStr string) {
	if strings.Contains(fName, "-x86_64") {
		return "amd64"
	}
	if strings.Contains(fName, "-aarch64") {
		return "arm64"
	}
	if strings.Contains(fName, "-arm64") {
		return "arm64"
	}
	if strings.Contains(fName, "-universal") {
		return "any"
	}
	return
}

func (c *CMake) vParser(tagName string) (vStr string) {
	return strings.TrimPrefix(tagName, "v")
}

func (c *CMake) insParser(fName string) (insStr string) {
	return version.Unarchiver
}

func (c *CMake) Start() {
	c.GhSearcher.Search(
		c.RepoName,
		c.tagFilter,
		c.fileFilter,
		c.vParser,
		c.archParser,
		c.osParser,
		c.insParser,
		nil,
	)
}

func (c *CMake) GetVersions() []byte {
	r, _ := c.Version.Marshal()
	return r
}

func TestCMake() {
	nn := NewCMake()
	nn.Start()

	ff := fmt.Sprintf(
		"/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/%s.json",
		nn.SDKName,
	)
	content, _ := json.MarshalIndent(nn.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

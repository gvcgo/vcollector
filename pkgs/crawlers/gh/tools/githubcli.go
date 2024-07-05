package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gvcgo/vcollector/internal/gh"
	"github.com/gvcgo/vcollector/internal/iconf"
	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
	"github.com/gvcgo/vcollector/pkgs/crawlers/gh/searcher"
	"github.com/gvcgo/vcollector/pkgs/version"
)

func init() {
	crawler.RegisterCrawler(NewGithubCli())
}

type GithubCli struct {
	SDKName  string
	RepoName string
	searcher.GhSearcher
}

func NewGithubCli() (g *GithubCli) {
	g = &GithubCli{
		SDKName:  "github-cli",
		RepoName: "cli/cli",
		GhSearcher: searcher.GhSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (g *GithubCli) GetSDKName() string {
	return g.SDKName
}

func (g *GithubCli) tagFilter(ri gh.ReleaseItem) bool {
	return searcher.GVersionRegexp.FindString(ri.TagName) != ""
}

func (g *GithubCli) fileFilter(a gh.Asset) bool {
	if strings.Contains(a.Url, "archive/refs/") {
		return false
	}
	if strings.HasSuffix(a.Url, "checksums.txt") {
		return false
	}
	if strings.HasSuffix(a.Url, ".deb") {
		return false
	}
	if strings.HasSuffix(a.Url, ".rpm") {
		return false
	}
	if strings.HasSuffix(a.Url, ".msi") {
		return false
	}
	if strings.Contains(a.Url, "_386") {
		return false
	}
	return true
}

func (g *GithubCli) osParser(fName string) (osStr string) {
	if strings.Contains(fName, "macOS") {
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

func (g *GithubCli) archParser(fName string) (archStr string) {
	if strings.Contains(fName, "amd64") {
		return "amd64"
	}
	if strings.Contains(fName, "arm64") {
		return "arm64"
	}
	return
}

func (g *GithubCli) vParser(tagName string) (vStr string) {
	return strings.TrimPrefix(tagName, "v")
}

func (g *GithubCli) insParser(fName string) (insStr string) {
	return version.Unarchiver
}

func (g *GithubCli) Start() {
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

func (g *GithubCli) GetVersions() []byte {
	rr, _ := g.Version.Marshal()
	return rr
}

func (g *GithubCli) HomePage() string {
	return "https://cli.github.com/"
}

func (g *GithubCli) GetInstallConf() (ic iconf.InstallerConfig) {
	return iconf.InstallerConfig{
		FlagFiles: &iconf.FileItems{
			Windows: []string{"bin"},
			MacOS:   []string{"bin"},
			Linux:   []string{"bin"},
		},
		FlagDirExcepted: false,
		BinaryDirs: &iconf.DirItems{
			Windows: []iconf.DirPath{{"bin"}},
			MacOS:   []iconf.DirPath{{"bin"}},
			Linux:   []iconf.DirPath{{"bin"}},
		},
	}
}

func TestGithubCli() {
	nn := NewGithubCli()
	nn.Start()

	ff := fmt.Sprintf(
		"/home/moqsien/projects/go/src/gvcgo/vcollector/test/%s.json",
		nn.SDKName,
	)
	content, _ := json.MarshalIndent(nn.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

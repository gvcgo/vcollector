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
	crawler.RegisterCrawler(NewTreeSitter())
}

type TreeSitter struct {
	SDKName  string
	RepoName string
	searcher.GhSearcher
}

func NewTreeSitter() (t *TreeSitter) {
	t = &TreeSitter{
		SDKName:  "tree-sitter",
		RepoName: "tree-sitter/tree-sitter",
		GhSearcher: searcher.GhSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (t *TreeSitter) GetSDKName() string {
	return t.SDKName
}

func (t *TreeSitter) tagFilter(ri gh.ReleaseItem) bool {
	return searcher.GVersionRegexp.FindString(ri.TagName) != ""
}

func (t *TreeSitter) fileFilter(a gh.Asset) bool {
	if strings.Contains(a.Url, "archive/refs/") {
		return false
	}
	if strings.HasSuffix(a.Name, ".wasm") {
		return false
	}
	if strings.HasSuffix(a.Name, ".js") {
		return false
	}
	return true
}

func (t *TreeSitter) osParser(fName string) (osStr string) {
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

func (t *TreeSitter) archParser(fName string) (archStr string) {
	if strings.Contains(fName, "-x64") {
		return "amd64"
	}
	if strings.Contains(fName, "-arm64") {
		return "arm64"
	}
	return
}

func (t *TreeSitter) vParser(tagName string) (vStr string) {
	return strings.TrimPrefix(tagName, "v")
}

func (t *TreeSitter) insParser(fName string) (insStr string) {
	return version.Unarchiver
}

func (t *TreeSitter) Start() {
	t.GhSearcher.Search(
		t.RepoName,
		t.tagFilter,
		t.fileFilter,
		t.vParser,
		t.archParser,
		t.osParser,
		t.insParser,
		nil,
	)
}

func (t *TreeSitter) GetVersions() []byte {
	r, _ := t.Version.Marshal()
	return r
}

func (t *TreeSitter) HomePage() string {
	return "https://tree-sitter.github.io/tree-sitter/"
}

func (t *TreeSitter) GetInstallConf() (ic iconf.InstallerConfig) {
	return iconf.InstallerConfig{
		FlagFiles: &iconf.FileItems{
			Windows: []string{"tree-sitter.exe"},
			MacOS:   []string{"tree-sitter"},
			Linux:   []string{"tree-sitter"},
		},
		FlagDirExcepted: true,
		BinaryDirs: &iconf.DirItems{
			Windows: []iconf.DirPath{},
			MacOS:   []iconf.DirPath{},
			Linux:   []iconf.DirPath{},
		},
		BinaryRename: &iconf.BinaryRename{
			NameFlag: "tree-sitter",
			RenameTo: "tree-sitter",
		},
	}
}

func TestTreeSitter() {
	nn := NewTreeSitter()
	nn.Start()

	ff := fmt.Sprintf(
		"/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/%s.json",
		nn.SDKName,
	)
	content, _ := json.MarshalIndent(nn.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

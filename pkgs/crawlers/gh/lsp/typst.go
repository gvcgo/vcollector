package lsp

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
	crawler.RegisterCrawler(NewTypstLsp())
}

type TypstLsp struct {
	SDKName  string
	RepoName string
	searcher.GhSearcher
}

func NewTypstLsp() (t *TypstLsp) {
	t = &TypstLsp{
		SDKName:  "typst-lsp",
		RepoName: "nvarner/typst-lsp",
		GhSearcher: searcher.GhSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (t *TypstLsp) GetSDKName() string {
	return t.SDKName
}

func (t *TypstLsp) tagFilter(ri gh.ReleaseItem) bool {
	return searcher.GhVersionRegexp.FindString(ri.TagName) != ""
}

func (t *TypstLsp) fileFilter(a gh.Asset) bool {
	if strings.Contains(a.Url, "archive/refs/") {
		return false
	}
	if strings.HasSuffix(a.Name, "linux-gnu") {
		return false
	}
	if strings.HasSuffix(a.Name, "-linux-gnueabihf") {
		return false
	}
	if strings.HasSuffix(a.Name, ".toml") {
		return false
	}
	if strings.HasSuffix(a.Name, ".png") {
		return false
	}
	if strings.HasSuffix(a.Name, ".md") {
		return false
	}
	if strings.HasSuffix(a.Name, ".vsix") {
		return false
	}
	return true
}

func (t *TypstLsp) osParser(fName string) (osStr string) {
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

func (t *TypstLsp) archParser(fName string) (archStr string) {
	if strings.Contains(fName, "-x86_64") {
		return "amd64"
	}
	if strings.Contains(fName, "-aarch64") {
		return "arm64"
	}
	return
}

func (t *TypstLsp) vParser(tagName string) (vStr string) {
	return strings.TrimPrefix(tagName, "v")
}

func (t *TypstLsp) insParser(fName string) (insStr string) {
	return version.Executable
}

func (t *TypstLsp) Start() {
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

func (t *TypstLsp) GetVersions() []byte {
	r, _ := t.Version.Marshal()
	return r
}

func (t *TypstLsp) HomePage() string {
	return "https://github.com/nvarner/typst-lsp"
}

func (t *TypstLsp) GetInstallConf() (ic iconf.InstallerConfig) {
	return iconf.InstallerConfig{
		FlagFiles: &iconf.FileItems{
			Windows: []string{"typst-lsp.exe"},
			MacOS:   []string{"typst-lsp"},
			Linux:   []string{"typst-lsp"},
		},
		FlagDirExcepted: true,
		BinaryDirs: &iconf.DirItems{
			Windows: []iconf.DirPath{},
			MacOS:   []iconf.DirPath{},
			Linux:   []iconf.DirPath{},
		},
		BinaryRename: &iconf.BinaryRename{
			NameFlag: "typst-lsp",
			RenameTo: "typst-lsp",
		},
	}
}

func TestTypstLsp() {
	nn := NewTypstLsp()
	nn.Start()

	ff := fmt.Sprintf(
		"/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/%s.json",
		nn.SDKName,
	)
	content, _ := json.MarshalIndent(nn.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

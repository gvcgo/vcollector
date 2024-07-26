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
	crawler.RegisterCrawler(NewTinyMist())
}

type TinyMist struct {
	SDKName  string
	RepoName string
	searcher.GhSearcher
}

func NewTinyMist() (t *TinyMist) {
	t = &TinyMist{
		SDKName:  "tinymist",
		RepoName: "Myriad-Dreamin/tinymist",
		GhSearcher: searcher.GhSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (t *TinyMist) GetSDKName() string {
	return t.SDKName
}

func (t *TinyMist) tagFilter(ri gh.ReleaseItem) bool {
	return searcher.GhVersionRegexp.FindString(ri.TagName) != ""
}

func (t *TinyMist) fileFilter(a gh.Asset) bool {
	if strings.Contains(a.Url, "archive/refs/") {
		return false
	}
	if strings.HasSuffix(a.Name, ".dwarf") {
		return false
	}
	if strings.HasSuffix(a.Name, ".debug") {
		return false
	}
	if strings.HasSuffix(a.Name, ".vsix") {
		return false
	}
	if strings.HasSuffix(a.Name, ".html") {
		return false
	}
	if strings.HasSuffix(a.Name, ".pdb") {
		return false
	}
	return true
}

func (t *TinyMist) osParser(fName string) (osStr string) {
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

func (t *TinyMist) archParser(fName string) (archStr string) {
	if strings.Contains(fName, "-x64") {
		return "amd64"
	}
	if strings.Contains(fName, "-arm64") {
		return "arm64"
	}
	return
}

func (t *TinyMist) vParser(tagName string) (vStr string) {
	return strings.TrimPrefix(tagName, "v")
}

func (t *TinyMist) insParser(fName string) (insStr string) {
	return version.Executable
}

func (t *TinyMist) Start() {
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

func (t *TinyMist) GetVersions() []byte {
	r, _ := t.Version.Marshal()
	return r
}

func (t *TinyMist) HomePage() string {
	return "https://github.com/Myriad-Dreamin/tinymist"
}

func (t *TinyMist) GetInstallConf() (ic iconf.InstallerConfig) {
	return iconf.InstallerConfig{
		FlagFiles: &iconf.FileItems{
			Windows: []string{"tinymist.exe"},
			MacOS:   []string{"tinymist"},
			Linux:   []string{"tinymist"},
		},
		FlagDirExcepted: true,
		BinaryDirs: &iconf.DirItems{
			Windows: []iconf.DirPath{},
			MacOS:   []iconf.DirPath{},
			Linux:   []iconf.DirPath{},
		},
		BinaryRename: &iconf.BinaryRename{
			NameFlag: "tinymist",
			RenameTo: "tinymist",
		},
	}
}

func TestTinyMist() {
	nn := NewTinyMist()
	nn.Start()

	ff := fmt.Sprintf(
		"/home/moqsien/projects/go/src/gvcgo/vcollector/test/%s.json",
		nn.SDKName,
	)
	content, _ := json.MarshalIndent(nn.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

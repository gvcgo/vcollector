package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gvcgo/vcollector/internal/gh"
	"github.com/gvcgo/vcollector/internal/iconf"
	"github.com/gvcgo/vcollector/pkgs/crawlers/gh/searcher"
	"github.com/gvcgo/vcollector/pkgs/version"
)

type Wasmer struct {
	SDKName  string
	RepoName string
	searcher.GhSearcher
}

func NewWasmer() (w *Wasmer) {
	w = &Wasmer{
		SDKName:  "wasmer",
		RepoName: "wasmerio/wasmer",
		GhSearcher: searcher.GhSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (w *Wasmer) GetSDKName() string {
	return w.SDKName
}

func (w *Wasmer) tagFilter(ri gh.ReleaseItem) bool {
	return searcher.GhVersionRegexp.FindString(ri.TagName) != ""
}

func (w *Wasmer) fileFilter(a gh.Asset) bool {
	if strings.Contains(a.Url, "archive/refs/") {
		return false
	}
	if strings.HasSuffix(a.Name, ".exe") {
		return false
	}
	if strings.Contains(a.Name, "-musl") {
		return false
	}
	if strings.Contains(a.Name, "-gnu64") {
		return false
	}
	return true
}

func (w *Wasmer) osParser(fName string) (osStr string) {
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

func (w *Wasmer) archParser(fName string) (archStr string) {
	if strings.Contains(fName, "-x86_64") {
		return "amd64"
	}
	if strings.Contains(fName, "-amd64") {
		return "amd64"
	}
	if strings.Contains(fName, "-aarch64") {
		return "arm64"
	}
	if strings.Contains(fName, "-arm64") {
		return "arm64"
	}
	return
}

func (w *Wasmer) vParser(tagName string) (vStr string) {
	return strings.TrimPrefix(tagName, "v")
}

func (w *Wasmer) insParser(fName string) (insStr string) {
	return version.Unarchiver
}

func (w *Wasmer) Start() {
	w.GhSearcher.Search(
		w.RepoName,
		w.tagFilter,
		w.fileFilter,
		w.vParser,
		w.archParser,
		w.osParser,
		w.insParser,
		nil,
	)
}

func (w *Wasmer) GetVersions() []byte {
	r, _ := w.Version.Marshal()
	return r
}

func (w *Wasmer) HomePage() string {
	return "https://wasmer.io/"
}

func (w *Wasmer) GetInstallConf() (ic iconf.InstallerConfig) {
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

func TestWasmer() {
	nn := NewWasmer()
	nn.Start()

	ff := fmt.Sprintf(
		"/home/moqsien/projects/go/src/gvcgo/vcollector/test/%s.json",
		nn.SDKName,
	)
	content, _ := json.MarshalIndent(nn.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

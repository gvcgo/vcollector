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

/*
https://github.com/WasmEdge/WasmEdge/releases
*/

type Wasmedge struct {
	SDKName  string
	RepoName string
	searcher.GhSearcher
}

func NewWasmedge() (w *Wasmedge) {
	w = &Wasmedge{
		SDKName:  "wasmedge",
		RepoName: "WasmEdge/WasmEdge",
		GhSearcher: searcher.GhSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (w *Wasmedge) GetSDKName() string {
	return w.SDKName
}

func (w *Wasmedge) tagFilter(ri gh.ReleaseItem) bool {
	return searcher.GVersionRegexp.FindString(ri.TagName) != ""
}

func (w *Wasmedge) fileFilter(a gh.Asset) bool {
	if strings.Contains(a.Url, "archive/refs/") {
		return false
	}
	if strings.HasPrefix(a.Name, "sbom") {
		return false
	}
	if strings.HasSuffix(a.Name, "static.tar.gz") {
		return false
	}
	if strings.Contains(a.Name, "-manylinux2014") {
		return false
	}
	if strings.Contains(a.Name, "-manylinux2010") {
		return false
	}
	if strings.HasSuffix(a.Name, ".exe") {
		return false
	}
	if strings.HasSuffix(a.Name, ".rpm") {
		return false
	}
	if strings.HasSuffix(a.Name, ".tar.xz") {
		return false
	}
	if strings.HasSuffix(a.Name, ".src.tar.gz") {
		return false
	}
	if strings.HasSuffix(a.Name, ".msi") {
		return false
	}
	if strings.Contains(a.Name, "-msvc") {
		return false
	}
	if strings.Contains(a.Name, "-plugin-") {
		return false
	}
	if strings.Contains(a.Name, "-runtime-only-") {
		return false
	}
	return true
}

func (w *Wasmedge) osParser(fName string) (osStr string) {
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

func (w *Wasmedge) archParser(fName string) (archStr string) {
	if strings.Contains(fName, "x86_64") {
		return "amd64"
	}
	if strings.Contains(fName, "windows") {
		return "amd64"
	}
	if strings.Contains(fName, "amd64") {
		return "amd64"
	}
	if strings.Contains(fName, "aarch64") {
		return "arm64"
	}
	if strings.Contains(fName, "arm64") {
		return "arm64"
	}
	return
}

func (w *Wasmedge) vParser(tagName string) (vStr string) {
	return strings.TrimPrefix(tagName, "v")
}

func (w *Wasmedge) insParser(fName string) (insStr string) {
	return version.Unarchiver
}

func (w *Wasmedge) Start() {
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

func (w *Wasmedge) GetVersions() []byte {
	r, _ := w.Version.Marshal()
	return r
}

func (w *Wasmedge) HomePage() string {
	return "https://wasmedge.org/"
}

func (w *Wasmedge) GetInstallConf() (ic iconf.InstallerConfig) {
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

func TestWasmedge() {
	nn := NewWasmedge()
	nn.Start()

	ff := fmt.Sprintf(
		"/home/moqsien/projects/go/src/gvcgo/vcollector/test/%s.json",
		nn.SDKName,
	)
	content, _ := json.MarshalIndent(nn.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

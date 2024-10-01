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
	crawler.RegisterCrawler(NewWasmtime())
}

type Wasmtime struct {
	SDKName  string
	RepoName string
	searcher.GhSearcher
}

func NewWasmtime() (w *Wasmtime) {
	w = &Wasmtime{
		SDKName:  "wasmtime",
		RepoName: "bytecodealliance/wasmtime",
		GhSearcher: searcher.GhSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (w *Wasmtime) GetSDKName() string {
	return w.SDKName
}

func (w *Wasmtime) tagFilter(ri gh.ReleaseItem) bool {
	return searcher.GhVersionRegexp.FindString(ri.TagName) != ""
}

func (w *Wasmtime) fileFilter(a gh.Asset) bool {
	if strings.Contains(a.Url, "archive/refs/") {
		return false
	}
	if strings.HasSuffix(a.Name, ".wasm") {
		return false
	}
	if strings.HasSuffix(a.Name, ".h") {
		return false
	}
	if strings.HasSuffix(a.Name, ".msi") {
		return false
	}
	if strings.Contains(a.Name, "-c-api") {
		return false
	}
	return true
}

func (w *Wasmtime) osParser(fName string) (osStr string) {
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

func (w *Wasmtime) archParser(fName string) (archStr string) {
	if strings.Contains(fName, "-x86_64") {
		return "amd64"
	}
	if strings.Contains(fName, "-aarch64") {
		return "arm64"
	}
	return
}

func (w *Wasmtime) vParser(tagName string) (vStr string) {
	return strings.TrimPrefix(tagName, "v")
}

func (w *Wasmtime) insParser(fName string) (insStr string) {
	return version.Unarchiver
}

func (w *Wasmtime) Start() {
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

func (w *Wasmtime) GetVersions() []byte {
	r, _ := w.Version.Marshal()
	return r
}

func (w *Wasmtime) HomePage() string {
	return "https://wasmtime.dev/"
}

func (w *Wasmtime) GetInstallConf() (ic iconf.InstallerConfig) {
	return iconf.InstallerConfig{
		FlagFiles: &iconf.FileItems{
			Windows: []string{"wasmtime.exe"},
			MacOS:   []string{"wasmtime"},
			Linux:   []string{"wasmtime"},
		},
		BinaryDirs: &iconf.DirItems{
			Windows: []iconf.DirPath{{}},
			MacOS:   []iconf.DirPath{{}},
			Linux:   []iconf.DirPath{{}},
		},
	}
}

func TestWasmTime() {
	nn := NewWasmtime()
	nn.Start()

	ff := fmt.Sprintf(
		"/home/moqsien/projects/go/src/gvcgo/vcollector/test/%s.json",
		nn.SDKName,
	)
	content, _ := json.MarshalIndent(nn.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

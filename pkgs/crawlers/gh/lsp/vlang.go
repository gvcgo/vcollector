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
	crawler.RegisterCrawler(NewVAnalyzer())
}

type VAnalyzer struct {
	SDKName  string
	RepoName string
	searcher.GhSearcher
}

func NewVAnalyzer() (v *VAnalyzer) {
	v = &VAnalyzer{
		SDKName:  "v-analyzer",
		RepoName: "vlang/v-analyzer",
		GhSearcher: searcher.GhSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (v *VAnalyzer) GetSDKName() string {
	return v.SDKName
}

func (v *VAnalyzer) tagFilter(ri gh.ReleaseItem) bool {
	if searcher.GVersionRegexp.FindString(ri.TagName) != "" {
		return true
	}
	if searcher.PreviewRegexp.FindString(ri.TagName) != "" {
		return true
	}
	return false
}

func (v *VAnalyzer) fileFilter(a gh.Asset) bool {
	if strings.HasSuffix(a.Name, "debug.zip") {
		return false
	}
	if strings.HasSuffix(a.Name, "dev.zip") {
		return false
	}
	if strings.HasSuffix(a.Name, ".vsix") {
		return false
	}
	return !strings.Contains(a.Url, "archive/refs/")
}

func (v *VAnalyzer) osParser(fName string) (osStr string) {
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

func (v *VAnalyzer) archParser(fName string) (archStr string) {
	if strings.Contains(fName, "-x86_64") {
		return "amd64"
	}
	if strings.Contains(fName, "-arm64") {
		return "arm64"
	}
	return
}

func (v *VAnalyzer) vParser(tagName string) (vStr string) {
	return strings.TrimPrefix(tagName, "v")
}

func (v *VAnalyzer) insParser(fName string) (insStr string) {
	return version.Unarchiver
}

func (v *VAnalyzer) Start() {
	v.GhSearcher.Search(
		v.RepoName,
		v.tagFilter,
		v.fileFilter,
		v.vParser,
		v.archParser,
		v.osParser,
		v.insParser,
		nil,
	)
}

func (v *VAnalyzer) GetVersions() []byte {
	r, _ := v.Version.Marshal()
	return r
}

func (v *VAnalyzer) HomePage() string {
	return "https://github.com/vlang/v-analyzer"
}

func (v *VAnalyzer) GetInstallConf() (ic iconf.InstallerConfig) {
	return iconf.InstallerConfig{
		FlagFiles: &iconf.FileItems{
			Windows: []string{"v-analyzer.exe"},
			MacOS:   []string{"v-analyzer"},
			Linux:   []string{"v-analyzer"},
		},
		FlagDirExcepted: true,
		BinaryDirs: &iconf.DirItems{
			Windows: []iconf.DirPath{},
			MacOS:   []iconf.DirPath{},
			Linux:   []iconf.DirPath{},
		},
	}
}

func TestVAnalyzer() {
	nn := NewVAnalyzer()
	nn.Start()

	ff := fmt.Sprintf(
		"/home/moqsien/projects/go/src/gvcgo/vcollector/test/%s.json",
		nn.SDKName,
	)
	content, _ := json.MarshalIndent(nn.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

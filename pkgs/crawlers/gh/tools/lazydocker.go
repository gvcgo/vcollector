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
	crawler.RegisterCrawler(NewLazydocker())
}

type Lazydocker struct {
	SDKName  string
	RepoName string
	searcher.GhSearcher
}

func NewLazydocker() (l *Lazydocker) {
	l = &Lazydocker{
		SDKName:  "lazydocker",
		RepoName: "jesseduffield/lazydocker",
		GhSearcher: searcher.GhSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (l *Lazydocker) GetSDKName() string {
	return l.SDKName
}

func (l *Lazydocker) tagFilter(ri gh.ReleaseItem) bool {
	return searcher.GhVersionRegexp.FindString(ri.TagName) != ""
}

func (l *Lazydocker) fileFilter(a gh.Asset) bool {
	if strings.Contains(a.Url, "archive/refs/") {
		return false
	}
	if strings.HasSuffix(a.Name, ".txt") {
		return false
	}
	return true
}

func (l *Lazydocker) osParser(fName string) (osStr string) {
	fName = strings.ToLower(fName)
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

func (l *Lazydocker) archParser(fName string) (archStr string) {
	if strings.Contains(fName, "x86_64") {
		return "amd64"
	}
	if strings.Contains(fName, "arm64") {
		return "arm64"
	}
	return
}

func (l *Lazydocker) vParser(tagName string) (vStr string) {
	return strings.TrimPrefix(tagName, "v")
}

func (l *Lazydocker) insParser(fName string) (insStr string) {
	return version.Unarchiver
}

func (l *Lazydocker) Start() {
	l.GhSearcher.Search(
		l.RepoName,
		l.tagFilter,
		l.fileFilter,
		l.vParser,
		l.archParser,
		l.osParser,
		l.insParser,
		nil,
	)
}

func (l *Lazydocker) GetVersions() []byte {
	r, _ := l.Version.Marshal()
	return r
}

func (l *Lazydocker) HomePage() string {
	return "https://github.com/jesseduffield/lazydocker"
}

func (l *Lazydocker) GetInstallConf() (ic iconf.InstallerConfig) {
	return iconf.InstallerConfig{
		FlagFiles: &iconf.FileItems{
			Windows: []string{"README.md", "LICENSE"},
			MacOS:   []string{"README.md", "LICENSE"},
			Linux:   []string{"README.md", "LICENSE"},
		},
		FlagDirExcepted: true,
		BinaryDirs: &iconf.DirItems{
			Windows: []iconf.DirPath{},
			MacOS:   []iconf.DirPath{},
			Linux:   []iconf.DirPath{},
		},
	}
}

func TestLazydocker() {
	nn := NewLazydocker()
	nn.Start()

	ff := fmt.Sprintf(
		"/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/%s.json",
		nn.SDKName,
	)
	content, _ := json.MarshalIndent(nn.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

package lans

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/gvcgo/vcollector/internal/gh"
	"github.com/gvcgo/vcollector/internal/iconf"
	"github.com/gvcgo/vcollector/pkgs/crawlers/gh/searcher"
	"github.com/gvcgo/vcollector/pkgs/version"
)

/*
https://github.com/mstorsjo/llvm-mingw/releases
*/
var (
	dateVersion = regexp.MustCompile(`\d{8}`)
)

type LlvmMingw struct {
	SDKName  string
	RepoName string
	searcher.GhSearcher
}

func NewLlvmMingw() (l *LlvmMingw) {
	l = &LlvmMingw{
		SDKName:  "llvm-mingw",
		RepoName: "mstorsjo/llvm-mingw",
		GhSearcher: searcher.GhSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (l *LlvmMingw) GetSDKName() string {
	return l.SDKName
}

func (l *LlvmMingw) tagFilter(ri gh.ReleaseItem) bool {
	return dateVersion.FindString(ri.TagName) != ""
}

func (l *LlvmMingw) fileFilter(aa gh.Asset) bool {
	if strings.Contains(aa.Url, "archive/refs/") {
		return false
	}
	if strings.Contains(aa.Url, "-i686") {
		return false
	}
	if strings.Contains(aa.Url, "-msvcrt-ubuntu") {
		return false
	}
	if strings.Contains(aa.Url, "-ucrt-") && !strings.Contains(aa.Url, "ubuntu") && !strings.Contains(aa.Url, "macos") {
		return false
	}
	return true
}

func (l *LlvmMingw) osParser(fName string) (osStr string) {
	if strings.Contains(fName, "-ucrt-ubuntu") {
		return "linux"
	}
	if strings.Contains(fName, "-ucrt-macos") {
		return "darwin"
	}
	return "windows"
}

func (l *LlvmMingw) archParser(fName string) (archStr string) {
	if strings.Contains(fName, "-x86_64") {
		return "amd64"
	}
	if strings.Contains(fName, "-aarch64") {
		return "arm64"
	}
	if strings.Contains(fName, "-universal") {
		return "any"
	}
	return
}

func (l *LlvmMingw) vParser(tagName string) (vStr string) {
	return dateVersion.FindString(tagName)
}

func (l *LlvmMingw) insParser(fName string) (insStr string) {
	return version.Unarchiver
}

func (l *LlvmMingw) Start() {
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

func (l *LlvmMingw) GetVersions() []byte {
	r, _ := l.Version.Marshal()
	return r
}

func (l *LlvmMingw) HomePage() string {
	return "https://github.com/mstorsjo/llvm-mingw"
}

func (l *LlvmMingw) GetInstallConf() (ic iconf.InstallerConfig) {
	return iconf.InstallerConfig{
		FlagFiles: &iconf.FileItems{
			Windows: []string{"bin", "lib", "include"},
			MacOS:   []string{"bin", "lib", "include"},
			Linux:   []string{"bin", "lib", "include"},
		},
		FlagDirExcepted: false,
		BinaryDirs: &iconf.DirItems{
			Windows: []iconf.DirPath{{"bin"}},
			MacOS:   []iconf.DirPath{{"bin"}},
			Linux:   []iconf.DirPath{{"bin"}},
		},
	}
}

func TestLlvmMingw() {
	nn := NewLlvmMingw()
	nn.Start()

	ff := fmt.Sprintf(
		"/home/moqsien/projects/go/src/gvcgo/vcollector/test/%s.json",
		nn.SDKName,
	)
	content, _ := json.MarshalIndent(nn.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

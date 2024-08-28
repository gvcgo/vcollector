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

/*
https://github.com/bazelbuild/bazel/releases
*/
func init() {
	crawler.RegisterCrawler(NewBazel())
}

type Bazel struct {
	SDKName  string
	RepoName string
	searcher.GhSearcher
}

func NewBazel() (b *Bazel) {
	b = &Bazel{
		SDKName:  "bazel",
		RepoName: "bazelbuild/bazel",
		GhSearcher: searcher.GhSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (b *Bazel) GetSDKName() string {
	return b.SDKName
}

func (b *Bazel) tagFilter(ri gh.ReleaseItem) bool {
	return searcher.GVersionRegexp.FindString(ri.TagName) != ""
}

func (b *Bazel) fileFilter(aa gh.Asset) bool {
	if strings.Contains(aa.Url, "archive/refs/") {
		return false
	}
	if !strings.Contains(aa.Url, "nojdk") {
		return false
	}
	if strings.HasSuffix(aa.Name, ".sha256") {
		return false
	}
	if strings.HasSuffix(aa.Name, ".sig") {
		return false
	}
	if strings.HasSuffix(aa.Name, ".zip") {
		return false
	}
	if strings.HasSuffix(aa.Name, ".sh") {
		return false
	}
	if strings.HasSuffix(aa.Name, ".sig") {
		return false
	}
	return true
}

func (b *Bazel) osParser(fName string) (osStr string) {
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

func (b *Bazel) archParser(fName string) (archStr string) {
	if strings.Contains(fName, "-x86_64") {
		return "amd64"
	}
	if strings.Contains(fName, "-arm64") {
		return "arm64"
	}
	return
}

func (b *Bazel) vParser(tagName string) (vStr string) {
	return strings.TrimPrefix(tagName, "v")
}

func (b *Bazel) insParser(fName string) (insStr string) {
	return version.Executable
}

func (b *Bazel) Start() {
	b.GhSearcher.Search(
		b.RepoName,
		b.tagFilter,
		b.fileFilter,
		b.vParser,
		b.archParser,
		b.osParser,
		b.insParser,
		nil,
	)
}

func (b *Bazel) GetVersions() []byte {
	r, _ := b.Version.Marshal()
	return r
}

func (b *Bazel) HomePage() string {
	return "https://github.com/bazelbuild/bazel"
}

func (b *Bazel) GetInstallConf() (ic iconf.InstallerConfig) {
	return iconf.InstallerConfig{
		FlagFiles: &iconf.FileItems{
			Windows: []string{"bazel.exe"},
			MacOS:   []string{"bazel"},
			Linux:   []string{"bazel"},
		},
		FlagDirExcepted: true,
		BinaryDirs: &iconf.DirItems{
			Windows: []iconf.DirPath{},
			MacOS:   []iconf.DirPath{},
			Linux:   []iconf.DirPath{},
		},
		BinaryRename: &iconf.BinaryRename{
			NameFlag: "bazel",
			RenameTo: "bazel",
		},
	}
}

func TestBazel() {
	nn := NewBazel()
	nn.Start()

	ff := fmt.Sprintf(
		"/home/moqsien/projects/go/src/gvcgo/vcollector/test/%s.json",
		nn.SDKName,
	)
	content, _ := json.MarshalIndent(nn.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

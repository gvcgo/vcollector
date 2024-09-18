package lans

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
	crawler.RegisterCrawler(NewKotlin())
}

type Kotlin struct {
	SDKName  string
	RepoName string
	searcher.GhSearcher
}

func NewKotlin() (k *Kotlin) {
	k = &Kotlin{
		SDKName:  "kotlin",
		RepoName: "JetBrains/kotlin",
		GhSearcher: searcher.GhSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (k *Kotlin) GetSDKName() string {
	return k.SDKName
}

func (k *Kotlin) tagFilter(ri gh.ReleaseItem) bool {
	return searcher.GhVersionRegexp.FindString(ri.TagName) != ""
}

func (k *Kotlin) fileFilter(a gh.Asset) bool {
	if strings.Contains(a.Url, "archive/refs/") {
		return false
	}
	if strings.HasPrefix(a.Name, "maven-") {
		return false
	}
	if strings.HasSuffix(a.Name, ".sha256") {
		return false
	}
	if strings.HasSuffix(a.Name, ".json") {
		return false
	}
	return true
}

func (k *Kotlin) osParser(fName string) (osStr string) {
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

func (k *Kotlin) osParserCompiler(fName string) (osStr string) {
	if strings.Contains(fName, "kotlin-compiler-") {
		return "any"
	}
	return
}

func (k *Kotlin) archParser(fName string) (archStr string) {
	if strings.Contains(fName, "-x86_64") {
		return "amd64"
	}
	if strings.Contains(fName, "-aarch64") {
		return "arm64"
	}
	return
}

func (k *Kotlin) archParseCompiler(fName string) (archStr string) {
	if strings.Contains(fName, "kotlin-compiler-") {
		return "any"
	}
	return
}

func (k *Kotlin) vParser(tagName string) (vStr string) {
	return "native-" + strings.TrimPrefix(tagName, "v")
}

func (k *Kotlin) vParserCompiler(tagName string) (vStr string) {
	return strings.TrimPrefix(tagName, "v")
}

func (k *Kotlin) insParser(fName string) (insStr string) {
	return version.Unarchiver
}

func (k *Kotlin) Start() {
	k.GhSearcher.Search(
		k.RepoName,
		k.tagFilter,
		k.fileFilter,
		k.vParser,
		k.archParser,
		k.osParser,
		k.insParser,
		nil,
	)
	k.GhSearcher.Search(
		k.RepoName,
		k.tagFilter,
		k.fileFilter,
		k.vParserCompiler,
		k.archParseCompiler,
		k.osParserCompiler,
		k.insParser,
		nil,
	)
}

func (k *Kotlin) GetVersions() []byte {
	r, _ := k.Version.Marshal()
	return r
}

func (k *Kotlin) HomePage() string {
	return "https://kotlinlang.org/"
}

func (k *Kotlin) GetInstallConf() (ic iconf.InstallerConfig) {
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

func TestKotlin() {
	nn := NewKotlin()
	nn.Start()

	ff := fmt.Sprintf(
		"/home/moqsien/projects/go/src/gvcgo/vcollector/test/%s.json",
		nn.SDKName,
	)
	content, _ := json.MarshalIndent(nn.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

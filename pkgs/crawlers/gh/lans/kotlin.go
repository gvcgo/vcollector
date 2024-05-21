package lans

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gvcgo/vcollector/internal/gh"
	"github.com/gvcgo/vcollector/pkgs/crawlers/gh/searcher"
	"github.com/gvcgo/vcollector/pkgs/version"
)

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
	return GhVersionRegexp.FindString(ri.TagName) != ""
}

func (k *Kotlin) fileFilter(a gh.Asset) bool {
	if strings.Contains(a.Url, "archive/refs/") {
		return false
	}
	if strings.HasPrefix(a.Name, "kotlin-compiler-") {
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

func (k *Kotlin) archParser(fName string) (archStr string) {
	if strings.Contains(fName, "-x86_64") {
		return "amd64"
	}
	if strings.Contains(fName, "-aarch64") {
		return "arm64"
	}
	return
}

func (k *Kotlin) vParser(tagName string) (vStr string) {
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
}

func TestKotlin() {
	nn := NewKotlin()
	nn.Start()

	ff := fmt.Sprintf(
		"/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/%s.json",
		nn.SDKName,
	)
	content, _ := json.MarshalIndent(nn.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

package lans

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gvcgo/vcollector/internal/gh"
	"github.com/gvcgo/vcollector/internal/iconf"
	"github.com/gvcgo/vcollector/internal/req"
	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
	"github.com/gvcgo/vcollector/pkgs/crawlers/gh/searcher"
	"github.com/gvcgo/vcollector/pkgs/version"
)

func init() {
	crawler.RegisterCrawler(NewBun())
}

type Bun struct {
	SDKName  string
	RepoName string
	searcher.GhSearcher
}

func NewBun() (b *Bun) {
	b = &Bun{
		SDKName:  "bun",
		RepoName: "oven-sh/bun",
		GhSearcher: searcher.GhSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (b *Bun) GetSDKName() string {
	return b.SDKName
}

func (b *Bun) tagFilter(ri gh.ReleaseItem) bool {
	return searcher.GhVersionRegexp.FindString(ri.TagName) != ""
}

func (b *Bun) fileFilter(a gh.Asset) bool {
	if strings.Contains(a.Url, "archive/refs/") {
		return false
	}
	if strings.Contains(a.Name, "profile") {
		return false
	}
	if strings.Contains(a.Name, "baseline") {
		return false
	}
	if strings.HasSuffix(a.Name, ".txt") {
		return false
	}
	if strings.HasSuffix(a.Name, ".txt.asc") {
		return false
	}
	return true
}

func (b *Bun) osParser(fName string) (osStr string) {
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

func (b *Bun) archParser(fName string) (archStr string) {
	if strings.Contains(fName, "-x64") {
		return "amd64"
	}
	if strings.Contains(fName, "-aarch64") {
		return "arm64"
	}
	return
}

func (b *Bun) vParser(tagName string) (vStr string) {
	return strings.TrimPrefix(searcher.GhVersionRegexp.FindString(tagName), "v")
}

func (b *Bun) insParser(fName string) (insStr string) {
	return version.Unarchiver
}

func (b *Bun) SumGetter(fName string, assets []gh.Asset) (sum, sumType string) {
	for _, a := range assets {
		if strings.HasSuffix(a.Name, "SHASUMS256.txt") {
			content := req.GetResp("https://gvc.1710717.xyz/proxy/"+a.Url, 1)
			for _, line := range strings.Split(content, "\n") {
				if strings.Contains(line, fName) {
					sList := strings.Split(line, fName)
					sum = strings.TrimSpace(sList[0])
					sumType = "sha256"
					return
				}
			}
		}
	}
	return
}

func (b *Bun) Start() {
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

func (b *Bun) GetVersions() []byte {
	r, _ := b.Version.Marshal()
	return r
}

func (b *Bun) HomePage() string {
	return "https://bun.sh/"
}

func (b *Bun) GetInstallConf() (ic iconf.InstallerConfig) {
	return iconf.InstallerConfig{
		FlagFiles: &iconf.FileItems{
			Windows: []string{"bun.exe"},
			MacOS:   []string{"bun"},
			Linux:   []string{"bun"},
		},
		FlagDirExcepted: true,
		BinaryDirs: &iconf.DirItems{
			Windows: []iconf.DirPath{},
			MacOS:   []iconf.DirPath{},
			Linux:   []iconf.DirPath{},
		},
	}
}

func TestBun() {
	bb := NewBun()
	bb.Start()

	ff := fmt.Sprintf(
		"/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/%s.json",
		bb.SDKName,
	)
	content, _ := json.MarshalIndent(bb.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gvcgo/vcollector/internal/gh"
	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
	"github.com/gvcgo/vcollector/pkgs/crawlers/gh/searcher"
	"github.com/gvcgo/vcollector/pkgs/version"
)

func init() {
	crawler.RegisterCrawler(NewProtoc())
}

type Protoc struct {
	SDKName  string
	RepoName string
	searcher.GhSearcher
}

func NewProtoc() (p *Protoc) {
	p = &Protoc{
		SDKName:  "protoc",
		RepoName: "protocolbuffers/protobuf",
		GhSearcher: searcher.GhSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (p *Protoc) GetSDKName() string {
	return p.SDKName
}

func (p *Protoc) tagFilter(ri gh.ReleaseItem) bool {
	return searcher.VVersionRegexp.FindString(ri.TagName) != ""
}

func (p *Protoc) fileFilter(a gh.Asset) bool {
	if strings.Contains(a.Url, "archive/refs/") {
		return false
	}
	if strings.HasPrefix(a.Name, "protobuf-") {
		return false
	}
	return true
}

func (p *Protoc) osParser(fName string) (osStr string) {
	if strings.Contains(fName, "-osx") {
		return "darwin"
	}
	if strings.Contains(fName, "linux") {
		return "linux"
	}
	if strings.Contains(fName, "-win") {
		return "windows"
	}
	return
}

func (p *Protoc) archParser(fName string) (archStr string) {
	if strings.Contains(fName, "-x86_64") {
		return "amd64"
	}
	if strings.Contains(fName, "-win64") {
		return "amd64"
	}
	if strings.Contains(fName, "-aarch_64") {
		return "arm64"
	}
	return
}

func (p *Protoc) vParser(tagName string) (vStr string) {
	return strings.TrimPrefix(tagName, "v")
}

func (p *Protoc) insParser(fName string) (insStr string) {
	return version.Unarchiver
}

func (p *Protoc) Start() {
	p.GhSearcher.Search(
		p.RepoName,
		p.tagFilter,
		p.fileFilter,
		p.vParser,
		p.archParser,
		p.osParser,
		p.insParser,
		nil,
	)
}

func (p *Protoc) GetVersions() []byte {
	r, _ := p.Version.Marshal()
	return r
}

func TestProtoc() {
	nn := NewProtoc()
	nn.Start()

	ff := fmt.Sprintf(
		"/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/%s.json",
		nn.SDKName,
	)
	content, _ := json.MarshalIndent(nn.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

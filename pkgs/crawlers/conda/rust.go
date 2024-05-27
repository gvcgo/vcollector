package conda

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gvcgo/vcollector/internal/iconf"
	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
	"github.com/gvcgo/vcollector/pkgs/version"
)

func init() {
	crawler.RegisterCrawler(NewRust())
}

type Rust struct {
	SDKName string
	CondaSearcher
}

func NewRust() (r *Rust) {
	r = &Rust{
		SDKName: "rust",
		CondaSearcher: CondaSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (r *Rust) GetSDKName() string {
	return r.SDKName
}

func (r *Rust) Start() {
	r.CondaSearcher.Search(r.SDKName)
}

func (r *Rust) GetVersions() []byte {
	result, _ := r.Version.Marshal()
	return result
}

func (r *Rust) HomePage() string {
	return "https://www.rust-lang.org/"
}

func (r *Rust) GetInstallConf() (ic iconf.InstallerConfig) {
	return iconf.InstallerConfig{
		BinaryDirs: &iconf.DirItems{
			Windows: []iconf.DirPath{{"cargo", "bin"}, {"Library", "bin"}},
			MacOS:   []iconf.DirPath{{"bin"}},
			Linux:   []iconf.DirPath{{"bin"}},
		},
	}
}

func TestRust() {
	rr := NewRust()
	rr.Start()

	ff := fmt.Sprintf(
		"/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/%s.json",
		rr.SDKName,
	)
	content, _ := json.MarshalIndent(rr.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

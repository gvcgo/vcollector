package conda

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
	"github.com/gvcgo/vcollector/pkgs/version"
)

func init() {
	crawler.RegisterCrawler(NewPerl())
}

type Perl struct {
	SDKName string
	CondaSearcher
}

func NewPerl() (p *Perl) {
	p = &Perl{
		SDKName: "perl",
		CondaSearcher: CondaSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (p *Perl) GetSDKName() string {
	return p.SDKName
}

func (p *Perl) Start() {
	p.CondaSearcher.Search(p.SDKName)
}

func (p *Perl) GetVersions() []byte {
	r, _ := p.Version.Marshal()
	return r
}

func TestPerl() {
	pp := NewPerl()
	pp.Start()

	ff := fmt.Sprintf(
		"/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/%s.json",
		pp.SDKName,
	)
	content, _ := json.MarshalIndent(pp.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

package conda

import (
	"encoding/json"
	"os"

	"github.com/gvcgo/vcollector/internal/iconf"
	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
	"github.com/gvcgo/vcollector/pkgs/version"
)

func init() {
	crawler.RegisterCrawler(NewPostgreSQL())
	crawler.RegisterCondaCrawler(NewPostgreSQL())
}

type PostgreSQL struct {
	SDKName string
	CondaSearcher
}

func NewPostgreSQL() (p *PostgreSQL) {
	p = &PostgreSQL{
		SDKName: "postgresql",
		CondaSearcher: CondaSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (p *PostgreSQL) GetSDKName() string {
	return p.SDKName
}

func (p *PostgreSQL) Start() {
	p.CondaSearcher.Search(p.SDKName)
}

func (p *PostgreSQL) GetVersions() []byte {
	r, _ := p.Version.Marshal()
	return r
}

func (p *PostgreSQL) HomePage() string {
	return "https://www.postgresql.org/"
}

func (p *PostgreSQL) GetInstallConf() (ic iconf.InstallerConfig) {
	return iconf.InstallerConfig{
		BinaryDirs: &iconf.DirItems{
			Windows: []iconf.DirPath{{"Library", "bin"}},
			MacOS:   []iconf.DirPath{{"bin"}},
			Linux:   []iconf.DirPath{{"bin"}},
		},
	}
}

func TestPostgreSQL() {
	cc := NewPostgreSQL()
	cc.Start()

	ff := "/home/moqsien/projects/go/src/gvcgo/vcollector/test/postgresql.json"
	content, _ := json.MarshalIndent(cc.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

package conda

import (
	"encoding/json"
	"os"

	"github.com/gvcgo/vcollector/internal/iconf"
	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
	"github.com/gvcgo/vcollector/pkgs/version"
)

func init() {
	crawler.RegisterCrawler(NewMySQL())
	crawler.RegisterCondaCrawler(NewMySQL())
}

type MySQL struct {
	SDKName string
	CondaSearcher
}

func NewMySQL() (m *MySQL) {
	m = &MySQL{
		SDKName: "mysql",
		CondaSearcher: CondaSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (m *MySQL) GetSDKName() string {
	return m.SDKName
}

func (m *MySQL) Start() {
	m.CondaSearcher.Search(m.SDKName)
}

func (m *MySQL) GetVersions() []byte {
	r, _ := m.Version.Marshal()
	return r
}

func (m *MySQL) HomePage() string {
	return "https://www.mysql.com/"
}

func (m *MySQL) GetInstallConf() (ic iconf.InstallerConfig) {
	return iconf.InstallerConfig{
		BinaryDirs: &iconf.DirItems{
			Windows: []iconf.DirPath{{"Library", "bin"}},
			MacOS:   []iconf.DirPath{{"bin"}},
			Linux:   []iconf.DirPath{{"bin"}},
		},
	}
}

func TestMySQL() {
	cc := NewMySQL()
	cc.Start()

	ff := "/home/moqsien/projects/go/src/gvcgo/vcollector/test/mysql.json"
	content, _ := json.MarshalIndent(cc.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

package conda

import (
	"encoding/json"
	"os"

	"github.com/gvcgo/vcollector/internal/iconf"
	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
	"github.com/gvcgo/vcollector/pkgs/version"
)

func init() {
	crawler.RegisterCrawler(NewMongoDB())
	crawler.RegisterCondaCrawler(NewMongoDB())
}

type MongoDB struct {
	SDKName string
	CondaSearcher
}

func NewMongoDB() (m *MongoDB) {
	m = &MongoDB{
		SDKName: "mongodb",
		CondaSearcher: CondaSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (m *MongoDB) GetSDKName() string {
	return m.SDKName
}

func (m *MongoDB) Start() {
	m.CondaSearcher.Search(m.SDKName)
}

func (m *MongoDB) GetVersions() []byte {
	r, _ := m.Version.Marshal()
	return r
}

func (m *MongoDB) HomePage() string {
	return "https://www.mongodb.com/"
}

func (m *MongoDB) GetInstallConf() (ic iconf.InstallerConfig) {
	return iconf.InstallerConfig{
		BinaryDirs: &iconf.DirItems{
			Windows: []iconf.DirPath{{"Library", "bin"}},
			MacOS:   []iconf.DirPath{{"bin"}, {"sbin"}},
			Linux:   []iconf.DirPath{{"bin"}, {"sbin"}},
		},
	}
}

func TestMongoDB() {
	cc := NewMongoDB()
	cc.Start()

	ff := "/home/moqsien/projects/go/src/gvcgo/vcollector/test/mongodb.json"
	content, _ := json.MarshalIndent(cc.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

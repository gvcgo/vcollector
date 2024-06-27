package conda

import (
	"github.com/gvcgo/vcollector/internal/iconf"
	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
	"github.com/gvcgo/vcollector/pkgs/version"
)

func init() {
	crawler.RegisterCrawler(NewFFMPeg())
	crawler.RegisterCondaCrawler(NewFFMPeg())
}

type FFMPeg struct {
	SDKName string
	CondaSearcher
}

func NewFFMPeg() (f *FFMPeg) {
	f = &FFMPeg{
		SDKName: "ffmpeg",
		CondaSearcher: CondaSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (f *FFMPeg) GetSDKName() string {
	return f.SDKName
}

func (f *FFMPeg) Start() {
	f.CondaSearcher.Search(f.SDKName)
}

func (f *FFMPeg) GetVersions() []byte {
	r, _ := f.Version.Marshal()
	return r
}

func (f *FFMPeg) HomePage() string {
	return "https://ffmpeg.org/"
}

func (f *FFMPeg) GetInstallConf() (ic iconf.InstallerConfig) {
	return iconf.InstallerConfig{
		BinaryDirs: &iconf.DirItems{
			Windows: []iconf.DirPath{{"Library", "bin"}},
			MacOS:   []iconf.DirPath{{"bin"}},
			Linux:   []iconf.DirPath{{"bin"}},
		},
	}
}

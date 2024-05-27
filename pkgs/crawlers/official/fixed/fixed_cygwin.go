package fixed

import (
	"github.com/gvcgo/vcollector/internal/iconf"
	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
	"github.com/gvcgo/vcollector/pkgs/version"
)

func init() {
	crawler.RegisterCrawler(NewCygwin())
}

type Cygwin struct {
	SDKName string
	Version version.VersionList
}

func NewCygwin() *Cygwin {
	return &Cygwin{
		SDKName: "cygwin",
		Version: make(version.VersionList),
	}
}

func (c *Cygwin) GetSDKName() string {
	return c.SDKName
}

func (c *Cygwin) Start() {
	c.Version["latest"] = version.Version{
		version.Item{
			Url:       "https://cygwin.com/setup-x86_64.exe",
			Arch:      "amd64",
			Os:        "windows",
			Installer: version.Executable,
		},
	}
}

func (c *Cygwin) GetVersions() []byte {
	result, _ := c.Version.Marshal()
	return result
}

func (c *Cygwin) HomePage() string {
	return "https://cygwin.com/"
}

func (c *Cygwin) GetInstallConf() (ic iconf.InstallerConfig) {
	return iconf.InstallerConfig{
		FlagFiles: &iconf.FileItems{
			Windows: []string{"setup-x86_64.exe"},
			MacOS:   []string{"setup-x86_64.exe"},
			Linux:   []string{"setup-x86_64.exe"},
		},
		BinaryRename: &iconf.BinaryRename{
			NameFlag: "setup-x86_64",
			RenameTo: "cygwin-installer",
		},
		BinaryDirs: &iconf.DirItems{
			Windows: []iconf.DirPath{},
			MacOS:   []iconf.DirPath{},
			Linux:   []iconf.DirPath{},
		},
	}
}

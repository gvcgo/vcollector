package fixed

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gvcgo/vcollector/internal/iconf"
	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
	"github.com/gvcgo/vcollector/pkgs/version"
)

func init() {
	crawler.RegisterCrawler(NewMoonBit())
}

var MoonBitVersions []string = []string{
	"latest",
}

var MoonBitCoreURLs []string = []string{
	"https://cli.moonbitlang.com/binaries/%s/moonbit-darwin-x86_64.tar.gz",
	"https://cli.moonbitlang.com/binaries/%s/moonbit-darwin-aarch64.tar.gz",
	"https://cli.moonbitlang.com/binaries/%s/moonbit-linux-x86_64.tar.gz",
	"https://cli.moonbitlang.com/binaries/%s/moonbit-windows-x86_64.zip",
}

var MoonBitLibURL string = "https://cli.moonbitlang.com/cores/core-%s.tar.gz"

type MoonBit struct {
	SDKName string
	Version version.VersionList
}

func NewMoonBit() (m *MoonBit) {
	m = &MoonBit{
		SDKName: "moonbit",
		Version: make(version.VersionList),
	}
	return
}

func (m *MoonBit) GetSDKName() string {
	return m.SDKName
}

func (m *MoonBit) osParser(coreUrl string) (osStr string) {
	if strings.Contains(coreUrl, "darwin") {
		osStr = "darwin"
	} else if strings.Contains(coreUrl, "linux") {
		osStr = "linux"
	} else if strings.Contains(coreUrl, "windows") {
		osStr = "windows"
	}
	return
}

func (m *MoonBit) archParser(coreUrl string) (archStr string) {
	if strings.Contains(coreUrl, "aarch64") {
		archStr = "arm64"
	} else if strings.Contains(coreUrl, "x86_64") {
		archStr = "amd64"
	}
	return
}

func (m *MoonBit) Start() {
	for _, vv := range MoonBitVersions {
		for _, uu := range MoonBitCoreURLs {
			dUrl := fmt.Sprintf(uu, vv)
			if m.Version[vv] == nil {
				m.Version[vv] = version.Version{}
			}
			osStr := m.osParser(dUrl)
			archStr := m.archParser(dUrl)
			if osStr == "" || archStr == "" {
				continue
			}
			m.Version[vv] = append(m.Version[vv], version.Item{
				Url:   dUrl,
				Os:    osStr,
				Arch:  archStr,
				Extra: fmt.Sprintf(MoonBitLibURL, vv),
			})
		}
	}
}

func (m *MoonBit) GetVersions() []byte {
	r, _ := m.Version.Marshal()
	return r
}

func (m *MoonBit) HomePage() string {
	return "https://www.moonbitlang.com/"
}

func (m *MoonBit) GetInstallConf() (ic iconf.InstallerConfig) {
	return iconf.InstallerConfig{
		FlagFiles: &iconf.FileItems{
			Windows: []string{"moon.exe", "moonc.exe"},
			MacOS:   []string{"moon", "moonc"},
			Linux:   []string{"moon", "moonc"},
		},
		FlagDirExcepted: true,
		BinaryDirs: &iconf.DirItems{
			Windows: []iconf.DirPath{},
			MacOS:   []iconf.DirPath{},
			Linux:   []iconf.DirPath{},
		},
	}
}

func TestMoonBit() {
	vv := NewMoonBit()
	vv.Start()

	ff := "/home/moqsien/projects/go/src/gvcgo/vcollector/test/moonbit.json"
	content, _ := json.MarshalIndent(vv.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

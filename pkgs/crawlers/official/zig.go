package official

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gvcgo/vcollector/internal/iconf"
	"github.com/gvcgo/vcollector/internal/req"
	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
	"github.com/gvcgo/vcollector/pkgs/version"
)

func init() {
	crawler.RegisterCrawler(NewZig())
}

/*
https://ziglang.org/download/
*/
type Zig struct {
	DownloadUrl string
	SDKName     string
	Version     version.VersionList
}

func NewZig() (z *Zig) {
	z = &Zig{
		DownloadUrl: "https://ziglang.org/download/",
		SDKName:     "zig",
		Version:     make(version.VersionList),
	}
	return
}

func (z *Zig) GetSDKName() string {
	return z.SDKName
}

func (z *Zig) parshArch(dUrl string) (arch string) {
	if strings.Contains(dUrl, "x86_64") {
		return "amd64"
	}
	if strings.Contains(dUrl, "aarch64") {
		return "arm64"
	}
	return
}

func (z *Zig) parseOS(dUrl string) (platform string) {
	if strings.Contains(dUrl, "windows") {
		return "windows"
	}
	if strings.Contains(dUrl, "linux") {
		return "linux"
	}
	if strings.Contains(dUrl, "macos") {
		return "darwin"
	}
	return
}

func (z *Zig) getResult() {
	doc := req.GetDocument(z.DownloadUrl)
	if doc == nil {
		return
	}
	doc.Find("h2").Each(func(i int, s *goquery.Selection) {
		vName := strings.TrimSpace(s.Text())
		doc.Find("table").Eq(i).Find("tr").Each(func(_ int, ss *goquery.Selection) {
			thStr := strings.ToLower(strings.TrimSpace(ss.Find("th").Text()))
			tdList := ss.Find("td").Nodes
			if thStr == "os" || len(tdList) < 4 {
				return
			}

			item := version.Item{}
			item.Url = ss.Find("td").Eq(1).Find("a").AttrOr("href", "")
			if item.Url == "" {
				return
			}
			item.Arch = z.parshArch(item.Url)
			item.Os = z.parseOS(item.Url)
			if item.Arch == "" || item.Os == "" {
				return
			}
			item.Installer = version.Unarchiver
			if _, ok := z.Version[vName]; !ok {
				z.Version[vName] = version.Version{}
			}
			z.Version[vName] = append(z.Version[vName], item)
		})
	})
}

func (z *Zig) Start() {
	z.getResult()
}

func (z *Zig) GetVersions() []byte {
	r, _ := z.Version.Marshal()
	return r
}

func (z *Zig) HomePage() string {
	return "https://ziglang.org/"
}

func (z *Zig) GetInstallConf() (ic iconf.InstallerConfig) {
	return iconf.InstallerConfig{
		FlagFiles: &iconf.FileItems{
			Windows: []string{"LICENSE"},
			MacOS:   []string{"LICENSE"},
			Linux:   []string{"LICENSE"},
		},
		BinaryDirs: &iconf.DirItems{
			Windows: []iconf.DirPath{},
			MacOS:   []iconf.DirPath{},
			Linux:   []iconf.DirPath{},
		},
	}
}

func TestZig() {
	zz := NewZig()
	zz.Start()

	ff := "/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/zig.json"
	content, _ := json.MarshalIndent(zz.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

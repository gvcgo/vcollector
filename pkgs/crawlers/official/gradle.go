package official

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gvcgo/vcollector/internal/iconf"
	"github.com/gvcgo/vcollector/internal/req"
	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
	"github.com/gvcgo/vcollector/pkgs/version"
)

func init() {
	crawler.RegisterCrawler(NewGradle())
}

var GradleExludeList = map[string]struct{}{
	"v0.7": {},
}

/*
https://gradle.org/releases/
https://gradle.org/release-checksums/
*/
type Gradle struct {
	DownloadUrl string
	SDKName     string
	Version     version.VersionList
	urlPattern  string
}

func NewGradle() (g *Gradle) {
	g = &Gradle{
		DownloadUrl: "https://gradle.org/release-checksums/",
		SDKName:     "gradle",
		Version:     make(version.VersionList),
		urlPattern:  "https://services.gradle.org/distributions/gradle-%s-all.zip",
	}
	return
}

func (g *Gradle) GetSDKName() string {
	return g.SDKName
}

func (g *Gradle) getResult() {
	doc := req.GetDocument(g.DownloadUrl)
	if doc == nil {
		return
	}
	doc.Find("h3.u-text-with-icon").Each(func(_ int, ss *goquery.Selection) {
		vStr := strings.TrimSpace(ss.Find("span").Eq(1).Text())
		if strings.HasPrefix(vStr, "v") {
			if _, ok := GradleExludeList[vStr]; ok {
				return
			}
			vStr = strings.TrimPrefix(vStr, "v")
			item := version.Item{}
			item.Arch = "any"
			item.Os = "any"
			item.Url = fmt.Sprintf(g.urlPattern, vStr)
			item.SumType = "sha256"
			item.Sum = ss.Next().Find("li").Eq(1).Find("code").Text()
			item.Installer = version.Unarchiver
			if _, ok := g.Version[vStr]; !ok {
				g.Version[vStr] = version.Version{}
			}
			g.Version[vStr] = append(g.Version[vStr], item)
		}
	})
}

func (g *Gradle) Start() {
	g.getResult()
}

func (g *Gradle) GetVersions() []byte {
	r, _ := g.Version.Marshal()
	return r
}

func (g *Gradle) HomePage() string {
	return "https://gradle.org/"
}

func (g *Gradle) GetInstallConf() (ic iconf.InstallerConfig) {
	return iconf.InstallerConfig{
		FlagFiles: &iconf.FileItems{
			Windows: []string{"LICENSE"},
			MacOS:   []string{"LICENSE"},
			Linux:   []string{"LICENSE"},
		},
		BinaryDirs: &iconf.DirItems{
			Windows: []iconf.DirPath{{"bin"}},
			MacOS:   []iconf.DirPath{{"bin"}},
			Linux:   []iconf.DirPath{{"bin"}},
		},
	}
}

func TestGradle() {
	gg := NewGradle()
	gg.Start()

	ff := "/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/gradle.json"
	content, _ := json.MarshalIndent(gg.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

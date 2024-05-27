package official

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gvcgo/vcollector/internal/iconf"
	"github.com/gvcgo/vcollector/internal/req"
	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
	"github.com/gvcgo/vcollector/pkgs/version"
)

func init() {
	crawler.RegisterCrawler(NewGroovy())
}

const (
	GroovyDistributionUrlPattern string = "https://archive.apache.org/dist/groovy/%s/distribution/"
)

var GroovyVersionRegexp = regexp.MustCompile(`\d+(.\d+){2}`)

/*
https://groovy.apache.org/download.html
https://archive.apache.org/dist/groovy/
*/
type Groovy struct {
	DownloadUrl string
	SDKName     string
	Version     version.VersionList
}

func NewGroovy() (g *Groovy) {
	g = &Groovy{
		DownloadUrl: "https://archive.apache.org/dist/groovy/",
		SDKName:     "groovy",
		Version:     make(version.VersionList),
	}
	return
}

func (g *Groovy) GetSDKName() string {
	return g.SDKName
}

func (g *Groovy) getSha256(sdkUrl string) (resp string) {
	sUrl := sdkUrl + ".sha256"
	resp = req.GetResp(sUrl)
	return
}

func (g *Groovy) getVersion(vhref string) {
	vName := strings.Trim(vhref, "/")
	dUrl := fmt.Sprintf(GroovyDistributionUrlPattern, vName)
	doc := req.GetDocument(dUrl)
	if doc == nil {
		return
	}
	var (
		sdkUrl    string
		sha256Str string
	)
	doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		hh := s.AttrOr("href", "")
		if !strings.Contains(hh, "-sdk-") {
			return
		}
		if strings.HasSuffix(hh, ".zip") {
			sdkUrl = fmt.Sprintf("%s/%s", strings.Trim(dUrl, "/"), strings.Trim(hh, "/"))
			sha256Str = g.getSha256(sdkUrl)
		}
	})

	if sdkUrl != "" && sha256Str != "" {
		item := version.Item{}
		item.Url = sdkUrl
		item.Arch = "any"
		item.Os = "any"
		item.Sum = strings.TrimSpace(sha256Str)
		item.SumType = "sha256"
		if strings.Contains(item.Sum, "apache-groovy-sdk") {
			sList := strings.Split(item.Sum, "apache-groovy-sdk")
			item.Sum = strings.TrimSpace(sList[0])
			item.SumType = "md5"
		}
		if vlist, ok := g.Version[vName]; !ok || vlist == nil {
			g.Version[vName] = version.Version{}
		}
		g.Version[vName] = append(g.Version[vName], item)
	}
}

func (g *Groovy) getResult() {
	doc := req.GetDocument(g.DownloadUrl)
	if doc == nil {
		return
	}
	doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		href := s.AttrOr("href", "")
		if GroovyVersionRegexp.FindString(href) != "" {
			g.getVersion(href)
		}
	})
}

func (g *Groovy) Start() {
	g.getResult()
}

func (g *Groovy) GetVersions() []byte {
	r, _ := g.Version.Marshal()
	return r
}

func (g *Groovy) HomePage() string {
	return "http://www.groovy-lang.org/"
}

func (g *Groovy) GetInstallConf() (ic iconf.InstallerConfig) {
	return iconf.InstallerConfig{
		FlagFiles: &iconf.FileItems{
			Windows: []string{"bin"},
			MacOS:   []string{"bin"},
			Linux:   []string{"bin"},
		},
		BinaryDirs: &iconf.DirItems{
			Windows: []iconf.DirPath{{"bin"}},
			MacOS:   []iconf.DirPath{{"bin"}},
			Linux:   []iconf.DirPath{{"bin"}},
		},
	}
}

func TestGroovy() {
	gg := NewGroovy()
	gg.Start()

	ff := "/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/groovy.json"
	content, _ := json.MarshalIndent(gg.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

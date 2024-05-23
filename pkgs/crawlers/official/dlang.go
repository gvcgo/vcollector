package official

import (
	"encoding/json"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gvcgo/vcollector/internal/req"
	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
	"github.com/gvcgo/vcollector/pkgs/version"
)

func init() {
	crawler.RegisterCrawler(NewDlang())
}

/*
https://downloads.dlang.org/releases/
*/
type Dlang struct {
	DownloadUrl string
	SDKName     string
	Version     version.VersionList
	host        string
}

func NewDlang() (d *Dlang) {
	d = &Dlang{
		DownloadUrl: "https://downloads.dlang.org/releases/2.x/",
		SDKName:     "dlang",
		Version:     make(version.VersionList),
		host:        "https://downloads.dlang.org",
	}
	return
}

func (d *Dlang) GetSDKName() string {
	return d.SDKName
}

func (d *Dlang) getVersion(vName, vHref string) {
	dUrl, _ := url.JoinPath(d.host, vHref)
	doc := req.GetDocument(dUrl)
	if doc == nil {
		return
	}
	doc.Find("div#content").Find("li").Find("a").Each(func(_ int, s *goquery.Selection) {
		href := s.AttrOr("href", "")
		item := version.Item{}
		if strings.HasSuffix(href, ".windows.zip") {
			item.Url, _ = url.JoinPath(d.host, href)
			item.Arch = "amd64"
			item.Os = "windows"
		} else if strings.HasSuffix(href, ".osx.zip") {
			item.Url, _ = url.JoinPath(d.host, href)
			item.Arch = "amd64"
			item.Os = "darwin"
		} else if strings.HasSuffix(href, ".linux.zip") {
			item.Url, _ = url.JoinPath(d.host, href)
			item.Arch = "amd64"
			item.Os = "linux"
		}

		if item.Url != "" {
			item.Installer = version.Unarchiver
			if vlist, ok := d.Version[vName]; !ok || vlist == nil {
				d.Version[vName] = version.Version{}
			}
			d.Version[vName] = append(d.Version[vName], item)
		}
	})
}

func (d *Dlang) getResult() {
	doc := req.GetDocument(d.DownloadUrl)
	if doc == nil {
		return
	}
	doc.Find("div#content").Find("li").Find("a").Each(func(_ int, s *goquery.Selection) {
		vName := s.Text()
		vHref := s.AttrOr("href", "")
		if vHref == "" {
			return
		}
		if strings.Count(vName, ".") < 2 {
			return
		}

		// higher than 2.065.0
		vList := strings.Split(vName, ".")
		vMinor, _ := strconv.Atoi(vList[1])
		if vMinor < 65 {
			return
		}
		d.getVersion(vName, vHref)
	})
}

func (d *Dlang) Start() {
	d.getResult()
}

func (d *Dlang) GetVersions() []byte {
	r, _ := d.Version.Marshal()
	return r
}

func (d *Dlang) HomePage() string {
	return "https://dlang.org/"
}

func TestDlang() {
	dd := NewDlang()
	dd.Start()

	ff := "/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/dlang.json"
	content, _ := json.MarshalIndent(dd.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

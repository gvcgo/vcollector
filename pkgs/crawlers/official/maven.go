package official

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gvcgo/vcollector/internal/req"
	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
	"github.com/gvcgo/vcollector/pkgs/version"
)

func init() {
	crawler.RegisterCrawler(NewMaven())
}

/*
https://dlcdn.apache.org/maven/
*/
type Maven struct {
	DownloadUrl string
	SDKName     string
	Version     version.VersionList
}

func NewMaven() (m *Maven) {
	m = &Maven{
		DownloadUrl: "https://dlcdn.apache.org/maven",
		SDKName:     "maven",
		Version:     make(version.VersionList),
	}
	return
}

func (m *Maven) GetSDKName() string {
	return m.SDKName
}

func (m *Maven) getSha(tarUrl string) string {
	shaUrl := tarUrl + ".sha512"
	resp := req.GetResp(shaUrl)
	return resp
}

func (m *Maven) getVersions(dUrl string) {
	doc := req.GetDocument(dUrl)
	if doc == nil {
		return
	}
	doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		text := s.Text()
		if strings.Contains(text, ".") && strings.HasSuffix(text, "/") {
			h := s.AttrOr("href", "")
			if h == "" {
				return
			}
			vUrl, _ := url.JoinPath(dUrl, h, "/binaries/")
			_doc := req.GetDocument(vUrl)
			if _doc == nil {
				return
			}
			_doc.Find("a").Each(func(_ int, ss *goquery.Selection) {
				t := ss.Text()
				if strings.HasSuffix(t, "bin.tar.gz") {
					href := ss.AttrOr("href", "")
					if href == "" {
						return
					}
					tarUrl, _ := url.JoinPath(vUrl, href)

					item := version.Item{}
					item.Url = tarUrl
					item.SumType = "sha512"
					item.Sum = m.getSha(tarUrl)
					item.Extra = text
					item.Arch = "any"
					item.Os = "any"
					item.Installer = version.Unarchiver
					vStr := strings.Trim(text, "/")
					if _, ok := m.Version[vStr]; !ok {
						m.Version[vStr] = version.Version{}
					}
					m.Version[vStr] = append(m.Version[vStr], item)
				}
			})
		}
	})
}

func (m *Maven) getResult() {
	routes := []string{
		"/maven-3/",
		"/maven-4/",
	}
	for _, route := range routes {
		dUrl, _ := url.JoinPath(m.DownloadUrl, route)
		fmt.Println(dUrl)
		m.getVersions(dUrl)
	}
}

func (m *Maven) Start() {
	m.getResult()
}

func (m *Maven) GetVersions() []byte {
	r, _ := m.Version.Marshal()
	return r
}

func TestMaven() {
	mm := NewMaven()
	mm.Start()

	ff := "/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/maven.json"
	content, _ := json.MarshalIndent(mm.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

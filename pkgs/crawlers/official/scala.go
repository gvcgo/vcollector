package official

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gvcgo/vcollector/internal/req"
	"github.com/gvcgo/vcollector/pkgs/version"
)

/*
https://www.scala-lang.org/download/all.html
*/
type Scala struct {
	DownloadUrl string
	SDKName     string
	Version     version.VersionList
}

func NewScala() (s *Scala) {
	s = &Scala{
		DownloadUrl: "https://www.scala-lang.org/download/all.html",
		SDKName:     "scala",
		Version:     make(version.VersionList),
	}
	return
}

func (s *Scala) GetSDKName() string {
	return s.SDKName
}

func (s *Scala) getResult() {
	doc := req.GetDocument(s.DownloadUrl)
	if doc == nil {
		return
	}
	doc.Find("div.download-elem").Find("a").Each(func(_ int, ss *goquery.Selection) {
		vName := strings.ReplaceAll(ss.Text(), "Scala ", "")
		vName = strings.TrimSpace(strings.ReplaceAll(vName, " ", "-"))
		if _, ok := s.Version[vName]; !ok {
			s.Version[vName] = version.Version{
				version.Item{
					Arch:      "any",
					Os:        "any",
					Installer: version.Coursier,
				},
			}
		}
	})
}

func (s *Scala) Start() {
	s.getResult()
}

func (s *Scala) GetVersions() version.VersionList {
	return s.Version
}

func TestScala() {
	ss := NewScala()
	ss.Start()

	ff := "/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/scala.json"
	content, _ := json.MarshalIndent(ss.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

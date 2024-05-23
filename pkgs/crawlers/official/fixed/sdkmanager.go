package fixed

import (
	"encoding/json"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gvcgo/vcollector/internal/req"
	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
	"github.com/gvcgo/vcollector/pkgs/version"
)

func init() {
	crawler.RegisterCrawler(NewSDKManager())
}

var SDKManagerPattern = regexp.MustCompile(`(\d+)`)

var SDKManagerOSMap = map[string]string{
	"windows": "windows",
	"linux":   "linux",
	"mac":     "darwin",
}

/*
https://developer.android.com/studio?hl=en
*/
type SDKManager struct {
	DownloadUrl string
	SDKName     string
	Version     version.VersionList
	host        string
}

func NewSDKManager() (s *SDKManager) {
	s = &SDKManager{
		DownloadUrl: "https://developer.android.com/studio?hl=en",
		SDKName:     "sdkmanager",
		Version:     make(version.VersionList),
		host:        "https://dl.google.com/android/repository",
	}
	return
}

func (s *SDKManager) GetSDKName() string {
	return s.SDKName
}

func (s *SDKManager) getResult() {
	doc := req.GetDocument(s.DownloadUrl)
	if doc == nil {
		return
	}
	doc.Find("table.download").Eq(1).Find("tr").Each(func(idx int, ss *goquery.Selection) {
		if idx == 0 {
			return
		}
		platform := strings.ToLower(ss.Find("td").Eq(0).Text())
		fName := ss.Find("td").Eq(1).Find("button").Text()
		if platform == "" || fName == "" {
			return
		}
		vName := SDKManagerPattern.FindString(fName)
		u, _ := url.JoinPath(s.host, fName)
		sha256Str := ss.Find("td").Eq(3).Text()

		item := version.Item{}
		item.Url = u
		item.Arch = "any"
		item.Os = SDKManagerOSMap[strings.TrimSpace(platform)]
		item.Sum = strings.TrimSpace(sha256Str)
		if item.Sum != "" {
			item.SumType = "sha256"
		}
		item.Installer = version.Unarchiver
		if _, ok := s.Version[vName]; !ok {
			s.Version[vName] = version.Version{}
		}
		s.Version[vName] = append(s.Version[vName], item)
	})
}

func (s *SDKManager) Start() {
	s.getResult()
}

func (s *SDKManager) GetVersions() []byte {
	r, _ := s.Version.Marshal()
	return r
}

func (s *SDKManager) HomePage() string {
	return "https://developer.android.com/studio"
}

func TestSDKManager() {
	os.Setenv(req.ProxyEnvName, "http://127.0.0.1:2023")
	ss := NewSDKManager()
	ss.Start()

	ff := "/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/sdkmanager.json"
	content, _ := json.MarshalIndent(ss.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}

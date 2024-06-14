package additional

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/gvcgo/vcollector/internal/req"
	"github.com/gvcgo/vcollector/pkgs/version"
)

/*
collect url patterns for cloudflare.
*/
const (
	GiteeHost      string = "https://gitee.com/moqsien/vsources/raw/main/"
	SDKListFile    string = "sdk-list.version.json"
	SDKVersionFile string = "%s.version.json"
)

var URLPatternList = map[string]struct{}{}

func ParseVersionList(sdkName string) {
	if sdkName == "conda-forge-pkgs" {
		return
	}
	vUrl := GiteeHost + fmt.Sprintf(SDKVersionFile, sdkName)
	versionList := version.VersionList{}
	req.GetJson(vUrl, &versionList, 180)
	for _, versions := range versionList {
		for _, item := range versions {
			itemUrl := item.Url
			if item.Url == "" {
				continue
			}
			var result string
			if u, err := url.Parse(itemUrl); err == nil {
				pList := strings.Split(u.Path, "/")
				if len(pList) > 1 {
					result = u.Host + "/" + pList[1]
				} else {
					result = u.Host + "/" + pList[0]
				}
			}
			if result != "" {
				URLPatternList[result] = struct{}{}
			}

		}
	}
}

func ParseSDKList() {
	sdkListUrl := GiteeHost + SDKListFile
	sl := &map[string]map[string]string{}
	req.GetJson(sdkListUrl, sl, 180)
	for sdk := range *sl {
		ParseVersionList(sdk)
	}
	result := []string{}
	for dUrl := range URLPatternList {
		result = append(result, dUrl)
	}
	data, _ := json.MarshalIndent(&result, "", "  ")
	os.WriteFile("/home/moqsien/projects/go/src/gvcgo/vcollector/url_patterns.txt", data, os.ModePerm)
}

func TestURLs() {
	ParseSDKList()
}

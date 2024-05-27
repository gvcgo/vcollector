package official

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gvcgo/vcollector/internal/iconf"
	"github.com/gvcgo/vcollector/internal/req"
	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
	"github.com/gvcgo/vcollector/pkgs/version"
)

func init() {
	crawler.RegisterCrawler(NewJDK())
}

/*
https://bell-sw.com/pages/downloads/
https://api.bell-sw.com/v1/liberica/releases/?&
*/

var JDKOsMap = map[string]string{
	"linux":   "linux",
	"windows": "windows",
	"macos":   "darwin",
}

var JDKArchMap = map[string]string{
	"x86": "amd64",
	"arm": "arm64",
}

type LibericaItem struct {
	Bitness        int    `json:"bitness"`
	UpdateVersion  int    `json:"updateVersion"`
	DownloadUrl    string `json:"downloadUrl"`
	BundleType     string `json:"bundleType"`
	FeatureVersion int    `json:"featureVersion"`
	PackageType    string `json:"packageType"`
	Architecture   string `json:"architecture"`
	ExtraVersion   int    `json:"extraVersion"`
	BuildVersion   int    `json:"buildVersion"`
	Os             string `json:"os"`
	InterimVersion int    `json:"interimVersion"`
	Version        string `json:"version"`
	Sha1           string `json:"sha1"`
	Size           int64  `json:"size"`
	PatchVersion   int    `json:"patchVersion"`
	LTS            bool   `json:"LTS"`
}

type LibericaResult []LibericaItem

type JDK struct {
	DownloadUrl string
	SDKName     string
	result      LibericaResult
	Verisons    version.VersionList
}

func NewJDK() (j *JDK) {
	j = &JDK{
		DownloadUrl: "https://api.bell-sw.com/v1/liberica/releases/?&",
		SDKName:     "jdk",
		result:      LibericaResult{},
		Verisons:    make(version.VersionList),
	}
	return
}

func (j *JDK) GetSDKName() string {
	return j.SDKName
}

func (j *JDK) getResult() {
	req.GetJson(j.DownloadUrl, &j.result)
}

func (j *JDK) filter() {
	for _, jItem := range j.result {
		if jItem.BundleType != "jdk-full" || jItem.Bitness != 64 {
			continue
		}

		if jItem.PackageType != "tar.gz" && jItem.PackageType != "zip" {
			continue
		}

		if jItem.Architecture != "x86" && jItem.Architecture != "arm" {
			continue
		}
		item := version.Item{}
		item.Os = JDKOsMap[jItem.Os]
		item.Arch = JDKArchMap[jItem.Architecture]
		item.Installer = version.Unarchiver
		item.Size = jItem.Size
		if jItem.LTS {
			item.LTS = "1"
		}
		item.SumType = "sha1"
		item.Sum = jItem.Sha1
		item.Url = jItem.DownloadUrl
		// featureVersion.extraVersion.updateVersion.patchVersion+buildVersion
		vStr := fmt.Sprintf(
			"%d.%d.%d.%d_%d",
			jItem.FeatureVersion,
			jItem.ExtraVersion,
			jItem.UpdateVersion,
			jItem.PatchVersion,
			jItem.BuildVersion,
		)
		item.Extra = jItem.Version
		if _, ok := j.Verisons[vStr]; !ok {
			j.Verisons[vStr] = version.Version{}
		}
		j.Verisons[vStr] = append(j.Verisons[vStr], item)
	}
}

func (j *JDK) Start() {
	j.getResult()
	j.filter()
}

func (j *JDK) GetVersions() []byte {
	r, _ := j.Verisons.Marshal()
	return r
}

func (j *JDK) HomePage() string {
	return "https://bell-sw.com/pages/downloads/"
}

func (j *JDK) GetInstallConf() (ic iconf.InstallerConfig) {
	return iconf.InstallerConfig{
		FlagFiles: &iconf.FileItems{
			Windows: []string{"bin", "lib", "include"},
			MacOS:   []string{"bin", "lib", "include"},
			Linux:   []string{"bin", "lib", "include"},
		},
		BinaryDirs: &iconf.DirItems{
			Windows: []iconf.DirPath{{"bin"}, {"jre", "bin"}},
			MacOS:   []iconf.DirPath{{"bin"}, {"jre", "bin"}},
			Linux:   []iconf.DirPath{{"bin"}, {"jre", "bin"}},
		},
		AdditionalEnvs: iconf.AdditionalEnvList{
			{
				Name:  "JAVA_HOME",
				Value: []iconf.DirPath{},
			},
			{
				Name:    "CLASSPATH",
				Value:   []iconf.DirPath{{"lib", "dt.jar"}, {"lib", "tools.jar"}, {"jre", "lib", "rt.jar"}},
				Version: "major<=8",
			},
		},
	}
}

func TestJDK() {
	jj := NewJDK()
	jj.Start()
	ff := "/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/jdk.json"
	content, _ := json.MarshalIndent(jj.Verisons, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)

	f := "/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/jdk_raw.json"
	content, _ = json.MarshalIndent(jj.result, "", "    ")
	os.WriteFile(f, content, os.ModePerm)
}

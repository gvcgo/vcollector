package conf

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	ConfFileName  string = "vconf.json"
	VersionSha256 string = "sha256.json"
)

func GetWorkDir() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".vcollector")
}

func GetVersionDir() string {
	workdir := GetWorkDir()
	return filepath.Join(workdir, "versions")
}

func GetShaFile() string {
	return filepath.Join(GetVersionDir(), VersionSha256)
}

type ShaList map[string]string

func GetSha256(sdkName string) (shaStr string) {
	shaList := ShaList{}
	shaPath := GetShaFile()
	content, _ := os.ReadFile(shaPath)
	json.Unmarshal(content, &shaList)
	shaStr = shaList[sdkName]
	return
}

func SaveSha256(sdkName string, shaStr string) {
	shaList := ShaList{}
	shaPath := GetShaFile()
	content, _ := os.ReadFile(shaPath)
	json.Unmarshal(content, &shaList)
	shaList[sdkName] = shaStr
	content, _ = json.Marshal(shaList)
	os.WriteFile(shaPath, content, os.ModePerm)
}

type Config struct {
	Proxy       string `json:"proxy"`
	GithubToken string `josn:"gh_token"`
	GithubRepo  string `json:"gh_repo"`
}

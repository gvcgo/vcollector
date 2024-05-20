package conf

import (
	"os"
	"path/filepath"
)

const (
	ConfFileName string = "vconf.json"
)

func GetWorkDir() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".vcollector")
}

func GetVersionDir() string {
	workdir := GetWorkDir()
	return filepath.Join(workdir, "versions")
}

type Config struct {
	Proxy       string `json:"proxy"`
	GithubToken string `josn:"gh_token"`
	GithubRepo  string `json:"gh_repo"`
}

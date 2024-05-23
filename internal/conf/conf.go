package conf

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gvcgo/goutils/pkgs/gutils"
	toml "github.com/pelletier/go-toml/v2"
)

const (
	ConfFileName string = "vconf.toml"
)

func GetWorkDir() string {
	homeDir, _ := os.UserHomeDir()
	p := filepath.Join(homeDir, ".vcollector")
	os.MkdirAll(p, os.ModePerm)
	return p
}

func GetVersionDir() string {
	workdir := GetWorkDir()
	p := filepath.Join(workdir, "versions")
	os.MkdirAll(p, os.ModePerm)
	return p
}

type Config struct {
	Proxy       string `json,toml:"proxy"`       // proxy: http or socks5
	GithubToken string `josn,toml:"githubToken"` // github access token
	GithubRepo  string `json,toml:"githubRepo"`  // github repository name, format: "user/repo"
}

func NewConfig() (c *Config) {
	c = &Config{
		Proxy:      "",
		GithubRepo: "gvcgo/vsources",
	}
	c.Load()
	return
}

func (c *Config) GetConfPath() string {
	return filepath.Join(GetWorkDir(), ConfFileName)
}

func (c *Config) Load() {
	p := c.GetConfPath()
	if ok, _ := gutils.PathIsExist(p); ok {
		content, _ := os.ReadFile(p)
		toml.Unmarshal(content, c)
	}
}

func (c *Config) save() {
	content, err := toml.Marshal(c)
	if len(content) > 0 && err == nil {
		os.WriteFile(c.GetConfPath(), content, os.ModePerm)
	}
}

func (c *Config) SetProxy(proxy string) {
	if strings.HasPrefix(proxy, "http://") || strings.HasPrefix(proxy, "socks5://") {
		c.Proxy = proxy
		c.save()
	}
}

func (c *Config) SetGithubToken(token string) {
	if token != "" {
		c.GithubToken = token
		c.save()
	}
}

func (c *Config) SetGithubRepo(repo string) {
	if repo != "" {
		c.GithubRepo = repo
		c.save()
	}
}

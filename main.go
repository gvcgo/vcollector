package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/gvcgo/vcollector/internal/conda"
	"github.com/gvcgo/vcollector/internal/conf"
	"github.com/gvcgo/vcollector/internal/utils"
	"github.com/gvcgo/vcollector/pkgs/crawlers/conda"
	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
	_ "github.com/gvcgo/vcollector/pkgs/crawlers/gh/tools"
	toml "github.com/pelletier/go-toml/v2"
)

type MyConfig struct {
	Version int               `toml:"version"`
	Name    string            `toml:"name"`
	Tags    []string          `toml:"tags"`
	Data    map[string]string `toml:"data"`
}

func TestToml() {
	mc := &MyConfig{
		Version: 10,
		Name:    "vmr",
		Tags:    []string{"v1", "v2"},
		Data: map[string]string{
			"darwin":  "macOS",
			"linux":   "Linux",
			"windows": "Windows",
		},
	}
	b, _ := toml.Marshal(mc)
	fmt.Println(string(b))
}

func UploadVSourceReadme() {
	uu := utils.NewUploader()
	localPath := "/Volumes/data/projects/go/src/gvcgo_org/vcollector/docs/README.md"
	uu.Github.UploadFile("README.md", localPath)
}

func UploadHomepageFile() {
	fName := "sdk-homepage.json"
	list := map[string]string{}
	for _, cc := range crawler.CrawlerList {
		list[cc.GetSDKName()] = cc.HomePage()
	}
	upl := utils.NewUploader()
	upl.DisableSaveSha256()
	content, _ := json.MarshalIndent(list, "", "  ")
	upl.Upload(fName, "", content)
}

func RunCrawler(cc crawler.Crawler) {
	if cc == nil {
		return
	}
	fmt.Println("start crawler:", cc.GetSDKName())
	cc.Start()
	uploader := utils.NewUploader()
	if cc.GetSDKName() == conda.CondaForgeSDKName {
		uploader.DisableSaveSha256()
	}
	uploader.Upload(cc.GetSDKName(), cc.HomePage(), cc.GetVersions())
}

func UploadMirrorsInChina() {
	var MirrorsInChina = map[string]string{
		"https://go.dev/dl/":                   "https://golang.google.cn/dl/",
		"https://nodejs.org/download/release/": "https://mirrors.huaweicloud.com/nodejs/",
		"https://storage.googleapis.com/flutter_infra_release/releases/stable/": "https://mirrors.tuna.tsinghua.edu.cn/flutter/flutter_infra/releases/stable/",
		"https://gradle.org/releases/next-steps/?version=":                      "https://mirrors.cloud.tencent.com/gradle/gradle-%s-all.zip",
		"https://dlcdn.apache.org/maven/":                                       "https://mirrors.aliyun.com/apache/maven/",
		"https://repo.anaconda.com/miniconda/":                                  "https://mirrors.tuna.tsinghua.edu.cn/anaconda/miniconda/",
		"https://mirrors.tuna.tsinghua.edu.cn/julia-releases/bin/":              "https://julialang-s3.julialang.org/bin/",
	}

	storageDir := conf.GetInstallConfigFileDir()
	mirrorFilePath := filepath.Join(storageDir, "customed_mirrors.toml")

	content, _ := toml.Marshal(MirrorsInChina)
	os.WriteFile(mirrorFilePath, content, os.ModePerm)
	uu := utils.NewUploader()
	uu.Github.UploadFile("mirrors/customed_mirrors.toml", mirrorFilePath)
}

func main() {
	// official.TestJDK()
	// official.TestGolang()
	// official.TestMaven()
	// official.TestGradle()
	// official.TestDotnet()
	// official.TestZig()
	// official.TestNode()
	// official.TestFlutter()
	// official.TestJulia()
	// official.TestDlang()
	// official.TestGroovy()
	// official.TestKubectl()
	// official.TestScala()
	// fixed.TestVSCode()
	// fixed.TestSDKManager()
	// fixed.TestMiniconda()
	// conda.TestPython()
	// conda.TestPyPy()
	// conda.TestGCC()
	// conda.TestClang()
	// conda.TestRuby()
	// conda.TestRust()
	// conda.TestLua()
	// conda.TestR()
	// conda.TestLFortran()
	// conda.TestPerl()
	// conda.TestTypst()
	conda.TestMySQL()
	conda.TestPostgreSQL()
	// RunCrawler(conda.NewNim())

	// gh.TestGithub()

	// lans.TestBun()
	// lans.TestCodon()
	// lans.TestDeno()
	// lans.TestGleam()
	// lans.TestKotlin()
	// lans.TestOdin()
	// lans.TestVlang()
	// lans.TestElixir()
	// RunCrawler(lans.NewClojure())
	// RunCrawler(lans.NewCrystal())

	// ttt := map[string]string{
	// 	"aaa": "bbb",
	// }
	// content, _ := toml.Marshal(ttt)
	// fmt.Println(string(content))
	// UploadMirrorsInChina()

	// lsp.TestDlangLsp()
	// lsp.TestTypstLsp()
	// lsp.TestVAnalyzer()
	// lsp.TestZls()
	// RunCrawler(lsp.NewTypstLsp())

	// tools.TestAsciinema()
	// tools.TestCMake()
	// tools.TestCoursier()
	// tools.TestFd()
	// tools.TestFzf()
	// tools.TestGitWin()
	// tools.TestGsudo()
	// tools.TestLazydocker()
	// tools.TestLazygit()
	// tools.TestProtoc()
	// tools.TestRipgrep()
	// tools.TestTreeSitter()
	// tools.TestTypstPreview()
	// tools.TestUpx()
	// tools.TestVhs()
	// RunCrawler(tools.NewNeovim())
	// RunCrawler(tools.NewAgg())

	// mix.TestErlang()
	// mix.TestPHP()

	// TestToml()

	// UploadVSourceReadme()
	// UploadHomepageFile()

	// RunCrawler(mix.NewPHP())
	// RunCrawler(official.NewJDK())
	// RunCrawler(mix.NewErlang())

	// m := fixed.NewMiniconda()
	// ic := m.GetInstallConf()
	// content, _ := toml.Marshal(ic)
	// localPath := filepath.Join("/Volumes/data/projects/go/src/gvcgo_org/vcollector", "miniconda.toml")
	// os.WriteFile(localPath, content, os.ModePerm)
	// uploader := utils.NewUploader()
	// uploader.Github.UploadFile("install/miniconda.toml", localPath)

	// additional.TestURLs()
}

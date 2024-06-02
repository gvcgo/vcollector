package main

import (
	"encoding/json"
	"fmt"

	_ "github.com/gvcgo/vcollector/internal/conda"
	"github.com/gvcgo/vcollector/internal/utils"
	"github.com/gvcgo/vcollector/pkgs/crawlers/conda"
	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
	_ "github.com/gvcgo/vcollector/pkgs/crawlers/gh/lans"
	_ "github.com/gvcgo/vcollector/pkgs/crawlers/gh/lsp"
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
	RunCrawler(conda.NewNim())

	// gh.TestGithub()

	// lans.TestBun()
	// lans.TestCodon()
	// lans.TestDeno()
	// lans.TestGleam()
	// lans.TestKotlin()
	// lans.TestOdin()
	// lans.TestVlang()
	// lans.TestElixir()

	// lsp.TestDlangLsp()
	// lsp.TestTypstLsp()
	// lsp.TestVAnalyzer()
	// lsp.TestZls()

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
}

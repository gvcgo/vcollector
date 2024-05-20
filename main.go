package main

import (
	_ "github.com/gvcgo/vcollector/internal/conda"
	"github.com/gvcgo/vcollector/pkgs/crawlers/conda"
	_ "github.com/gvcgo/vcollector/pkgs/crawlers/official"
	_ "github.com/gvcgo/vcollector/pkgs/crawlers/official/fixed"
)

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
	conda.TestTypst()
}

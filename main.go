package main

import (
	_ "github.com/gvcgo/vcollector/pkgs/crawlers/official"
	"github.com/gvcgo/vcollector/pkgs/crawlers/official/fixed"
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
	fixed.TestMiniconda()
}

package fixed

import (
	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
	"github.com/gvcgo/vcollector/pkgs/version"
)

func init() {
	crawler.RegisterCrawler(NewRustup())
}

type Rustup struct {
	SDKName string
	Version version.VersionList
}

func NewRustup() *Rustup {
	return &Rustup{
		SDKName: "rustup",
		Version: make(version.VersionList),
	}
}

func (r *Rustup) GetSDKName() string {
	return r.SDKName
}

func (r *Rustup) Start() {
	r.Version["latest"] = version.Version{
		version.Item{
			Url:       "https://static.rust-lang.org/rustup/dist/x86_64-apple-darwin/rustup-init",
			Arch:      "amd64",
			Os:        "darwin",
			Installer: version.Executable,
		},
		version.Item{
			Url:       "https://static.rust-lang.org/rustup/dist/aarch64-apple-darwin/rustup-init",
			Arch:      "arm64",
			Os:        "darwin",
			Installer: version.Executable,
		},
		version.Item{
			Url:       "https://static.rust-lang.org/rustup/dist/x86_64-unknown-linux-gnu/rustup-init",
			Arch:      "amd64",
			Os:        "linux",
			Installer: version.Executable,
		},
		version.Item{
			Url:       "https://static.rust-lang.org/rustup/dist/aarch64-unknown-linux-gnu/rustup-init",
			Arch:      "arm64",
			Os:        "linux",
			Installer: version.Executable,
		},
		version.Item{
			Url:       "https://static.rust-lang.org/rustup/dist/x86_64-pc-windows-msvc/rustup-init.exe",
			Arch:      "amd64",
			Os:        "windows",
			Installer: version.Executable,
		},
		version.Item{
			Url:       "https://static.rust-lang.org/rustup/dist/aarch64-pc-windows-msvc/rustup-init.exe",
			Arch:      "arm64",
			Os:        "windows",
			Installer: version.Executable,
		},
	}
}

func (r *Rustup) GetVersions() []byte {
	rr, _ := r.Version.Marshal()
	return rr
}

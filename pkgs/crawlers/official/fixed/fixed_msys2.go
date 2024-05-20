package fixed

import "github.com/gvcgo/vcollector/pkgs/version"

type Msys2 struct {
	SDKName string
	Version version.VersionList
}

func NewMsys2(sdkName string) *Msys2 {
	return &Msys2{
		SDKName: "msys2",
		Version: make(version.VersionList),
	}
}

func (m *Msys2) GetSDKName() string {
	return m.SDKName
}

func (m *Msys2) Start() {
	m.Version["latest"] = version.Version{
		version.Item{
			Url:       "https://github.com/msys2/msys2-installer/releases/download/nightly-x86_64/msys2-x86_64-latest.exe",
			Arch:      "amd64",
			Os:        "windows",
			Installer: version.Executable,
		},
	}
}

func (m *Msys2) GetVersions() version.VersionList {
	return m.Version
}

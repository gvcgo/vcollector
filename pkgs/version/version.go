package version

import "encoding/json"

const (
	Darwin  string = "darwin"
	Linux   string = "linux"
	Windows string = "windows"
)

const (
	Conda      string = "conda"
	CondaForge string = "conda-forge"
	Coursier   string = "coursier"
	Unarchiver string = "unarchiver"
	Executable string = "executable"
	Dpkg       string = "dpkg"
	Rpm        string = "rpm"
)

type Item struct {
	Url       string `json:"url"`       // download url
	Arch      string `json:"arch"`      // amd64 | arm64
	Os        string `json:"os"`        // linux | darwin | windows
	Sum       string `json:"sum"`       // Checksum
	SumType   string `json:"sum_type"`  // sha1 | sha256 | sha512 | md5
	Size      int64  `json:"size"`      // Size in bytes
	Installer string `json:"installer"` // conda | conda-forge | coursier | unarchiver | executable | dpkg | rpm
	LTS       string `json:"lts"`       // Long Term Support
	Extra     string `json:"extra"`     // Extra Info
}

type Version []Item

type VersionList map[string]Version

func (v *VersionList) Unmarshal(data []byte) error {
	return json.Unmarshal(data, v)
}

func (v VersionList) Marshal() ([]byte, error) {
	return json.Marshal(v)
}

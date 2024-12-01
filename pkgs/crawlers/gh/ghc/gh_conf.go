package ghc

import (
	"github.com/gvcgo/vcollector/internal/iconf"
)

type GhConfig struct {
	SDK           string                `json:"sdk"`
	Repo          string                `json:"repo"`
	HomePage      string                `json:"homepage"`
	Exclude       map[string][]string   `json:"exclude"`
	Include       map[string][]string   `json:"include"`
	Os            map[string]string     `json:"os"`
	Arch          map[string]string     `json:"arch"`
	Type          string                `json:"type"`
	VersionRegExp string                `json:"versionRegExp"`
	VersionOTP    string                `json:"versionOTP"`
	Install       iconf.InstallerConfig `json:"install"`
}

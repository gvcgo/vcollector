package official

import "github.com/gvcgo/vcollector/pkgs/version"

/*
https://storage.googleapis.com/flutter_infra_release/releases/releases_{linux/macos/windows}.json
https://storage.flutter-io.cn/flutter_infra_release/releases/releases_{linux/macos/windows}.json
*/
type Flutter struct {
	DownloadUrl string
	SDKName     string
	Version     version.VersionList
}

package iconf

type FileItems struct {
	Windows []string `toml,json:"windows"`
	Linux   []string `toml,json:"linux"`
	MacOS   []string `toml,json:"darwin"`
}

type (
	DirPath  []string
	DirItems struct {
		Windows []DirPath `toml,json:"windows"` // <symbolLinkPath>/<filepath.Join(List)>, ...
		Linux   []DirPath `toml,json:"linux"`
		MacOS   []DirPath `toml,json:"darwin"`
	}
)

type AdditionalEnv struct {
	Name    string
	Value   []DirPath // <symbolLinkPath>/<filepath.Join(Value)>
	Version string    // major>8 or major<=8(for JDK)
}

type AdditionalEnvList []AdditionalEnv

type BinaryRename struct {
	NameFlag string `toml,json:"name_flag"`
	RenameTo string `toml,json:"rename_to"`
}

/*
Installation configs
*/
type InstallerConfig struct {
	FlagFiles       *FileItems        `toml,json:"flag_files"`
	FlagDirExcepted bool              `toml,json:"flag_dir_excepted"`
	BinaryDirs      *DirItems         `toml,json:"binary_dirs"`
	BinaryRename    *BinaryRename     `toml,json:"binary_rename"`
	AdditionalEnvs  AdditionalEnvList `toml,json:"additional_envs"`
}

package ghc

import (
	"embed"
	"fmt"

	"encoding/json"
)

//go:embed from_github/*
var jsonFS embed.FS

func TestEmbed() {
	content, _ := jsonFS.ReadFile("from_github/gh_projects/bun.json")
	cc := &GhConfig{}
	json.Unmarshal(content, cc)
	fmt.Println(cc.Install)
}

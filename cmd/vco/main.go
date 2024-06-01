package main

import (
	"github.com/gvcgo/goutils/pkgs/gtea/gprint"
	"github.com/gvcgo/goutils/pkgs/gutils"
	"github.com/gvcgo/vcollector/internal/conf"
	"github.com/spf13/cobra"
)

type Cli struct {
	rootCmd *cobra.Command
}

func NewCli() (c *Cli) {
	c = &Cli{
		rootCmd: &cobra.Command{
			Use:   "vco",
			Short: "vco is a tool for crawling version list for SDKs.",
			Long:  "vco <Command> <SubCommand> --flags args...",
		},
	}
	c.initiate()
	return
}

func (c *Cli) initiate() {
	c.rootCmd.AddCommand(&cobra.Command{
		Use:     "set-repo",
		Short:   "set remote repository.",
		Long:    "vco set-repo <remote-repo: user/reponame>",
		Aliases: []string{"sr"},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}
			cfg := conf.NewConfig()
			cfg.SetGithubRepo(args[0])
		},
	})

	c.rootCmd.AddCommand(&cobra.Command{
		Use:     "set-proxy",
		Short:   "set local proxy.",
		Long:    "vco set-proxy <http/socks5>",
		Aliases: []string{"sp"},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}
			cfg := conf.NewConfig()
			cfg.SetProxy(args[0])
		},
	})

	c.rootCmd.AddCommand(&cobra.Command{
		Use:     "set-token",
		Short:   "set github access token.",
		Long:    "vco set-token <github token>",
		Aliases: []string{"st"},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}
			cfg := conf.NewConfig()
			cfg.SetGithubToken(args[0])
		},
	})

	c.rootCmd.AddCommand(&cobra.Command{
		Use:     "crawl",
		Short:   "crawl version list for SDKs.",
		Long:    "vco crawl.",
		Aliases: []string{"c"},
		Run: func(cmd *cobra.Command, args []string) {
			start()
		},
	})
}

func (c *Cli) Run() {
	if c.rootCmd != nil {
		if err := c.rootCmd.Execute(); err != nil {
			gprint.PrintError("%+v", err)
		}
	}
}

func main() {
	cs := gutils.CtrlCSignal{}
	cs.ListenSignal()
	cli := NewCli()
	cli.Run()

	// GenerateSDKHompePageListForDocs()
}

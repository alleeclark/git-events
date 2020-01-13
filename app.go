package main

import (
	"github.com/urfave/cli"
)

//New cli application for git-events commands
func New() *cli.App {
	app := cli.NewApp()
	app.Name = "git-events"
	app.Version = "0.0.1"
	app.Usage = "stream git changes over grpc"

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "git-user", Value: "gitevents", Usage: "git user for ssh", Required: false},
		cli.StringFlag{Name: "git-url", Usage: "git url to clone", Required: false},
		cli.StringFlag{Name: "git-branch", Value: "master", Usage: "git branch to run syncing on"},
		cli.StringFlag{Name: "git-dir", Value: "/var/gitevents/data", Usage: "directory to pull to"},
		cli.BoolFlag{Name: "git-ssh-enabled", Usage: "enable ssh agent usage", Hidden: true},
		cli.BoolFlag{Name: "git-ssh-file", Usage: "read public, private, and passpharse for ssh agent", Hidden: true},
		cli.StringFlag{Name: "git-fingerprint-path", FilePath: "/var/gitevents/.ssh/fingerprint", Usage: "git RSA finerprint id", Required: false},
		cli.BoolFlag{Name: "metrics", Usage: "send metrics to pushgateway", EnvVar: "GITEVENTS_METRICS", Hidden: true},
		cli.StringFlag{Name: "pushgateway-addr", Value: "localhost:9091", Usage: "push gateway address for metrics", Hidden: true},
		cli.StringFlag{Name: "port", Value: "9000", Usage: "grpc port to run on"},
		cli.Uint64Flag{Name: "interval", Value: 300, Usage: "interval in seconds"},
	}
	app.Commands = []cli.Command{startCommand}
	return app
}

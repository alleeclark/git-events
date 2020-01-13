package main

import (
	"fmt"
	"git-events/git"
	"net"
	"os"

	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

var startCommand = cli.Command{
	Name:        "start",
	Usage:       "start a streaming frequency",
	ArgsUsage:   "[flags] <ref>",
	Description: "fetch contents changes and sync to consul",
	Action: func(c *cli.Context) error {
		lis, err := net.Listen("tcp", ":"+c.GlobalString("port"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to listen: %v", err)
			return err
		}
		s := grpc.NewServer()
		git.RegisterEventsServiceServer(s, &server{
			interval: c.GlobalUint64("interval"),
			branch:   c.GlobalString("branch"),
			name:     c.GlobalString("git-dir"),
		})
		fmt.Fprintf(os.Stdout, "Getting ready to start server on port %s", c.GlobalString("port"))

		if err := s.Serve(lis); err != nil {
			fmt.Fprintf(os.Stderr, "failed to serve: %v", err)
			return err
		}
		return nil
	},
}

package main

import (
	"fmt"
	"os"

	"StandBy-Web/routes"

	mingfuflags "github.com/Fu-XDU/mingfu_go_common/flags"
	"github.com/labstack/gommon/log"
	"github.com/urfave/cli/v2"
)

const (
	clientIdentifier = "StandBy-Web"
	clientVersion    = "1.0.0"
	clientUsage      = "A minimal Web Clock"
)

var (
	app = cli.NewApp()
)

func init() {
	app.Action = ServerApp
	app.Name = clientIdentifier
	app.Version = clientVersion
	app.Usage = clientUsage
	app.Commands = []*cli.Command{}
	app.Flags = append(app.Flags, mingfuflags.GinFlags...)
}

func ServerApp(ctx *cli.Context) error {
	if args := ctx.Args(); args.Len() > 0 {
		return fmt.Errorf("invalid command: %q", args.First())
	}
	err := prepare()
	if err != nil {
		log.Error(err)
	}
	return err
}

func prepare() (err error) {
	routes.Run()
	return
}

func main() {
	if err := app.Run(os.Args); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

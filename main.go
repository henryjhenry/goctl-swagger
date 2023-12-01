package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/aishuchen/goctl-swagger/render"
	"github.com/urfave/cli/v2"
)

var (
	version  = "0.0.1"
	commands = []*cli.Command{
		{
			Name:  "swagger",
			Usage: "generates swagger json file",
			//Action: action.Generator,
			Action: render.Do,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "host",
					Usage: "api request address",
				},
				&cli.StringFlag{
					Name:  "basePath",
					Usage: "url request prefix",
				},
				&cli.StringFlag{
					Name:  "target",
					Usage: "swagger save file name",
				},
				&cli.StringFlag{
					Name:  "schemes",
					Usage: "swagger support schemes: http, https, ws, wss",
				},
				&cli.StringFlag{
					Name:  "tagPrefix",
					Usage: "add prefix on operation's tags",
				},
			},
		},
	}
)

func main() {
	app := cli.NewApp()
	app.Usage = "a plugin of goctl to generate swagger json file"
	app.Version = fmt.Sprintf("%s %s/%s", version, runtime.GOOS, runtime.GOARCH)
	app.Commands = commands
	if err := app.Run(os.Args); err != nil {
		fmt.Printf("goctl-swagger: %+v\n", err)
	}
}

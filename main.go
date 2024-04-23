package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/henryjhenry/goctl-swagger/render"
	"github.com/urfave/cli/v2"
)

var (
	version  = "0.0.2"
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
					Usage: "api request prefix",
				},
				&cli.StringFlag{
					Name:  "target",
					Usage: "swagger save file name, default: ./swagger.json",
				},
				&cli.StringFlag{
					Name:  "schemes",
					Usage: "swagger support schemes: http, https, ws, wss",
				},
				&cli.StringFlag{
					Name:  "tagPrefix",
					Usage: "add prefix on operation's tags",
				},
				&cli.StringFlag{
					Name:  "outsideSchema",
					Usage: "add outside schema api file",
				},
				&cli.StringFlag{
					Name:  "responseKey",
					Usage: "special response data key when outsideSchema is set, default: data",
				},
			},
		},
	}
)

func main() {
	logger := log.Default()
	app := cli.NewApp()
	app.Usage = "a plugin of goctl to generate swagger json file"
	app.Version = fmt.Sprintf("%s %s/%s", version, runtime.GOOS, runtime.GOARCH)
	app.Commands = commands
	//cwd, err := os.Getwd()
	//if err != nil {
	//	panic(err)
	//}
	if err := app.Run(os.Args); err != nil {
		logger.Fatalf("goctl-swagger: %+v\n", err)
	}
}

package main

import (
	"os"

	"github.com/dafsic/gambler/utils"
	"github.com/dafsic/gambler/version"
	"github.com/urfave/cli/v2"
)

// 这里mylog还没有初始化，不能使用
var logger = utils.NewLogger(os.Stdout, "main", utils.LDebug, utils.Ldefault)

func main() {
	app := &cli.App{
		Name:    "gambler",
		Usage:   "gambler command [args]",
		Version: version.GamblerVersion.String(),
		Commands: []*cli.Command{
			runCmd,
			//...
		},
	}
	if err := app.Run(os.Args); err != nil {
		logger.Warnf("%+v", err)
		return
	}
}

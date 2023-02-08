package main

import (
	"github.com/dafsic/gambler/config"
	"github.com/dafsic/gambler/lib/mylog"
	"github.com/dafsic/gambler/modules/channels"
	"github.com/dafsic/gambler/modules/client"
	"github.com/dafsic/gambler/modules/listent"
	"github.com/dafsic/gambler/modules/robot"
	"github.com/dafsic/gambler/modules/server"
	"github.com/dafsic/gambler/modules/store"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
)

var runCmd = &cli.Command{
	Name:  "run",
	Usage: "Auto bet robot",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			EnvVars: []string{"GAMBLER_CFG"},
			Value:   "~/gambler/config.toml",
			Usage:   "Load configuration from `FILE`",
		},
		&cli.StringFlag{
			Name:    "pool",
			Aliases: []string{"p"},
			EnvVars: []string{"GAMBLER_POOL"},
			Value:   "",
			Usage:   "41开头的池子地址",
		},
		&cli.StringFlag{
			Name:    "refund",
			Aliases: []string{"r"},
			EnvVars: []string{"GAMBLER_REFUND"},
			Value:   "",
			Usage:   "池子回款地址",
		},
		&cli.StringFlag{
			Name:    "address",
			Aliases: []string{"addr"},
			EnvVars: []string{"GAMBLER_ADDR"},
			Value:   "",
			Usage:   "下注使用的地址",
		},
	},
	Action: func(cctx *cli.Context) error {
		// 要测试下多个模块都需要的情况下是不是初始化多个实例还是一个，文章中看是多个，只能有一个实例的话就需要sync.Once.
		// 目前看只会调用一次构造函数，不知道是不是因为返回值都是指针
		fx.New(
			fx.Supply(cctx), //config 模块需要从命令行参数中获取配置文件路径
			config.CfgModule,
			mylog.Module,
			channels.ChanManagerModule,
			listent.ListenModule,
			client.ClientModule,
			robot.SchedModule,
			robot.RegisterRouters,
			server.ServerModule,
			store.StoreModule,
			fx.NopLogger,
		).Run() // 模块里不能有阻塞的协程，都要用go开启一个新的线程，Run()函数会在app.start后卡住等信号，收到中断信号会调用app.stop

		return nil
	},
}

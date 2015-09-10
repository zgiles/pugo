package boot

import (
	"github.com/codegangsta/cli"
	"gopkg.in/inconshreveable/log15.v2"
	"pugo/src/core"
	"pugo/src/middle"
	"pugo/src/service"
)

var (
	serverCommand = cli.Command{
		Name:   "server",
		Usage:  "start pugo blog web server",
		Action: Server,
	}
)

func Server(ctx *cli.Context) {
	opt := service.BootstrapInitOption{true, false, false}
	if err := service.Call(service.Bootstrap.Init, opt); err != nil {
		log15.Crit("Server.start.fail", "error", err)
	}
	// if not installed,try to install
	if core.Cfg.Install == "0" {
		Install(ctx)
	}

	// re-init all things
	opt = service.BootstrapInitOption{true, true, true}
	if err := service.Call(service.Bootstrap.Init, opt); err != nil {
		log15.Crit("Server.start.fail", "error", err)
	}
	// bootstrap service, preload data
	opt2 := service.BootstrapOption{true, true}
	if err := service.Call(service.Bootstrap.Bootstrap, opt2); err != nil {
		log15.Crit("Server.start.fail", "error", err)
	}
	log15.Info("Server.prepare")

	// set middleware and routers
	core.Server.Use(
		middle.Logger(),
		middle.Authorizor(),
		middle.Themer())

	// start server
	log15.Info("Server.start." + core.Cfg.Http.Host + ":" + core.Cfg.Http.Port)
	core.Server.Run(core.Cfg.Http.Host + ":" + core.Cfg.Http.Port)
}

package main

import (
	"flag"
	"github.com/deeptest-com/deeptest-next/cmd/server/serve"
	"github.com/deeptest-com/deeptest-next/cmd/server/v1/router"
	"github.com/deeptest-com/deeptest-next/internal/pkg/core/auth"
	"github.com/deeptest-com/deeptest-next/internal/pkg/inits"
	"github.com/deeptest-com/deeptest-next/internal/pkg/serve/cache"
	"github.com/deeptest-com/deeptest-next/internal/pkg/serve/cron_server"
	"github.com/deeptest-com/deeptest-next/internal/pkg/serve/database"
	"github.com/deeptest-com/deeptest-next/internal/pkg/serve/viper_server"
	"github.com/deeptest-com/deeptest-next/internal/pkg/serve/web"
	"github.com/deeptest-com/deeptest-next/internal/pkg/serve/web/web_iris"
	"github.com/deeptest-com/deeptest-next/internal/pkg/serve/zap_server"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/repo"
	_logUtils "github.com/deeptest-com/deeptest-next/pkg/libs/log"
	"github.com/facebookgo/inject"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	flagSet *flag.FlagSet
)

func main() {
	flagSet = flag.NewFlagSet("deeptest", flag.ContinueOnError)
	flagSet.Parse(os.Args[1:])

	viper_server.Init(database.GetViperConfig())

	zap_server.Init()
	inits.Init()

	webServer := web_iris.Init()
	serve.AddStatic(webServer.App)

	cache.Init()
	err := auth.InitDriver(&auth.Config{
		DriverType:      "redis",
		HmacSecret:      nil,
		UniversalClient: cache.Instance(),
	})
	if err != nil {
		_logUtils.Zap.Panic(err.Error())
	}

	// inject objects
	var g inject.Graph
	g.Logger = logrus.StandardLogger()

	cronServer := cron_server.NewCronServer()
	indexModule := router.NewIndexModule()

	err = g.Provide(
		&inject.Object{Value: repo.GetDbInstance()},
		&inject.Object{Value: cronServer},

		&inject.Object{Value: indexModule},
	)
	if err != nil {
		_logUtils.Errorf("provide usecase objects to the Graph: %v", err)
	}
	err = g.Populate()
	if err != nil {
		_logUtils.Errorf("populate the incomplete Objects: %v", err)
	}
	cronServer.Start()

	webServer.AddModule(web_iris.Party{
		Perfix:    "/api/v1",
		PartyFunc: indexModule.ApiParty(),
	})

	web.Start(webServer)
}

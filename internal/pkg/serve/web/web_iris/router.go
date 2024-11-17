package web_iris

import (
	"github.com/deeptest-com/deeptest-next/internal/pkg/config"
	"github.com/deeptest-com/deeptest-next/internal/pkg/libs/arr"
	"github.com/deeptest-com/deeptest-next/internal/pkg/serve/web/web_iris/middleware"
	"net/http"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/pprof"
	"github.com/kataras/iris/v12/middleware/rate"
	"github.com/kataras/iris/v12/middleware/recover"
)

// InitRouter
func (ws *WebServer) InitRouter() error {
	app := ws.App.Party("/").AllowMethods(iris.MethodOptions)
	{
		app.Get("/v0/version", func(ctx iris.Context) {
			ctx.WriteString("IRIS-ADMIN is running!!!")
		})

		app.UseRouter(middleware.CrsAuth())
		app.UseRouter(recover.New())
		if !config.CONFIG.Limit.Disable {
			limitV1 := rate.Limit(config.CONFIG.Limit.Limit, config.CONFIG.Limit.Burst, rate.PurgeEvery(time.Minute, 5*time.Minute))
			app.Use(limitV1)
		}
		if config.CONFIG.System.Level == "debug" {
			debug := func(index iris.Party) {
				index.Get("/", func(ctx iris.Context) {
					ctx.HTML("<h1>请点击<a href='/debug/pprof'>这里</a>打开调试页面")
				})
				index.Any("/pprof", pprof.New())
				index.Any("/pprof/{action:path}", pprof.New())
			}
			app.PartyFunc("/debug", debug)
		}

		for _, party := range ws.parties {
			app.PartyFunc(party.Perfix, party.PartyFunc)
		}
	}

	// http test must build
	if err := ws.App.Build(); err != nil {
		return err
	}

	return nil
}

// GetSources
// - PermRoutes
// - NoPermRoutes
func (ws *WebServer) GetSources() ([]map[string]string, []map[string]string) {
	methodExcepts := strings.Split(config.CONFIG.Except.Method, ";")
	uris := strings.Split(config.CONFIG.Except.Uri, ";")
	routeLen := len(ws.App.GetRoutes())
	permRoutes := make([]map[string]string, 0, routeLen)
	noPermRoutes := make([]map[string]string, 0, routeLen)

	for _, r := range ws.App.GetRoutes() {
		route := map[string]string{
			"path": r.Path,
			"name": r.Name,
			"act":  r.Method,
		}
		httpStatusType := arr.NewCheckArrayType(4)
		httpStatusType.AddMutil(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete)
		if !httpStatusType.Check(r.Method) {
			noPermRoutes = append(noPermRoutes, route)
			continue
		}

		if len(methodExcepts) > 0 && len(uris) > 0 && len(methodExcepts) == len(uris) {
			for i := 0; i < len(methodExcepts); i++ {
				if strings.EqualFold(r.Method, strings.ToLower(methodExcepts[i])) && strings.EqualFold(r.Path, strings.ToLower(uris[i])) {
					noPermRoutes = append(noPermRoutes, route)
					continue
				}
			}
		}

		permRoutes = append(permRoutes, route)
	}
	return permRoutes, noPermRoutes
}

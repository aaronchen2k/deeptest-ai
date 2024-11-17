package web_iris

import (
	stdContext "context"
	"errors"
	"github.com/deeptest-com/deeptest-next/internal/pkg/config"
	"github.com/deeptest-com/deeptest-next/internal/pkg/libs/arr"
	"github.com/deeptest-com/deeptest-next/internal/pkg/serve/web/web_iris/middleware"
	_str "github.com/deeptest-com/deeptest-next/pkg/libs/string"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/recover"
)

var ErrAuthDriverEmpty = errors.New("auth driver initialize fail")

// WebServer
// - App iris application
// - idleConnsClosed
// - addr
// - timeFormat
// - staticPrefix

type WebServer struct {
	App             *iris.Application
	idleConnsClosed chan struct{}
	parties         []Party
	addr            string
	timeFormat      string
}

// Party
// - perfix
// - partyFunc
type Party struct {
	Perfix    string
	PartyFunc func(index iris.Party)
}

// Init
func Init() *WebServer {
	app := iris.New()
	if config.CONFIG.System.Tls {
		app.Use(middleware.LoadTls())
	}

	app.Use(recover.New())

	app.Validator = validator.New()
	app.Logger().SetLevel(config.CONFIG.System.Level)
	idleConnsClosed := make(chan struct{})

	iris.RegisterOnInterrupt(func() {
		timeout := 10 * time.Second
		ctx, cancel := stdContext.WithTimeout(stdContext.Background(), timeout)
		defer cancel()
		app.Shutdown(ctx) // close all hosts
		close(idleConnsClosed)
	})

	config.SetDefaultAddrAndTimeFormat()

	return &WebServer{
		App:             app,
		addr:            config.CONFIG.System.Addr,
		timeFormat:      config.CONFIG.System.TimeFormat,
		idleConnsClosed: idleConnsClosed,
	}
}

// GetEngine
func (ws *WebServer) GetEngine() *iris.Application {
	return ws.App
}

// AddModule
func (ws *WebServer) AddModule(parties ...Party) {
	ws.parties = append(ws.parties, parties...)
}

// AddWebStatic
func (ws *WebServer) AddWebStatic(staticAbsPath, webPrefix string, paths ...string) {
	webPrefixs := strings.Split(config.CONFIG.System.WebPrefix, ",")
	wp := arr.NewCheckArrayType(2)
	for _, webPrefix := range webPrefixs {
		wp.Add(webPrefix)
	}
	if wp.Check(webPrefix) {
		return
	}

	fsOrDir := iris.Dir(staticAbsPath)
	opt := iris.DirOptions{
		IndexName: "index.html",
		SPA:       true,
	}
	ws.App.HandleDir(webPrefix, fsOrDir, opt)
	config.CONFIG.System.WebPrefix = _str.Join(config.CONFIG.System.WebPrefix, ",", webPrefix)
}

// AddUploadStatic
func (ws *WebServer) AddUploadStatic(webPrefix, staticAbsPath string) {
	fsOrDir := iris.Dir(staticAbsPath)
	ws.App.HandleDir(webPrefix, fsOrDir)
	config.CONFIG.System.StaticPrefix = webPrefix
}

// Run
func (ws *WebServer) Run() {
	ws.App.Listen(
		ws.addr,
		iris.WithoutInterruptHandler,
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
		iris.WithTimeFormat(ws.timeFormat),
	)
	<-ws.idleConnsClosed
}

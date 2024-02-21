package run

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"projects/LDmitryLD/geo-nginx/geo/config"
	"projects/LDmitryLD/geo-nginx/geo/internal/db"
	"projects/LDmitryLD/geo-nginx/geo/internal/infrastructure/cache"
	"projects/LDmitryLD/geo-nginx/geo/internal/infrastructure/component"
	"projects/LDmitryLD/geo-nginx/geo/internal/infrastructure/responder"
	"projects/LDmitryLD/geo-nginx/geo/internal/infrastructure/router"
	"projects/LDmitryLD/geo-nginx/geo/internal/infrastructure/server"
	"projects/LDmitryLD/geo-nginx/geo/internal/modules"
	"projects/LDmitryLD/geo-nginx/geo/internal/storages"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type Runner interface {
	Run() error
}

type App struct {
	conf   config.AppConf
	logger *zap.Logger
	srv    server.Server
	Sig    chan os.Signal
}

func NewApp(conf config.AppConf, logger *zap.Logger) *App {
	return &App{conf: conf, logger: logger, Sig: make(chan os.Signal, 1)}
}

func (a *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())

	errGroup, ctx := errgroup.WithContext(ctx)

	errGroup.Go(func() error {
		sigIng := <-a.Sig
		a.logger.Info("signal interrupt recieved", zap.Stringer("os_signal", sigIng))
		cancel()
		return nil
	})

	errGroup.Go(func() error {
		err := a.srv.Serve(ctx)
		if err != nil && err != http.ErrServerClosed {
			a.logger.Error("app: server error:", zap.Error(err))
			return err
		}
		return nil
	})

	if err := errGroup.Wait(); err != nil {
		return err
	}

	return nil
}

func (a *App) Bootstrap(options ...interface{}) Runner {

	_, sqlAdapter, err := db.NewSqlDB(a.conf.DB, a.logger)
	if err != nil {
		a.logger.Fatal("error init db", zap.Error(err))
	}

	cacheClient, err := cache.NewRedisClient(a.conf.Cache.Host, a.conf.Cache.Port, a.logger)
	if err != nil {
		a.logger.Fatal("error init cache", zap.Error(err))
	}

	responseManages := responder.NewResponder(a.logger)

	components := component.NewComponents(a.conf, responseManages, a.logger)

	newStorages := storages.NewStorages(sqlAdapter, cacheClient, a.logger)

	services := modules.NewServices(newStorages, components)

	controllers := modules.NewControllers(services, components)

	r := router.NewRouter(controllers)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", a.conf.Server.Port),
		Handler: r,
	}

	a.srv = server.NewHTTPServer(a.conf.Server, srv, a.logger)

	return a
}

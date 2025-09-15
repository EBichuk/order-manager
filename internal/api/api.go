package api

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"order-manager/internal/cache"
	"order-manager/internal/config"
	"order-manager/internal/controller/http"
	"order-manager/internal/controller/kafka"
	"order-manager/internal/repository"
	"order-manager/internal/service"
	"order-manager/pkg/db"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	cfg         config.Config
	logger      *slog.Logger
	repo        *repository.Repository
	s           *service.Service
	httpServer  *http.Server
	pool        *pgxpool.Pool
	cache       *cache.Cache
	kafkaReader *kafka.Consumer
}

func NewApp() *App {
	app := &App{}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to init configs %v", err)
	}

	app.logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	app.pool, err = db.InitPool(fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		cfg.Db.User, cfg.Db.Password, cfg.Db.Host, cfg.Db.Port, cfg.Db.Name))
	if err != nil {
		log.Fatalf("Failed to connect to db %v", err)
	}

	app.cache = cache.NewCache(cfg.Cache.Size)

	app.repo = repository.NewRepository(app.pool)

	app.s = service.NewService(app.repo, app.cache, app.logger)

	app.kafkaReader = kafka.NewConsumer(app.s, app.logger, cfg.Topic, cfg.Brokers)

	handlerOrder := http.NewHandler(app.s, app.logger)
	app.httpServer = http.NewServer(handlerOrder, cfg.Addr)
	app.cfg = cfg

	return app
}

func (a *App) RunApp() {
	err := a.s.FillCache(a.cfg.Cache.Size)
	if err != nil {
		log.Fatalf("Failed to fill cache %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go a.kafkaReader.Start(ctx)

	go func() {
		err = a.httpServer.StartHttpServer()
		if err != nil {
			log.Fatalf("Failed to start http server %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	if err := a.kafkaReader.Stop(); err != nil {
		log.Fatalf("Failed to stop kafka %v", err)
	}
	a.pool.Close()

	a.logger.Info("application stopped")
}

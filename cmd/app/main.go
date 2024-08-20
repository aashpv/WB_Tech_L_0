package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"wb-order-service/internal/config"
	"wb-order-service/internal/http-server/handlers"
	"wb-order-service/internal/stan"
	"wb-order-service/internal/storage/memcache"
	"wb-order-service/internal/storage/postgres"
)

func main() {
	cfg := config.MustLoad()

	log := newLogger()

	log.Info("starting wb-order-service")

	storage, err := postgres.NewDB(cfg.StorageConn)
	if err != nil {
		log.Error("failed to init storage")
		os.Exit(1)
	}
	log.Info("storage initialized")

	cache := memcache.NewCache(storage)
	log.Info("cache initialized")

	err = cache.LoadOrders()
	if err != nil {
		log.Error("failed to load orders", "err", err)
		os.Exit(1)
	}
	log.Info("cache is loaded")

	stn := &stan.Stan{
		ClusterId: cfg.NatsStreaming.ClusterId,
		ClientId:  cfg.NatsStreaming.ClientId,
		Subject:   cfg.NatsStreaming.Subject,
	}

	if err := stn.NewStan(cache); err != nil {
		log.Error("failed to init stan", err)
		os.Exit(1)
	}
	log.Info("stan initialized")

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	fs := http.FileServer(http.Dir("web/static"))
	router.Handle("/*", http.StripPrefix("/", fs))

	router.Get("/", handlers.OrderPage)
	router.Get("/api/{id}", handlers.GetOrderHandler(log, cache))

	srv := &http.Server{
		Addr:         cfg.HttpServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}

	log.Info("server started on port: ", slog.String("address", cfg.HttpServer.Address))

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server is stopped")
}

func newLogger() *slog.Logger {
	var log *slog.Logger

	log = slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	return log
}

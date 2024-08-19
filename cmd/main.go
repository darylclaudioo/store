package main

import (
	"context"
	"hexagon-architecture/config"
	"hexagon-architecture/internal/api"
	productsRepo "hexagon-architecture/internal/domain/products/repository"
	productsDB "hexagon-architecture/internal/domain/products/repository/db"
	"hexagon-architecture/internal/infrastructure"
	"hexagon-architecture/internal/service"
	"hexagon-architecture/pkg/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	configuration := config.GetConfig()

	database, _ := config.NewDB(configuration.DB, configuration.App.Env)

	var products productsRepo.Repository = productsDB.New(database, "products")

	service := service.New(
		products,
	)

	contexts := context.Background()

	otelShutdown, err := infrastructure.SetupOTelSDK(contexts, configuration.App.Name, configuration.App.Version, configuration.Otel.Host, configuration.App.Env)
	if err != nil {
		return
	}
	defer func() {
		err = otelShutdown(contexts)
	}()

	httpServer := http.New(http.Config{
		Port:            configuration.App.Port,
		ReadTimeout:     configuration.App.ReadTimeout,
		WriteTimeout:    configuration.App.WriteTimeout,
		GracefulTimeout: configuration.App.GracefulTimeout,
	})
	router := httpServer.Router()
	router.Use(recover.New())
	router.Use(logger.New(logger.Config{
		Format: "[${time}] ${ip}  ${status} - ${latency} ${method} ${path}\n",
	}))
	router.Use(otelfiber.Middleware())

	api.New(service, configuration.App.Env).Register(router)

	httpServerChan := httpServer.Run()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case errors := <-httpServerChan:
		if errors != nil {
			return
		}
	case <-signalChan:
	}

}

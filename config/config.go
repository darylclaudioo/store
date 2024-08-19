package config

import (
	"context"
	"errors"
	"fmt"
	"hexagon-architecture/internal/infrastructure"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

type config struct {
	App  appConfig  `envconfig:"APP"`
	Otel otelConfig `envconfig:"OTEL"`
	DB   dbConfig   `envconfig:"DB"`
}

type appConfig struct {
	Env             string        `envconfig:"ENV" default:"development" validate:"required,oneof=local development staging production" mod:"no_space,lcase"` // local
	Port            string        `envconfig:"PORT" default:"8006" validate:"required" mod:"no_space"`
	ReadTimeout     time.Duration `envconfig:"READ_TIMEOUT" default:"5s" validate:"required,gt=0"`
	WriteTimeout    time.Duration `envconfig:"WRITE_TIMEOUT" default:"5s" validate:"required,gt=0"`
	GracefulTimeout time.Duration `envconfig:"GRACEFUL_TIMEOUT" default:"10s" validate:"required,gt=0"`
	Host            string        `envconfig:"HOST" validate:"required,url"`
	Name            string        `envconfig:"NAME" validate:"required"`
	Version         string        `envconfig:"VERSION" validate:"required"`
}

type otelConfig struct {
	Host string `envconfig:"GRPC_HOST" validate:"required"`
}

type dbConfig struct {
	Address         string        `envconfig:"ADDRESS" default:"localhost:27018" validate:"required"`
	Name            string        `envconfig:"NAME" default:"store" validate:"required"`
	SslMode         string        `envconfig:"SSL_MODE" default:"disable"`
	MaxConnOpen     int           `envconfig:"MAX_CONN_OPEN" default:"10" validate:"required,gt=0"`
	MaxConnIdle     int           `envconfig:"MAX_CONN_IDLE" default:"10" validate:"required,gt=0"`
	MaxConnLifetime time.Duration `envconfig:"MAX_CONN_LIFETIME" default:"60s" validate:"required,gt=0"`
}

const envPrefix = "STORE"

func GetConfig() *config {
	var cfg config

	_ = godotenv.Load()

	if errors := envconfig.Process(envPrefix, &cfg); errors != nil {
		return nil
	}

	infrastructure.InitSlog()

	return &cfg
}

func NewDB(cfg dbConfig, env string) (*mongo.Database, error) {
	var dns string
	switch env {
	case "local", "development":
		dns = fmt.Sprintf("mongodb://%s", cfg.Address)
	default:
	}
	dbOptions := options.Client()
	dbOptions.Monitor = otelmongo.NewMonitor()
	dbOptions.ApplyURI(dns)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbClient, err := mongo.Connect(ctx, dbOptions)
	if err != nil {
		fmt.Println("=========== Failed to connect to database ", err)
		return nil, errors.Join(err)
	}
	fmt.Println("=========== success to connect to database ")
	err = dbClient.Ping(context.Background(), nil)
	if err != nil {
		fmt.Println("=========== Failed to ping database: ", err)
		return nil, err
	}
	fmt.Println("=========== success ping to connect to database ")
	return dbClient.Database(cfg.Name), nil
}

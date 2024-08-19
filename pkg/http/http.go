package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Server interface {
	Run() chan error
	Router() *fiber.App
	Close() error
}

type server struct {
	server *http.Server
	router *fiber.App
	cfg    Config
}

type Config struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	GracefulTimeout time.Duration
}

func New(cfg Config) Server {
	return &server{
		router: fiber.New(),
		cfg:    cfg,
	}
}

func (s *server) Router() *fiber.App {
	return s.router
}

func (s *server) Run() chan error {
	var ch = make(chan error)
	go s.run(ch)
	return ch
}

func (s *server) run(ch chan error) {
	err := s.router.Listen(":" + s.cfg.Port)
	if err != nil {
		ch <- err
		return
	}

	s.server = &http.Server{
		ReadTimeout:  s.cfg.ReadTimeout,
		WriteTimeout: s.cfg.WriteTimeout,
	}
}

func (s *server) Close() error {
	if s.server == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.GracefulTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}

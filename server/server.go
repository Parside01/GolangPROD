package server

import (
	"log/slog"
	"solution/pkg/database/postgres"
	"solution/server/controllers"
	"solution/server/controllers/auth"
	"solution/server/controllers/me"
	"solution/server/controllers/me/friends"
	"solution/server/controllers/posts"
	"solution/server/controllers/profiles"

	"github.com/labstack/echo/v4"
)

type Server struct {
	router      *echo.Echo
	address     string
	controllers []controllers.Controller
	logger      *slog.Logger
}

func NewServer(address string, logger *slog.Logger) *Server {
	return &Server{
		address: address,
		logger:  logger,
		router:  echo.New(),
		controllers: []controllers.Controller{
			controllers.NewCountriesController(logger, postgres.New(logger)),
			controllers.NewPingController(),
			auth.NewRegisterController(logger, postgres.New(logger)),
			auth.NewSignInController(logger, postgres.New(logger)),
			me.NewProfileController(logger, postgres.New(logger)),
			me.NewUpdatePasswordController(logger, postgres.New(logger)),
			profiles.NewProfilesController(logger, postgres.New(logger)),
			posts.NewPostController(logger, postgres.New(logger)),
			friends.NewFriendsController(postgres.New(logger), logger),
		},
	}
}

func (s *Server) Start() error {
	s.registerRoutesrs()
	s.logger.Info("server has been started", "address", s.address)
	return s.router.Start(s.address)
}

func (s *Server) registerRoutesrs() {
	for _, route := range s.controllers {
		group := s.router.Group(route.GetGroup())
		for _, middleware := range route.GetMiddleware() {
			group.Use(middleware)
		}
		for _, handler := range route.GetHandlers() {
			group.Add(handler.GetMethod(), handler.GetPath(), handler.GetHandler())
		}
	}
}

package httpserver

import (
	"fmt"

	"gameapp/config"
	"gameapp/service/authservice"
	"gameapp/service/userservice"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	config  config.Config
	authSvc authservice.Service
	userSvc userservice.Service
}

func New(
	config config.Config,
	authSvc authservice.Service,
	userSvc userservice.Service,
) Server {
	return Server{
		config:  config,
		authSvc: authSvc,
		userSvc: userSvc,
	}
}

func (s Server) Serve() {
	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.GET("/health-check", s.healthCheck)

	// User routes
	userGroup := e.Group("/users")
	userGroup.POST("/register", s.userRegisterHandler)
	userGroup.POST("/login", s.userLoginHandler)
	// e.GET("/users/profile", s.userProfileHandler)

	fmt.Println("server is listening on port 8080...")

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", s.config.HTTPServer.Port)))
}

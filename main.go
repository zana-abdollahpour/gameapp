package main

import (
	"time"

	"gameapp/config"
	"gameapp/delivery/httpserver"
	"gameapp/repository/mysql"
	"gameapp/service/authservice"
	"gameapp/service/userservice"
)

const (
	JwtSignKey                 = "BZ0niKtToA4TwoNjP1na"
	AccessTokenSubject         = "at"
	RefreshTokenSubject        = "rt"
	AccessTokenExpiryDuration  = time.Hour * 24
	RefreshTokenExpiryDuration = time.Hour * 24 * 7
)

func setupServices(cfg config.Config) (authservice.Service, userservice.Service) {
	authSvc := authservice.New(cfg.Auth)

	MysqlRepo := mysql.New(cfg.Mysql)
	userSvc := userservice.New(authSvc, MysqlRepo)

	return authSvc, userSvc
}

func main() {
	cfg := config.Config{
		HTTPServer: config.HTTPServer{Port: 8080},
		Mysql: mysql.Config{
			Username: "gameapp",
			Password: "gameappt0lk2o20",
			Host:     "localhost",
			Port:     3308,
			DBName:   "gameapp_db",
		},
		Auth: authservice.Config{
			SignKey:               JwtSignKey,
			AccessSubject:         AccessTokenSubject,
			RefreshSubject:        RefreshTokenSubject,
			AccessExpirationTime:  AccessTokenExpiryDuration,
			RefreshExpirationTime: RefreshTokenExpiryDuration,
		},
	}

	// TODO: add command for migrations
	// mgr := migrator.New(cfg.Mysql)
	// mgr.Up(0)

	authSvc, userSvc := setupServices(cfg)

	server := httpserver.New(cfg, authSvc, userSvc)
	server.Serve()

}

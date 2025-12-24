package mysql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Username string
	Password string
	Port     int
	Host     string
	DBName   string
}

type MySQLDB struct {
	db     *sql.DB
	config Config
}

func New(cfg Config) *MySQLDB {
	db, err := sql.Open(
		"mysql",
		fmt.Sprintf(
			"%s:%s@(%s:%d)/%s",
			cfg.Username,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.DBName,
		),
	)
	if err != nil {
		panic(fmt.Errorf("can not open mysql db: %v", err))
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return &MySQLDB{db: db, config: cfg}
}

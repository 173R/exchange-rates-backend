package db

import (
	"embed"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"github.com/wolframdeus/exchange-rates-backend/configs"
	gormpg "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
)

//go:embed migrations/*
var migrations embed.FS

// New создает новый экземпляр для работы с БД PostgreSQL.
func New(host string, port uint, user, pass, name string) (*gorm.DB, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		user, pass, host, port, name,
	)

	return gorm.Open(gormpg.Open(connStr), nil)
}

// NewByConfig работает аналогично New, но использует данные из
// конфига проекта.
func NewByConfig() (*gorm.DB, error) {
	c := configs.Db
	return New(c.Host, c.Port, c.User, c.Pass, c.Name)
}

// RunMigrations запускает миграции проекта.
func RunMigrations() error {
	gormDB, err := NewByConfig()
	if err != nil {
		return err
	}

	sqldb, err := gormDB.DB()
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(sqldb, &postgres.Config{})
	if err != nil {
		return err
	}

	source, err := httpfs.New(http.FS(migrations), "migrations")
	if err != nil {
		return err
	}

	mg, err := migrate.NewWithInstance("httpfs", source, configs.Db.Name, driver)
	if err != nil {
		panic(err)
	}

	if err = mg.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
	}
	return nil
}

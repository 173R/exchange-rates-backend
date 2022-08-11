package db

import (
	"fmt"
	"github.com/wolframdeus/exchange-rates-backend/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewClient создает новый экземпляр для работы с БД PostgreSQL.
func NewClient(host string, port uint, user, pass, name string) (*gorm.DB, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		user, pass, host, port, name,
	)

	return gorm.Open(postgres.Open(connStr), nil)
}

// NewClientByConfig работает аналогично NewClient, но использует данные из
// конфига проекта.
func NewClientByConfig() (*gorm.DB, error) {
	c := configs.Db
	return NewClient(c.Host, c.Port, c.User, c.Pass, c.Name)
}

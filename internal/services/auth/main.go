package auth

import (
	"errors"
	"github.com/wolframdeus/exchange-rates-backend/configs"
	"github.com/wolframdeus/exchange-rates-backend/internal/services/users"
	"github.com/wolframdeus/exchange-rates-backend/internal/tg"
	"net/url"
)

type Jwt struct {
}

type Service struct {
	// Сервис для работы с пользователями.
	uSrv *users.Service
}

func (s *Service) AuthenticateTg(initData string) (*Jwt, error) {
	// Валидируем параметры запуска.
	if ok, err := tg.ValidateInitData(initData, configs.Tg.SecretKey); err != nil {
		return nil, err
	} else if !ok {
		return nil, errors.New("init data is invalid")
	}

	// Парсим параметры запуска как query-параметры.
	q, err := url.Parse(initData)
	if err != nil {
		return nil, err
	}
}

// NewService создает указатель на новый экземпляр Service.
func NewService(uSrv *users.Service) *Service {
	return &Service{uSrv: uSrv}
}

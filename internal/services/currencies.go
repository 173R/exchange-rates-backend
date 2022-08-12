package services

import (
	"github.com/wolframdeus/exchange-rates-backend/internal/db/models"
	"github.com/wolframdeus/exchange-rates-backend/internal/repositories"
)

type Currencies struct {
	curRep *repositories.Currencies
}

// FindAll возвращает информацию о всех валютах.
func (c *Currencies) FindAll() ([]models.Currency, error) {
	return c.curRep.FindAll()
}

// NewCurrencies возвращает ссылку на новый экземпляр Currencies.
func NewCurrencies(curRep *repositories.Currencies) *Currencies {
	return &Currencies{curRep: curRep}
}

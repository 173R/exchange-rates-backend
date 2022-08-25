package http

import (
	"github.com/gin-gonic/gin"
	"github.com/wolframdeus/exchange-rates-backend/internal/context"
	creppkg "github.com/wolframdeus/exchange-rates-backend/internal/repositories/currencies"
	exratesreppkg "github.com/wolframdeus/exchange-rates-backend/internal/repositories/exrates"
	ureppkg "github.com/wolframdeus/exchange-rates-backend/internal/repositories/users"
	csrvpkg "github.com/wolframdeus/exchange-rates-backend/internal/services/currencies"
	exratessrvpkg "github.com/wolframdeus/exchange-rates-backend/internal/services/exrates"
	usrvpkg "github.com/wolframdeus/exchange-rates-backend/internal/services/users"
	"gorm.io/gorm"
)

// NewContextConfigMiddleware возвращает новый промежуточный обработчик,
// который внедряет в контекст запроса список известных сервисов.
func NewContextConfigMiddleware(db *gorm.DB) gin.HandlerFunc {
	curSrv := csrvpkg.New(creppkg.NewCurrencies(db))
	uSrv := usrvpkg.NewUsers(ureppkg.NewUsers(db))
	exRatesSrv := exratessrvpkg.NewService(exratesreppkg.NewRepository(db))

	return context.NewGinHandler(func(c *context.Gin) {
		// Инджектим сервисы.
		c.InjectServices(curSrv, uSrv, exRatesSrv)

		// Инджектим параметры запуска.
		if err := c.InjectLaunchParams(); err != nil {
			c.SendError(err)
		}
	})
}

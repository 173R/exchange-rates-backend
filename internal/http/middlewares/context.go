package http

import (
	"github.com/gin-gonic/gin"
	"github.com/wolframdeus/exchange-rates-backend/internal/context"
	"github.com/wolframdeus/exchange-rates-backend/internal/repositories"
	"github.com/wolframdeus/exchange-rates-backend/internal/services"
	"github.com/wolframdeus/exchange-rates-backend/internal/services/currencies"
	"gorm.io/gorm"
)

// NewContextConfigMiddleware возвращает новый промежуточный обработчик,
// который внедряет в контекст запроса список известных сервисов.
func NewContextConfigMiddleware(db *gorm.DB) gin.HandlerFunc {
	curRep := repositories.NewCurrencies(db)
	curSrv := currencies.New(curRep)

	uRep := repositories.NewUsers(db)
	uSrv := services.NewUsers(uRep)

	return context.NewGinHandler(func(c *context.Gin) {
		// Инджектим сервисы.
		c.InjectServices(curSrv, uSrv)

		// Инджектим параметры запуска.
		if err := c.InjectLaunchParams(); err != nil {
			c.SendError(err)
		}
	})
}

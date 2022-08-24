package http

import (
	"github.com/gin-gonic/gin"
	"github.com/wolframdeus/exchange-rates-backend/internal/context"
	creppkg "github.com/wolframdeus/exchange-rates-backend/internal/repositories/currencies"
	ureppkg "github.com/wolframdeus/exchange-rates-backend/internal/repositories/users"
	csrvpkg "github.com/wolframdeus/exchange-rates-backend/internal/services/currencies"
	usrvpkg "github.com/wolframdeus/exchange-rates-backend/internal/services/users"
	"gorm.io/gorm"
)

// NewContextConfigMiddleware возвращает новый промежуточный обработчик,
// который внедряет в контекст запроса список известных сервисов.
func NewContextConfigMiddleware(db *gorm.DB) gin.HandlerFunc {
	curRep := creppkg.NewCurrencies(db)
	curSrv := csrvpkg.New(curRep)

	uRep := ureppkg.NewUsers(db)
	uSrv := usrvpkg.NewUsers(uRep)

	return context.NewGinHandler(func(c *context.Gin) {
		// Инджектим сервисы.
		c.InjectServices(curSrv, uSrv)

		// Инджектим параметры запуска.
		if err := c.InjectLaunchParams(); err != nil {
			c.SendError(err)
		}
	})
}

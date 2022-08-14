package http

import (
	"github.com/gin-gonic/gin"
	"github.com/wolframdeus/exchange-rates-backend/internal/context"
	"github.com/wolframdeus/exchange-rates-backend/internal/repositories"
	"github.com/wolframdeus/exchange-rates-backend/internal/services"
	services2 "github.com/wolframdeus/exchange-rates-backend/internal/services/translations"
	"gorm.io/gorm"
)

// NewContextConfigMiddleware возвращает новый промежуточный обработчик,
// который внедряет в контекст запроса список известных сервисов.
func NewContextConfigMiddleware(db *gorm.DB) gin.HandlerFunc {
	curRep := repositories.NewCurrencies(db)
	trlRep := repositories.NewTranslations(db)

	curSrv := services.NewCurrencies(curRep)
	trlSrv := services2.NewTranslations(trlRep)

	return context.NewGinHandler(func(c *context.Gin) {
		c.InjectServices(curSrv, trlSrv)
		c.InjectLaunchParams()
	})
}

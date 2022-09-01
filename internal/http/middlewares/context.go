package http

import (
	"github.com/gin-gonic/gin"
	"github.com/wolframdeus/exchange-rates-backend/internal/context"
)

// NewContextConfigMiddleware возвращает новый промежуточный обработчик,
// который внедряет в контекст запроса список известных сервисов.
func NewContextConfigMiddleware() gin.HandlerFunc {
	return context.NewGinHandler(func(c *context.Gin) {
		// Инджектим параметры запуска.
		if err := c.InjectLaunchParams(); err != nil {
			c.SendError(err)
		}
	})
}

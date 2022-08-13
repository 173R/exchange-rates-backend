package http

import (
	"github.com/gin-gonic/gin"
	"github.com/wolframdeus/exchange-rates-backend/internal/context"
)

// NewCustomRecoveryMiddleware возвращает новый обработчик для процессинга
// паник.
func NewCustomRecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(gc *gin.Context, err any) {
		c := context.NewGin(gc)

		if e, ok := err.(error); ok {
			c.SendError(e.Error())
		} else {
			c.SendError("Произошла неизвестная ошибка.")
		}
	})
}

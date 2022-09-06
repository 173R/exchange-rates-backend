package context

import (
	"github.com/gin-gonic/gin"
)

type Gin struct {
	Context
	// Оригинальный контекст gin.
	Gin *gin.Context
}

// SendData успешно отправляет указанные данные по единому формату.
func (c *Gin) SendData(data any) {
	c.Gin.JSON(200, map[string]interface{}{
		"ok":   true,
		"data": data,
	})
}

// SendError отправляет указанную ошибку удаленному клиенту.
func (c *Gin) SendError(data any) {
	if err, ok := data.(error); ok {
		data = err.Error()
	}
	c.Gin.JSON(400, map[string]interface{}{
		"ok":    false,
		"error": data,
	})
	c.Gin.Abort()
}

// NewGin возвращает указатель на новый экземпляр Gin.
func NewGin(ctx *gin.Context) *Gin {
	return &Gin{
		Context: Context{},
		Gin:     ctx,
	}
}

package context

import (
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

type Context struct {
	Gin *gin.Context
}

// CaptureError захватывает ошибку и отправляет её в Sentry.
func (c *Context) CaptureError(err error) {
	// TODO: Установить IP клиента.
	hub := sentrygin.GetHubFromContext(c.Gin)
	if hub == nil {
		return
	}
	hub.CaptureException(err)
}

// Send успешно отправляет указанные данные по единому формату.
func (c *Context) Send(data any) {
	c.Gin.JSON(200, map[string]interface{}{
		"ok":   true,
		"data": data,
	})
}

// SendError отправляет указанную ошибку удаленному клиенту.
func (c *Context) SendError(data any) {
	c.Gin.JSON(400, map[string]interface{}{
		"ok":    false,
		"error": data,
	})
}

// New возвращает ссылку на новый экземпляр Context.
func New(c *gin.Context) *Context {
	return &Context{Gin: c}
}

// NewHandler возвращает новый обработчик и в нем вызывает переданную функцию
// с уже обернутым контекстом.
func NewHandler(f func(c *Context)) gin.HandlerFunc {
	return func(gc *gin.Context) {
		f(New(gc))
	}
}

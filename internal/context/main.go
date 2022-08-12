package context

import (
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/wolframdeus/exchange-rates-backend/internal/services"
)

const (
	// Ключ контекста, в котором хранятся сервисы.
	contextKeyServices = "__services"
)

type Services struct {
	// Сервис для работы с валютами.
	Currencies *services.Currencies
}

type Context struct {
	// Оригинальный контекст gin.
	Gin *gin.Context
	// Список доступных сервисов.
	Services *Services
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

// SendData успешно отправляет указанные данные по единому формату.
func (c *Context) SendData(data any) {
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

// InjectServices помещает в контекст gin список сервисов.
func (c *Context) InjectServices(curSrv *services.Currencies) {
	c.Gin.Set(contextKeyServices, &Services{Currencies: curSrv})
}

// New возвращает ссылку на новый экземпляр Context.
func New(gc *gin.Context) *Context {
	c := &Context{Gin: gc}

	// Восстанавливаем список сервисом из контекста Gin.
	if val, ok := gc.Get(contextKeyServices); ok {
		if srv, ok := val.(*Services); ok {
			c.Services = srv
		}
	}

	return c
}

// NewHandler возвращает новый обработчик и в нем вызывает переданную функцию
// с уже обернутым контекстом.
func NewHandler(f func(c *Context)) gin.HandlerFunc {
	return func(gc *gin.Context) {
		f(New(gc))
	}
}

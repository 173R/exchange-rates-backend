package context

import (
	"context"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/wolframdeus/exchange-rates-backend/internal/launchparams"
	"github.com/wolframdeus/exchange-rates-backend/internal/services"
	"github.com/wolframdeus/exchange-rates-backend/internal/services/currencies"
)

type Gin struct {
	Context
	// Оригинальный контекст gin.
	Gin *gin.Context
}

// CaptureError захватывает ошибку и отправляет её в Sentry.
func (c *Gin) CaptureError(err error) {
	// TODO: Установить IP клиента.
	hub := sentrygin.GetHubFromContext(c.Gin)
	if hub == nil {
		return
	}
	hub.CaptureException(err)
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

// InjectServices помещает в контекст gin список сервисов.
func (c *Gin) InjectServices(curSrv *currencies.Currencies, uSrv *services.Users) {
	c.inject(contextKeyServices, &Services{
		Currencies: curSrv,
		Users:      uSrv,
	})
}

// InjectLaunchParams извлекает из текущего запроса параметры запуска и
// помещает их в контекст.
func (c *Gin) InjectLaunchParams() error {
	// Извлекаем данные о параметрах запуска.
	h := c.Gin.GetHeader(HeaderLaunchParams)
	if h == "" {
		return nil
	}

	// Пытаемся извлечь параметры запуска из заголовка.
	params, err := launchparams.Derive(h)
	if err != nil {
		return err
	}

	c.inject(contextKeyLaunchParams, params)
	return nil
}

func (c *Gin) inject(key string, value any) {
	ctx := context.WithValue(c.Gin.Request.Context(), key, value)
	c.Gin.Request = c.Gin.Request.WithContext(ctx)
}

// NewGin возвращает ссылку на новый экземпляр Gin.
func NewGin(ctx *gin.Context) *Gin {
	return &Gin{
		Context: *NewContext(ctx.Request.Context()),
		Gin:     ctx,
	}
}

// NewGinHandler возвращает новый обработчик и в нем вызывает переданную функцию
// с уже обернутым контекстом.
func NewGinHandler(f func(c *Gin)) gin.HandlerFunc {
	return func(gc *gin.Context) {
		f(NewGin(gc))
	}
}

package context

import (
	"context"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/wolframdeus/exchange-rates-backend/internal/launchparams"
	"github.com/wolframdeus/exchange-rates-backend/internal/services"
	services2 "github.com/wolframdeus/exchange-rates-backend/internal/services/translations"
)

type Gin struct {
	// Оригинальный контекст gin.
	Gin *gin.Context
	// Список доступных сервисов.
	Services *Services
	// Список параметров запуска.
	LaunchParams *launchparams.Params
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
	c.Gin.JSON(400, map[string]interface{}{
		"ok":    false,
		"error": data,
	})
}

// InjectServices помещает в контекст gin список сервисов.
func (c *Gin) InjectServices(
	curSrv *services.Currencies,
	trlSrv *services2.Translations,
) {
	c.inject(contextKeyServices, &Services{
		Currencies:   curSrv,
		Translations: trlSrv,
	})
}

// InjectLaunchParams извлекает из текущего запроса параметры запуска и
// помещает их в контекст.
func (c *Gin) InjectLaunchParams() {
	// Пытаемся извлечь параметры запуска из заголовка.
	params, err := launchparams.Derive(c.Gin.GetHeader(HeaderLaunchParams))
	if err != nil {
		return
	}
	c.inject(contextKeyLaunchParams, params)
}

func (c *Gin) inject(key string, value any) {
	ctx := context.WithValue(c.Gin.Request.Context(), key, value)
	c.Gin.Request = c.Gin.Request.WithContext(ctx)
}

// NewGin возвращает ссылку на новый экземпляр Gin.
func NewGin(gc *gin.Context) *Gin {
	c := &Gin{Gin: gc}
	ctx := gc.Request.Context()

	// Восстанавливаем список сервисов из контекста.
	if srv := getServicesFromContext(ctx); srv != nil {
		c.Services = srv
	}

	// Восстанавливаем параметры запуска.
	if params := getLaunchParamsFromContext(ctx); params != nil {
		c.LaunchParams = params
	}

	return c
}

// NewGinHandler возвращает новый обработчик и в нем вызывает переданную функцию
// с уже обернутым контекстом.
func NewGinHandler(f func(c *Gin)) gin.HandlerFunc {
	return func(gc *gin.Context) {
		f(NewGin(gc))
	}
}

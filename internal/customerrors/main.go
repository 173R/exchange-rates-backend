package customerrors

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"time"
)

type Error struct {
	Original       error
	PublicErrorMsg string
	Breadcrumbs    []*sentry.Breadcrumb
}

// AddBreadcrumb добавляет новую хлебную крошку.
func (e *Error) AddBreadcrumb(bc *sentry.Breadcrumb) *Error {
	e.Breadcrumbs = append(e.Breadcrumbs, bc)
	return e
}

// Error возвращает текст ошибки для внутреннего использования.
func (e *Error) Error() string {
	return e.Original.Error()
}

// PublicError возвращает текст ошибки, готовый для публичного
// представления.
func (e *Error) PublicError() string {
	if len(e.PublicErrorMsg) > 0 {
		return e.PublicErrorMsg
	}
	if e, ok := e.Original.(*Error); ok {
		return e.PublicError()
	}
	return "Произошла неизвестная ошибка."
}

// SetPublicError изменяет текст ошибки для публичного представления, и
// добавляет лог об этом в список хлебных крошек.
func (e *Error) SetPublicError(msg string) *Error {
	e.PublicErrorMsg = msg
	e.Breadcrumbs = append(e.Breadcrumbs, &sentry.Breadcrumb{
		Type:      "debug",
		Category:  "error.message",
		Message:   fmt.Sprintf("Error message was changed to %q", msg),
		Data:      nil,
		Level:     sentry.LevelDebug,
		Timestamp: time.Now(),
	})
	return e
}

// New стандартный конструктор ошибки.
func New(err error) *Error {
	return &Error{Original: err}
}

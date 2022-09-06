package sentry

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"time"
)

// RepErrorBreadcrumb создает стандартный breadcrumb для ошибки в
// репозитории.
func RepErrorBreadcrumb(
	repName string,
	funcName string,
	params map[string]interface{},
) *sentry.Breadcrumb {
	return &sentry.Breadcrumb{
		Type:     "error",
		Category: fmt.Sprintf("repositories.%s", repName),
		Message:  "An error occurred while calling rep function.",
		Data: map[string]interface{}{
			"funcName": funcName,
			"params":   params,
		},
		Level:     sentry.LevelError,
		Timestamp: time.Now(),
	}
}

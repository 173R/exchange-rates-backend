package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wolframdeus/exchange-rates-backend/configs"
)

// Run запускает HTTP-сервер проекта.
func Run() error {
	g := gin.Default()

	g.GET("/hello-world", func(c *gin.Context) {
		c.JSON(200, map[string]string{
			"message": "its ok!",
		})
	})

	if err := g.Run(fmt.Sprintf("0.0.0.0:%d", configs.App.Port)); err != nil {
		return err
	}
	return nil
}

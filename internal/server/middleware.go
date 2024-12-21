package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
)

func (a *EchoApp) LoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			c.Error(err)
		}

		a.logger.Infof(fmt.Sprintf("%s | %s | %s", c.Request().Method, c.Request().RequestURI, c.Request().Header.Get("X-Real-IP")))
		return nil
	}
}

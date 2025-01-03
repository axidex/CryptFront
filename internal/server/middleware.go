package server

import (
	"github.com/labstack/echo/v4"
)

func (a *EchoApp) LoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			a.logger.Errorf("%s | %s | %s | %s", c.Request().Method, c.Request().RequestURI, c.Request().Header.Get("X-Real-IP"), err)
			c.Error(err)
		}

		a.logger.Infof("%s | %s | %s", c.Request().Method, c.Request().RequestURI, c.Request().Header.Get("X-Real-IP"))
		return nil
	}
}

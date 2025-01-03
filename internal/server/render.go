package server

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func (a *EchoApp) Render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		a.logger.Errorf("Got error when rendering: %s", err)
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}

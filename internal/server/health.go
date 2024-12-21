package server

import (
	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (a *EchoApp) health(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, gin.H{"status": "OK"})
}

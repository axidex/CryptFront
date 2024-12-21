package server

import (
	"fmt"
	"front/internal/models"
	"front/internal/template"
	"github.com/labstack/echo/v4/middleware"
	"net/http"

	"github.com/axidex/CryptBot/pkg/logger"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
)

type App interface {
	Run() error
	Stop(error)
}

type EchoApp struct {
	serv      *echo.Echo
	logger    logger.Logger
	apiClient *resty.Client
	port      int
	appRoutes map[string]models.Route
}

func NewEchoApp(port int, logger logger.Logger, apiClient *resty.Client, appRoutes map[string]models.Route) App {

	app := EchoApp{
		serv:      echo.New(),
		logger:    logger,
		port:      port,
		apiClient: apiClient,
		appRoutes: appRoutes,
	}

	app.InitRoutes()

	return &app
}

func (a *EchoApp) InitRoutes() {
	a.serv.Use(middleware.Recover())
	a.serv.Use(a.LoggerMiddleware)

	a.serv.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost},
		AllowHeaders: []string{
			"Content-Type",
			"Authorization",
			"Accept",
			"Origin",
		},
		MaxAge: 3600,
	}))

	a.serv.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		Skipper:            middleware.DefaultSkipper,
		XSSProtection:      "1; mode=block",
		ContentTypeNosniff: "nosniff",
		XFrameOptions:      "SAMEORIGIN",
	}))

	a.serv.GET("/", a.HomeHandler)
	a.serv.POST("/model-fields", a.ModelFieldsHandler)
	a.serv.POST("/send-to-api", a.SendToAPIHandler)

	a.serv.GET("/ping", a.health)
}

func (a *EchoApp) Run() error {
	route := fmt.Sprintf(":%d", a.port)
	a.logger.Infof("Starting server on %s", route)
	return a.serv.Start(route)
}

func (a *EchoApp) Stop(err error) {
	a.logger.Infof("stopping server: %s", err)
	err = a.serv.Close()
	if err != nil {
		a.logger.Errorf("error stopping server: %s", err)
		return
	}
	a.logger.Infof("server stopped")
}

func (a *EchoApp) HomeHandler(c echo.Context) error {
	return Render(c, http.StatusOK, template.Home(a.appRoutes))
}

// ModelFieldsHandler возвращает HTML для динамических полей
func (a *EchoApp) ModelFieldsHandler(c echo.Context) error {
	var request struct {
		Model string `json:"model"`
	}
	if err := c.Bind(&request); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request")
	}

	route, exist := a.appRoutes[request.Model]
	if !exist {
		return c.String(http.StatusBadRequest, "Unknown model")
	}

	return Render(c, http.StatusOK, template.Problem(route))
}

func (a *EchoApp) SendToAPIHandler(c echo.Context) error {
	var request struct {
		Model string            `json:"model"`
		Data  map[string]string `json:"data"`
	}
	if err := c.Bind(&request); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request")
	}

	a.logger.Infof("Got Handler request model - %s | %+v", request.Model, request.Data)

	model, exist := a.appRoutes[request.Model]
	if !exist {
		a.logger.Errorf("model %s not found", request.Model)
		return c.String(http.StatusBadRequest, "Unknown model")
	}

	resp, err := a.apiClient.R().SetQueryParams(request.Data).Post(model.Handler)
	if err != nil {
		a.logger.Errorf("error sending message to API: %s", err)
		return c.String(http.StatusBadRequest, "Invalid request")
	}

	if resp.StatusCode() != http.StatusOK {
		a.logger.Errorf("bad response status code %d", resp.StatusCode())
		return c.String(http.StatusBadRequest, "Invalid request")
	}

	return c.String(http.StatusOK, resp.String())
}

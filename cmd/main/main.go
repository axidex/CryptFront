package main

import (
	"context"
	"flag"
	"fmt"
	"front/internal/server"
	"front/internal/swagger"
	"front/internal/utils"
	"github.com/axidex/CryptBot/pkg/logger"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"github.com/oklog/run"
	"os"
	"strconv"
	"syscall"
)

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "Enable debug mode")
	flag.Parse()

	appLogger := logger.NewZapLogger()

	err := godotenv.Load()
	if err != nil {
		appLogger.Fatal("Error loading .env file")
	}

	frontPort, _ := strconv.Atoi(os.Getenv("FRONT_PORT"))
	apiHost := os.Getenv("API_HOST")
	apiPort, _ := strconv.Atoi(os.Getenv("API_PORT"))
	openApiRoute := os.Getenv("OPENAPI_ROUTE")

	apiClient := resty.New().SetBaseURL(fmt.Sprintf("http://%s:%d", apiHost, apiPort))

	//appRoutes := map[string]models.Route{
	//	"des3": {
	//		Handler: "/des3",
	//		Params: map[string]string{
	//			"l":  "5",
	//			"r":  "1",
	//			"k1": "3",
	//			"k2": "7",
	//			"k3": "5",
	//		},
	//	},
	//}

	data, err := swagger.RequestOpenApi(apiClient, openApiRoute)
	if err != nil {
		appLogger.Errorf("Got error when tried to get open api: %s", err)
		return
	}

	appRoutes, err := swagger.GetRoutes(data)
	if err != nil {
		appLogger.Errorf("Got error when parsing open api: %s", err)
		return
	}

	appLogger.Infof("Routes: %s", utils.PrettyStruct(appRoutes))

	serv := server.NewEchoApp(frontPort, appLogger, apiClient, appRoutes)

	runG := run.Group{}
	runG.Add(serv.Run, serv.Stop)
	runG.Add(run.SignalHandler(context.TODO(), syscall.SIGINT, syscall.SIGTERM))
	err = runG.Run()
	if err != nil {
		appLogger.Fatal(err)
	}
}

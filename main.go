package main

import (
	distributorhandler "challenge2016/internal/handler/distributor"
	"challenge2016/internal/service/distributor"
	"challenge2016/internal/util"
	"gofr.dev/pkg/gofr"
	"os"
)

func main() {
	app := gofr.New()

	inputFile := app.Config.Get("INPUT_FILE")

	file, err := os.Open(inputFile)
	if err != nil {
		app.Logger().Errorf("failed to open input file")
		return
	}

	err = util.LoadLocations(app.Logger(), file)
	if err != nil {
		app.Logger().Errorf("failed to process input file")
		return
	}

	file.Close()

	distributorSvc := distributor.New()
	distributorHandler := distributorhandler.New(distributorSvc)

	app.POST("/distributor", distributorHandler.Add)
	app.GET("/distributor/{name}", distributorHandler.Get)
	app.GET("/distributor/{name}/permission", distributorHandler.Check)

	app.Run()
}

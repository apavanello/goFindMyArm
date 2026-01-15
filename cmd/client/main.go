package main

import (
	"github.com/apavanello/goFindMyArm/frontend"
	"github.com/apavanello/goFindMyArm/internal/client"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

func main() {
	// Create an instance of the app structure
	app := client.NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "tmp_client",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: frontend.Assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.Startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

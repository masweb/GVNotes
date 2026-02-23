package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed db/migrations
var migrations embed.FS

func main() {
	appMigrations = migrations

	app := NewApp()

	err := wails.Run(&options.App{
		Title:                    "GVNotes",
		Width:                    1600,
		Height:                   1200,
		MinWidth:                 533,
		MinHeight:                400,
		HideWindowOnClose:        true,
		EnableDefaultContextMenu: true,

		BackgroundColour: &options.RGBA{R: 34, G: 37, B: 49, A: 1},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		CSSDragProperty:  "--wails-draggable",
		CSSDragValue:     "drag",
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Windows: &windows.Options{
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			CustomTheme: &windows.ThemeSettings{
				DarkModeTitleBar: windows.RGB(34, 37, 49),
				DarkModeBorder:   windows.RGB(34, 37, 49),
			},
		},
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: true,
				HideTitle:                  true,
				UseToolbar:                 true,
				HideToolbarSeparator:       false,
			},
			About: &mac.AboutInfo{
				Title:   "GVNotes",
				Message: "© 2026 Guillermo Valentín",
			},
		},
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

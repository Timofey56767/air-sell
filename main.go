package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/TOMMy-Net/air-sell/cache"
	"github.com/TOMMy-Net/air-sell/db"
	"github.com/TOMMy-Net/air-sell/internal"
	"github.com/TOMMy-Net/air-sell/tools"
)

func main() {
	storage, err := db.ConnectPostgres(db.PostgresConnector{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "770948",
		Database: "air_sell",
		SSL:      "disable",
	})
	if err != nil {
		panic(err)
	}

	tools.NewValidator() // активация валидатора 
	cache.NewCache() // инициализация кеша
	myApp := app.New()
	
	mainWindow := myApp.NewWindow("Air Buy")
	mainWindow.SetIcon(internal.GetIcon())
	mainWindow.Resize(fyne.Size{Width: 1000, Height: 700})
	mainWindow.SetMaster()
	mainWindow.CenterOnScreen()

	var newSettings = internal.NewSettings()
	newSettings.Window = mainWindow
	newSettings.Storage = storage
	newSettings.App = myApp
	newSettings.SignInWindow()

	mainWindow.Show()
	myApp.Run()

	go cacheKeeper(newSettings)

}


func cacheKeeper(s *internal.Settings)  {
	for {
		time.Sleep(1 * time.Second)
		a, err := s.Storage.AllAirports()
		if err != nil {
			continue
		}
		for _, v := range a {
			cache.AddAirport(v)
		}
	}
}


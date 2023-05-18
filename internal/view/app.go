package view

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/Tyz3/nymgraph/cmd/app/config"
	"github.com/Tyz3/nymgraph/internal/controller"
)

var (
	logo = fyne.NewStaticResource("logo", config.NymLogoBin)
)

type App struct {
	app        fyne.App
	controller *controller.Controller

	ChooseClientWindow *ChooseClientWindow
	HomeWindow         *HomeWindow
}

func NewApp(appName string, controller *controller.Controller) *App {
	a := &App{
		app:        app.NewWithID(appName),
		controller: controller,
	}

	// Окно выбора аккаунта
	a.ChooseClientWindow = NewChooseClientWindow(controller, a.app, "Nymgraph", logo)
	a.HomeWindow = NewHomeWindow(controller, a.app, "Nymgraph", logo)

	// Окно взаимодействия
	a.ChooseClientWindow.OnSubmit = func() {
		if err := a.HomeWindow.Load(); err != nil {
			fmt.Println(err)
		}
		a.HomeWindow.Window.Show()
		a.ChooseClientWindow.Window.Close()
	}
	a.HomeWindow.Window.SetOnClosed(func() {
		if err := a.controller.NymClient.Close(); err != nil {
			fmt.Println(err)
		}
		a.app.Quit()
	})

	return a
}

func (a *App) Run() error {
	if err := a.ChooseClientWindow.Load(); err != nil {
		return err
	} else {
		a.ChooseClientWindow.Window.ShowAndRun()
	}

	return nil
}

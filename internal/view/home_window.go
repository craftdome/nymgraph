package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/Tyz3/nymgraph/internal/controller"
)

type HomeWindow struct {
	App        fyne.App
	Controller *controller.Controller
	Window     fyne.Window
}

func NewHomeWindow(controller *controller.Controller, app fyne.App, title string, icon fyne.Resource) *HomeWindow {
	w := &HomeWindow{
		App:        app,
		Controller: controller,
		Window:     app.NewWindow(title),
	}

	w.Window.SetIcon(icon)
	w.Window.Resize(fyne.NewSize(350, 450))
	w.Window.CenterOnScreen()

	w.Window.SetContent(widget.NewLabel("asdasd"))

	return w
}

func (w *HomeWindow) Load() error {

	return nil
}

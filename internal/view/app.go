package view

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Tyz3/nymgraph/cmd/app/config"
	"github.com/Tyz3/nymgraph/internal/entity"
	"github.com/Tyz3/nymgraph/internal/model"
	"github.com/Tyz3/nymgraph/internal/service"
	"github.com/pkg/errors"
	"os"
	"sync"
	"time"
)

var (
	logo = fyne.NewStaticResource("logo", config.NymLogoBin)
)

type App struct {
	app        fyne.App
	controller *service.Service

	models []*model.Pseudonym
	menu   *fyne.Menu

	openedChat map[string]*HomeWindow
}

func NewApp(appName string, controller *service.Service) *App {
	a := &App{
		app:        app.NewWithID(appName),
		controller: controller,

		menu:       fyne.NewMenu("Nymgraph"),
		openedChat: make(map[string]*HomeWindow),
	}

	a.app.Settings().SetTheme(theme.DarkTheme())

	return a
}

func (a *App) Run() {
	a.Load()

	if desk, ok := a.app.(desktop.App); ok {
		desk.SetSystemTrayMenu(a.menu)
		desk.SetSystemTrayIcon(logo)
	}

	a.app.Run()
}

func (a *App) Close() {
	var wg sync.WaitGroup
	for _, m := range a.models {
		m := m
		wg.Add(1)
		go func() {
			defer wg.Done()
			if m.NymClient != nil && m.NymClient.IsOnline() {
				m.NymClient.Close()
			}
		}()
	}
	wg.Wait()
}

func (a *App) Load() {
	settItem := fyne.NewMenuItem("Settings", func() {

	})
	settItem.Icon = theme.SettingsIcon()
	a.menu.Items = append(a.menu.Items,
		fyne.NewMenuItemSeparator(),
		settItem,
	)

	a.Reload()

	go func() {
		for ; ; time.Sleep(time.Second) {
			a.update()
		}
	}()
}

func (a *App) Reload() {
	pseudonyms, err := a.controller.Pseudonyms.GetAll()
	if err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrapf(err, "controller.Pseudonyms.GetAll"))
		return
	}

	for _, pseudonym := range pseudonyms {
		a.addMenuClient(pseudonym)
	}
}

func (a *App) update() {
	for _, m := range a.models {
		m := m
		label := fmt.Sprintf("%s (%s)", m.Pseudonym.Name, m.Pseudonym.Server)

		// Если добавлен новый клиент, то у него не будет MenuItem
		if m.MenuItem == nil {
			m.MenuItem = fyne.NewMenuItem(label, nil)
			m.MenuItem.Icon = theme.NewPrimaryThemedResource(theme.LoginIcon())
			m.MenuItem.Disabled = true
			a.menu.Items = append([]*fyne.MenuItem{m.MenuItem}, a.menu.Items...)
			//fmt.Println(label, "Добавлен элемент меню")
		}

		if m.NymClient.IsOnline() {
			//fmt.Println(label, "IsOnline")
			continue
		}

		// Если клиент не онлайн, то пробуем подключиться
		if err := m.NymClient.Dial(); err != nil {
			//fmt.Println(label, "не удалось подключиться")
			m.MenuItem.Action = nil
			m.MenuItem.Disabled = true
			continue
		}

		// Включаем прослушку входящих сообщений
		if err := m.NymClient.ListenAndServe(); err != nil {
			m.MenuItem.Action = nil
			m.MenuItem.Disabled = true
			continue
		}

		// Если подключение к клиенту установлено, активируем кнопку в меню
		m.MenuItem.Disabled = false
		m.MenuItem.Action = func() {
			// Если Checked = true, то чат-клиент уже открыт, вызываем на него фокус
			if m.MenuItem.Checked {
				//fmt.Println(label, "окно уже открыто - фокус")
				a.openedChat[m.Pseudonym.Name].Window.RequestFocus()
			} else {
				var w *HomeWindow
				if opened, exists := a.openedChat[m.Pseudonym.Name]; exists {
					//fmt.Println(label, "открываем созданное окно")
					w = opened
				} else {
					//fmt.Println(label, "создаём новое окно")
					w = NewHomeWindow(a.controller, a.app, fmt.Sprintf("Connected - %s", label), logo, m)
					w.Load()
					a.openedChat[m.Pseudonym.Name] = w
				}
				w.Window.Show()
				w.Window.SetCloseIntercept(func() {
					//fmt.Println(label, "скрываем окно")
					w.Window.Hide()
					m.MenuItem.Checked = false
					a.menu.Refresh()
				})
				m.MenuItem.Checked = true
				a.menu.Refresh()
			}
		}
	}
	a.menu.Refresh()
}

func (a *App) addMenuClient(e *entity.Pseudonym) {
	a.models = append(a.models, &model.Pseudonym{
		Pseudonym: e,
		NymClient: a.controller.NymClient.New(e),
	})
}

func (a *App) MessageError(err error) {
	w := a.app.NewWindow("Error")
	w.SetIcon(theme.InfoIcon())
	w.Resize(fyne.NewSize(400, 100))
	w.SetContent(container.NewBorder(
		nil,
		nil,
		container.NewGridWrap(
			fyne.NewSize(70, 70),
			widget.NewIcon(theme.ErrorIcon()),
		),
		nil,
		widget.NewLabel(err.Error()),
	))
	w.Show()
}

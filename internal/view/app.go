package view

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/craftdome/nymgraph/cmd/app/config"
	"github.com/craftdome/nymgraph/internal/entity"
	"github.com/craftdome/nymgraph/internal/service"
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

	pseudonyms []*entity.Pseudonym
	menu       *fyne.Menu

	settingsWindow *SettingsWindow
	openedChats    map[string]*HomeWindow
	mu             sync.Mutex
}

func NewApp(appName string, controller *service.Service) *App {
	a := &App{
		app:        app.NewWithID(appName),
		controller: controller,

		menu:        fyne.NewMenu("Nymgraph"),
		openedChats: make(map[string]*HomeWindow),
	}

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
	for _, chat := range a.openedChats {
		chat := chat
		wg.Add(1)
		go func() {
			defer wg.Done()
			chat.Close()
		}()
	}
	wg.Wait()

	if a.controller.Config.DeleteHistoryAfterQuit() {
		if err := a.controller.Sent.Truncate(); err == nil {
			fmt.Println("Sent data cleaned")
		}

		if err := a.controller.Received.Truncate(); err == nil {
			fmt.Println("Received/Replies data cleaned")
		}
	}
}

func (a *App) Load() {
	settItem := fyne.NewMenuItem("Settings", func() {
		if a.settingsWindow == nil {
			a.settingsWindow = NewSettingsWindow(a.controller, a.app, "Settings", theme.SettingsIcon())
			a.settingsWindow.Window.Show()
			a.settingsWindow.Load()
		} else {
			a.settingsWindow.Window.RequestFocus()
		}

		a.settingsWindow.Window.SetCloseIntercept(func() {
			a.settingsWindow.Unload()
			a.settingsWindow.Window.Close()
			a.settingsWindow = nil
		})
		a.settingsWindow.OnCreate = func(pseudonym *entity.Pseudonym) {
			a.pseudonyms = append(a.pseudonyms, pseudonym)
			a.update()
		}
		a.settingsWindow.OnUpdate = func(pseudonym *entity.Pseudonym) {
			for i, p := range a.pseudonyms {
				if p.ID == pseudonym.ID {
					a.pseudonyms[i].Name = pseudonym.Name
					a.pseudonyms[i].Server = pseudonym.Server
					openedChat := a.openedChats[pseudonym.Name]
					a.mu.Lock()
					delete(a.openedChats, pseudonym.Name)
					openedChat.Close()
					for i, mi := range a.menu.Items {
						if mi == openedChat.menuItem {
							a.menu.Items = append(a.menu.Items[:i], a.menu.Items[i+1:]...)
							break
						}
					}
					a.mu.Unlock()
					break
				}
			}
			a.update()
		}
		a.settingsWindow.OnDelete = func(pseudonym *entity.Pseudonym) {
			for i, p := range a.pseudonyms {
				if p.ID == pseudonym.ID {
					a.pseudonyms = append(a.pseudonyms[:i], a.pseudonyms[i+1:]...)
					break
				}
			}
			openedChat := a.openedChats[pseudonym.Name]
			a.mu.Lock()
			delete(a.openedChats, pseudonym.Name)
			openedChat.Close()
			for i, mi := range a.menu.Items {
				if mi == openedChat.menuItem {
					a.menu.Items = append(a.menu.Items[:i], a.menu.Items[i+1:]...)
					break
				}
			}
			a.mu.Unlock()
			a.update()
		}
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

	a.pseudonyms = pseudonyms
}

func (a *App) update() {
	var wg sync.WaitGroup
	for _, pseudonym := range a.pseudonyms {
		pseudonym := pseudonym

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer a.menu.Refresh()

			// Search an opened chat
			openedChat, exists := a.openedChats[pseudonym.Name]
			if !exists {
				// Create connect to the nym-client by pseudonym settings
				connect := a.controller.NymClient.New(pseudonym)
				// Create Window for opening when MenuItem will be tapped
				openedChat = NewHomeWindow(
					a.controller, a.app,
					fmt.Sprintf("Connected - %s (%s)", pseudonym.Name, pseudonym.Server),
					logo,
					pseudonym,
					connect,
				)
				// Handle socket closed
				connect.OnCloseCallback = func() {
					fmt.Println("Connection closed")
				}
				// Update map and menu
				a.mu.Lock()
				a.openedChats[pseudonym.Name] = openedChat
				a.mu.Unlock()
				a.menu.Items = append([]*fyne.MenuItem{openedChat.MenuItem()}, a.menu.Items...)
				a.menu.Refresh()
			}

			// Check connect by online
			if openedChat.ClientConnect().IsOnline() {
				openedChat.UpdateSelfAddress()
				return
			}

			// Try to connect with the nym-client
			if err := openedChat.ClientConnect().Dial(); err != nil {
				// Disable MenuItem when failure
				openedChat.MenuItem().Action = nil
				openedChat.MenuItem().Disabled = true
				return
			}

			// Start listening incoming message from the nym-client
			if err := openedChat.ClientConnect().ListenAndServe(); err != nil {
				openedChat.MenuItem().Action = nil
				openedChat.MenuItem().Disabled = true
				return
			}

			openedChat.MenuItem().Disabled = false
			openedChat.MenuItem().Action = func() {
				if openedChat.MenuItem().Checked {
					openedChat.Window.RequestFocus()
				} else {
					openedChat.Load()
					openedChat.Window.Show()
					openedChat.Window.SetCloseIntercept(func() {
						openedChat.Window.Hide()
						openedChat.MenuItem().Checked = false
						a.menu.Refresh()
						openedChat.Unload()
					})
					openedChat.MenuItem().Checked = true
					a.menu.Refresh()
				}
			}
		}()
		wg.Wait()
	}
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

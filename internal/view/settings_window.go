package view

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Tyz3/nymgraph/internal/entity"
	"github.com/Tyz3/nymgraph/internal/service"
	"github.com/Tyz3/nymgraph/internal/view/custom_widget"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"regexp"
)

var (
	proxyRegex = regexp.MustCompile(`((.+):(.+)@)?(\d+.\d+.\d+.\d+):(\d+)`)
)

type SettingsWindow struct {
	App        fyne.App
	Controller *service.Service
	Window     fyne.Window

	pseudonyms             []*entity.Pseudonym
	addPseudonym           *entity.Pseudonym
	useProxy               bool
	proxyCredentials       string
	deleteHistoryAfterQuit bool

	list                        *custom_widget.TList[*entity.Pseudonym]
	addButton                   *widget.Button
	deleteHistoryAfterQuitCheck *widget.Check
	useProxyCheck               *widget.Check
	proxyCredentialsEntry       *widget.Entry

	OnUpdate func(*entity.Pseudonym)
	OnCreate func(*entity.Pseudonym)
	OnDelete func(*entity.Pseudonym)
}

func NewSettingsWindow(controller *service.Service, app fyne.App, title string, icon fyne.Resource) *SettingsWindow {
	w := &SettingsWindow{
		App:        app,
		Controller: controller,
		Window:     app.NewWindow(title),

		addPseudonym: &entity.Pseudonym{},
	}

	w.Window.SetIcon(icon)
	w.Window.Resize(fyne.NewSize(480, 450))
	w.Window.CenterOnScreen()

	{
		validate := validator.New()

		clientNameEntry := widget.NewEntry()
		clientNameEntry.PlaceHolder = "max"
		clientNameEntry.Validator = func(s string) error {
			w.addPseudonym.Name = s
			if err := validate.StructPartial(w.addPseudonym, "Name"); err != nil {
				return err
			}

			return nil
		}

		serverNameEntry := widget.NewEntry()
		serverNameEntry.Validator = func(s string) error {
			w.addPseudonym.Server = s
			if err := validate.StructPartial(w.addPseudonym, "Server"); err != nil {
				return err
			}

			return nil
		}
		serverNameEntry.PlaceHolder = "127.0.0.1:1977"

		w.deleteHistoryAfterQuitCheck = widget.NewCheck("Delete chat history after quit", func(b bool) {
			w.deleteHistoryAfterQuit = b
			w.Controller.Config.SetDeleteHistoryAfterQuit(b)
		})

		w.list = new(custom_widget.TList[*entity.Pseudonym])
		w.list.Init(6, func() fyne.CanvasObject {
			return custom_widget.NewClientEntry()
		})
		w.list.SetUpdateItemFunc(func(id int, pseudonym *entity.Pseudonym, object fyne.CanvasObject) {
			entry := object.(*custom_widget.ClientEntry)

			entry.SetModel(pseudonym)
			entry.OnDeleteTapped = func(pseudonym *entity.Pseudonym) {
				dialog.ShowConfirm(
					"Confirmation",
					fmt.Sprintf("Confirm client %s deletion", pseudonym.Name),
					func(b bool) {
						if !b {
							return
						}
						deleted, err := w.Controller.Pseudonyms.Delete(pseudonym.ID)
						if err != nil {
							dialog.ShowError(errors.Wrapf(err, "Controller.Pseudonyms.Delete id=%d", pseudonym.ID), w.Window)
							return
						}

						for i, p := range w.pseudonyms {
							if p.ID == deleted.ID {
								w.pseudonyms = append(w.pseudonyms[:i], w.pseudonyms[i+1:]...)
								break
							}
						}
						w.list.Reload()
						w.list.UnselectAll()
						w.OnDelete(deleted)
					},
					w.Window,
				)
			}
			entry.OnEditTapped = func(pseudonym *entity.Pseudonym) {
				w.addPseudonym = pseudonym
				clientNameEntry.SetText(pseudonym.Name)
				serverNameEntry.SetText(pseudonym.Server)
				dialog.ShowForm(
					"Editing",
					"Save",
					"Cancel",
					[]*widget.FormItem{
						{
							Text:     "Pseudonym Name",
							Widget:   clientNameEntry,
							HintText: "give a name to your nym-client",
						},
						{
							Text:     "Server",
							Widget:   serverNameEntry,
							HintText: "nym-client endpoint like 127.0.0.1:1977",
						},
					},
					func(b bool) {
						if !b {
							return
						}

						updated, err := w.Controller.Pseudonyms.Update(w.addPseudonym.ID, w.addPseudonym.Name, w.addPseudonym.Server)
						if err != nil {
							dialog.ShowError(errors.Wrap(err, "Controller.Pseudonyms.Update"), w.Window)
							return
						}

						for i, p := range w.pseudonyms {
							if p.ID == updated.ID {
								w.pseudonyms[i] = updated
								break
							}
						}

						w.list.Reload()
						w.OnUpdate(updated)
					},
					w.Window,
				)
			}
			entry.OnLeftClick = func(event *fyne.PointEvent) {
				w.list.Select(id)
			}
		})
		w.list.SetDataSourceFunc(func() []*entity.Pseudonym {
			return w.pseudonyms
		})

		w.addButton = widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
			dialog.ShowForm(
				"Creating",
				"Create",
				"Cancel",
				[]*widget.FormItem{
					{
						Text:     "Pseudonym name",
						Widget:   clientNameEntry,
						HintText: "give a name to your nym-client",
					},
					{
						Text:     "Сервер",
						Widget:   serverNameEntry,
						HintText: "nym-client endpoint like 127.0.0.1:1977",
					},
				},
				func(b bool) {
					if !b {
						return
					}

					created, err := w.Controller.Pseudonyms.Create(w.addPseudonym.Name, w.addPseudonym.Server)
					if err != nil {
						dialog.ShowError(err, w.Window)
						return
					}

					w.pseudonyms = append(w.pseudonyms, created)
					w.list.Reload()
					w.OnCreate(created)
				},
				w.Window,
			)
		})

		w.useProxyCheck = widget.NewCheck("Use SOCKS5", func(b bool) {
			w.useProxy = b
			w.Controller.Config.UseProxy(b)
			if b {
				w.proxyCredentialsEntry.Enable()
			} else {
				w.proxyCredentialsEntry.Disable()
			}
		})
		// TODO
		w.useProxyCheck.Disable()

		w.proxyCredentialsEntry = widget.NewEntry()
		w.proxyCredentialsEntry.PlaceHolder = "[user:pass@]addr:port"
		w.proxyCredentialsEntry.Disable()
		w.proxyCredentialsEntry.Validator = func(s string) error {
			if s == "" && !w.useProxy {
				return nil
			}

			if !proxyRegex.MatchString(s) {
				return errors.New("invalid format: [user:pass@]addr:port")
			}

			w.Controller.Config.SetProxy(s)
			return nil
		}

		w.Window.SetContent(
			container.NewBorder(
				container.NewVBox(
					container.NewBorder(
						nil,
						nil,
						w.useProxyCheck,
						nil,
						w.proxyCredentialsEntry,
					),
					w.deleteHistoryAfterQuitCheck,
				),
				w.addButton,
				nil,
				nil,
				w.list,
			),
		)
	}

	return w
}

func (w *SettingsWindow) Load() {
	pseudonyms, err := w.Controller.Pseudonyms.GetAll()
	if err != nil {
		dialog.ShowError(err, w.Window)
		return
	}

	w.useProxy = w.Controller.UsingProxy()
	w.proxyCredentials = w.Controller.GetProxy()
	w.deleteHistoryAfterQuit = w.Controller.DeleteHistoryAfterQuit()

	w.pseudonyms = pseudonyms
	w.update()
}

func (w *SettingsWindow) Unload() {
	w.pseudonyms = nil
	w.addPseudonym = &entity.Pseudonym{}
	w.proxyCredentials = ""
}

func (w *SettingsWindow) reload() {
	w.Unload()
	w.Load()
}

func (w *SettingsWindow) update() {
	w.list.Reload()

	w.useProxyCheck.SetChecked(w.useProxy)
	w.proxyCredentialsEntry.SetText(w.proxyCredentials)
	w.deleteHistoryAfterQuitCheck.SetChecked(w.deleteHistoryAfterQuit)
}

package view

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Tyz3/fyne-extra/xwidget"
	"github.com/Tyz3/nymgraph/internal/controller"
	"github.com/Tyz3/nymgraph/internal/entity"
	"github.com/Tyz3/nymgraph/internal/model"
	"github.com/Tyz3/nymgraph/internal/view/custom_widget"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

type ChooseClientWindow struct {
	App        fyne.App
	Controller *controller.Controller
	Window     fyne.Window

	selectedClient *model.Client

	clientsList   *xwidget.PagedList
	addPseudonym  *entity.Pseudonym
	addButton     *widget.Button
	connectButton *widget.Button

	OnSubmit func()
}

func NewChooseClientWindow(controller *controller.Controller, app fyne.App, title string, icon fyne.Resource) *ChooseClientWindow {
	w := &ChooseClientWindow{
		App:        app,
		Controller: controller,
		Window:     app.NewWindow(title),
	}

	w.Window.SetIcon(icon)
	w.Window.Resize(fyne.NewSize(350, 450))
	w.Window.CenterOnScreen()

	{
		validate := validator.New()

		clientNameEntry := widget.NewEntry()
		clientNameEntry.Validator = func(s string) error {
			return validate.StructPartial(w.addPseudonym, "Name")
		}
		clientNameEntry.OnChanged = func(s string) {
			w.addPseudonym.Name = s
		}
		serverNameEntry := widget.NewEntry()
		serverNameEntry.Validator = func(s string) error {
			return validate.StructPartial(w.addPseudonym, "Server")
		}
		serverNameEntry.OnChanged = func(s string) {
			w.addPseudonym.Server = s
		}

		w.addPseudonym = &entity.Pseudonym{}
		w.clientsList = xwidget.NewPagedList(6, func() fyne.CanvasObject {
			return custom_widget.NewClientEntry()
		})
		w.clientsList.SetUpdateItemFunc(func(id widget.ListItemID, object fyne.CanvasObject) {
			client := w.clientsList.GetFilteredItems()[id].(*model.Client)
			entry := object.(*custom_widget.ClientEntry)

			entry.SetModel(client)
			entry.OnDeleteTapped = func(client *model.Client) {
				if _, err := w.Controller.Pseudonyms.Delete(client.Pseudonym.ID); err != nil {
					dialog.ShowError(errors.Wrapf(err, "Controller.Pseudonyms.Delete id=%d", client.Pseudonym.ID), w.Window)
					return
				}
				w.clientsList.Refresh()
			}
			entry.OnEditTapped = func(client *model.Client) {
				w.addPseudonym = client.Pseudonym
				clientNameEntry.SetText(client.Pseudonym.Name)
				serverNameEntry.SetText(client.Pseudonym.Server)
				dialog.ShowForm(
					"Редактирование клиента",
					"Сохранить",
					"Отменить",
					[]*widget.FormItem{
						{
							Text:     "Название клиента",
							Widget:   clientNameEntry,
							HintText: "Дайте название своему nym-client",
						},
						{
							Text:     "Сервер",
							Widget:   serverNameEntry,
							HintText: "Конечная точка подключения к nym-client",
						},
					},
					func(b bool) {
						if b {
							if _, err := w.Controller.Pseudonyms.Update(w.addPseudonym.ID, w.addPseudonym.Name, w.addPseudonym.Server); err != nil {
								dialog.ShowError(errors.Wrap(err, "Controller.Pseudonyms.Update"), w.Window)
								//return
							}
							w.clientsList.Reload()
						}
					},
					w.Window,
				)
			}
			entry.OnLeftClick = func(event *fyne.PointEvent) {
				w.clientsList.Select(id)
			}
		})
		w.clientsList.OnSelected = func(id widget.ListItemID) {
			w.selectedClient = w.clientsList.GetFilteredItems()[id].(*model.Client)
			w.connectButton.Enable()
		}
		w.clientsList.SetDataSourceFunc(w.dataSource)
		w.addButton = widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
			dialog.ShowForm(
				"Новый клиент",
				"Добавить",
				"Отменить",
				[]*widget.FormItem{
					{
						Text:     "Название клиента",
						Widget:   clientNameEntry,
						HintText: "Дайте название своему nym-client",
					},
					{
						Text:     "Сервер",
						Widget:   serverNameEntry,
						HintText: "Конечная точка подключения к nym-client",
					},
				},
				func(b bool) {
					if b {
						if _, err := w.Controller.Pseudonyms.Create(w.addPseudonym.Name, w.addPseudonym.Server); err != nil {
							dialog.ShowError(errors.Wrapf(err, "Controller.Pseudonyms.Create Name=%s, Server=%s", w.addPseudonym.Name, w.addPseudonym.Server), w.Window)
							return
						}
						w.clientsList.Reload()
					}
				},
				w.Window,
			)
		})
		w.connectButton = widget.NewButtonWithIcon("Подключиться", theme.LoginIcon(), func() {
			fmt.Println(1)
			if err := w.Controller.NymClient.Dial(w.selectedClient.Pseudonym); err != nil {
				dialog.ShowError(errors.Wrap(err, "Controller.Pseudonyms.GetAll"), w.Window)
			}
			fmt.Println(2)

			if w.OnSubmit != nil {
				w.OnSubmit()
			}
		})
		w.connectButton.Importance = widget.HighImportance
		w.connectButton.Disable()

		w.Window.SetContent(
			container.NewBorder(
				nil,
				container.NewGridWithColumns(2,
					w.addButton,
					w.connectButton,
				),
				nil,
				nil,
				w.clientsList,
			),
		)
	}

	return w
}

func (w *ChooseClientWindow) dataSource() []any {
	clients, err := w.Controller.Pseudonyms.GetAll()
	if err != nil {
		dialog.ShowError(errors.Wrap(err, "Controller.Pseudonyms.GetAll"), w.Window)
		return []any{}
	}

	entries := make([]any, 0, len(clients))

	for _, b := range clients {
		entries = append(entries, &model.Client{Pseudonym: b})
	}

	return entries
}

func (w *ChooseClientWindow) Load() error {
	w.clientsList.Reload()

	return nil
}

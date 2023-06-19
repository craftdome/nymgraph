package view

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Tyz3/fyne-extra/xwidget"
	"github.com/Tyz3/nymgraph/internal/entity"
	"github.com/Tyz3/nymgraph/internal/model"
	"github.com/Tyz3/nymgraph/internal/service"
	"github.com/Tyz3/nymgraph/internal/view/custom_widget"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

type ChooseClientWindow struct {
	App        fyne.App
	Controller *service.Service
	Window     fyne.Window

	selectedClient *model.Pseudonym
	addPseudonym   *entity.Pseudonym

	list          *xwidget.PagedList
	addButton     *widget.Button
	connectButton *widget.Button

	OnSubmit func(*model.Pseudonym)
}

func NewChooseClientWindow(controller *service.Service, app fyne.App, title string, icon fyne.Resource) *ChooseClientWindow {
	w := &ChooseClientWindow{
		App:        app,
		Controller: controller,
		Window:     app.NewWindow(title),

		addPseudonym: &entity.Pseudonym{},
		OnSubmit:     func(*model.Pseudonym) {},
	}

	w.Window.SetIcon(icon)
	w.Window.Resize(fyne.NewSize(480, 450))
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

		w.list = xwidget.NewPagedList(6, func() fyne.CanvasObject {
			return custom_widget.NewClientEntry()
		})
		w.list.SetUpdateItemFunc(func(id widget.ListItemID, object fyne.CanvasObject) {
			client := w.list.GetFilteredItems()[id].(*model.Pseudonym)
			entry := object.(*custom_widget.ClientEntry)

			entry.SetModel(client)
			entry.OnDeleteTapped = func(client *model.Pseudonym) {
				dialog.ShowConfirm(
					"Confirmation",
					fmt.Sprintf("Confirm client %s deletion", client.Pseudonym.Name),
					func(b bool) {
						if !b {
							return
						}
						if _, err := w.Controller.Pseudonyms.Delete(client.Pseudonym.ID); err != nil {
							dialog.ShowError(errors.Wrapf(err, "Controller.Pseudonyms.Delete id=%d", client.Pseudonym.ID), w.Window)
							return
						}
						w.list.Reload()
						w.list.UnselectAll()
						if len(w.list.GetItems()) == 0 {
							w.connectButton.Disable()
						}
					},
					w.Window,
				)
			}
			entry.OnEditTapped = func(client *model.Pseudonym) {
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
							w.list.Reload()
						}
					},
					w.Window,
				)
			}
			entry.OnLeftClick = func(event *fyne.PointEvent) {
				w.list.Select(id)
			}
		})
		w.list.OnSelected = func(id widget.ListItemID) {
			w.selectedClient = w.list.GetFilteredItems()[id].(*model.Pseudonym)
			if w.selectedClient.NymClient != nil {
				w.connectButton.Enable()
			} else {
				w.connectButton.Disable()
			}
		}
		w.list.SetDataSourceFunc(w.dataSource)

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
						w.list.Reload()
					}
				},
				w.Window,
			)
		})

		w.connectButton = widget.NewButtonWithIcon("Подключиться", theme.LoginIcon(), func() {
			w.connectButton.Disable()
			defer w.connectButton.Enable()

			w.OnSubmit(w.selectedClient)
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
				w.list,
			),
		)
	}

	return w
}

func (w *ChooseClientWindow) dataSource() []any {
	pseudonyms, err := w.Controller.Pseudonyms.GetAll()
	if err != nil {
		dialog.ShowError(errors.Wrap(err, "Controller.Pseudonyms.GetAll"), w.Window)
		return []any{}
	}

	entries := make([]any, 0, len(pseudonyms))
	for _, pseudonym := range pseudonyms {
		entries = append(entries, &model.Pseudonym{
			Pseudonym: pseudonym,
			NymClient: w.Controller.NymClient.New(pseudonym),
		})
	}

	return entries
}

func (w *ChooseClientWindow) Load() error {
	w.list.Reload()

	return nil
}

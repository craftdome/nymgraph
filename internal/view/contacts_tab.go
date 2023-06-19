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
	"github.com/Tyz3/nymgraph/pkg/utils"
	"github.com/pkg/errors"
	"regexp"
)

var (
	NymAddrRegexp = regexp.MustCompile("([A-z0-9]{44}).([A-z0-9]{44})@([A-z0-9]{44})")
)

type ContactsTab struct {
	App        fyne.App
	Controller *service.Service
	Window     fyne.Window
	TabItem    *container.TabItem

	model *model.Pseudonym

	list      *xwidget.PagedList
	addButton *widget.Button

	OnUpdate func()
}

func NewContactsTab(app fyne.App, controller *service.Service, window fyne.Window, title string, icon fyne.Resource, pseudonym *model.Pseudonym) *ContactsTab {
	t := &ContactsTab{
		App:        app,
		Controller: controller,
		Window:     window,
		TabItem:    container.NewTabItemWithIcon(title, icon, nil),

		model: pseudonym,
	}

	{
		t.list = xwidget.NewPagedList(100, func() fyne.CanvasObject {
			entry := custom_widget.NewContactEntry()

			entry.AddContextMenuItem("Edit", theme.DocumentCreateIcon(), func() {
				editAliasEntry := widget.NewEntry()
				editAliasEntry.SetText(entry.GetModel().Alias)
				editAliasEntry.PlaceHolder = "type a contact name"
				editAliasEntry.Validator = func(s string) error {
					if s == entry.GetModel().Alias {
						return nil
					}
					for _, item := range t.list.GetItems() {
						contact := item.(*entity.Contact)
						if contact.Alias == s {
							return errors.New("same alias already exists")
						}
					}
					return nil
				}

				editAddressEntry := widget.NewEntry()
				editAddressEntry.SetText(entry.GetModel().Address)
				editAddressEntry.PlaceHolder = "type a nym-address"
				editAddressEntry.Validator = func(s string) error {
					if !NymAddrRegexp.MatchString(s) {
						return errors.New("incorrect nym-address format (must be like a44.b44@c44)")
					}

					if entry.GetModel().Address == s {
						return nil
					}

					for _, item := range t.list.GetItems() {
						contact := item.(*entity.Contact)
						if contact.Address == s {
							return errors.New(fmt.Sprintf("this nym-address already belongs to %s", contact.Alias))
						}
					}

					return nil
				}

				editForm := dialog.NewForm(
					"Edit contact",
					"Keep changes",
					"Cancel",
					[]*widget.FormItem{
						widget.NewFormItem("Alias", editAliasEntry),
						widget.NewFormItem("Address", editAddressEntry),
					},
					func(b bool) {
						if !b {
							return
						}

						if _, err := t.Controller.Contacts.Update(entry.GetModel().ID, editAddressEntry.Text, editAliasEntry.Text); err != nil {
							dialog.ShowError(errors.Wrap(err, "Controller.Contacts.Update"), t.Window)
							return
						}

						utils.ShowSplash(fmt.Sprintf("Contact %s is edited", entry.GetModel().Alias))
						t.list.Reload()
						t.OnUpdate()
					},
					t.Window,
				)
				editForm.Resize(fyne.NewSize(400, 0))
				editForm.Show()
			})
			entry.AddContextMenuItem("Delete", theme.DeleteIcon(), func() {
				if _, err := t.Controller.Contacts.Delete(entry.GetModel().ID); err != nil {
					dialog.ShowError(errors.Wrapf(err, "Controller.Contacts.Delete %d", entry.GetModel().ID), t.Window)
					return
				}
				utils.ShowSplash(fmt.Sprintf("Contact %s is deleted", entry.GetModel().Alias))
				t.list.Reload()
			})
			return entry
		})
		t.list.SetDataSourceFunc(t.dataSource)
		t.list.SetUpdateItemFunc(func(id widget.ListItemID, object fyne.CanvasObject) {
			contact := t.list.GetFilteredItems()[id].(*entity.Contact)
			entry := object.(*custom_widget.ContactEntry)

			entry.SetModel(contact)
			entry.OnLeftClick = func(event *fyne.PointEvent) {
				t.list.Select(id)
			}
			entry.OnCopyButtonClick = func() {
				if err := utils.CopyToClipboard(contact.Address); err != nil {
					dialog.ShowError(err, t.Window)
					return
				}
				utils.ShowSplash("copied")
			}
		})

		addAliasEntry := widget.NewEntry()
		addAliasEntry.PlaceHolder = "type a contact name"
		addAliasEntry.Validator = func(s string) error {
			for _, item := range t.list.GetItems() {
				contact := item.(*entity.Contact)
				if contact.Alias == s {
					return errors.New("same alias already exists")
				}
			}
			return nil
		}

		addAddressEntry := widget.NewEntry()
		addAddressEntry.PlaceHolder = "type a nym-address"
		addAddressEntry.Validator = func(s string) error {
			if !NymAddrRegexp.MatchString(s) {
				return errors.New("incorrect nym-address format (must be like a44.b44@c44)")
			}

			for _, item := range t.list.GetItems() {
				contact := item.(*entity.Contact)
				if contact.Address == s {
					return errors.New("same nym-address already exists")
				}
			}

			return nil
		}

		addForm := dialog.NewForm(
			"Create contact",
			"Create",
			"Cancel",
			[]*widget.FormItem{
				widget.NewFormItem("Alias", addAliasEntry),
				widget.NewFormItem("Address", addAddressEntry),
			},
			func(b bool) {
				if !b {
					return
				}

				if _, err := t.Controller.Contacts.Create(t.model.Pseudonym.ID, addAddressEntry.Text, addAliasEntry.Text); err != nil {
					dialog.ShowError(errors.Wrap(err, "Controller.Contacts.Create"), t.Window)
					return
				}

				utils.ShowSplash(fmt.Sprintf("Contact %s is created", addAliasEntry.Text))
				t.list.Reload()
				t.OnUpdate()
			},
			t.Window,
		)
		addForm.Resize(fyne.NewSize(400, 0))

		t.addButton = widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
			t.addButton.Disable()
			addForm.Show()
			t.addButton.Enable()
		})

		t.TabItem.Content = container.NewBorder(
			nil,
			t.addButton,
			nil,
			nil,
			t.list,
		)
	}

	return t
}

func (t *ContactsTab) Load() {
	t.list.Reload()
}

func (t *ContactsTab) dataSource() []any {
	contacts, err := t.Controller.Contacts.GetAll(t.model.Pseudonym.ID)
	if err != nil {
		dialog.ShowError(errors.Wrapf(err, "Controller.Contacts.GetAll %d", t.model.Pseudonym.ID), t.Window)
		return []any{}
	}

	entries := make([]any, 0, len(contacts))
	for _, c := range contacts {
		entries = append(entries, c)
	}

	return entries
}

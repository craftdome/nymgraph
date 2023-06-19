package view

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Tyz3/nymgraph/internal/entity"
	"github.com/Tyz3/nymgraph/internal/model"
	"github.com/Tyz3/nymgraph/internal/service"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

type ChatTab struct {
	App        fyne.App
	Controller *service.Service
	Window     fyne.Window
	TabItem    *container.TabItem

	model *model.Pseudonym

	contacts map[string]*entity.Contact

	selectedContact *entity.Contact
	replySurbs      int

	contactSelectEntry *widget.Select
	textEntry          *widget.Entry
	replySurbsEntry    *widget.Entry
	sendButton         *widget.Button
	messageArea        *fyne.Container
}

func NewChatTab(app fyne.App, controller *service.Service, window fyne.Window, title string, icon fyne.Resource, pseudonym *model.Pseudonym) *ChatTab {
	t := &ChatTab{
		App:        app,
		Controller: controller,
		Window:     window,
		TabItem:    container.NewTabItemWithIcon(title, icon, nil),

		model: pseudonym,
	}

	{
		t.contactSelectEntry = widget.NewSelect(nil, func(s string) {
			t.selectedContact = t.contacts[s]
			if t.textEntry.Text != "" {
				t.sendButton.Enable()
			}
		})
		t.contactSelectEntry.PlaceHolder = "choose a contact"

		t.textEntry = widget.NewMultiLineEntry()
		t.textEntry.OnChanged = func(s string) {
			if s != "" {
				t.sendButton.Enable()
			} else {
				t.sendButton.Disable()
			}
		}
		t.textEntry.PlaceHolder = "start txting..."
		t.textEntry.Wrapping = fyne.TextWrapWord
		t.textEntry.SetMinRowsVisible(3)

		t.replySurbsEntry = widget.NewEntry()
		t.replySurbsEntry.Wrapping = fyne.TextWrapOff
		t.replySurbsEntry.Validator = func(s string) (err error) {
			if s == "" {
				return nil
			}
			if i, err := strconv.Atoi(s); err != nil {
				return err
			} else if i < 0 {
				return errors.New("")
			}
			return nil
		}

		t.sendButton = widget.NewButtonWithIcon("", theme.MailSendIcon(), func() {
			t.sendButton.Disable()
			defer t.sendButton.Enable()

			if obj, err := t.sendMessage(t.selectedContact, t.textEntry.Text, t.replySurbs); err != nil {
				dialog.ShowError(err, t.Window)
				return
			} else {
				t.messageArea.Add(obj)
			}
		})
		t.sendButton.Importance = widget.LowImportance

		t.messageArea = container.NewVBox()
		scroll := container.NewVScroll(t.messageArea)

		split := container.NewVSplit(
			scroll,
			container.NewPadded(
				container.NewVBox(
					container.NewHBox(
						widget.NewLabel("Address:"),
						t.contactSelectEntry,
						layout.NewSpacer(),
						widget.NewLabel("Reply SURBs:"),
						container.NewGridWrap(
							fyne.NewSize(80, 34),
							t.replySurbsEntry,
						),
					),
					container.NewBorder(
						nil,
						nil,
						nil,
						t.sendButton,
						t.textEntry,
					),
				),
			),
		)

		split.SetOffset(0.9)

		t.TabItem.Content = split
	}

	return t
}

func (t *ChatTab) Load() {
	t.Reload()
}

func (t *ChatTab) Reload() {
	t.contacts = make(map[string]*entity.Contact, 1)

	contacts, err := t.Controller.Contacts.GetAll(t.model.Pseudonym.ID)
	if err != nil {
		dialog.ShowError(errors.Wrap(err, "Controller.Contacts.GetAll"), t.Window)
		return
	}

	t.sendButton.Disable()
	t.contactSelectEntry.Selected = ""
	t.contactSelectEntry.Options = nil
	for _, c := range contacts {
		key := fmt.Sprintf("(%s) %s...", c.Alias, c.Address[:7])
		t.contacts[key] = c
		t.contactSelectEntry.Options = append(t.contactSelectEntry.Options, key)
	}
	t.contactSelectEntry.Refresh()
}

func (t *ChatTab) sendMessage(contact *entity.Contact, text string, replySurbs int) (fyne.CanvasObject, error) {
	if err := t.model.NymClient.SendMessage(text, contact.Address, replySurbs); err != nil {
		return nil, err
	}

	sent, err := t.Controller.Sent.Create(contact.ID, text)
	if err != nil {
		return nil, err
	}

	icon := widget.NewIcon(theme.MoveUpIcon())
	timestamp := widget.NewLabel(sent.CreateAt.Format(time.RFC822Z))
	alias := widget.NewLabel(contact.Alias)
	alias.Alignment = fyne.TextAlignTrailing
	alias.TextStyle = fyne.TextStyle{Bold: true}
	msg := widget.NewLabel(sent.Text)
	msg.Wrapping = fyne.TextWrapWord

	size := t.messageArea.Size()

	content := container.NewVBox(
		container.NewBorder(
			nil,
			nil,
			nil,
			container.NewHBox(
				alias,
				timestamp,
				icon,
			),
		),
		msg,
	)
	content.Resize(fyne.NewSize(size.Width*0.66, content.Size().Height))

	return container.NewBorder(
		nil,
		nil,
		nil,
		content,
	), nil
}

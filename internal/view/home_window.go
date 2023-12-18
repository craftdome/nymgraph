package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/craftdome/go-nym/response"
	"github.com/craftdome/nymgraph/internal/entity"
	"github.com/craftdome/nymgraph/internal/model"
	"github.com/craftdome/nymgraph/internal/nym_client"
	"github.com/craftdome/nymgraph/internal/service"
	"github.com/craftdome/nymgraph/pkg/utils"
)

type HomeWindow struct {
	App        fyne.App
	Controller *service.Service
	Window     fyne.Window

	menuItem  *fyne.MenuItem
	pseudonym *entity.Pseudonym
	connect   *nym_client.ClientConnect

	selfAddress     *widget.Label
	copySelfAddress *widget.Button

	contactsTab *ContactsTab
	chatTab     *ChatTab

	Tabs *container.AppTabs

	OnClose func()
}

func NewHomeWindow(controller *service.Service, app fyne.App, title string, icon fyne.Resource, pseudonym *entity.Pseudonym, connect *nym_client.ClientConnect) *HomeWindow {
	w := &HomeWindow{
		App:        app,
		Controller: controller,
		Window:     app.NewWindow(title),

		pseudonym: pseudonym,
		connect:   connect,
	}

	w.Window.SetIcon(icon)
	w.Window.Resize(fyne.NewSize(500, 680))
	w.Window.CenterOnScreen()

	{
		w.menuItem = fyne.NewMenuItem(w.pseudonym.Pretty(), nil)
		w.menuItem.Icon = theme.NewPrimaryThemedResource(theme.LoginIcon())
		w.menuItem.Disabled = true

		w.selfAddress = widget.NewLabel("")
		w.selfAddress.Wrapping = fyne.TextTruncate

		w.copySelfAddress = widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
			if err := utils.CopyToClipboard(connect.SelfAddress()); err != nil {
				dialog.ShowError(err, w.Window)
				return
			}
			utils.ShowSplash("copied")
		})
		w.copySelfAddress.Importance = widget.LowImportance

		w.chatTab = NewChatTab(app, controller, w.Window, "Chat", theme.RadioButtonIcon(), pseudonym)
		w.chatTab.OnSendMessageCallback = func(text, address string, replySurbs int) error {
			if err := w.connect.SendMessage(text, address, replySurbs); err != nil {
				return err
			}

			return nil
		}
		w.chatTab.OnReplyMessageCallback = func(text, senderTag string) error {
			if err := w.connect.SendReply(text, senderTag); err != nil {
				return err
			}

			return nil
		}
		connect.OnReceiveCallback = func(received *response.Received) {
			w.chatTab.HandleReceivedMessage(received.Message, received.SenderTag)
		}
		connect.OnErrorCallback = func(error *response.Error) {
			w.chatTab.HandleErrorMessage(error.Message)
		}

		w.contactsTab = NewContactsTab(app, controller, w.Window, "Contacts", theme.AccountIcon(), pseudonym)
		w.contactsTab.OnUpdateCallback = func(last, new *model.Contact) {
			w.chatTab.UpdateContact(last, new)
		}
		w.contactsTab.OnCreateCallback = func(contact *model.Contact) {
			w.chatTab.CreateContact(contact)
		}
		w.contactsTab.OnDeleteCallback = func(contact *model.Contact) {
			w.chatTab.DeleteContact(contact)
		}

		w.Tabs = container.NewAppTabs(
			w.chatTab.TabItem,
			w.contactsTab.TabItem,
		)

		w.Window.SetContent(container.NewBorder(
			container.NewBorder(
				nil,
				nil,
				widget.NewLabel("Self Address:"),
				w.copySelfAddress,
				w.selfAddress,
			),
			nil,
			nil,
			nil,
			w.Tabs,
		))
	}

	return w
}

func (w *HomeWindow) Load() {
	w.contactsTab.Load()
	w.chatTab.Load()
}

func (w *HomeWindow) Unload() {
	w.contactsTab.Unload()
	w.chatTab.Unload()
}

func (w *HomeWindow) MenuItem() *fyne.MenuItem {
	return w.menuItem
}

func (w *HomeWindow) ClientConnect() *nym_client.ClientConnect {
	return w.connect
}

func (w *HomeWindow) UpdateSelfAddress() {
	if w.connect.SelfAddress() == "" {
		w.selfAddress.SetText("")
	} else {
		w.selfAddress.SetText(w.connect.SelfAddress()[:21] + "...")
	}
}

func (w *HomeWindow) Close() {
	if w.connect.IsOnline() {
		if err := w.connect.Close(); err != nil {
			return
		}
	}
}

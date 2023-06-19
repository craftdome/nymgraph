package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Tyz3/nymgraph/internal/model"
	"github.com/Tyz3/nymgraph/internal/service"
	"github.com/Tyz3/nymgraph/pkg/utils"
)

type HomeWindow struct {
	App        fyne.App
	Controller *service.Service
	Window     fyne.Window

	model *model.Pseudonym

	selfAddress     *widget.Label
	copySelfAddress *widget.Button

	contactsTab *ContactsTab
	chatTab     *ChatTab

	Tabs *container.AppTabs

	OnClose func()
}

func NewHomeWindow(controller *service.Service, app fyne.App, title string, icon fyne.Resource, model *model.Pseudonym) *HomeWindow {
	w := &HomeWindow{
		App:        app,
		Controller: controller,
		Window:     app.NewWindow(title),

		model: model,
	}

	w.Window.SetIcon(icon)
	w.Window.Resize(fyne.NewSize(500, 680))
	w.Window.CenterOnScreen()

	{
		w.selfAddress = widget.NewLabel(model.NymClient.SelfAddress()[:21] + "...")
		w.selfAddress.Wrapping = fyne.TextTruncate

		w.copySelfAddress = widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
			if err := utils.CopyToClipboard(model.NymClient.SelfAddress()); err != nil {
				dialog.ShowError(err, w.Window)
				return
			}
			utils.ShowSplash("copied")
		})
		w.copySelfAddress.Importance = widget.LowImportance

		w.chatTab = NewChatTab(app, controller, w.Window, "Chat", theme.RadioButtonIcon(), model)
		w.contactsTab = NewContactsTab(app, controller, w.Window, "Contacts", theme.AccountIcon(), model)
		w.contactsTab.OnUpdate = func() {
			w.chatTab.Reload()
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

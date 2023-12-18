package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/craftdome/nymgraph/internal/entity"
	"github.com/craftdome/nymgraph/internal/model"
	"github.com/craftdome/nymgraph/internal/service"
	"github.com/craftdome/nymgraph/internal/view/custom_widget"
	"github.com/pkg/errors"
	"sort"
	"strconv"
	"time"
)

type ChatTab struct {
	App        fyne.App
	Controller *service.Service
	Window     fyne.Window
	TabItem    *container.TabItem

	pseudonym *entity.Pseudonym

	sentModel *model.Sent

	contacts map[string]*model.Contact

	contactSelectEntry *widget.Select
	textEntry          *widget.Entry
	replySurbsEntry    *widget.Entry
	sendButton         *widget.Button
	chat               *custom_widget.Chat

	OnSendMessageCallback  func(text, address string, replySurbs int) error
	OnReplyMessageCallback func(text, senderTag string) error
}

func NewChatTab(app fyne.App, controller *service.Service, window fyne.Window, title string, icon fyne.Resource, pseudonym *entity.Pseudonym) *ChatTab {
	t := &ChatTab{
		App:        app,
		Controller: controller,
		Window:     window,
		TabItem:    container.NewTabItemWithIcon(title, icon, nil),

		pseudonym: pseudonym,

		sentModel: &model.Sent{
			Sent: &entity.Sent{},
		},
	}

	{
		t.contactSelectEntry = widget.NewSelect(nil, func(s string) {
			t.sentModel.Contact = t.contacts[s].Contact

			// Enabling send condition
			if t.textEntry.Text != "" {
				t.sendButton.Enable()
			}
		})
		t.contactSelectEntry.PlaceHolder = "choose a contact"

		t.textEntry = widget.NewMultiLineEntry()
		t.textEntry.OnChanged = func(s string) {
			t.sentModel.Sent.Text = s
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
				t.sentModel.Sent.ReplySurbs = 0
				return nil
			}
			if i, err := strconv.Atoi(s); err != nil {
				return err
			} else if i < 0 {
				return errors.New("")
			} else {
				t.sentModel.Sent.ReplySurbs = i
				return nil
			}
		}

		t.sendButton = widget.NewButtonWithIcon("", theme.MailSendIcon(), func() {
			t.sendButton.Disable()
			defer t.sendButton.Enable()

			// Send message to mixnet
			if err := t.OnSendMessageCallback(t.sentModel.Sent.Text, t.sentModel.Contact.Address, t.sentModel.Sent.ReplySurbs); err != nil {
				dialog.ShowError(err, t.Window)
				return
			}

			// Save message to db
			sent, err := t.Controller.Sent.Create(t.sentModel.Contact.ID, t.sentModel.Sent.Text, t.sentModel.Sent.ReplySurbs)
			if err != nil {
				dialog.ShowError(err, t.Window)
				return
			}

			// Draw message at chat
			t.AddChatSent(sent)

			t.textEntry.SetText("")
			t.replySurbsEntry.SetText("")
		})
		t.sendButton.Importance = widget.LowImportance

		t.chat = custom_widget.NewChat()
		t.chat.SetAutoScroll(true)

		split := container.NewVSplit(
			t.chat,
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

func (t *ChatTab) Unload() {
	t.contactSelectEntry.SetSelected("")
	t.contactSelectEntry.Options = nil
	t.textEntry.SetText("")
	t.replySurbsEntry.SetText("")
	t.chat.Reset()
}

func (t *ChatTab) UpdateContact(old, new *model.Contact) {
	oldKey := old.Pretty()
	newKey := new.Pretty()
	for i, key := range t.contactSelectEntry.Options {
		if oldKey == key {
			if t.contactSelectEntry.Selected == oldKey {
				t.contactSelectEntry.Selected = newKey
			}
			t.contactSelectEntry.Options[i] = newKey
			delete(t.contacts, oldKey)
			t.contacts[newKey] = new
			t.contactSelectEntry.Refresh()
			return
		}
	}
}

func (t *ChatTab) CreateContact(new *model.Contact) {
	newKey := new.Pretty()
	t.contacts[newKey] = new
	t.contactSelectEntry.Options = append(t.contactSelectEntry.Options, newKey)
	t.contactSelectEntry.Refresh()
}

func (t *ChatTab) DeleteContact(old *model.Contact) {
	oldKey := old.Pretty()
	for i, key := range t.contactSelectEntry.Options {
		if key == oldKey {
			t.contactSelectEntry.Selected = ""
			t.contactSelectEntry.Options = append(t.contactSelectEntry.Options[:i], t.contactSelectEntry.Options[i+1:]...)
			delete(t.contacts, oldKey)
			t.contactSelectEntry.Refresh()

			// TODO update chat messages
			return
		}
	}
}

func (t *ChatTab) Reload() {
	t.contacts = make(map[string]*model.Contact, 1)

	contacts, err := t.Controller.Contacts.GetAll(t.pseudonym.ID)
	if err != nil {
		dialog.ShowError(errors.Wrap(err, "Controller.Contacts.GetAll"), t.Window)
		return
	}

	t.sendButton.Disable()
	t.contactSelectEntry.Selected = ""
	t.contactSelectEntry.Options = nil
	for _, c := range contacts {
		key := c.Pretty()
		t.contacts[key] = c
		t.contactSelectEntry.Options = append(t.contactSelectEntry.Options, key)
	}
	t.contactSelectEntry.Refresh()

	chatElements := make(map[time.Time]fyne.CanvasObject, 1)

	// Preparing message history for load
	for _, v := range contacts {
		sent, err := t.Controller.Sent.GetAll(v.Contact.ID)
		if err != nil {
			dialog.ShowError(errors.Wrap(err, "Controller.Sent.GetAll"), t.Window)
			return
		}

		for _, m := range sent {
			chatElements[m.Sent.CreateAt] = custom_widget.NewChatSent(m)
		}
	}

	received, err := t.Controller.Received.GetAll(t.pseudonym.ID)
	if err != nil {
		dialog.ShowError(errors.Wrap(err, "Controller.Received.GetAll"), t.Window)
		return
	}

	for _, m := range received {
		chatElements[m.Received.CreateAt] = custom_widget.NewChatReceived(m)
		replies, err := t.Controller.Replies.GetAll(m.Received.ID)
		if err != nil {
			dialog.ShowError(errors.Wrap(err, "Controller.Replies.GetAll"), t.Window)
			return
		}

		for _, r := range replies {
			chatElements[r.Reply.CreateAt] = custom_widget.NewChatReply(r)
		}
	}

	// Collecting timestamps for sort
	keys := make([]time.Time, 0, len(chatElements))
	for k := range chatElements {
		keys = append(keys, k)
	}

	// Sort by timestamp
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Before(keys[j])
	})

	for _, k := range keys {
		element := chatElements[k]
		if _, ok := element.(*custom_widget.ChatReceived); ok {
			chatReceived := element.(*custom_widget.ChatReceived)
			chatReceived.ReplyAction = func() {
				var text string
				replyTextEntry := widget.NewMultiLineEntry()
				replyTextEntry.Wrapping = fyne.TextWrapWord
				replyTextEntry.SetMinRowsVisible(3)
				replyTextEntry.OnChanged = func(s string) {
					text = s
				}
				dialog.ShowCustomConfirm("Reply", "Sent", "Cancel",
					replyTextEntry,
					func(b bool) {
						if !b || text == "" {
							return
						}

						if err := t.OnReplyMessageCallback(text, chatReceived.GetModel().Received.SenderTag); err != nil {
							dialog.ShowError(err, t.Window)
							return
						}

						// Save message to db
						reply, err := t.Controller.Replies.Create(chatReceived.GetModel().Received.ID, text)
						if err != nil {
							dialog.ShowError(err, t.Window)
							return
						}

						// Draw message at chat
						chatSent := custom_widget.NewChatReply(reply)
						t.chat.AddChatElement(chatSent)
					},
					t.Window,
				)
			}
		}
		t.chat.LoadChatElement(element)
	}
	t.chat.ScrollToBottom()
}

func (t *ChatTab) HandleReceivedMessage(text, senderTag string) {
	received, err := t.Controller.Received.Create(t.pseudonym.ID, text, senderTag)
	if err != nil {
		dialog.ShowError(err, t.Window)
		return
	}

	t.AddChatReceived(received)
}

func (t *ChatTab) HandleErrorMessage(text string) {
	dialog.ShowError(errors.New(text), t.Window)
}

func (t *ChatTab) AddChatReceived(received *model.Received) {
	chatReceived := custom_widget.NewChatReceived(received)
	chatReceived.ReplyAction = func() {
		var text string
		replyTextEntry := widget.NewMultiLineEntry()
		replyTextEntry.Wrapping = fyne.TextWrapWord
		replyTextEntry.SetMinRowsVisible(3)
		replyTextEntry.OnChanged = func(s string) {
			text = s
		}
		dialog.ShowCustomConfirm("Reply", "Sent", "Cancel",
			replyTextEntry,
			func(b bool) {
				if !b || text == "" {
					return
				}

				if err := t.OnReplyMessageCallback(text, received.Received.SenderTag); err != nil {
					dialog.ShowError(err, t.Window)
					return
				}

				// Save message to db
				reply, err := t.Controller.Replies.Create(received.Received.ID, text)
				if err != nil {
					dialog.ShowError(err, t.Window)
					return
				}

				// Draw message at chat
				chatSent := custom_widget.NewChatReply(reply)
				t.chat.AddChatElement(chatSent)
			},
			t.Window,
		)
	}
	t.chat.AddChatElement(chatReceived)
}

func (t *ChatTab) AddChatReply(reply *model.Reply) {
	chatReply := custom_widget.NewChatReply(reply)
	t.chat.AddChatElement(chatReply)
	t.chat.ScrollToBottom()
}

func (t *ChatTab) AddChatSent(sent *model.Sent) {
	chatSent := custom_widget.NewChatSent(sent)
	t.chat.AddChatElement(chatSent)
	t.chat.ScrollToBottom()
}

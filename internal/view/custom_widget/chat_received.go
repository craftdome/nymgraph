package custom_widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Tyz3/nymgraph/internal/model"
	"image/color"
	"time"
)

type ChatReceived struct {
	widget.BaseWidget

	model *model.Received

	contextMenu    *fyne.Menu
	timestampText  *canvas.Text
	senderTagText  *canvas.Text
	msgRichText    *widget.RichText
	senderNameText *canvas.Text

	area *fyne.Container

	ReplyAction func()
}

func NewChatReceived(model *model.Received) *ChatReceived {
	w := &ChatReceived{
		model: model,

		senderTagText: canvas.NewText("", subColor),
	}

	item := fyne.NewMenuItem("Delete", func() {

	})
	item.Icon = theme.DeleteIcon()

	w.timestampText = canvas.NewText("", subColor)
	w.timestampText.TextSize = 10

	w.senderTagText = canvas.NewText("", subColor)
	w.senderTagText.TextSize = 10

	w.senderNameText = canvas.NewText("Anonymous", color.Black)
	w.senderNameText.TextStyle = fyne.TextStyle{Bold: true}
	w.senderNameText.TextSize = 14

	w.area = container.NewVBox()

	w.update()

	w.ExtendBaseWidget(w)
	return w
}

func (w *ChatReceived) GetModel() *model.Received {
	return w.model
}

func (w *ChatReceived) update() {
	w.timestampText.Text = w.model.Received.CreateAt.Format(time.DateTime)
	w.senderTagText.Text = "senderTag: " + w.model.Received.SenderTag
	w.msgRichText = widget.NewRichTextFromMarkdown(w.model.Received.Text)
	w.msgRichText.Wrapping = fyne.TextWrapWord

	replyButton := widget.NewButtonWithIcon("", theme.MailReplyIcon(), func() {
		w.ReplyAction()
	})
	replyButton.Importance = widget.LowImportance
	if w.model.Received.SenderTag == "" {
		replyButton.Disable()
	}

	infoArea := container.NewVBox()
	if w.model.Received.SenderTag != "" {
		infoArea.Add(w.senderTagText)
	}
	infoArea.Add(container.NewBorder(
		nil,
		nil,
		nil,
		w.timestampText,
	))

	w.area.RemoveAll()
	w.area.Add(
		container.NewBorder(
			container.NewBorder(
				nil,
				nil,
				container.NewHBox(
					widget.NewIcon(theme.MoveDownIcon()),
					w.senderNameText,
					replyButton,
				),
				nil,
			),
			nil,
			nil,
			container.NewGridWrap(fyne.NewSize(100, 0)),
			container.NewVBox(
				widget.NewCard("", "", w.msgRichText),
				container.NewBorder(
					nil,
					nil,
					nil,
					infoArea,
				),
			),
		),
	)
}

func (w *ChatReceived) CreateRenderer() fyne.WidgetRenderer {
	w.ExtendBaseWidget(w)
	return &chatReceivedRenderer{objects: []fyne.CanvasObject{
		w.area,
	}}
}

type chatReceivedRenderer struct {
	objects []fyne.CanvasObject
}

func (r *chatReceivedRenderer) Destroy() {
}

func (r *chatReceivedRenderer) Layout(size fyne.Size) {
	for _, o := range r.objects {
		o.Resize(size)
	}
}

func (r *chatReceivedRenderer) MinSize() fyne.Size {
	minSize := fyne.NewSize(0, 0)
	for _, child := range r.objects {
		minSize = minSize.Max(child.MinSize())
	}
	return minSize
}

func (r *chatReceivedRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *chatReceivedRenderer) Refresh() {
	for _, o := range r.objects {
		o.Refresh()
	}
}

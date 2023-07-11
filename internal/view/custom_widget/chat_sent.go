package custom_widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Tyz3/nymgraph/internal/model"
	"image/color"
	"strconv"
	"time"
)

var (
	subColor = color.RGBA{R: 135, G: 135, B: 135, A: 255}
)

type ChatSent struct {
	widget.BaseWidget

	model *model.Sent

	contextMenu    *fyne.Menu
	timestampText  *canvas.Text
	toText         *canvas.Text
	surbsReplyText *canvas.Text
	msgRichText    *widget.RichText
	senderNameText *canvas.Text

	area *fyne.Container
}

func NewChatSent(model *model.Sent) *ChatSent {
	w := &ChatSent{
		model: model,

		surbsReplyText: canvas.NewText("", subColor),
	}

	item := fyne.NewMenuItem("Delete", func() {

	})
	item.Icon = theme.DeleteIcon()

	w.timestampText = canvas.NewText("", subColor)
	w.timestampText.TextSize = 10

	w.toText = canvas.NewText("", subColor)
	w.toText.TextSize = 10

	w.surbsReplyText = canvas.NewText("", subColor)
	w.surbsReplyText.TextSize = 10

	w.senderNameText = canvas.NewText("You", color.Black)
	w.senderNameText.TextStyle = fyne.TextStyle{Bold: true}
	w.senderNameText.TextSize = 14

	w.area = container.NewVBox()

	w.update()

	w.ExtendBaseWidget(w)
	return w
}

func (w *ChatSent) GetModel() *model.Sent {
	return w.model
}

func (w *ChatSent) update() {
	w.timestampText.Text = w.model.Sent.CreateAt.Format(time.DateTime)
	w.toText.Text = "to: " + w.model.Contact.Alias
	w.surbsReplyText.Text = "surbs: " + strconv.Itoa(w.model.Sent.ReplySurbs)
	w.msgRichText = widget.NewRichTextFromMarkdown(w.model.Sent.Text)
	w.msgRichText.Wrapping = fyne.TextWrapWord

	w.area.RemoveAll()
	w.area.Add(
		container.NewBorder(
			container.NewBorder(
				nil,
				nil,
				nil,
				container.NewHBox(
					w.senderNameText,
					widget.NewIcon(theme.MoveUpIcon()),
				),
			),
			nil,
			container.NewGridWrap(fyne.NewSize(100, 0)),
			nil,
			container.NewVBox(
				widget.NewCard("", "", w.msgRichText),
				container.NewBorder(
					nil,
					nil,
					container.NewVBox(
						container.NewHBox(
							w.toText,
							w.surbsReplyText,
						),
						container.NewBorder(
							nil,
							nil,
							w.timestampText,
							nil,
						),
					),
					nil,
				),
			),
		),
	)
}

func (w *ChatSent) CreateRenderer() fyne.WidgetRenderer {
	w.ExtendBaseWidget(w)
	return &chatSentRenderer{objects: []fyne.CanvasObject{
		w.area,
	}}
}

type chatSentRenderer struct {
	objects []fyne.CanvasObject
}

func (r *chatSentRenderer) Destroy() {
}

func (r *chatSentRenderer) Layout(size fyne.Size) {
	for _, o := range r.objects {
		o.Resize(size)
	}
}

func (r *chatSentRenderer) MinSize() fyne.Size {
	minSize := fyne.NewSize(0, 0)
	for _, child := range r.objects {
		minSize = minSize.Max(child.MinSize())
	}
	return minSize
}

func (r *chatSentRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *chatSentRenderer) Refresh() {
	for _, o := range r.objects {
		o.Refresh()
	}
}

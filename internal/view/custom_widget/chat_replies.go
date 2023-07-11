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

type ChatReply struct {
	widget.BaseWidget

	model *model.Reply

	contextMenu    *fyne.Menu
	timestampText  *canvas.Text
	toText         *canvas.Text
	msgRichText    *widget.RichText
	senderNameText *canvas.Text

	area *fyne.Container
}

func NewChatReply(model *model.Reply) *ChatReply {
	w := &ChatReply{
		model: model,
	}

	item := fyne.NewMenuItem("Delete", func() {

	})
	item.Icon = theme.DeleteIcon()

	w.timestampText = canvas.NewText("", subColor)
	w.timestampText.TextSize = 10

	w.toText = canvas.NewText("", subColor)
	w.toText.TextSize = 10

	w.senderNameText = canvas.NewText("You", color.Black)
	w.senderNameText.TextStyle = fyne.TextStyle{Bold: true}
	w.senderNameText.TextSize = 14

	w.area = container.NewVBox()

	w.update()

	w.ExtendBaseWidget(w)
	return w
}

func (w *ChatReply) GetModel() *model.Reply {
	return w.model
}

func (w *ChatReply) update() {
	w.timestampText.Text = w.model.Reply.CreateAt.Format(time.DateTime)
	w.toText.Text = "to: " + w.model.Received.Received.SenderTag
	w.msgRichText = widget.NewRichTextFromMarkdown(w.model.Reply.Text)
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
						w.toText,
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

func (w *ChatReply) CreateRenderer() fyne.WidgetRenderer {
	w.ExtendBaseWidget(w)
	return &chatReplyRenderer{objects: []fyne.CanvasObject{
		w.area,
	}}
}

type chatReplyRenderer struct {
	objects []fyne.CanvasObject
}

func (r *chatReplyRenderer) Destroy() {
}

func (r *chatReplyRenderer) Layout(size fyne.Size) {
	for _, o := range r.objects {
		o.Resize(size)
	}
}

func (r *chatReplyRenderer) MinSize() fyne.Size {
	minSize := fyne.NewSize(0, 0)
	for _, child := range r.objects {
		minSize = minSize.Max(child.MinSize())
	}
	return minSize
}

func (r *chatReplyRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *chatReplyRenderer) Refresh() {
	for _, o := range r.objects {
		o.Refresh()
	}
}

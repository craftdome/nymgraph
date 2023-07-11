package custom_widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Chat struct {
	widget.BaseWidget

	autoScroll bool

	scrollArea *container.Scroll
	baseArea   *fyne.Container
}

func NewChat() *Chat {
	w := &Chat{}

	w.baseArea = container.NewVBox()
	padded := container.NewPadded(w.baseArea)
	w.baseArea.Move(w.baseArea.Position().SubtractXY(6, 0))
	w.scrollArea = container.NewVScroll(padded)

	return w
}

func (w *Chat) SetAutoScroll(b bool) {
	w.autoScroll = b
}

func (w *Chat) AddChatElement(e fyne.CanvasObject) {
	w.baseArea.Add(e)
	w.baseArea.Refresh()
	if w.autoScroll {
		w.ScrollToBottom()
	}
}

func (w *Chat) LoadChatElement(e fyne.CanvasObject) {
	w.baseArea.Add(e)
}

func (w *Chat) ScrollToBottom() {
	w.scrollArea.ScrollToBottom()
}

func (w *Chat) Reset() {
	w.baseArea.RemoveAll()
}

func (w *Chat) CreateRenderer() fyne.WidgetRenderer {
	w.ExtendBaseWidget(w)
	return &chatRenderer{objects: []fyne.CanvasObject{
		w.scrollArea,
	}}
}

type chatRenderer struct {
	objects []fyne.CanvasObject
}

func (r *chatRenderer) Destroy() {
}

func (r *chatRenderer) Layout(size fyne.Size) {
	for _, o := range r.objects {
		o.Resize(size)
	}
}

func (r *chatRenderer) MinSize() fyne.Size {
	minSize := fyne.NewSize(0, 0)
	for _, child := range r.objects {
		minSize = minSize.Max(child.MinSize())
	}
	return minSize
}

func (r *chatRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *chatRenderer) Refresh() {
	for _, o := range r.objects {
		o.Refresh()
	}
}

package xwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type TappedLabel struct {
	widget.BaseWidget

	label   *widget.Label
	minSize fyne.Size

	OnLeftClick  func(*fyne.PointEvent)
	OnRightClick func(*fyne.PointEvent)
}

func NewTappedLabel(minSize fyne.Size) *TappedLabel {
	w := &TappedLabel{
		label:   widget.NewLabel(""),
		minSize: minSize,
	}
	w.ExtendBaseWidget(w)
	return w
}

func (w *TappedLabel) SetText(text string) {
	w.label.SetText(text)
}

func (w *TappedLabel) SetStyle(style fyne.TextStyle) {
	w.label.TextStyle = style
}

func (w *TappedLabel) SetAlignment(align fyne.TextAlign) {
	w.label.Alignment = align
}

func (w *TappedLabel) SetWrapping(wrap fyne.TextWrap) {
	w.label.Wrapping = wrap
}

func (w *TappedLabel) Tapped(event *fyne.PointEvent) {
	if w.OnLeftClick != nil {
		w.OnLeftClick(event)
	}
}

func (w *TappedLabel) TappedSecondary(event *fyne.PointEvent) {
	if w.OnRightClick != nil {
		w.OnRightClick(event)
	}
}

func (w *TappedLabel) CreateRenderer() fyne.WidgetRenderer {
	w.ExtendBaseWidget(w)
	return &TappedLabelRenderer{objects: []fyne.CanvasObject{w.label}, minSize: w.minSize}
}

type TappedLabelRenderer struct {
	objects []fyne.CanvasObject
	minSize fyne.Size
}

func (r *TappedLabelRenderer) Destroy() {
}

func (r *TappedLabelRenderer) Layout(size fyne.Size) {
	for _, o := range r.objects {
		o.Resize(size)
	}
}

func (r *TappedLabelRenderer) MinSize() fyne.Size {
	return r.minSize
}

func (r *TappedLabelRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *TappedLabelRenderer) Refresh() {
	for _, o := range r.objects {
		o.Refresh()
	}
}

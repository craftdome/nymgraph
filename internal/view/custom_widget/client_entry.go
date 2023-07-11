package custom_widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Tyz3/fyne-extra/xwidget"
	"github.com/Tyz3/nymgraph/internal/entity"
)

type ClientEntry struct {
	widget.BaseWidget

	model *entity.Pseudonym

	buttonMenu *fyne.Menu

	nameLabel    *xwidget.TappedLabel
	serverLabel  *xwidget.TappedLabel
	editButton   *widget.Button
	deleteButton *widget.Button

	OnLeftClick    func(*fyne.PointEvent)
	OnRightClick   func(*fyne.PointEvent)
	OnEditTapped   func(*entity.Pseudonym)
	OnDeleteTapped func(*entity.Pseudonym)
}

func NewClientEntry() *ClientEntry {
	w := &ClientEntry{
		buttonMenu:   fyne.NewMenu(""),
		nameLabel:    xwidget.NewTappedLabel(fyne.NewSize(120, 34)),
		serverLabel:  xwidget.NewTappedLabel(fyne.NewSize(120, 34)),
		editButton:   widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), nil),
		deleteButton: widget.NewButtonWithIcon("", theme.ContentClearIcon(), nil),
	}

	w.nameLabel.SetStyle(fyne.TextStyle{Bold: true})

	w.nameLabel.OnLeftClick = func(event *fyne.PointEvent) {
		if w.OnLeftClick != nil {
			w.OnLeftClick(event)
		}
	}
	w.serverLabel.OnLeftClick = func(event *fyne.PointEvent) {
		if w.OnLeftClick != nil {
			w.OnLeftClick(event)
		}
	}

	w.ExtendBaseWidget(w)
	return w
}

func (w *ClientEntry) GetModel() *entity.Pseudonym {
	return w.model
}

func (w *ClientEntry) SetModel(model *entity.Pseudonym) {
	w.model = model

	w.nameLabel.SetText(w.model.Name)
	w.serverLabel.SetText(w.model.Server)
	w.editButton.OnTapped = func() {
		if w.OnEditTapped != nil {
			w.OnEditTapped(w.model)
		}
	}
	w.deleteButton.OnTapped = func() {
		if w.OnDeleteTapped != nil {
			w.OnDeleteTapped(w.model)
		}
	}
}

func (w *ClientEntry) Tapped(event *fyne.PointEvent) {
	if w.OnLeftClick != nil {
		w.OnLeftClick(event)
	}
}

func (w *ClientEntry) TappedSecondary(event *fyne.PointEvent) {
	if w.OnRightClick != nil {
		w.OnRightClick(event)
	}
}

func (w *ClientEntry) CreateRenderer() fyne.WidgetRenderer {
	w.ExtendBaseWidget(w)
	return &ClientEntryRenderer{objects: []fyne.CanvasObject{
		container.NewHBox(
			w.nameLabel,
			layout.NewSpacer(),
			w.serverLabel,
			layout.NewSpacer(),
			w.editButton,
			w.deleteButton,
		),
	}}
}

type ClientEntryRenderer struct {
	objects []fyne.CanvasObject
}

func (r *ClientEntryRenderer) Destroy() {
}

func (r *ClientEntryRenderer) Layout(size fyne.Size) {
	for _, o := range r.objects {
		o.Resize(size)
	}
}

func (r *ClientEntryRenderer) MinSize() fyne.Size {
	minSize := fyne.NewSize(0, 0)
	for _, child := range r.objects {
		minSize = minSize.Max(child.MinSize())
	}
	return minSize
}

func (r *ClientEntryRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *ClientEntryRenderer) Refresh() {
	for _, o := range r.objects {
		o.Refresh()
	}
}

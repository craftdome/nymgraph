package custom_widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/craftdome/nymgraph/internal/model"
	"github.com/craftdome/nymgraph/pkg/fyne-extra/xwidget"
)

type ContactEntry struct {
	widget.BaseWidget

	model *model.Contact

	contextMenu *fyne.Menu

	icon         *widget.Icon
	aliasLabel   *xwidget.TappedLabel
	addressLabel *xwidget.TappedLabel
	copyButton   *widget.Button

	OnLeftClick       func(*fyne.PointEvent)
	OnRightClick      func(*fyne.PointEvent)
	OnCopyButtonClick func()
}

func NewContactEntry() *ContactEntry {
	w := &ContactEntry{
		contextMenu:  fyne.NewMenu(""),
		icon:         widget.NewIcon(theme.AccountIcon()),
		aliasLabel:   xwidget.NewTappedLabel(fyne.NewSize(100, 34)),
		addressLabel: xwidget.NewTappedLabel(fyne.NewSize(120, 34)),
	}

	w.aliasLabel.SetStyle(fyne.TextStyle{Bold: true})

	w.aliasLabel.OnLeftClick = func(event *fyne.PointEvent) {
		if w.OnLeftClick != nil {
			w.OnLeftClick(event)
		}
	}
	w.addressLabel.OnLeftClick = func(event *fyne.PointEvent) {
		if w.OnLeftClick != nil {
			w.OnLeftClick(event)
		}
	}

	w.aliasLabel.OnRightClick = func(event *fyne.PointEvent) {
		if w.OnRightClick != nil {
			w.OnRightClick(event)
		}
	}
	w.addressLabel.OnRightClick = func(event *fyne.PointEvent) {
		if w.OnRightClick != nil {
			w.OnRightClick(event)
		}
	}

	w.copyButton = widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		if w.OnCopyButtonClick != nil {
			w.OnCopyButtonClick()
		}
	})
	w.copyButton.Importance = widget.LowImportance

	w.OnRightClick = func(event *fyne.PointEvent) {
		widget.ShowPopUpMenuAtPosition(w.contextMenu, fyne.CurrentApp().Driver().CanvasForObject(w), event.AbsolutePosition)
	}

	w.ExtendBaseWidget(w)
	return w
}

func (w *ContactEntry) GetModel() *model.Contact {
	return w.model
}

func (w *ContactEntry) SetModel(model *model.Contact) {
	w.model = model

	w.aliasLabel.SetText(w.model.Contact.Alias)
	w.addressLabel.SetText(w.model.Contact.Address[:21] + "...")
}

func (w *ContactEntry) AddContextMenuItem(label string, icon fyne.Resource, action func()) {
	item := fyne.NewMenuItem(label, action)
	item.Icon = icon
	w.contextMenu.Items = append(w.contextMenu.Items, item)
}

func (w *ContactEntry) Tapped(event *fyne.PointEvent) {
	if w.OnLeftClick != nil {
		w.OnLeftClick(event)
	}
}

func (w *ContactEntry) TappedSecondary(event *fyne.PointEvent) {
	if w.OnRightClick != nil {
		w.OnRightClick(event)
	}
}

func (w *ContactEntry) CreateRenderer() fyne.WidgetRenderer {
	w.ExtendBaseWidget(w)
	return &ContactEntryRenderer{objects: []fyne.CanvasObject{
		container.NewHBox(
			w.icon,
			w.aliasLabel,
			layout.NewSpacer(),
			w.addressLabel,
			layout.NewSpacer(),
			w.copyButton,
		),
	}}
}

type ContactEntryRenderer struct {
	objects []fyne.CanvasObject
}

func (r *ContactEntryRenderer) Destroy() {
}

func (r *ContactEntryRenderer) Layout(size fyne.Size) {
	for _, o := range r.objects {
		o.Resize(size)
	}
}

func (r *ContactEntryRenderer) MinSize() fyne.Size {
	minSize := fyne.NewSize(0, 0)
	for _, child := range r.objects {
		minSize = minSize.Max(child.MinSize())
	}
	return minSize
}

func (r *ContactEntryRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *ContactEntryRenderer) Refresh() {
	for _, o := range r.objects {
		o.Refresh()
	}
}

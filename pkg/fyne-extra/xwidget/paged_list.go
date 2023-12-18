package xwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/craftdome/nymgraph/pkg/fyne-extra/utils"
)

type filteredEntry struct {
	OriginalID widget.ListItemID
	Object     any
}

type PagedList struct {
	widget.List

	reversed   bool
	filterMode bool

	pageSize int
	page     int

	originalEntries []any
	filteredEntries []*filteredEntry

	dataSourceFunc func() []any
	filterItemFunc func(any) bool
	updateItemFunc func(widget.ListItemID, fyne.CanvasObject)

	OnFilterItemFunc func(any, bool)
}

func NewPagedList(pageSize int, createItemFunc func() fyne.CanvasObject) *PagedList {
	w := &PagedList{
		pageSize: pageSize,
		page:     1,

		dataSourceFunc: func() []any { return []any{} },
		filterItemFunc: func(any) bool { return true },
	}
	w.Length = w.length
	w.CreateItem = createItemFunc
	w.UpdateItem = w.updateItem

	w.ExtendBaseWidget(w)
	return w
}

// SetReversed
// Устанавливает флаг переворота отображаемых элементов списка
func (w *PagedList) SetReversed(reversed bool) {
	w.reversed = reversed
	w.update()
}

// SetPage
// Устанавливает текущую страницу и перерисовывает её
func (w *PagedList) SetPage(page int) {
	w.page = page
	w.Refresh()
}

// SetFilterMode
// Устанавливает флаг фильтрации загруженных данных
func (w *PagedList) SetFilterMode(mode bool) {
	w.filterMode = mode
	w.update()
}

// SetFilterFunc
// Устанавливает функцию фильтрации загруженных данных
func (w *PagedList) SetFilterFunc(f func(any) bool) {
	w.filterItemFunc = f
}

// SetDataSourceFunc
// Устанавливает функцию загрузки данных
func (w *PagedList) SetDataSourceFunc(dataSourceFunc func() []any) {
	w.dataSourceFunc = dataSourceFunc
}

// SetUpdateItemFunc
// Устанавливает функцию инициализации элемента списка
func (w *PagedList) SetUpdateItemFunc(f func(widget.ListItemID, fyne.CanvasObject)) {
	w.updateItemFunc = f
}

// GetPages
// Возвращает количество страниц в списке
func (w *PagedList) GetPages() int {
	pages := len(w.filteredEntries) / w.pageSize
	if len(w.filteredEntries)%w.pageSize != 0 {
		pages++
	}
	if pages == 0 {
		pages++
	}
	return pages
}

func (w *PagedList) GetItems() []any {
	return w.originalEntries
}

func (w *PagedList) GetFilteredItems() []any {
	res := make([]any, 0, len(w.filteredEntries))
	for _, e := range w.filteredEntries {
		res = append(res, e.Object)
	}
	return res
}

// Reload
// 1. Очищает все данные из памяти
// 2. Запрашивает новые с помощью dataSourceFunc()
// 3. Выполняет фильтрацию если FilterMode = true
func (w *PagedList) Reload() {
	w.originalEntries = nil

	entries := w.dataSourceFunc()
	if entries == nil {
		return
	}

	w.originalEntries = append(w.originalEntries, entries...)
	w.update()
}

// ApplyActionForFiltered
// Применяет функцию к отображаемым данным и перерисовывает их
// Если action вернёт false, то ApplyActionForFiltered также вернёт false
func (w *PagedList) ApplyActionForFiltered(action func(any) bool) bool {
	for i := 0; i < len(w.filteredEntries); i++ {
		if !action(w.filteredEntries[i].Object) {
			return false
		}
	}

	w.update()
	return true
}

// DeleteItem
// Метод удаления элемента в списке
func (w *PagedList) DeleteItem(item any) {
	var id int
	for i, e := range w.filteredEntries {
		if e.Object == item {
			id = i
			break
		}
	}

	// Удалить элемент в сыром списке
	origID := w.filteredEntries[id].OriginalID
	w.originalEntries[origID] = nil
	w.originalEntries = append(w.originalEntries[:origID], w.originalEntries[origID+1:]...)

	// Удалить элемент в фильтрованном списке
	w.filteredEntries[id] = nil
	w.filteredEntries = append(w.filteredEntries[:id], w.filteredEntries[id+1:]...)

	for i := id; i < len(w.filteredEntries); i++ {
		if w.reversed {
			w.filteredEntries[i].OriginalID++
		} else {
			w.filteredEntries[i].OriginalID--
		}
	}
	w.update()
}

func (w *PagedList) length() int {
	if len(w.filteredEntries) < w.pageSize {
		return len(w.filteredEntries)
	}

	remained := len(w.filteredEntries) - w.pageSize*(w.page-1)
	if remained < w.pageSize {
		return remained
	} else {
		return w.pageSize
	}
}

func (w *PagedList) updateItem(item widget.ListItemID, object fyne.CanvasObject) {
	// Переопределяем ID для смещения по страницам списка
	id := item + (w.page-1)*w.pageSize
	if w.updateItemFunc != nil {
		w.updateItemFunc(id, object)
	}
}

func (w *PagedList) update() {
	// Очистка отображаемого списка
	w.filteredEntries = nil

	for i := 0; i < len(w.originalEntries); i++ {
		// Режим фильтрации
		if w.filterMode && !w.filterItemFunc(w.originalEntries[i]) {
			if w.OnFilterItemFunc != nil {
				w.OnFilterItemFunc(w.originalEntries[i], false)
			}
			continue
		}

		if w.OnFilterItemFunc != nil {
			w.OnFilterItemFunc(w.originalEntries[i], true)
		}
		// Сохраняем ID в отображаемом списке
		w.filteredEntries = append(w.filteredEntries, &filteredEntry{
			OriginalID: i,
			Object:     w.originalEntries[i],
		})
	}

	// Переворот списка
	if w.reversed {
		utils.Reverse(w.filteredEntries)
	}

	// Обновляем виджет
	w.Refresh()
}

package custom_widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type DisplayedItem[T comparable] struct {
	Linked T
	itemID int
}

type TList[T comparable] struct {
	widget.List

	reversed   bool
	filterMode bool

	pageSize int
	page     int

	items          []T
	displayedItems []*DisplayedItem[T]

	dataSourceFunc func() []T
	filterItemFunc func(T) bool
	updateItemFunc func(int, T, fyne.CanvasObject)
}

func (w *TList[T]) Init(pageSize int, createItemFunc func() fyne.CanvasObject) *TList[T] {
	w.pageSize = pageSize
	w.page = 1
	w.dataSourceFunc = func() []T { return []T{} }
	w.filterItemFunc = func(T) bool { return true }

	w.Length = w.length
	w.CreateItem = createItemFunc
	w.UpdateItem = w.updateItem

	w.ExtendBaseWidget(w)
	return w
}

// SetReversed
// Устанавливает флаг переворота отображаемых элементов списка
func (w *TList[T]) SetReversed(reversed bool) {
	w.reversed = reversed
	w.update()
}

// SetPage
// Устанавливает текущую страницу и перерисовывает её
func (w *TList[T]) SetPage(page int) {
	w.page = page
	w.Refresh()
}

// SetFilterMode
// Устанавливает флаг фильтрации загруженных данных
func (w *TList[T]) SetFilterMode(mode bool) {
	w.filterMode = mode
	w.update()
}

// SetFilterFunc
// Устанавливает функцию фильтрации загруженных данных
func (w *TList[T]) SetFilterFunc(f func(T) bool) {
	w.filterItemFunc = f
}

// SetDataSourceFunc
// Устанавливает функцию загрузки данных
func (w *TList[T]) SetDataSourceFunc(dataSourceFunc func() []T) {
	w.dataSourceFunc = dataSourceFunc
}

// SetUpdateItemFunc
// Устанавливает функцию инициализации элемента списка
func (w *TList[T]) SetUpdateItemFunc(f func(int, T, fyne.CanvasObject)) {
	w.updateItemFunc = f
}

// GetPages
// Возвращает количество страниц в списке
func (w *TList[T]) GetPages() int {
	pages := len(w.displayedItems) / w.pageSize
	if len(w.displayedItems)%w.pageSize != 0 {
		pages++
	}
	if pages == 0 {
		pages++
	}
	return pages
}

func (w *TList[T]) GetItems() []T {
	return w.items
}

func (w *TList[T]) GetDisplayedItems() []*DisplayedItem[T] {
	return w.displayedItems
}

// Reload
// 1. Очищает все данные из памяти
// 2. Запрашивает новые с помощью dataSourceFunc()
// 3. Выполняет фильтрацию если FilterMode = true
func (w *TList[T]) Reload() {
	w.items = nil

	entries := w.dataSourceFunc()
	if entries == nil {
		return
	}

	w.items = append(w.items, entries...)
	w.update()
}

// DeleteItem
// Метод удаления элемента в списке
func (w *TList[T]) DeleteItem(item T) {
	var id int
	for i, e := range w.displayedItems {
		if e.Linked == item {
			id = i
			break
		}
	}

	// Удалить элемент в сыром списке
	origID := w.displayedItems[id].itemID
	//w.items[origID] = nil
	w.items = append(w.items[:origID], w.items[origID+1:]...)

	// Удалить элемент в фильтрованном списке
	w.displayedItems[id] = nil
	w.displayedItems = append(w.displayedItems[:id], w.displayedItems[id+1:]...)

	for i := id; i < len(w.displayedItems); i++ {
		if w.reversed {
			w.displayedItems[i].itemID++
		} else {
			w.displayedItems[i].itemID--
		}
	}
	w.update()
}

func (w *TList[T]) length() int {
	if len(w.displayedItems) < w.pageSize {
		return len(w.displayedItems)
	}

	remained := len(w.displayedItems) - w.pageSize*(w.page-1)
	if remained < w.pageSize {
		return remained
	} else {
		return w.pageSize
	}
}

func (w *TList[T]) updateItem(item widget.ListItemID, object fyne.CanvasObject) {
	// Переопределяем ID для смещения по страницам списка
	// item - это тот id элемента списка, который мы видим, то есть
	// первый элемент 1 стр будет 0 и на 2 стр тоже будет 0
	id := item + (w.page-1)*w.pageSize
	//fmt.Printf("%d + (%d-1)*%d = %d\n", item, w.page, w.pageSize, id)
	if w.updateItemFunc != nil {
		w.updateItemFunc(id, w.displayedItems[id].Linked, object)
	}
}

func (w *TList[T]) update() {
	// Очистка отображаемого списка
	w.displayedItems = nil

	for i := 0; i < len(w.items); i++ {
		// Режим фильтрации
		if w.filterMode && !w.filterItemFunc(w.items[i]) {
			continue
		}

		// Сохраняем ID в отображаемом списке
		w.displayedItems = append(w.displayedItems, &DisplayedItem[T]{
			itemID: i,
			Linked: w.items[i],
		})
	}

	// Переворот списка
	if w.reversed {
		Reverse(w.displayedItems)
	}

	// Обновляем виджет
	w.Refresh()
}

func Reverse[S ~[]E, E any](s S) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

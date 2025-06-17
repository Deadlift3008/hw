package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	firstElem *ListItem
	lastElem  *ListItem
	length    int
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.firstElem
}

func (l *list) Back() *ListItem {
	return l.lastElem
}

func (l *list) PushFront(v interface{}) *ListItem {
	l.length++
	if l.firstElem == nil {
		l.firstElem = &ListItem{
			Value: v,
			Next:  nil,
			Prev:  nil,
		}

		l.lastElem = l.firstElem

		return l.firstElem
	}

	newFirstElem := &ListItem{
		Value: v,
		Next:  l.firstElem,
		Prev:  nil,
	}

	l.firstElem.Prev = newFirstElem
	l.firstElem = newFirstElem

	return l.firstElem
}

func (l *list) PushBack(v interface{}) *ListItem {
	l.length++

	if l.lastElem == nil {
		l.firstElem = &ListItem{
			Value: v,
			Next:  nil,
			Prev:  nil,
		}

		l.lastElem = l.firstElem

		return l.firstElem
	}

	newLastElem := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  l.lastElem,
	}

	l.lastElem.Next = newLastElem
	l.lastElem = newLastElem

	return l.lastElem
}

func (l *list) Remove(i *ListItem) {
	prev := i.Prev
	next := i.Next

	if prev != nil {
		prev.Next = next
	} else {
		l.firstElem = next
	}

	if next != nil {
		next.Prev = prev
	} else {
		l.lastElem = prev
	}

	l.length--
}

func (l *list) MoveToFront(i *ListItem) {
	prev := i.Prev
	next := i.Next

	if prev != nil {
		prev.Next = next
	} else {
		l.firstElem = next
	}

	if next != nil {
		next.Prev = prev
	} else {
		l.lastElem = prev
	}

	l.firstElem.Prev = i
	i.Prev = nil
	i.Next = l.firstElem
	l.firstElem = i
}

func NewList() List {
	return new(list)
}

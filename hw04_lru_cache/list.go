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
	elemDict  map[*ListItem]*ListItem
	length    int
}

func (l *list) Len() int {
	return l.length
}

func NewList() List {
	return new(list)
}

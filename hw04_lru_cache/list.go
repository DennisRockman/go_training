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

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	newItem := &ListItem{Value: v}
	if l.head == nil {
		l.head = newItem
		l.tail = newItem
	} else {
		l.insertBefore(newItem, l.head)
	}
	l.len++
	return newItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newItem := &ListItem{Value: v}
	if l.tail == nil {
		newItem = l.PushFront(v)
	} else {
		l.insertAfter(newItem, l.tail)
		l.len++
	}
	return newItem
}

func (l *list) Remove(i *ListItem) {
	l.lazyRemove(i)
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if i.Prev != nil {
		l.lazyRemove(i)
		l.insertBefore(i, l.head)
	}
}

func (l *list) insertAfter(newItem, i *ListItem) {
	newItem.Prev = i
	if i.Next == nil {
		l.tail = newItem
	} else {
		newItem.Next = i.Next
		i.Next.Prev = newItem
	}
	i.Next = newItem
}

func (l *list) insertBefore(newItem, i *ListItem) {
	newItem.Next = i
	if i.Prev == nil {
		l.head = newItem
	} else {
		newItem.Prev = i.Prev
		i.Prev.Next = newItem
	}
	i.Prev = newItem
}

func (l *list) lazyRemove(i *ListItem) {
	if i.Prev == nil {
		l.head = i.Next
	} else {
		i.Prev.Next = i.Next
	}

	if i.Next == nil {
		l.tail = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}

	i.Prev = nil
	i.Next = nil
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	head, tail *ListItem
	len        int
}

func NewList() List {
	return new(list)
}

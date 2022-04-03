package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	item, isPresent := l.items[key]
	cItem := &cacheItem{key: key, value: value}
	if isPresent {
		item.Value = cItem
		l.queue.MoveToFront(item)
	} else {
		l.items[key] = l.queue.PushFront(cItem)
		cItem, ok := l.queue.Back().Value.(*cacheItem)
		if l.queue.Len() > l.capacity && ok {
			delete(l.items, cItem.key)
			l.queue.Remove(l.queue.Back())
		}
	}
	return isPresent
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	item, isPresent := l.items[key]
	if isPresent {
		cItem, ok := item.Value.(*cacheItem)
		if ok {
			l.queue.MoveToFront(item)
			return cItem.value, true
		}
	}
	return nil, false
}

func (l *lruCache) Clear() {
	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

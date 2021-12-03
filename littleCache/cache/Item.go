package cache

import (
	"sync"
	"time"
)

type Item struct {
	sync.RWMutex

	key interface{}

	value interface{}

	duration time.Duration

	storTime time.Time
}

func NewItem(key interface{}, value interface{}, duration time.Duration) *Item {
	t := time.Now()
	return &Item{
		key:      key,
		value:    value,
		duration: duration,
		storTime: t,
	}
}

func (item *Item) Value() interface{} {
	return item.value
}

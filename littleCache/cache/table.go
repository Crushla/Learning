package cache

import (
	"errors"
	"sync"
	"time"
)

type Table struct {
	sync.RWMutex

	name string // table name

	items map[interface{}]*Item

	nextCheckTime time.Duration

	cleaner *time.Timer
}

func (t *Table) Add(key interface{}, value interface{}, time time.Duration) *Item {
	item := NewItem(key, value, time)
	t.Lock()
	t.AddItem(item)
	return item
}

func (t *Table) AddItem(item *Item) {
	t.items[item.key] = item
	checkTime := t.nextCheckTime
	t.Unlock()

	if item.duration > 0 && (checkTime == 0 || item.duration < checkTime) {
		t.expirationCheck()
	}
}

func (t *Table) DeleteItem(key interface{}) (*Item, error) {
	item, ok := t.items[key]
	if !ok {
		return nil, errors.New("error key not found")
	}
	t.Unlock()
	item.RLock()
	defer item.RUnlock()
	t.Lock()
	delete(t.items, key)
	return item, nil
}

//过期处理
func (t *Table) expirationCheck() {
	t.Lock()
	//关闭
	if t.cleaner != nil {
		t.cleaner.Stop()
	}

	//当前时间
	now := time.Now()
	waitTime := 0 * time.Second
	for key, item := range t.items {
		item.RLock()
		//item可以存放的时间
		duration := item.duration
		//存放的时间
		storTime := item.storTime
		item.RUnlock()
		if now.Sub(storTime) > duration {
			t.DeleteItem(key)
		} else {
			//找出最小等待时间，等时间到了继续做过期检验
			if waitTime == 0 || duration-now.Sub(storTime) < waitTime {
				waitTime = duration - now.Sub(storTime)
			}
		}
	}
	t.nextCheckTime = waitTime
	if waitTime > 0 {
		t.cleaner = time.AfterFunc(waitTime, func() {
			go t.expirationCheck()
		})
	}
	t.Unlock()
}

func (t *Table) Value(key interface{}) (*Item, error) {
	t.RLock()
	r, ok := t.items[key]
	t.RUnlock()
	if ok {
		return r, nil
	}
	return nil, errors.New("error key not found")

}

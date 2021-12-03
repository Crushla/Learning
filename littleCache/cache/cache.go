package cache

import "sync"

var (
	cache = make(map[string]*Table)
	mutex sync.RWMutex
)

func Cache(name string) *Table {
	mutex.RLock()
	table, ok := cache[name]
	mutex.RUnlock()
	if !ok {
		//如果没有这个名字的cache表
		mutex.Lock()
		table, ok = cache[name]
		if !ok {
			table = &Table{
				name:  name,
				items: make(map[interface{}]*Item),
			}
			cache[name] = table
		}
		mutex.Unlock()
	}

	return table
}

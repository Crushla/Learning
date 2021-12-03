package main

import (
	"Learning/littleCache/cache"
	"fmt"
	"time"
)

func main() {
	table := cache.Cache("myCache")
	table.Add("key", "value", 5*time.Second)
	time.Sleep(6 * time.Second)
	res, err := table.Value("key")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(res.Value().(string))
}

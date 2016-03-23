package cache

import (
	"errors"
	"sync"
)

var (
	TableMap map[string]Table
	lock     sync.Mutex
)

func GetTable(name string) Table {
	lock.Lock()
	defer lock.Unlock()
	value, ok := TableMap[name]
	if ok {
		return value
	}
	value = CreateTable(name)
	return value
}

func DeleteTable(name string) error {
	lock.Lock()
	defer lock.Unlock()
	_, ok := TableMap[name]
	if ok {
		delete(TableMap, name)
	} else {
		return errors.New("Not Founded")
	}
	return nil
}

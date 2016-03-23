package cache

import (
	"sync"
	"time"
)

type DeleteCallBackFunc func(interface{}, interface{})
type Item struct {
	sync.Mutex
	Key        string
	Value      []byte
	Createtime time.Time
	Lasttime   time.Time
	Cas        uint64
	Flag       uint16
	expiretime time.Duration
	deleteFunc DeleteCallBackFunc
}

func CreateItem(key string, value []byte, cas uint64, flag uint16, expire time.Duration, callback DeleteCallBackFunc) *Item {
	t := time.Now()
	return &Item{
		Key:        key,
		Value:      value,
		expiretime: expire,
		Cas:        cas,
		Flag:       flag,
		Createtime: t,
		Lasttime:   t,
		deleteFunc: callback,
	}
}

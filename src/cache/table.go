package cache

import (
	"errors"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type Table struct {
	sync.Mutex
	name       string
	createtime time.Time
	items      map[string]*Item
	interval   time.Duration
}

func CreateTable(name string) Table {
	db := Table{
		name:       name,
		createtime: time.Now(),
		items:      make(map[string]*Item),
	}
	return db
}

func (t *Table) Insert(
	key string,
	value []byte,
	flag uint16,
	cas *uint64,
	expire time.Duration,
	callback DeleteCallBackFunc) error {

	atomic.AddUint64(cas, 1)
	t.Lock()
	// value, err := d.items[key]
	item := CreateItem(key, value, *cas, flag, expire, callback)
	t.items[key] = item
	t.Unlock()
	if expire > 0 && (t.interval == 0 || expire < t.interval) {
		t.expireFunc()
	}
	return nil
}
func (t *Table) Add(
	key string,
	value []byte,
	flag uint16,
	cas *uint64,
	expire time.Duration,
	callback DeleteCallBackFunc) error {

	t.Lock()
	_, ok := t.items[key]
	if ok {
		return errors.New("exists")
	}
	atomic.AddUint64(cas, 1)
	item := CreateItem(key, value, *cas, flag, expire, callback)
	t.items[key] = item
	t.Unlock()
	if expire > 0 && (t.interval == 0 || expire < t.interval) {
		t.expireFunc()
	}
	return nil
}

func (t *Table) Cas(key string,
	value []byte,
	flag uint16,
	cas *uint64,
	excas uint64,
	expire time.Duration,
	callback DeleteCallBackFunc) error {

	t.Lock()
	it, ok := t.items[key]
	if !ok {
		return errors.New("unexists")
	}
	if it.Cas != excas {
		return errors.New("exists")
	}
	atomic.AddUint64(cas, 1)
	item := CreateItem(key, value, *cas, flag, expire, callback)
	t.items[key] = item
	t.Unlock()
	if expire > 0 && (t.interval == 0 || expire < t.interval) {
		t.expireFunc()
	}
	return nil
}

func (t *Table) Replace(
	key string,
	value []byte,
	flag uint16,
	cas *uint64,
	expire time.Duration,
	callback DeleteCallBackFunc) error {

	t.Lock()
	_, ok := t.items[key]
	if !ok {
		return errors.New("unexists")
	}
	atomic.AddUint64(cas, 1)
	item := CreateItem(key, value, *cas, flag, expire, callback)
	t.items[key] = item
	t.Unlock()
	if expire > 0 && (t.interval == 0 || expire < t.interval) {
		t.expireFunc()
	}
	return nil
}

func (t *Table) Delete(key string) error {
	t.Lock()
	defer t.Unlock()
	item, ok := t.items[key]
	if !ok {
		return errors.New("Not Founded")
	}
	if item.deleteFunc != nil {
		item.deleteFunc(item.Key, item.Value)
	}
	delete(t.items, key)
	return nil
}

func (t *Table) Get(key string) (*Item, error) {
	t.Lock()
	defer t.Unlock()
	item, ok := t.items[key]
	if !ok {
		return nil, errors.New("Not Founded")
	}
	return item, nil
}
func (t *Table) Incr(key string, incr uint32) (uint32, error) {
	t.Lock()
	defer t.Unlock()
	item, ok := t.items[key]
	if !ok {
		return 0, errors.New("Not Founded")
	}
	var val uint32
	if len(item.Value) <= 4 {
		v, err := strconv.ParseUint(string(item.Value), 10, 32)
		if err != nil {
			return 0, errors.New("Not Uint32")
		} else {
			val = uint32(v) + incr
			item.Value = []byte(strconv.FormatUint(uint64(val), 10))
		}

	} else {
		return 0, errors.New("Not Uint32")
	}
	t.items[key] = item
	return val, nil
}
func (t *Table) Decr(key string, decr uint32) (uint32, error) {
	t.Lock()
	defer t.Unlock()
	item, ok := t.items[key]
	if !ok {
		return 0, errors.New("Not Founded")
	}
	var val uint32
	if len(item.Value) <= 4 {
		v, err := strconv.ParseUint(string(item.Value), 10, 32)
		if err != nil {
			return 0, errors.New("Not Uint32")
		} else {
			val = uint32(v) - decr
			item.Value = []byte(strconv.FormatUint(uint64(val), 10))
		}

	} else {
		return 0, errors.New("Not Uint32")
	}
	t.items[key] = item
	return val, nil
}

func (t *Table) expireFunc() {
	t.Lock()
	items := t.items
	t.Unlock()
	now := time.Now()
	smallesttime := t.interval

	for key, item := range items {
		item.Lock()
		expire := item.expiretime
		lasttime := item.Lasttime
		item.Unlock()
		if expire > 0 && expire <= now.Sub(lasttime) {
			t.Delete(key)
		}
		if expire > 0 && expire < smallesttime {
			smallesttime = expire
		}
	}
	t.Lock()
	if smallesttime < t.interval {
		t.interval = smallesttime
	}
	time.AfterFunc(t.interval, func() {
		go t.expireFunc()
	})
	t.Unlock()
}

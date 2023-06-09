package krkstops

import (
	"sync"
	"time"
)

const (
	airlyExpire = time.Minute * 10
	depsExpire  = time.Second * 10
)

type entry[T any] struct {
	v     T
	setAt time.Time
}

type cache[T any] struct {
	d   map[uint]entry[T]
	m   sync.RWMutex
	ttl time.Duration
}

func (c *cache[T]) get(k uint) (T, *time.Time, bool) {
	c.m.RLock()
	v, ok := c.d[k]
	c.m.RUnlock()

	if ok && v.setAt.Add(c.ttl).After(time.Now()) {
		return v.v, &v.setAt, true
	} else {
		var empty T
		return empty, nil, false
	}
}

func (c *cache[T]) set(k uint, v T) {
	c.m.Lock()
	c.d[k] = entry[T]{
		v:     v,
		setAt: time.Now(),
	}
	c.m.Unlock()
}

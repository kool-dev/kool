package shell

import "sync"

type lookupCache struct {
	mtx   *sync.RWMutex
	cache map[string]error
}

func newLookupCache() *lookupCache {
	var mtx sync.RWMutex
	return &lookupCache{&mtx, nil}
}

func (l *lookupCache) fetch(key string) (exists bool, err error) {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	if l.cache != nil {
		err, exists = l.cache[key]
	}

	return
}

func (l *lookupCache) set(key string, err error) {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	if l.cache == nil {
		l.cache = make(map[string]error)
	}

	l.cache[key] = err
}

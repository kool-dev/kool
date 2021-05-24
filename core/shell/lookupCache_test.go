package shell

import (
	"errors"
	"testing"
)

func TestLookupCache(t *testing.T) {
	cache := newLookupCache()

	if cache.mtx == nil {
		t.Error("missing mtx on lookupCache")
	}

	if cache.cache != nil {
		t.Errorf("unexpected default non-nil value for cache: %v", cache.cache)
	}

	if exists, err := cache.fetch("unknown"); err != nil || exists {
		t.Errorf("bad return for fetching unknown key from cache")
	}

	var err = errors.New("error")
	cache.set("key", err)

	if exists, cached := cache.fetch("key"); !exists || !errors.Is(cached, err) {
		t.Errorf("failed to fetch error from cache - exists: %v got: %v expected: %v", exists, cached, err)
	}
}

package forms

import (
	"sync"

	"github.com/ipoluianov/altsysinfo/system"
)

// The system information takes up to half a second to collect (system_profiler,
// WMI...), so it is loaded in the background once and kept: switching between
// the categories does not wait for it. Refresh loads it again.

// infoKey is the cache key of system.GetInfo, shown by the common and the RAM pages
const infoKey = "info"

// loadResult is the information of a page
type loadResult struct {
	info   system.Info
	tables []*system.DetailTable
	err    error
}

// cacheEntry is the information of one key, loaded or being loaded
type cacheEntry struct {
	done chan struct{} // closed when res is set
	res  loadResult
}

func (e *cacheEntry) ready() bool {
	select {
	case <-e.done:
		return true
	default:
		return false
	}
}

type dataCache struct {
	mtx     sync.Mutex
	entries map[string]*cacheEntry
}

var cache = &dataCache{entries: make(map[string]*cacheEntry)}

// cacheKey returns the key of the information the page shows; ok is false
// for a page that loads nothing slow
func cacheKey(mode string) (key string, ok bool) {
	switch mode {
	case "common", "ram":
		return infoKey, true
	case "pcidev":
		return "", false
	}
	return mode, true
}

// get returns the entry of the key, starting to load it if it is not there
func (c *dataCache) get(key string) *cacheEntry {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if e, ok := c.entries[key]; ok {
		return e
	}
	e := &cacheEntry{done: make(chan struct{})}
	c.entries[key] = e
	go func() {
		e.res = load(key)
		close(e.done)
	}()
	return e
}

// drop forgets the information of the key, so the next get loads it again
func (c *dataCache) drop(key string) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	delete(c.entries, key)
}

// prefetch loads all the pages one by one in the background, so they open at once
func (c *dataCache) prefetch() {
	keys := []string{infoKey}
	for _, category := range system.DetailCategories() {
		keys = append(keys, category.ID)
	}
	go func() {
		for _, key := range keys {
			<-c.get(key).done
		}
	}()
}

func load(key string) loadResult {
	if key == infoKey {
		info, err := system.GetInfo()
		return loadResult{info: info, err: err}
	}
	for _, category := range system.DetailCategories() {
		if category.ID == key {
			tables, err := category.Load()
			return loadResult{tables: tables, err: err}
		}
	}
	return loadResult{}
}

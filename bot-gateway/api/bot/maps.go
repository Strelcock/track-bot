package bot

import "sync"

type waitMap struct {
	mu           sync.RWMutex
	waitForInput map[int64]bool
}

func newWaitMap() *waitMap {
	wait := make(map[int64]bool)
	return &waitMap{waitForInput: wait}
}

type infoMap struct {
	mu   sync.RWMutex
	info map[int64][]string
}

func NewInfoMap() *infoMap {
	info := make(map[int64][]string)
	return &infoMap{info: info}
}

func (i *infoMap) Add(key int64, val []string) {
	i.mu.Lock()
	i.info[key] = val
	i.mu.Unlock()
}

func (i *infoMap) Get(key int64) []string {
	return i.info[key]
}

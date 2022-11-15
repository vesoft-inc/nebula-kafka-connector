package manager

import "sync"

type WaitGroupMap struct {
	mu sync.Mutex
	m  map[string]*sync.WaitGroup
}

func NewWaitGroups() *WaitGroupMap {
	return &WaitGroupMap{
		m: make(map[string]*sync.WaitGroup),
	}
}

func (w *WaitGroupMap) Add(key string, delta int) {
	w.mu.Lock()
	wg := w.m[key]
	if wg == nil {
		wg = &sync.WaitGroup{}
		w.m[key] = wg
	}
	w.mu.Unlock()
	wg.Add(delta)
}

func (w *WaitGroupMap) Done(key string) {
	w.mu.Lock()
	wg := w.m[key]
	w.mu.Unlock()
	if wg != nil {
		wg.Done()
	}
}
func (w *WaitGroupMap) Wait(key string) {
	w.mu.Lock()
	wg := w.m[key]
	w.mu.Unlock()
	if wg != nil {
		wg.Wait()
	}
}

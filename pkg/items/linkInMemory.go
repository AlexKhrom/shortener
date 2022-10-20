package items

import (
	"fmt"
	"sync"
)

type LinkRepoMemory struct {
	newLinks  map[string]string
	origLinks map[string]string
	mutNew    sync.RWMutex
	mutOrig   sync.RWMutex
}

func (r *LinkRepoMemory) NewLink(origLink string) string {
	r.mutNew.RLock()
	newLink, ok := r.origLinks[origLink]
	r.mutNew.RUnlock()
	if ok {
		return newLink
	}
	key := origLink
	for {
		newLink = CreatNewLinkString(key)
		r.mutNew.RLock()
		_, ok := r.newLinks[newLink]
		r.mutNew.RUnlock()
		if !ok {
			r.mutNew.Lock()
			r.newLinks[newLink] = origLink
			r.mutNew.Unlock()
			r.mutNew.Lock()
			r.origLinks[origLink] = newLink
			r.mutNew.Unlock()
			return newLink
		}
		key = newLink
	}
}

func (r *LinkRepoMemory) GetOrigLink(unifiedLink string) string {
	r.mutNew.RLock()
	origLink, ok := r.newLinks[unifiedLink]
	r.mutNew.RUnlock()
	if !ok {
		fmt.Println("some err with get orig string")
		return ""
	}
	return origLink
}

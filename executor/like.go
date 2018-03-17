package executor

import (
	"regexp"
	"sync"
)

var (
	regCache = map[likeStr]*regexp.Regexp{}
	regLock  = sync.Mutex{}
)

type likeStr string

// TODO : hell, this is more complicated than this, later :)
func (ls likeStr) regexp() *regexp.Regexp {
	regLock.Lock()
	defer regLock.Unlock()
	if r, ok := regCache[ls]; ok {
		return r
	}

	re := "^"
	for i := 0; i < len(ls); i++ {
		if ls[i] == '%' {
			re += ".*"
		} else if ls[i] == '_' {
			re += "."
		} else {
			re += string(ls[i])
		}
	}
	re += "$"
	r := regexp.MustCompile(re)
	regCache[ls] = r
	return r
}

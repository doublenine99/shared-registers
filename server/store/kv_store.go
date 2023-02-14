package store

import (
	"sync"
)

var noop = ""
var s = sync.Map{} // use concurrentMap for simplicity first

func Get(key string) (string, bool) {
	v, ok := s.Load(key)
	if !ok {
		return noop, false
	}
	// should be ok if not map is private so that no one can put other type of value
	return v.(string), true
}

func Set(key, value string) {
	s.Store(key, value)
}

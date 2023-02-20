package store

import (
	"shared-registers/common/proto"
	"sync"
)

var s = sync.Map{} // use concurrentMap for simplicity first

func Get(key string) (*proto.StoredValue, error) {
	v, ok := s.Load(key)
	if !ok {
		return nil, nil
	}
	// should be ok if not map is private so that no one can put other type of value
	return v.(*proto.StoredValue), nil
}

func Set(key string, value *proto.StoredValue) {
	s.Store(key, value)
	//log.Printf("Stored %s %v\n", key, value)
}

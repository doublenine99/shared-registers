package store

import (
	"shared-registers/common/proto"
	"sync"
)

type valueType struct {
	value     string
	timestamp *proto.TimeStamp
}

var noop = ""
var s = sync.Map{} // use concurrentMap for simplicity first

func Get(key string) (string, *proto.TimeStamp, bool) {
	v, ok := s.Load(key)
	var empty *proto.TimeStamp
	if !ok {
		return noop, empty, false
	}
	// should be ok if not map is private so that no one can put other type of value
	var out = v.(valueType)
	return out.value, out.timestamp, true
}

func Set(key, value string, ts *proto.TimeStamp) {
	var val valueType
	val.value = value
	val.timestamp = ts

	// TODO: only insert if my ts > old ts
	s.Store(key, val)
}

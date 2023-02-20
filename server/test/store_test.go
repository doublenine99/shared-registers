package test

import (
	"math/rand"
	"shared-registers/common/proto"
	"shared-registers/server/store"
	"strconv"
	"testing"
)

func TestConcurrentSet(t *testing.T) {
	n := 100000
	collideChance := 10000
	// concurrently set n times with the collideChance to test the performance in concurrent cases
	for i := 0; i < n; i++ {
		go func(i int) {
			idx := rand.Intn(n / collideChance)
			s := strconv.Itoa(idx)
			store.Set(s, &proto.StoredValue{
				Val: s,
				Ts: &proto.TimeStamp{
					ClientID:      "cid",
					RequestNumber: uint64(i),
				},
			})
		}(i)
	}
}

func TestSet(t *testing.T) {
	n := 100000
	collideChance := 10000
	values := make([]string, n)
	for i := 0; i < n; i++ {
		idx := rand.Intn(n / collideChance)
		s := strconv.Itoa(idx)
		values[idx] = s

		store.Set(s, &proto.StoredValue{
			Val: s,
			Ts: &proto.TimeStamp{
				ClientID:      "cid",
				RequestNumber: uint64(i),
			},
		})
	}
	for i, v := range values {
		if v != "" {
			val, err := store.Get(strconv.Itoa(i))
			if err != nil || v != val.GetVal() {
				t.Errorf("value not match for key %d, %s %s", i, v, val)
			}
		}
	}
}

func TestServiceHandler(t *testing.T) {
	// TODO: test that series of sets and gets gives correct output
}

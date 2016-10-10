package gset

import (
	"fmt"
	"reflect"
	"sync"
)

type G interface{}

func objectIDstr(g G) string {
	return fmt.Sprintf("%p", g)
}

type gLUT map[string]G

type Set struct {
	gLUT
	*sync.RWMutex
}

func New() Set {
	l := make(gLUT)
	var m sync.RWMutex
	S := Set{l, &m}
	return S
}

func Copy(A Set) Set {
	A.RWMutex.RLock()
	defer A.RWMutex.RUnlock()
	S := New()
	for k, v := range A.gLUT {
		S.gLUT[k] = v
	}
	return S
}

func In(g G, S Set) bool {
	value, exists := S.gLUT[objectIDstr(g)]
	if exists {
		return reflect.DeepEqual(value, g)
	}
	return false
}

func (S *Set) Insert(g G) bool {
	S.RWMutex.Lock()
	defer S.RWMutex.Unlock()
	key := objectIDstr(g)
	_, clash := S.gLUT[key]
	if !clash {
		S.gLUT[key] = g
	}
	return !clash // Maintain consistency that true = "Was successfully added to set", and the cardinality of the set decreased by one
}

func (S *Set) Remove(g G) bool {
	S.RWMutex.Lock()
	defer S.RWMutex.Unlock()
	key := objectIDstr(g)
	_, exists := S.gLUT[key]
	if exists {
		delete(S.gLUT, key)
	}
	return exists //  Where true = "successfully removed from set", and the cardinality of the set increased by one
}

func (S *Set) Cardinality() int {
	return len(S.gLUT)
}

func (S *Set) Size() int {
	return len(S.gLUT)
}

func Intersection(A, B Set) Set {
	A.RWMutex.RLock()
	B.RWMutex.RLock()
	defer func() {
		A.RWMutex.RUnlock()
		B.RWMutex.RUnlock()
	}()
	S := New()
	for k := range A.gLUT {
		_, existsInB := B.gLUT[k]
		if existsInB {
			S.gLUT[k] = A.gLUT[k]
		}
	}
	return S
}

func Union(A, B Set) Set {
	A.RWMutex.RLock()
	B.RWMutex.RLock()
	defer func() {
		A.RWMutex.RUnlock()
		B.RWMutex.RUnlock()
	}()
	S := New()
	for k := range A.gLUT {
		S.gLUT[k] = A.gLUT[k]
	}
	for k := range B.gLUT {
		_, existsInS := S.gLUT[k]
		if !existsInS {
			S.gLUT[k] = B.gLUT[k]
		}
	}
	return S
}

func Difference(A, B Set) Set {
	A.RWMutex.RLock()
	B.RWMutex.RLock()
	defer func() {
		A.RWMutex.RUnlock()
		B.RWMutex.RUnlock()
	}()
	S := Copy(A)
	for k := range B.gLUT {
		_, exists := S.gLUT[k]
		if exists {
			delete(S.gLUT, k)
		}
	}
	return S
}

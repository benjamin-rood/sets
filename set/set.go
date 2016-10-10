package set

import "sync"

type G interface{}

type Equality interface {
	Equals(b G) bool
}

type Item interface {
	Hashable
	Equality
	Less(than Item) bool
}

type Hashable interface {
	Key() key
}

type key string

type lut map[key]Hashable

type Set struct {
	lut
	*sync.RWMutex
}

func Intersection(A, B Set) Set {
	A.RWMutex.RLock()
	B.RWMutex.RLock()
	defer func() {
		A.RWMutex.RUnlock()
		B.RWMutex.RUnlock()
	}()
	S := New()
	for k := range A.lut {
		_, existsInB := B.lut[k]
		if existsInB {
			S.lut[k] = A.lut[k]
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
	for k := range A.lut {
		S.lut[k] = A.lut[k]
	}
	for k := range B.lut {
		_, existsInS := S.lut[k]
		if !existsInS {
			S.lut[k] = B.lut[k]
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
	for k := range B.lut {
		_, exists := S.lut[k]
		if exists {
			delete(S.lut, k)
		}
	}
	return S
}

func New() Set {
	l := make(lut)
	var m sync.RWMutex
	S := Set{l, &m}
	return S
}

func Copy(A Set) Set {
	A.RWMutex.RLock()
	defer A.RWMutex.RUnlock()
	S := New()
	for k, v := range A.lut {
		S.lut[k] = v
	}
	return S
}

func (S *Set) Insert(h Hashable) bool {
	S.RWMutex.Lock()
	defer S.RWMutex.Unlock()
	_, clash := S.lut[h.Key()]
	if !clash {
		S.lut[h.Key()] = h
	}
	return clash
}

func (S *Set) Remove(h Hashable) bool {
	S.RWMutex.Lock()
	defer S.RWMutex.Unlock()
	_, exists := S.lut[h.Key()]
	if exists {
		delete(S.lut, h.Key())
	}
	return exists
}

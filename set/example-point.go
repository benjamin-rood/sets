package set

import "fmt"

type point struct {
	x, y    int
	hashkey key
}

func (p *point) Key() key {
	return p.hashkey
}

func (p *point) Define(S *Set) (key, bool) {
	hk := key(fmt.Sprintf("%p", p))
	p.hashkey = hk
	_, clash := S.lut[hk]
	if clash {
		return hk, false
	}
	return hk, true
}

func (p *point) IsDefined() bool {
	return p.hashkey == ""
}

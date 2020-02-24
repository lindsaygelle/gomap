package lex

import "sync"

// var _ Stringer = (&stringer{})(nil)

// Stringer is the interface that manages key value pairs
//
// Stringer accepts any interface as a key but expects a string as its value.
// Stringer is safe for concurrent use by multiple goroutines without additional locking or coordination.
type Stringer interface {
	Add(interface{}, string) Stringer
	AddLength(interface{}, string) int
	AddOK(interface{}, string) bool
	Del(interface{}) Stringer
	DelAll() Stringer
	DelLength(interface{}) int
	DelSome(...interface{}) Stringer
	DelSomeLength(...interface{}) int
	DelOK(interface{}) bool
	Each(func(interface{}, string)) Stringer
	EachBreak(func(interface{}, string) bool) Stringer
	EachKey(func(interface{})) Stringer
	EachValue(func(interface{})) Stringer
	Fetch(interface{}) interface{}
	FetchSome(...interface{}) []interface{}
	FetchSomeLength(...interface{}) ([]interface{}, int)
	Get(interface{}) (interface{}, bool)
	GetLength(interface{}) (interface{}, int, bool)
	Has(interface{}) bool
	Keys() []interface{}
	Len() int
	Lock() Stringer
	Map(func(interface{}, string) interface{}) Stringer
	MapBreak(func(interface{}, string) (interface{}, bool)) Stringer
	MapOK(func(interface{}, string) (interface{}, bool)) Stringer
	Mutate(func()) Stringer
	Not(interface{}) bool
	NotSome(...interface{}) bool
	NotSomeLength(...interface{})
	Unlock() Stringer
	Values() []interface{}
}

type stringer struct {
	mu sync.Mutex
	l  *Lex
}

func (stringer *stringer) Add(k interface{}, v string) Stringer {
	return stringer.Mutate(func() { stringer.l.Add(k, v) })
}
func (stringer *stringer) AddLength(k interface{}, v string) int {
	var l int
	stringer.Mutate(func() { l = stringer.l.AddLength(k, v) })
	return l
}
func (stringer *stringer) AddOK(k interface{}, v string) bool {
	var ok bool
	stringer.Mutate(func() { ok = stringer.l.AddOK(k, v) })
	return ok
}
func (stringer *stringer) Del(k interface{}) Stringer {
	return stringer.Mutate(func() { stringer.l.Del(k) })
}
func (stringer *stringer) DelAll() Stringer {
	return stringer.Mutate(func() { stringer.l.DelAll() })
}

func (stringer *stringer) DelLength(k interface{}) int {
	stringer.mu.Lock()
	var l = stringer.l.DelLength(k)
	stringer.mu.Unlock()
	return l
}

func (stringer *stringer) DelSome(k ...interface{}) Stringer {
	stringer.mu.Lock()
	stringer.l.DelSome(k...)
	stringer.mu.Unlock()
	return stringer
}
func (stringer *stringer) DelSomeLength(k ...interface{}) int {
	return stringer.l.DelSomeLength()
}
func (stringer *stringer) DelOK(k interface{}) bool {
	return stringer.l.DelOK(k)
}
func (stringer *stringer) Each(fn func(interface{}, string)) Stringer {
	stringer.l.Each(func(k, v interface{}) {
		fn(k, v.(string))
	})
	return stringer
}
func (stringer *stringer) EachBreak(fn func(k, v interface{}) bool) Stringer {
	stringer.l.EachBreak(func(k, v interface{}) bool {
		return fn(k, v.(string))
	})
	return stringer
}
func (stringer *stringer) EachKey(fn func(k interface{})) Stringer {
	stringer.l.EachKey(fn)
	return stringer
}
func (stringer *stringer) EachValue(fn func(v string)) Stringer {
	stringer.l.EachValue(func(v interface{}) {
		fn(v.(string))
	})
	return stringer
}
func (stringer *stringer) Lock() Stringer {
	stringer.mu.Lock()
	return stringer
}

func (stringer *stringer) Mutate(fn func()) Stringer {
	stringer.mu.Lock()
	fn()
	stringer.mu.Unlock()
	return stringer
}

func (stringer *stringer) Unlock() Stringer {
	stringer.mu.Unlock()
	return stringer
}

// 28 july 2014

package ui

import (
	"fmt"
	"reflect"
	"sync"
)

// Table is a Control that displays a list of like-structured data in a grid where each row represents an item and each column represents a bit of data.
// As such, a Table renders a []struct{...} where each field of the struct is rendered using package fmt's %v rule.
// Tables maintain their own storage behind a sync.RWMutex-compatible sync.Locker; use Table.Lock()/Table.Unlock() to make changes and Table.RLock()/Table.RUnlock() to merely read values.
// TODO headers
type Table interface {
	Control

	// Lock and Unlock lock and unlock Data for reading or writing.
	// RLock and RUnlock lock and unlock Data for reading only.
	// These methods have identical semantics to the analogous methods of sync.RWMutex.
	// In addition, Unlock() will request an update of the Table to account for whatever was changed.
	Lock()
	Unlock()
	RLock()
	RUnlock()

	// Data returns the internal data.
	// The returned value will contain an object of type pointer to slice of some structure; use a type assertion to get the properly typed object out.
	// Do not call this outside a Lock()..Unlock() or RLock()..RUnlock() pair.
	Data() interface{}
}

type tablebase struct {
	lock		sync.RWMutex
	data		interface{}
}

// NewTable creates a new Table.
// Currently, the argument must be a reflect.Type representing the structure that each item in the Table will hold, and the Table will be initially empty.
// This will change in the future.
func NewTable(ty reflect.Type) Table {
	if ty.Kind() != reflect.Struct {
		panic(fmt.Errorf("unknown or unsupported type %v given to NewTable()", ty))
	}
	b := new(tablebase)
	// we want a pointer to a slice
	b.data = reflect.New(reflect.SliceOf(ty)).Interface()
	return finishNewTable(b, ty)
}

func (b *tablebase) Lock() {
	b.lock.Lock()
}

// Unlock() is defined on each backend implementation of Table
// they should all call this, however
func (b *tablebase) unlock() {
	b.lock.Unlock()
}

func (b *tablebase) RLock() {
	b.lock.RLock()
}

func (b *tablebase) RUnlock() {
	b.lock.RUnlock()
}

func (b *tablebase) Data() interface{} {
	return b.data
}

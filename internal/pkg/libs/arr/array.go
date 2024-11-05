package arr

import "sync"

type ArrayType interface {
	Add(value interface{})
	Check(value interface{}) bool
	Len() int
	Values() map[interface{}]bool
}

type CheckArrayType struct {
	values map[interface{}]bool
	sm     sync.Mutex
	len    int
}

func NewCheckArrayType(len int) *CheckArrayType {
	return &CheckArrayType{values: make(map[interface{}]bool, len)}
}

func (ct *CheckArrayType) Add(value interface{}) {
	defer ct.sm.Unlock()
	ct.sm.Lock()
	ct.values[value] = true
	ct.len++
}

func (ct *CheckArrayType) AddMutil(values ...interface{}) {
	for _, v := range values {
		v := v
		ct.Add(v)
	}
}

func (ct *CheckArrayType) Check(value interface{}) bool {
	defer ct.sm.Unlock()
	ct.sm.Lock()
	if b, ok := ct.values[value]; ok && b {
		return true
	}
	return false
}

func (ct *CheckArrayType) Len() int {
	defer ct.sm.Unlock()
	ct.sm.Lock()
	return ct.len
}

func (ct *CheckArrayType) Values() map[interface{}]bool {
	defer ct.sm.Unlock()
	ct.sm.Lock()
	return ct.values
}

package utils

import (
	"container/list"
	"sync"
)

type ConcurrentList struct {
	lock *sync.RWMutex
	list *list.List
}

func NewConcurrentList() *ConcurrentList {
	return &ConcurrentList{
		lock: &sync.RWMutex{},
		list: list.New(),
	}
}

func (cl *ConcurrentList) PushBack(v interface{}) *list.Element {
	cl.lock.Lock()

	ele := cl.list.PushBack(v)

	cl.lock.Unlock()

	return ele
}

func (cl *ConcurrentList) Pop() interface{} {
	cl.lock.Lock()

	value := cl.list.Remove(cl.list.Back())

	cl.lock.Unlock()

	return value
}

func (cl *ConcurrentList) PopAll() []interface{} {
	datas := make([]interface{}, 0, cl.Len())

	if cap(datas) > 0 {
		cl.lock.Lock()

		for i := 0; i < cap(datas); i++ {
			datas = append(datas, cl.list.Remove(cl.list.Front()))
		}
		cl.lock.Unlock()

		return datas
	} else {
		return nil
	}

}

func (cl *ConcurrentList) Back() *list.Element {
	cl.lock.RLock()

	ele := cl.list.Back()

	cl.lock.RUnlock()

	return ele
}

func (cl *ConcurrentList) Front() *list.Element {
	cl.lock.RLock()

	ele := cl.list.Front()

	cl.lock.RUnlock()

	return ele
}

func (cl *ConcurrentList) Len() int {
	cl.lock.RLock()

	length := cl.list.Len()

	cl.lock.RUnlock()

	return length
}

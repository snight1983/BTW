package bitcoinvipsvr

import "sync"

type iElement interface{}

type syncQueue struct {
	lock    *sync.RWMutex
	element []iElement
}

func newQueue() *syncQueue {
	return &syncQueue{
		lock: new(sync.RWMutex),
	}
}

func (q *syncQueue) Push(e iElement) {
	q.lock.RLock()
	defer q.lock.RUnlock()
	q.element = append(q.element, e)
}

func (q *syncQueue) Pop() iElement {
	if q.IsEmpty() {
		return nil
	}
	q.lock.Lock()
	defer q.lock.Unlock()
	firstElement := q.element[0]
	q.element = q.element[1:]
	return firstElement
}

func (q *syncQueue) Clear() bool {
	if q.IsEmpty() {
		return false
	}
	q.lock.Lock()
	defer q.lock.Unlock()
	for i := 0; i < q.size(); i++ {
		q.element[i] = nil
	}
	q.element = nil
	return true
}

func (q *syncQueue) IsEmpty() bool {
	q.lock.RLock()
	defer q.lock.RUnlock()
	if len(q.element) == 0 {
		return true
	}
	return false
}

func (q *syncQueue) size() int {
	q.lock.RLock()
	defer q.lock.RUnlock()
	return len(q.element)
}

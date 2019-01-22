package bitcoinvipsvr

import "sync"

type sSyncMap struct {
	lock     *sync.RWMutex
	bm       map[interface{}]interface{}
	isUpdate bool
	isInsert bool
}

func newSyncMap() *sSyncMap {
	return &sSyncMap{
		lock: new(sync.RWMutex),
		bm:   make(map[interface{}]interface{}),
	}
}

//Get from maps return the k's value
func (m *sSyncMap) Get(k interface{}) interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if val, ok := m.bm[k]; ok {
		return val
	}
	return nil
}

func (m *sSyncMap) Set(k interface{}, v interface{}) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	if val, ok := m.bm[k]; !ok {
		m.isInsert = true
		m.bm[k] = v
	} else if val != v {
		m.bm[k] = v
		m.isUpdate = true
	} else {
		return false
	}
	return true
}

// Returns true if k is exist in the map.
func (m *sSyncMap) Check(k interface{}) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if _, ok := m.bm[k]; !ok {
		return false
	}
	return true
}

func (m *sSyncMap) Delete(k interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.bm, k)
}

func (m *sSyncMap) EachItem(eachFun func(interface{}, interface{})) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	for key, value := range m.bm {
		eachFun(key, value)
	}
}

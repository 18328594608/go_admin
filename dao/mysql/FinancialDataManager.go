package mysql

import "sync"

var instance *FinancialDataManager
var once sync.Once

type FinancialDataManager struct {
	data sync.Map //[group FinancialData]
	once sync.Once
}

func GetMemoryInstance() *FinancialDataManager {
	once.Do(func() {
		instance = &FinancialDataManager{}
	})
	return instance
}

func (m *FinancialDataManager) SetData(key string, value FinancialData) {
	m.data.Store(key, value)
}

func (m *FinancialDataManager) GetData(key string) (FinancialData, bool) {
	value, ok := m.data.Load(key)
	if !ok {
		return FinancialData{}, false
	}
	return value.(FinancialData), true
}

func (m *FinancialDataManager) DeleteData(key string) {
	m.data.Delete(key)
}

// UpdateData 更新sync.Map中的FinancialData对象
func (m *FinancialDataManager) UpdateData(key string, updateFunc func(*FinancialData)) {
	value, ok := m.data.Load(key)
	if !ok {
		return // Key不存在，直接返回
	}
	data := value.(FinancialData)
	updateFunc(&data)       // 调用更新函数修改FinancialData对象
	m.data.Store(key, data) // 重新存储更新后的对象
}

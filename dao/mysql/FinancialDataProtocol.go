package mysql

import "sync"

type SymbolData struct {
	Fees    Fee
	Symbols Symbol
}

func NewSymbolData(fee Fee, symbol Symbol) SymbolData {
	return SymbolData{
		Fees:    fee,
		Symbols: symbol,
	}
}

type FinancialData struct {
	Groups     Group
	symbolData sync.Map //[symbol SymbolData]
}

func NewFinancialData(group Group) *FinancialData {
	return &FinancialData{
		Groups: group,
	}
}

func (fd *FinancialData) SetGroup(group Group) {
	fd.Groups = group
}

func (fd *FinancialData) GetGroup() Group {
	return fd.Groups
}

func (fd *FinancialData) DeleteGroup() {
	fd.Groups = Group{}
}

// SetSymbolData 增加或更新symbolData
func (fd *FinancialData) SetSymbolData(symbol string, data SymbolData) {
	fd.symbolData.Store(symbol, data)
}

// GetSymbolData 查询symbolData
func (fd *FinancialData) GetSymbolData(symbol string) (SymbolData, bool) {
	value, ok := fd.symbolData.Load(symbol)
	if !ok {
		return SymbolData{}, false
	}
	return value.(SymbolData), true
}

func (fd *FinancialData) DeleteSymbolData(symbol string) {
	fd.symbolData.Delete(symbol)
}

func (fd *FinancialData) UpdateSymbolData(symbol string, updateFunc func(*SymbolData)) {
	data, ok := fd.GetSymbolData(symbol)
	if !ok {
		return
	}
	updateFunc(&data)
	fd.SetSymbolData(symbol, data) // 重新存储修改后的数据
}

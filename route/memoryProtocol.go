package route

import (
	"go_web_app/dao/mysql"
	"sync"
)

type SymbolStore struct {
	symbols sync.Map
	once    sync.Once
}

var symbolStoreInstance *SymbolStore
var symbolStoreOnce sync.Once

func GetSymbolStore() *SymbolStore {
	symbolStoreOnce.Do(func() {
		symbolStoreInstance = &SymbolStore{}
	})
	return symbolStoreInstance
}

func (store *SymbolStore) AddSymbol(symbol string, sym mysql.Symbol) {
	store.symbols.Store(symbol, sym)
}

func (store *SymbolStore) GetSymbol(symbol string) (mysql.Symbol, bool) {
	val, ok := store.symbols.Load(symbol)
	if !ok {
		return mysql.Symbol{}, false
	}
	return val.(mysql.Symbol), true
}

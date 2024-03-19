package mysql

import (
	"fmt"
	"go.uber.org/zap"
	"strings"
)

type Symbol struct {
	ID           int     `db:"id" json:"-"`
	Symbol       string  `db:"symbol" json:"symbol"`
	Security     string  `db:"security" json:"security"`
	Digit        int     `db:"digit" json:"digit"`
	Currency     string  `db:"currency" json:"currency"`
	ContractSize int     `db:"contract_size" json:"contract_size"`
	Percentage   int     `db:"percentage" json:"percentage"`
	MarginCalc   int     `db:"margin_calc" json:"margin_calc"`
	ProfitCalc   int     `db:"profit_calc" json:"profit_calc"`
	SwapCalc     int     `db:"swap_calc" json:"swap_calc"`
	TickSize     float64 `db:"tick_size" json:"tick_size"`
	TickPrice    int     `db:"tick_price" json:"tick_price"`
	Monday       string  `db:"monday" json:"monday"`
	Tuesday      string  `db:"tuesday" json:"tuesday"`
	Wednesday    string  `db:"wednesday" json:"wednesday"`
	Thursday     string  `db:"thursday" json:"thursday"`
	Friday       string  `db:"friday" json:"friday"`
}

func GetSymbols() ([]Symbol, error) {
	var symbols []Symbol
	query := `SELECT * FROM symbol`
	err := db.Select(&symbols, query)
	if err != nil {
		zap.L().Error("query symbols failed", zap.Error(err))
		return nil, err
	}
	return symbols, nil
}

func GetSymbolsFromSymbol(symbol string) ([]Symbol, error) {
	var symbols []Symbol
	query := `SELECT * FROM symbol WHERE  symbol = ?`
	err := db.Select(&symbols, query, symbol)
	if err != nil {
		zap.L().Error("query symbols from symbol failed", zap.Error(err))
		return nil, err
	}
	return symbols, nil
}

func GetDistinctSecurity() ([]string, error) {
	var securities []string
	query := `SELECT DISTINCT security FROM symbol`
	err := db.Select(&securities, query)
	if err != nil {
		zap.L().Error("query distinct securities failed", zap.Error(err))
		return nil, err
	}
	return securities, nil
}

func InsertSymbol(symbol Symbol) error {
	var count int
	query := `SELECT COUNT(*) FROM symbol WHERE symbol = ?`
	err := db.Get(&count, query, symbol.Symbol)
	if err != nil {
		zap.L().Error("check symbol existence failed", zap.Error(err))
		return err
	}
	if count > 0 {
		return fmt.Errorf("symbol '%s' already exists", symbol.Symbol)
	}

	query = `
        INSERT INTO symbol (
            symbol, security, digit, currency, contract_size, percentage,
            margin_calc, profit_calc, swap_calc, tick_size, tick_price,
            monday, tuesday, wednesday, thursday, friday
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `
	_, err = db.Exec(
		query,
		symbol.Symbol, symbol.Security, symbol.Digit, symbol.Currency,
		symbol.ContractSize, symbol.Percentage, symbol.MarginCalc,
		symbol.ProfitCalc, symbol.SwapCalc, symbol.TickSize, symbol.TickPrice,
		symbol.Monday, symbol.Tuesday, symbol.Wednesday, symbol.Thursday, symbol.Friday,
	)
	if err != nil {
		zap.L().Error("insert symbol failed", zap.Error(err))
		return err
	}
	return nil
}

//根据品种删除Symbol
func DeleteSymbolBySymbol(symbol string) error {
	deleteQuery := "DELETE FROM symbol WHERE symbol = ?"
	_, err := db.Exec(deleteQuery, symbol)
	if err != nil {
		zap.L().Error("[delete symbol]  failed", zap.Error(err))
		return err
	}
	return nil
}

// UpdateSymbol 更新 Symbol 记录
func UpdateSymbol(s *Symbol) error {
	query := "UPDATE symbol SET "
	var params []interface{}

	// 构建更新SQL语句，考虑到字段类型对应
	if s.Security != "" {
		query += "security=?, "
		params = append(params, s.Security)
	}
	if s.Digit != 0 {
		query += "digit=?, "
		params = append(params, s.Digit)
	}
	if s.Currency != "" {
		query += "currency=?, "
		params = append(params, s.Currency)
	}
	if s.ContractSize != 0 {
		query += "contract_size=?, "
		params = append(params, s.ContractSize)
	}
	if s.Percentage != 0 {
		query += "percentage=?, "
		params = append(params, s.Percentage)
	}
	if s.MarginCalc != 0 {
		query += "margin_calc=?, "
		params = append(params, s.MarginCalc)
	}
	if s.ProfitCalc != 0 {
		query += "profit_calc=?, "
		params = append(params, s.ProfitCalc)
	}
	if s.SwapCalc != 0 {
		query += "swap_calc=?, "
		params = append(params, s.SwapCalc)
	}
	if s.TickSize != 0.0 {
		query += "tick_size=?, "
		params = append(params, s.TickSize)
	}
	if s.TickPrice != 0 {
		query += "tick_price=?, "
		params = append(params, s.TickPrice)
	}
	if s.Monday != "" {
		query += "monday=?, "
		params = append(params, s.Monday)
	}
	if s.Tuesday != "" {
		query += "tuesday=?, "
		params = append(params, s.Tuesday)
	}
	if s.Wednesday != "" {
		query += "wednesday=?, "
		params = append(params, s.Wednesday)
	}
	if s.Thursday != "" {
		query += "thursday=?, "
		params = append(params, s.Thursday)
	}
	if s.Friday != "" {
		query += "friday=?, "
		params = append(params, s.Friday)
	}

	// 移除最后的逗号和空格，并添加WHERE子句
	query = strings.TrimSuffix(query, ", ") + " WHERE symbol=?"
	params = append(params, s.Symbol)

	// 执行SQL更新操作
	_, err := db.Exec(query, params...)
	if err != nil {
		zap.L().Error("Failed to update symbol", zap.Error(err))
		return err
	}
	zap.L().Info("Executing SQL", zap.String("query", query))

	return nil
}

//检查symbol是否是唯一行数据
func CheckSingleSymbolRow(symbol string) bool {
	var rowCount int

	query := "SELECT COUNT(*) FROM symbol where  symbol=?"
	err := db.Get(&rowCount, query, symbol)
	if err != nil {
		zap.L().Error("[symbol check] query symbol table row count", zap.Error(err))
		return false
	}
	return rowCount == 1
}

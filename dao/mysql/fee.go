package mysql

import (
	"fmt"
	"go.uber.org/zap"
	"strings"
)

type Fee struct {
	ID         int     `db:"id" json:"-"`
	Symbol     string  `db:"symbol" json:"symbol"`
	Group      string  `db:"group" json:"group"`
	Percentage int     `db:"percentage" json:"percentage"`
	Fee        float64 `db:"fee" json:"fee"`
	SwapLong   float64 `db:"swap_long" json:"swap_long"`
	SwapShort  float64 `db:"swap_short" json:"swap_short"`
	BidSpread  int     `db:"bid_spread" json:"bid_spread"`
	AskSpread  int     `db:"ask_spread" json:"ask_spread"`
}

//获取所有品种的费率
func GetFees() ([]Fee, error) {
	var fees []Fee
	query := `SELECT * FROM fee`
	err := db.Select(&fees, query)
	if err != nil {
		zap.L().Error("query fee failed", zap.Error(err))
		return nil, err
	}
	return fees, nil
}

//根据品种名获取品种费率
func GetFeesFromSymbol(symbol string) ([]Fee, error) {
	var fees []Fee
	query := `SELECT * FROM fee WHERE symbol = ?`
	err := db.Select(&fees, query, symbol)
	if err != nil {
		zap.L().Error("query fee by symbol failed", zap.Error(err))
		return nil, err
	}
	return fees, nil
}

//根据组获取该组下所有品种的费率
func GetFeesFromGroup(group string) ([]Fee, error) {
	var fees []Fee
	query := `SELECT * FROM fee WHERE ` + "`group`" + ` = ?`
	err := db.Select(&fees, query, group)
	if err != nil {
		zap.L().Error("query fee by group failed", zap.Error(err))
		return nil, err
	}
	return fees, nil
}

//在费率表里面新增组
func InsertGroupForFee(newGroup string) error {
	query := `
        INSERT INTO fee (symbol, "group", percentage, fee, swap_long, swap_short, bid_spread, ask_spread)
        SELECT symbol, $1, percentage, fee, swap_long, swap_short, bid_spread, ask_spread
        FROM fee
        WHERE "group" = 'APP_DEMO_RAW17'
    `
	_, err := db.Exec(query, newGroup)
	if err != nil {
		zap.L().Error("copy fee data to new group failed", zap.Error(err))
		return err
	}
	return nil
}

//新增费率表数据
func InsertFeeData(fee Fee) error {
	if fee.Group == "" || fee.Group == "*/" {
		//不写组则默认创建所有
		groups, err := GetGroups()
		if err != nil {
			zap.L().Error("[get group]Failed to get groups:", zap.Error(err))
			return err
		}

		for _, group := range groups {
			fee.Group = group.Group
			err := InsertFee(fee)
			if err != nil {
				zap.L().Error("[insert fee]Failed :", zap.Error(err))
				return err
			}
		}
		return nil
	} else {
		return InsertFee(fee)
	}
	return fmt.Errorf("failed to insert fee data")

}

//新增费率表 onece
func InsertFee(fee Fee) error {

	//空类型不执行
	if fee.Symbol == "" {
		return nil
	}

	var count int
	checkQuery := "SELECT COUNT(*) FROM fee WHERE symbol = ? AND `group` = ?"

	err := db.Get(&count, checkQuery, fee.Symbol, fee.Group)
	if err != nil {
		zap.L().Error("query count failed", zap.Error(err))
		return err
	}

	// 如果已经存在相同的记录，则直接返回
	if count > 0 {
		zap.L().Error("[sql error] symbol already exists", zap.String("symbol", fee.Symbol))
		return fmt.Errorf("symbol '%s' already exists", fee.Symbol)
	}

	query := "INSERT INTO fee (symbol, `group`, percentage, fee, swap_long, swap_short, bid_spread, ask_spread) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	_, err = db.Exec(query, fee.Symbol, fee.Group, fee.Percentage, fee.Fee, fee.SwapLong, fee.SwapShort, fee.BidSpread, fee.AskSpread)
	if err != nil {
		zap.L().Error("insert fee failed", zap.Error(err))
		return err
	}
	return nil
}

//根据品种删除Fee
func DeleteFeeBySymbol(symbol string) error {
	deleteQuery := "DELETE FROM fee WHERE symbol = ?"
	_, err := db.Exec(deleteQuery, symbol)
	if err != nil {
		zap.L().Error("insert fee failed", zap.Error(err))
		return err
	}
	return nil
}

func CheckFeeRow(symbol string) bool {
	var rowCount int

	query := "SELECT COUNT(*) FROM fee where  symbol=?"
	err := db.Get(&rowCount, query, symbol)
	if err != nil {
		zap.L().Error("[fee check] query fee table row count", zap.Error(err))
		return false
	}
	return rowCount >= 1
}

// UpdateFee 更新 Fee 记录
func UpdateFee(f *Fee) error {
	query := "UPDATE fee SET "
	var params []interface{}

	if f.Symbol != "" {
		query += "symbol=?, "
		params = append(params, f.Symbol)
	}
	if f.Group != "" {
		query += "`group`=?, "
		params = append(params, f.Group)
	}
	if f.Percentage != 0 {
		query += "percentage=?, "
		params = append(params, f.Percentage)
	}
	if f.Fee != 0.0 {
		query += "fee=?, "
		params = append(params, f.Fee)
	}
	if f.SwapLong != 0.0 {
		query += "swap_long=?, "
		params = append(params, f.SwapLong)
	}
	if f.SwapShort != 0.0 {
		query += "swap_short=?, "
		params = append(params, f.SwapShort)
	}
	if f.BidSpread != 0 {
		query += "bid_spread=?, "
		params = append(params, f.BidSpread)
	}
	if f.AskSpread != 0 {
		query += "ask_spread=?, "
		params = append(params, f.AskSpread)
	}

	// 移除最后的逗号和空格，并添加WHERE子句
	query = strings.TrimSuffix(query, ", ") + " WHERE symbol=?"
	params = append(params, f.Symbol)

	// 执行SQL更新操作
	_, err := db.Exec(query, params...)
	if err != nil {
		zap.L().Error("Failed to update fee", zap.Error(err))
		return err
	}
	zap.L().Info("Executing SQL", zap.String("query", query))

	return nil
}

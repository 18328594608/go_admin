package mysql

import (
	"fmt"
	"go_web_app/setting"

	"go.uber.org/zap"

	"github.com/spf13/viper"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func Close() {
	_ = db.Close()
}

func Init(config *setting.MysqlConfig) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		config.User, config.Password,
		config.Host, config.Port,
		config.DbName,
	)
	// 也可以使用MustConnect连接不成功就panic
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("connect DB failed, err:%v\n", zap.Error(err))
		return
	}
	db.SetMaxOpenConns(viper.GetInt("mysql.max_open_connection"))
	db.SetMaxIdleConns(viper.GetInt("mysql.max_idle_connection"))
	return
}

func LoadConfig() {
	symbols, err := GetSymbols()
	if err != nil {
		zap.L().Error("Failed to get symbols:", zap.Error(err))
		return
	}

	//获取详细symbol数据
	symbolMap := make(map[string]Symbol)
	for _, symbol := range symbols {
		zap.L().Info("Symbol details", zap.Any("symbol", symbol))
		symbolMap[symbol.Symbol] = symbol
	}
	fees, err := GetFees()
	if err != nil {
		zap.L().Error("Failed to get fees:", zap.Error(err))
		return
	}

	//获取详细fee数据
	feeMap := make(map[string]Fee)
	for _, fee := range fees {
		zap.L().Info("Fee details", zap.Any("fee", fee))
		feeMap[fee.Symbol] = fee
	}

	groups, err := GetGroups()
	if err != nil {
		zap.L().Error("Failed to get symbols:", zap.Error(err))
		return
	}

	for _, group := range groups {
		zap.L().Info("group details", zap.Any("group", group))
		financialData := NewFinancialData(group)

		for _, fee := range fees {
			if fee.Group == group.Group {
				zap.L().Info("Fee details", zap.Any("fee", fee))
				symbol, ok := symbolMap[fee.Symbol]
				if ok {
					symbolData := NewSymbolData(fee, symbol)
					financialData.SetSymbolData(symbol.Symbol, symbolData)
					GetMemoryInstance().SetData(fee.Group, *financialData)

				} else {
					zap.L().Error("symbolMap not Found", zap.Any("symbol", fee.Symbol))
				}
			}
		}
	}
}

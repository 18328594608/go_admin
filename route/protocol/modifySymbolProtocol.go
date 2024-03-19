package protocol

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go_web_app/dao/mysql"
	"go_web_app/route/reply"
)

type ModifySymbolHandler struct {
	Data ModifySymbolData
}

type ModifySymbolData struct {
	Symbol       string  `json:"symbol"`
	Security     *string `db:"security" json:"security,omitempty"`
	Digit        *int    `db:"digit" json:"digit,omitempty"`
	Currency     *string `db:"currency" json:"currency,omitempty"`
	ContractSize *int    `db:"contract_size" json:"contract_size,omitempty"`
	Percentage   *int    `db:"percentage" json:"percentage,omitempty"`
	MarginCalc   *int    `db:"margin_calc" json:"margin_calc,omitempty"`
	ProfitCalc   *int    `db:"profit_calc" json:"profit_calc,omitempty"`
	SwapCalc     *int    `db:"swap_calc" json:"swap_calc,omitempty"`
}

func (info *ModifySymbolData) NewCopyToSymbol() *mysql.Symbol {
	symbol := &mysql.Symbol{
		Symbol: info.Symbol,
	}

	if info.Security != nil {
		symbol.Security = *info.Security
	}

	if info.Digit != nil {
		symbol.Digit = *info.Digit
	}

	if info.Currency != nil {
		symbol.Currency = *info.Currency
	}

	if info.ContractSize != nil {
		symbol.ContractSize = *info.ContractSize
	}

	if info.Percentage != nil {
		symbol.Percentage = *info.Percentage
	}

	if info.MarginCalc != nil {
		symbol.MarginCalc = *info.MarginCalc
	}

	if info.ProfitCalc != nil {
		symbol.ProfitCalc = *info.ProfitCalc
	}

	if info.SwapCalc != nil {
		symbol.SwapCalc = *info.SwapCalc
	}

	return symbol
}

func (info *ModifySymbolData) NewCopyToFee() *mysql.Fee {
	fee := &mysql.Fee{
		Symbol: info.Symbol,
	}

	if info.Percentage != nil {
		fee.Percentage = *info.Percentage
	}

	return fee
}

func (h *ModifySymbolHandler) Handle(context *gin.Context) {

	symbol := h.Data.NewCopyToSymbol()
	if mysql.CheckSingleSymbolRow(symbol.Symbol) {
		err := mysql.UpdateSymbol(symbol)
		if err != nil {
			zap.L().Error("[update symbol] failed:", zap.Error(err))
			reply.NewResponseData(reply.CodeDatabaseError, "modifySymbol").Reply(context)
			return
		}
	} else {
		zap.L().Error("[update symbol] failed: symbol count error")
		reply.NewResponseData(reply.CodeDbCountError, "modifySymbol").Reply(context)
		return
	}

	if h.Data.Percentage != nil {
		fee := h.Data.NewCopyToFee()
		if mysql.CheckFeeRow(fee.Symbol) {
			err := mysql.UpdateFee(fee)
			if err != nil {
				zap.L().Error("[update fee] failed:", zap.Error(err))
				reply.NewResponseData(reply.CodeDatabaseError, "modifySymbol").Reply(context)
				return
			}
		} else {
			zap.L().Error("[update fee] failed: fee count error")
			reply.NewResponseData(reply.CodeDbCountError, "modifySymbol").Reply(context)
			return
		}
	}
	reply.NewResponseData(reply.CodeSuccess, "modifySymbol").Reply(context)

}

package protocol

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go_web_app/dao/mysql"
	"go_web_app/route/reply"
)

type ModifyFeeHandler struct {
	Data ModifyFeeData
}

type ModifyFeeData struct {
	Symbol     string   `db:"symbol" json:"symbol"`
	Percentage *int     `db:"percentage" json:"percentage,omitempty"`
	Fee        *float64 `db:"fee" json:"fee,omitempty"`
	SwapLong   *float64 `db:"swap_long" json:"swap_long,omitempty"`
	SwapShort  *float64 `db:"swap_short" json:"swap_short,omitempty"`
	BidSpread  *int     `db:"bid_spread" json:"bid_spread,omitempty"`
	AskSpread  *int     `db:"ask_spread" json:"ask_spread,omitempty"`
}

func (data *ModifyFeeData) NewFeeFromModifyFeeData() *mysql.Fee {
	fee := &mysql.Fee{
		Symbol: data.Symbol,
	}

	if data.Percentage != nil {
		fee.Percentage = *data.Percentage
	}
	if data.Fee != nil {
		fee.Fee = *data.Fee
	}
	if data.SwapLong != nil {
		fee.SwapLong = *data.SwapLong
	}
	if data.SwapShort != nil {
		fee.SwapShort = *data.SwapShort
	}
	if data.BidSpread != nil {
		fee.BidSpread = *data.BidSpread
	}
	if data.AskSpread != nil {
		fee.AskSpread = *data.AskSpread
	}

	return fee
}

func (data *ModifyFeeData) NewSymbolFromModifyFeeData() *mysql.Symbol {
	symbol := &mysql.Symbol{
		Symbol: data.Symbol,
	}

	if data.Percentage != nil {
		symbol.Percentage = *data.Percentage
	}

	return symbol
}

func (h *ModifyFeeHandler) Handle(context *gin.Context) {
	if h.Data.Percentage != nil {
		symbol := h.Data.NewSymbolFromModifyFeeData()
		if mysql.CheckSingleSymbolRow(symbol.Symbol) {
			err := mysql.UpdateSymbol(symbol)
			if err != nil {
				zap.L().Error("[update symbol] failed:", zap.Error(err))
				reply.NewResponseData(reply.CodeDatabaseError, "modifyFee").Reply(context)
				return
			}
		} else {
			zap.L().Error("[update symbol] failed: symbol count error")
			reply.NewResponseData(reply.CodeDbCountError, "modifyFee").Reply(context)
			return
		}
	}

	fee := h.Data.NewFeeFromModifyFeeData()
	if mysql.CheckFeeRow(fee.Symbol) {
		err := mysql.UpdateFee(fee)
		if err != nil {
			zap.L().Error("[update fee] failed:", zap.Error(err))
			reply.NewResponseData(reply.CodeDatabaseError, "modifyFee").Reply(context)
			return
		}
	} else {
		zap.L().Error("[update fee] failed: fee count error")
		reply.NewResponseData(reply.CodeDbCountError, "modifyFee").Reply(context)
		return
	}

	reply.NewResponseData(reply.CodeSuccess, "modifyFee").Reply(context)
}

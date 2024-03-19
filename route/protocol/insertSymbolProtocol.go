package protocol

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go_web_app/dao/mysql"
	"go_web_app/route/reply"
)

type InsertSymbolHandler struct {
	Data InsertSymbolData
}

type InsertSymbolData struct {
	Fee    mysql.Fee
	Symbol mysql.Symbol
}

func (h *InsertSymbolHandler) Handle(context *gin.Context) {
	symbol := mysql.Symbol{}
	symbol = h.Data.Symbol
	if symbol.Symbol != "" {
		err := mysql.InsertSymbol(symbol)
		if err != nil {
			zap.L().Error("Failed to Insert symbol:", zap.Error(err))
			reply.NewResponseData(reply.CodeDatabaseError, "insertSymbol").Reply(context)
			return
		}
	}

	fee := h.Data.Fee
	if fee.Symbol != "" {
		err := mysql.InsertFeeData(fee)
		if err != nil {
			zap.L().Error("Failed to Insert fee:", zap.Error(err))
			reply.NewResponseData(reply.CodeDatabaseError, "insertSymbol").Reply(context)
			return
		}
		reply.NewResponseData(reply.CodeSuccess, "insertSymbol").Reply(context)
	}
}

package protocol

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go_web_app/dao/mysql"
	"go_web_app/route/reply"
)

type DeleteSymbolHandler struct {
	Data DeleteSymbolData
}

type DeleteSymbolData struct {
	Symbol string `json:"symbol"`
}

func (h *DeleteSymbolHandler) Handle(context *gin.Context) {

	err := mysql.DeleteFeeBySymbol(h.Data.Symbol)
	if err != nil {
		zap.L().Error("[delete fee] failed:", zap.Error(err))
		reply.NewResponseData(reply.CodeDatabaseError, "deleteSymbol").Reply(context)
		return
	}

	err = mysql.DeleteSymbolBySymbol(h.Data.Symbol)
	if err != nil {
		zap.L().Error("[delete symbol] failed:", zap.Error(err))
		reply.NewResponseData(reply.CodeDatabaseError, "deleteSymbol").Reply(context)
		return
	}

	reply.NewResponseData(reply.CodeSuccess, "deleteSymbol").Reply(context)
}

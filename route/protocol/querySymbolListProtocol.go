package protocol

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go_web_app/dao/mysql"
	"go_web_app/route/reply"
)

type SymbolListHandler struct {
	Data SymbolListData
}

type SymbolListData struct {
}

func (h *SymbolListHandler) Handle(context *gin.Context) {
	symbols, err := mysql.GetSymbols()
	if err != nil {
		zap.L().Error("Failed to get symbols:", zap.Error(err))
		reply.NewResponseData(reply.CodeInternalServerError, "symbolList").Reply(context)
		return
	}

	//获取详细symbol数据
	var symbolNames []string
	for _, sym := range symbols {
		symbolNames = append(symbolNames, sym.Symbol)
	}
	reply.NewResponseData(reply.CodeSuccess, "symbolList", symbolNames).Reply(context)
}

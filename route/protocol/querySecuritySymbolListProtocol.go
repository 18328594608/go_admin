package protocol

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go_web_app/dao/mysql"
	"go_web_app/route/reply"
)

type SecuritySymbolListHandler struct {
	Data SecuritySymbolListData
}

type SecuritySymbolListData struct {
}

func (h *SecuritySymbolListHandler) Handle(context *gin.Context) {
	symbols, err := mysql.GetSymbols()
	if err != nil {
		zap.L().Error("Failed to get symbols:", zap.Error(err))
		reply.NewResponseData(reply.CodeInternalServerError, "querySecuritySymbolList").Reply(context)
		return
	}

	// 获取详细symbol数据
	symbolMap := make(map[string][]string)
	for _, symbol := range symbols {
		if _, ok := symbolMap[symbol.Security]; !ok {
			symbolMap[symbol.Security] = []string{}
		}
		symbolMap[symbol.Security] = append(symbolMap[symbol.Security], symbol.Symbol)
	}
	reply.NewResponseData(reply.CodeSuccess, "querySecuritySymbolList", symbolMap).Reply(context)
}

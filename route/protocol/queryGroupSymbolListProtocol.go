package protocol

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go_web_app/dao/mysql"
	"go_web_app/route/reply"
)

type GroupSymbolListHandler struct {
	Data GroupSymbolListData
}

type GroupSymbolListData struct {
}

func (h *GroupSymbolListHandler) Handle(context *gin.Context) {
	fees, err := mysql.GetFees()
	if err != nil {
		zap.L().Error("Failed to get symbols:", zap.Error(err))
		reply.NewResponseData(reply.CodeInternalServerError, "queryGroupSymbolList").Reply(context)
		return
	}

	// 获取详细symbol数据
	symbolMap := make(map[string][]string)
	for _, symbol := range fees {
		if _, ok := symbolMap[symbol.Group]; !ok {
			symbolMap[symbol.Group] = []string{}
		}
		symbolMap[symbol.Group] = append(symbolMap[symbol.Group], symbol.Symbol)
	}
	reply.NewResponseData(reply.CodeSuccess, "queryGroupSymbolList", symbolMap).Reply(context)
}

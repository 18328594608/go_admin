package protocol

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go_web_app/dao/mysql"
	"go_web_app/route/reply"
)

type SecurityListHandler struct {
	Data SecurityListData
}

type SecurityListData struct {
}

func (h *SecurityListHandler) Handle(context *gin.Context) {
	securityList, err := mysql.GetDistinctSecurity()
	if err != nil {
		zap.L().Error("Failed to get securityList:", zap.Error(err))
		reply.NewResponseData(reply.CodeInternalServerError, "querySecurityList").Reply(context)
		return
	}
	reply.NewResponseData(reply.CodeSuccess, "querySecurityList", securityList).Reply(context)
}

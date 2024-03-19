package protocol

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go_web_app/dao/mysql"
	"go_web_app/route/reply"
)

type GroupListHandler struct {
	Data GroupListData
}

type GroupListData struct {
}

func (h *GroupListHandler) Handle(context *gin.Context) {
	groups, err := mysql.GetGroups()
	if err != nil {
		zap.L().Error("Failed to get symbols:", zap.Error(err))
		reply.NewResponseData(reply.CodeInternalServerError, "groupList").Reply(context)
		return
	}

	//获取详细symbol数据
	var groupNames []string
	for _, group := range groups {
		groupNames = append(groupNames, group.Group)
	}
	reply.NewResponseData(reply.CodeSuccess, "groupList", groupNames).Reply(context)
}

package protocol

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go_web_app/dao/mysql"
	"go_web_app/route/reply"
)

type GroupInfoHandler struct {
	Data GroupInfoData
}

type GroupInfoData struct {
	Group string `json:"group"`
}

type GroupInfo struct {
	Group      string  `json:"group"`
	Point      int     `json:"point"`
	ModifyTime int64   `json:"modify_time"`
	Leverage   int     `json:"leverage"`
	Fee        float64 `db:"fee" json:"fee"`
	SwapLong   float64 `db:"swap_long" json:"swap_long"`
	SwapShort  float64 `db:"swap_short" json:"swap_short"`
}

func (gi *GroupInfo) CopyFromGroup(group mysql.Group) {
	gi.Leverage = group.Leverage
	gi.Group = group.Group
}

func (gi *GroupInfo) CopyFromFee(fee mysql.Fee) {
	gi.Group = fee.Group
	gi.Fee = fee.Fee
	gi.SwapLong = fee.SwapLong
	gi.SwapShort = fee.SwapShort
}

func (h *GroupInfoHandler) Handle(context *gin.Context) {

	group, err := mysql.GetGroupByGroup(h.Data.Group)
	if err != nil {
		zap.L().Error("Failed to get groups:", zap.Error(err))
		return
	}

	groupInfo := GroupInfo{}
	groupInfo.CopyFromGroup(group)

	fees, err := mysql.GetFeesFromGroup(h.Data.Group)
	if err != nil {
		zap.L().Error("Failed to get fee:", zap.Error(err))
		return
	}
	//获取详细fee数据
	if len(fees) > 0 {
		firstFee := mysql.Fee{}
		firstFee = fees[0]
		groupInfo.CopyFromFee(firstFee)
	}
	reply.NewResponseData(reply.CodeSuccess, "groupInfo", groupInfo).Reply(context)
}

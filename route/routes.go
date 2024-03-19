package route

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go_web_app/logger"
	"go_web_app/route/protocol"
	"net/http"
)

type RequestBody struct {
	Cmd  string          `json:"cmd"`
	Data json.RawMessage `json:"data"`
}

func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(mode)
	}
	r := gin.New()
	// 最重要的就是这个日志库
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	r.GET("/test", func(context *gin.Context) {
		context.String(http.StatusOK, "11111")
	})
	r.POST("/", func(context *gin.Context) {
		var requestBody RequestBody
		// 解析JSON请求体到结构体
		if err := context.BindJSON(&requestBody); err != nil {
			// 如果解析出错，返回400状态码和错误信息
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var handler CmdHandler
		switch requestBody.Cmd {
		case "querySymbol":
			var data protocol.SymbolInfoData
			if err := json.Unmarshal(requestBody.Data, &data); err != nil {
				context.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "Invalid data for symbolInfo"})
				return
			}
			handler = &protocol.SymbolInfoHandler{Data: data}
		case "queryGroup":
			var data protocol.GroupInfoData
			if err := json.Unmarshal(requestBody.Data, &data); err != nil {
				context.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "Invalid data for groupInfo"})
				return
			}
			handler = &protocol.GroupInfoHandler{Data: data}
		case "insertSymbol":
			var data protocol.InsertSymbolData
			if err := json.Unmarshal(requestBody.Data, &data); err != nil {
				context.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "Invalid data for insertSymbol"})
				return
			}
			handler = &protocol.InsertSymbolHandler{Data: data}
		case "deleteSymbol":
			var data protocol.DeleteSymbolData
			if err := json.Unmarshal(requestBody.Data, &data); err != nil {
				context.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "Invalid data for deleteSymbol"})
				return
			}
			handler = &protocol.DeleteSymbolHandler{Data: data}

		case "modifyFee":
			var data protocol.ModifyFeeData
			if err := json.Unmarshal(requestBody.Data, &data); err != nil {
				context.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "Invalid data for modifyFee"})
				return
			}
			handler = &protocol.ModifyFeeHandler{Data: data}
		case "symbolList":
			var data protocol.SymbolListData
			if err := json.Unmarshal(requestBody.Data, &data); err != nil {
				context.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "Invalid data for modifyFee"})
				return
			}
			handler = &protocol.SymbolListHandler{Data: data}
		case "querySecurityList":
			var data protocol.SecurityListData
			if err := json.Unmarshal(requestBody.Data, &data); err != nil {
				context.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "Invalid data for querySecurityList"})
				return
			}
			handler = &protocol.SecurityListHandler{Data: data}
		case "querySecuritySymbolList":
			var data protocol.SecuritySymbolListData
			if err := json.Unmarshal(requestBody.Data, &data); err != nil {
				context.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "Invalid data for querySecuritySymbolList"})
				return
			}
			handler = &protocol.SecuritySymbolListHandler{Data: data}
		case "modifySymbol":
			var data protocol.ModifySymbolData
			if err := json.Unmarshal(requestBody.Data, &data); err != nil {
				context.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "Invalid data for querySecuritySymbolList"})
				return
			}
			handler = &protocol.ModifySymbolHandler{Data: data}
		case "groupList":
			var data protocol.GroupListData
			if err := json.Unmarshal(requestBody.Data, &data); err != nil {
				context.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "Invalid data for querySecuritySymbolList"})
				return
			}
			handler = &protocol.GroupListHandler{Data: data}
		case "queryGroupSymbolList":
			var data protocol.GroupSymbolListData
			if err := json.Unmarshal(requestBody.Data, &data); err != nil {
				context.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "Invalid data for querySecuritySymbolList"})
				return
			}
			handler = &protocol.GroupSymbolListHandler{Data: data}

		default:
			context.JSON(http.StatusBadRequest, gin.H{"code": 3, "message": "Unknown cmd"})
			return
		}

		// 使用多态处理命令
		handler.Handle(context)
	})
	return r
}

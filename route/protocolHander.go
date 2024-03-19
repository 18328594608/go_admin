package route

import "github.com/gin-gonic/gin"

// CmdHandler 接口定义了处理命令的方法
type CmdHandler interface {
	Handle(context *gin.Context)
}

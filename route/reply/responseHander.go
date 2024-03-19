package reply

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ErrorCode int

const (
	CodeSuccess             ErrorCode = iota // 0
	CodeInvalidInput                         // 1
	CodeDatabaseError                        // 2
	CodeNotFound                             // 3
	CodeUnauthorized                         // 4
	CodeForbidden                            // 5
	CodeTooManyRequests                      // 6
	CodeInternalServerError                  // 7
	CodeServiceUnavailable                   // 8
	CodeTimeout                              // 9
	CodeConflict                             // 10
	CodeDbCountError
)

var errorCodeMessages = map[ErrorCode]string{
	CodeSuccess:             "Success",
	CodeInvalidInput:        "Invalid input provided",
	CodeDatabaseError:       "Database operation failed",
	CodeNotFound:            "Requested resource not found",
	CodeUnauthorized:        "Unauthorized access",
	CodeForbidden:           "Forbidden operation",
	CodeTooManyRequests:     "Too many requests",
	CodeInternalServerError: "Internal server error",
	CodeServiceUnavailable:  "Service unavailable",
	CodeTimeout:             "Operation timed out",
	CodeConflict:            "Resource conflict",
	CodeDbCountError:        "Database count error",
}

// GetErrorMessage 返回给定错误代码对应的错误消息
func GetErrorMessage(code ErrorCode) string {
	message, exists := errorCodeMessages[code]
	if !exists {
		return "Unknown error"
	}
	return message
}

type ResponseData struct {
	Code    int         `json:"code"`    // 错误码，0表示成功，非0表示失败
	Cmd     string      `json:"cmd"`     // 操作命令或类型
	Data    interface{} `json:"data"`    // 返回的数据
	Message string      `json:"message"` // 成功或错误消息
}

func NewResponseData(code ErrorCode, cmd string, args ...interface{}) *ResponseData {
	response := &ResponseData{
		Code: int(code),
		Cmd:  cmd,
	}

	if code != 0 {
		// 获取错误消息
		response.Message = GetErrorMessage(ErrorCode(code))
	} else if len(args) > 0 {
		// 如果code为0且提供了额外参数，将第一个参数视为Data
		response.Data = args[0]
		response.Message = "Success"
	}

	return response
}

func (responseData *ResponseData) Reply(context *gin.Context) {
	context.JSON(http.StatusOK, responseData)
}

package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/outsstill/go-kit/logger"
)

func Success(c *gin.Context, msg ...string) {
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"success": true,
		"message": defaultMessage("操作成功", msg...),
	})
}

func Data(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "操作成功！",
		"data":    data,
		"code":    http.StatusOK,
	})
}
func Fail(c *gin.Context, msg ...string) {
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusBadRequest,
		"success": false,
		"message": defaultMessage("请求处理失败，请查看 error 的值", msg...),
	})
}

func Error(c *gin.Context, err error, msg ...string) {
	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"code":    http.StatusInternalServerError,
		"success": false,
		"message": defaultMessage("请求处理失败，请查看 error 的值", msg...),
		"error":   err.Error(),
	})
}

func Abort(c *gin.Context, code int, msg ...string) {
	c.AbortWithStatusJSON(code, gin.H{
		"code":    code,
		"message": defaultMessage("数据不存在，请确定请求正确", msg...),
	})
}

// Abort404 响应 404，未传参 msg 时使用默认消息
func Abort404(c *gin.Context, msg ...string) {
	c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
		"code":    http.StatusNotFound,
		"message": defaultMessage("数据不存在，请确定请求正确", msg...),
	})
}

// Abort403 响应 403，未传参 msg 时使用默认消息
func Abort403(c *gin.Context, msg ...string) {
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		"code":    http.StatusForbidden,
		"message": defaultMessage("权限不足，请确定您有对应的权限", msg...),
	})
}

// Abort403 响应 403，未传参 msg 时使用默认消息
func Abort401(c *gin.Context, msg ...string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"code":    http.StatusUnauthorized,
		"message": defaultMessage("权限不足，请确定您有对应的权限", msg...),
	})
}

// Abort500 响应 500，未传参 msg 时使用默认消息
func Abort500(c *gin.Context, msg ...string) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"code":    http.StatusInternalServerError,
		"message": defaultMessage("服务器内部错误，请稍后再试", msg...),
	})
}

// BadRequest 响应 400，传参 err 对象，未传参 msg 时使用默认消息
// 在解析用户请求，请求的格式或者方法不符合预期时调用
func BadRequest(c *gin.Context, err error, msg ...string) {
	logger.LogIf(err)
	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"code":    http.StatusInternalServerError,
		"message": defaultMessage("请求解析错误，请确认请求格式是否正确。上传文件请使用 multipart 标头，参数请使用 JSON 格式。", msg...),
		"error":   err.Error(),
	})
}

// Unauthorized 响应 401，未传参 msg 时使用默认消息
// 登录失败、jwt 解析失败时调用
func Unauthorized(c *gin.Context, msg ...string) {
	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"code":    http.StatusUnauthorized,
		"message": defaultMessage("请求解析错误，请确认请求格式是否正确。上传文件请使用 multipart 标头，参数请使用 JSON 格式。", msg...),
	})
}

// token认证失败，返回特定code让前端重新登录
func AuthFail(c *gin.Context, msg ...string) {
	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"code":    401,
		"message": defaultMessage("请求解析错误，请确认请求格式是否正确。上传文件请使用 multipart 标头，参数请使用 JSON 格式。", msg...),
	})
}

func ValidationFields(c *gin.Context, errors map[string][]string) {
	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"code":    http.StatusUnprocessableEntity,
		"message": "请求验证不通过，具体请查看 errors",
		"errors":  errors,
	})
}

func ValidationError(c *gin.Context, errors map[string][]string) {
	c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
		"code":    http.StatusUnprocessableEntity,
		"message": "请求验证不通过，具体请查看 errors",
		"errors":  errors,
	})
}

// defaultMessage 内用的辅助函数，用以支持默认参数默认值
// Go 不支持参数默认值，只能使用多变参数来实现类似效果
func defaultMessage(defaultMsg string, msg ...string) (message string) {
	if len(msg) > 0 {
		message = msg[0]
	} else {
		message = defaultMsg
	}
	return
}

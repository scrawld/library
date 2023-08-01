package ginx

type HttpStatus int

const (
	HttpStatusOk           HttpStatus = 200
	HttpStatusParamsErr    HttpStatus = 400 // 参数验证错误
	HttpStatusAuthErr      HttpStatus = 401 // 授权过期
	HttpStatusForbidden    HttpStatus = 403 // 无权限
	HttpStatusNotFound     HttpStatus = 404 // 内容不存在
	HttpStatusNotAllowed   HttpStatus = 405 // 不允许使用
	HttpStatusPayloadLarge HttpStatus = 413 // 请求体过大
	HttpStatusUnMediaType  HttpStatus = 415 // 不支持的媒体类型
	HttpStatusRetryWith    HttpStatus = 449 // 稍后重试
	HttpStatusServerErr    HttpStatus = 500 // 服务器内部错误
)

type UserType string

const (
	UserTypeAny UserType = "any" // 允许全部用户访问
)

package ginx

// 通用错误码
const (
	ErrCodeParams               = 400 // 参数验证错误
	ErrCodeAuthExpire           = 401 // 授权过期
	ErrCodeForbidden            = 403 // 无权限
	ErrCodeNotFound             = 404 // 内容不存在
	ErrCodeNotAllowed           = 405 // 不允许使用
	ErrCodePayloadTooLarge      = 413 // 请求体过大
	ErrCodeUnsupportedMediaType = 415 // 不支持的媒体类型
	ErrCodeRetryLater           = 449 // 稍后重试
	ErrCodeServer               = 500 // 服务器内部错误
)

type Role string

const (
	RoleAny Role = ""
)

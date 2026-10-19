# ginx

## NewContext 请求上下文封装, 实现了日志链路追踪、参数解析、参数校验

### 使用示例

```go
package main

import (
	_ "yourpackage/library/ginx"
	ginxMiddleware "yourpackage/library/library/ginx/middleware"
)

router := gin.Default()
router.Use(ginxMiddleware.RequestId()).Use(ginxMiddleware.Logger())

type UserIndexReq struct {
   UserId    int      `json:"userId" binding:"required"`
   UserName  string   `json:"userName" binding:"required"`
}

func (cont *UserController) Index(c *gin.Context) {
   p := &UserIndexReq{}
   ctx, ok := ginx.NewContext(c, ginx.RoleAny, p)
	if !ok {
		return
	}
	r := UserIndexResp{}

	ctx.Log.Infof("userId %s", p.UserId)
	ctx.Log.Errorf("userId %s", p.UserId)

	// ...
	if err != nil {
		ctx.RenderServerError(fmt.Errorf("search user report error: %s", err))
		return
	}
	ctx.Render(r)
}
```

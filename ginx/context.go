package ginx

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/scrawld/library/zaplog"

	"github.com/gin-gonic/gin"
)

type Context struct {
	ctx    *gin.Context
	Log    *zaplog.TracingLogger
	Params interface{}
}

/**
 * NewContext 请求上下文封装, 实现了日志链路追踪、参数解析、参数校验
 *
 * Example:
 *
 * import (
 * 	_ "github.com/scrawld/library/ginx"
 * 	ginxMiddleware "github.com/scrawld/library/library/ginx/middleware"
 * )
 *
 * router := gin.Default()
 * router.Use(ginxMiddleware.RequestId()).Use(ginxMiddleware.Logger())
 *
 * type UserIndexReq struct {
 *	UserId    int      `json:"userId" binding:"required"`
 *	UserName  string   `json:"userName" binding:"required"`
 * }
 *
 * func (cont *UserController) Index(c *gin.Context) {
 *	p := &UserIndexReq{}
 *	ctx, ok := ginx.NewContext(c, ginx.UserTypeAny, p)
 * 	if !ok {
 * 		return
 * 	}
 * 	r := UserIndexResp{}
 *
 * 	ctx.Log.Infof("userId %s", p.UserId)
 * 	ctx.Log.Errorf("userId %s", p.UserId)
 *
 * 	// ...
 * 	if err != nil {
 * 		ctx.RenderServerError(fmt.Errorf("search user report error: %s", err))
 * 		return
 * 	}
 * 	ctx.Render(r)
 * }
 */
func NewContext(ctx *gin.Context, utype UserType, params interface{}) (r *Context, isAllow bool) {
	r = &Context{
		ctx:    ctx,
		Log:    zaplog.New(ctx.GetHeader("X-Request-Id")).Named("API"),
		Params: params,
	}
	r.ctx.Set("logger", r.Log)

	defer func() {
		// request log
		query := r.ctx.Request.URL.RawQuery
		if query != "" {
			query = "?" + query
		}
		r.Log.Infof("%s | %s %s | %+v", r.ctx.ClientIP(), r.ctx.Request.Method, r.ctx.Request.URL.Path+query, params)
	}()

	// validate params
	if err := r.ShouldBind(params); err != nil {
		r.Log.Infof("Params validate error: %s", err)
		r.RenderError(HttpStatusParamsErr, fmt.Errorf("Params validate error: %s", err))
		return
	}
	isAllow = true
	return
}

func (c *Context) ShouldBind(obj interface{}) error {
	if obj == nil {
		return nil
	}
	return c.ctx.ShouldBind(obj)
}

func (c *Context) ClientIP() string {
	return c.ctx.ClientIP()
}

/************ Render **************/
type RenderStruct struct {
	Code HttpStatus  `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// WriteExcel 将excel文件返回
func (c *Context) WriteExcel(h []byte, filename string) {
	if len(h) == 0 {
		return
	}
	c.ctx.Header("Content-Type", "application/vnd.ms-excel;charset=utf8")
	c.ctx.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.xlsx"`, url.QueryEscape(filename)))
	c.ctx.Writer.Write(h)
}

// Render 成功返回并打印日志
func (c *Context) Render(data interface{}) {
	if data == nil {
		data = []string{}
	}
	h := RenderStruct{
		Code: HttpStatusOk,
		Data: data,
	}
	c.Log.Infof("Response: %+v", h)
	c.ctx.JSON(http.StatusOK, h)
}

// RenderNotLog 成功返回无日志
func (c *Context) RenderNotLog(data interface{}) {
	if data == nil {
		data = []string{}
	}
	h := RenderStruct{
		Code: HttpStatusOk,
		Data: data,
	}
	c.ctx.JSON(http.StatusOK, h)
}

// RenderError 错误返回
func (c *Context) RenderError(code HttpStatus, e error) {
	h := RenderStruct{
		Code: code,
		Msg:  e.Error(),
		Data: []string{},
	}
	c.Log.Infof("Response: %+v", h)
	c.ctx.JSON(http.StatusOK, h)
}

// RenderRealError 错误返回并将realE以Error日志打印
func (c *Context) RenderRealError(code HttpStatus, userE error, realE error) {
	if realE != nil {
		c.Log.Warnf(realE.Error())
	}
	if userE == nil {
		userE = errors.New("Server exception")
	}
	h := RenderStruct{
		Code: code,
		Msg:  userE.Error(),
		Data: []string{},
	}
	c.Log.Infof("Response: %+v", h)
	c.ctx.JSON(http.StatusOK, h)
}

// RenderServerError 返回系统错误500并打印e
func (c *Context) RenderServerError(e error) {
	c.RenderRealError(HttpStatusServerErr, e, e)
}

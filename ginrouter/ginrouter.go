package ginrouter

import (
	"github.com/Techit-Kakaew/urlshorten/httprouter"
	"github.com/gin-gonic/gin"
)

func NewHandler(handler func(httprouter.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := NewContext(c)
		handler(ctx)
	}
}

type context struct {
	ctx *gin.Context
}

func NewContext(c *gin.Context) *context {
	return &context{ctx: c}
}

func (c *context) PathValue(key string) string {
	return c.ctx.Param(key)
}

func (c *context) Bind(v any) error {
	return c.ctx.Bind(v)
}

func (c *context) Status(code int) {
	c.ctx.Status(code)
}

func (c *context) JSON(code int, v any) {
	c.ctx.JSON(code, v)
}

func (c *context) Redirect(code int, url string) {
	c.ctx.Redirect(code, url)
}

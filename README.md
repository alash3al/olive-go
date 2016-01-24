# olive.go
Just a lightweight golang web application middleware

# Author
[Mohammed Al Ashaal, `a full-stack developer`](http://www.alash3al.xyz)

# Install
`go get github.com/alash3al/olive-go`

# Docs
[Godoc](http://godoc.org/github.com/alash3al/olive-go)

# Quick overview:
```go
package main

import "net/http"
import "github.com/alash3al/olive-go"

func main() {
	olive.New().GET("/", func(ctx *olive.Context) bool {
		ctx.SetBody("index")
		// return false = "don't run the next matched route with the same method and pattern if any"
		// this feature allows you to run multiple routes with the same properties
		return false
	}).ANY("/page/?(.*?)", func(ctx *olive.Context) bool {
		var body []byte
		ctx.LimitBody(20)
		err := ctx.GetBody(&body)
		ctx.SetBody("this is your input \n")
		ctx.SetBody(body)
		_ = err
		return true
	}).GET("/page", func(ctx *olive.Context) bool {
		ctx.SetBody([]byte("hi !"))
		return false
	}).POST("/page/([^/]+)/and/([^/]+)", func(ctx *olive.Context) bool {
		var input map[string]string
		ctx.GetBody(&input, 512) // parse the request body into {input} and returns error if any
		ctx.SetBody(ctx.Params)
		return false
	}).GroupBy("path", "/api/v1", func(ApiV1 *olive.App){
		ApiV1.GET("/ok", func(ctx *olive.Context) bool {
			ctx.Res.Write([]byte("api/v1/ok"))
			return false
		}).GET("/page/([^/]+)/and/([^/]+)", func(ctx *olive.Context) bool {
			ctx.Res.Write([]byte("api/v1/ " + ctx.Params[0] + " " + ctx.Params[1]))
			return false
		})
	}).ANY("?.*?", olive.Handler(http.NotFoundHandler(), false)).Listen(":80")
}
```

# Changes

**Version 3.0**
- `Context.GetQuery` now accepts new param called `body` and its type is bool, so you can get the request body as url-decoded as url.Values
- `Context.GetBody` now accepts one paramater, and you don't need to `make([]byte, ...)` just pass a `&v` where `v` is `[]byte`
- added `Context.LimitBody` to limit the request body to prevent any memory-leaks attacks while reading it using `Context.GetBody` .

**Version 2.0**
- removed panics handler
- removed `Context.AddHeaders()` and `Context.SetHeaders()`
- added `Context.DelHeader()`
- renamed `Context.Query()` to `Context.GetQuery()`
- renamed `Context.Body()` to `Context.GetBody()`
- renamed `Context.Send()` to `Context.SetBody()`
- added support for html templates in `Context.SetBody()`
- renamed `App.Group()` to `App.GroupBy`
- add support for custom vhost routing

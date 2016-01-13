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
			ctx.Res.Write([]byte("index"))
			// return false = "don't run the next matched route with the same method and pattern if any"
			// this feature allows yout to run multiple routes with the same properties
			return false
		}).CONNECT("/", func(ctx *olive.Context) bool {
			// olive automatically catch any panic and recover it to the 
			// "stdout" using "log" package .
			panic("connect method !")
			return false
		}).ANY("/page/?(.*?)", func(ctx *olive.Context) bool {
			ctx.Res.Write([]byte("i'm the parent \n"))
			return true
		}).GET("/page", func(ctx *olive.Context) bool {
			ctx.Res.Write([]byte("page"))
			return false
		}).GET("/page/([^/]+)/and/([^/]+)", func(ctx *olive.Context) bool {
			ctx.Res.Write([]byte(ctx.Params[0] + " " + ctx.Params[1]))
			return false
		}).Group("/api/v1", func(ApiV1 *olive.App){
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

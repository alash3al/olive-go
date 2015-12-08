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

	import ( "github.com/alash3al/olive-go"; "net/http" )

	func main(){
		// lets set our base "localhost.me"
		// this app will only work on "localhost.me/*"
		// By default the hostname is "*"
		// so it will work on any hostname, but here
		// we change it to 'localhost.me' which is a free service
		// that routes all requests to your own localhost server on port 80,
		// you can use regex strings too .
		app := olive.NewApp().SetHostname("localhost.me")

		// just like the main http library
		// but we have more features such as "regex strings" .
		// You can chain multiple handleFunc calls from this one easily .
		app.HandleFunc(`/hello`, func(ctx *olive.Context){
			ctx.WriteJSON(map[string]string{
				"message": "hello world",
			})
		})

			// creating sub apps (routes) is so easy .
			api := app.NewSubApp().SetHostname(`api.localtest.me`)

				// nested sub-apps
				apiAuth := api.NewSubApp().SetPath(`/auth/?.*`)

				// login
				apiAuth.HandleFunc(`/login`, func(ctx *Context){
					ctx.WriteJSON(map[string]string{
						"message": "login",
					})
				}).SetMethod(`POST`)

		app.ListenAndServe(":80")
	}
```

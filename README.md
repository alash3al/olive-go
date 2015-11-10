# olive.go
Just a very simple and elegant http(s) middleware (router)

# Author
[Mohammed Al Ashaal, `a full-stack developer`](http://www.alash3al.xyz)

# Quick overview
```go
	package main

	import ( "github.com/alash3al/olive-go"; "net/http" )

	func main(){
		app := olive.NewApp()

		// create new route
		// Path = "*" --> any path
		// Method = "*" --> any method
		// Host = "*" --> any host
		// "*" -> is the default value for each property
		// Exclusive =  true, "true is the default"
		// 		it means that once it matches the current request,
		//		stop and don't run similar routes with the same request properties .

		app.Factory().SetPath("/tst").SetHost("*").SetMethod("*").SetFunc(func(c *olive.Context){
			c.Res.Write([]byte("tst"))
		})

		app.Factory().SetPath(`/page/([^/]+)`).SetHost("*").SetMethod("*").SetFunc(func(c *olive.Context){
			c.Res.Write([]byte("current-page: " + c.Args[0]))
		})

		app.Factory().SetPath("*").SetHost("cdn.mysite.com").SetHandler(http.FileServer(http.Dir(`/root/cdn/`)))

		app.Add(
			olive.NewRoute().SetPath("/new"), // ... and so on, multiple routes are supported
		)

		app.ListenAndServe(":80")
		// or
		// app.ListenAndServeTLS(":433", ........)
	}
```

# Documentation
Documentation is available on [Godoc](https://godoc.org/github.com/alash3al/olive-go) 

# Installation
`go get github.com/alash3al/olive-go`

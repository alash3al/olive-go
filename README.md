# olive.go
Just a lightweight golang web application middleware

# Author
[Mohammed Al Ashaal, `a full-stack developer`](http://www.alash3al.xyz)

# Structures
`olive.go` is just a middleware that implements `http.Handler` interface, so you can use it as you want .
in `olive` there are just 3 structures:

### 1)- olive.Context
> contains the `*http.Request`, `http.ResponseWriter` and `Args []string`,  
> `Args` are the arguments from the matched routes `regexp results`

```go

type Context struct {
	Req	*http.Request
	Res 	http.ResponseWriter
	Args	[]string
}

```

### 2)- olive.Route
> contains the data we need in each route.  

```go

type Route struct {
	method		string
	vhost		string
	path		*regexp.Regexp
	callback	func(Context)
}

```

> It also contains the following methods, to controle when and where the route will be dispatched .

```go

  // --> set the route's path `or posix regexp string`
  Route.SetPath(p string)
  
  // --> set the route's method `or posix regexp string`
  // --> `*` means any request method
  Route.SetMethod(m string)

  // --> set the route's vhost 'subdomain' `or posix regexp string`
  Route.SetVhost(v string)

  // --> set the route's callback
  Route.SetHandler(fn func(olive.Context))

```

### 3)- olive.Handler
> the main wrapper that implements the `http.Handler` interface, this is the structure  

```go

type Handler struct {
	routes []*Route
}

```

> and the following methods

```go
 
 // --> add new route ?
 // --> it returns a `*Route` so you can customize it as descriped above
 Handler.HandleFunc(func(olive.Context))
 
 // --> ServeHTTP
 // its just an implementation of http.Handler
 
 // --> ListenAndServe
 Handler.ListenAndServe(addr string) error
 
 // --> listenAndServeTLS
 Handler.ListenAndServeTLS(addr string, certFile string, keyFile string) error
```

# Lets learn it

```go

import(
	"github.com/alash3al/olive-go"
	"net/http" // just for custom handlers
)

func main() {
	// initialize it
	app := olive.NewHandler()
	
	// new handler for `localtest.me/hello-world`
	// >> `localtest.me` is a free service that routes all requests to your own `localhost`
	app.HandleFunc(func(o olive.Context){
		o.Res.Write([]byte(`Hello World`))
	}).SetPath(`hello-world`)

	// new handler for `api.localtest.me/<anything>`
	app.HandleFunc(func(o olive.Context){
		o.Res.Write([]byte(`current path is ` + o.Args[0]))
	}).SetPath(`?(.*?)`.SetVhost(`api.localhost.me`))

	// new handler for `<anything>.localtest.me/<anything>`
	app.HandleFunc(func(o olive.Context){
		// vhost args will be the first in the args array
		// path args will be the last in the args array
		o.Res.Write([]byte(`current vhost is ` + o.Args[0] + `, path is ` + o.Args[1]))
	}).SetPath(`?(.*?)`.SetVhost(`?(.*?).localhost.me`))

	// new handler for `localtest.me/assets/` 
	// to handle static files from `/root/www/`
	app.HandleFunc(func(o olive.Context){
		http.StripPrefix(`/assets/`, http.FileServer(http.Dir(`/root/www/`))).ServeHTTP(c.Res, c.Req)
	}).SetPath(`assets`)

	// listen on port '80'
	app.ListenAndServe(`:80`)
}

```

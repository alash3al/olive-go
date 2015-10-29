# olive.go
Just a lightweight golang web application middleware

# Author
(Mohammed Al Ashaal, `a full-stack developer`)[http://www.alash3al.xyz]

# Structures
`olive.go` is just a middleware that implements `http.Handler` interface, so you can use it as you want .
in `olive` there are just 3 structures:

### 1)- olive.Context
> contains the `*http.Request`, `http.ResponseWriter` and `Args []string`,  
> `Args` are the arguments from the matched routes `regexp results`

```go

type Context struct {
	Req		*http.Request
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

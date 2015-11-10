// Olive - a lightweight simple Golang web application framework .
/*
	// Example
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
			olive.NewRoute().SetPath("*"), // ... and so on
		)

		app.ListenAndServe(":80")
		// or
		// app.ListenAndServeTLS(":433", ........)
	}
*/
package olive

// Our Requirements
import (
	"net/http"
	"regexp"
	"strings"
)

// A Context is the a request + response + route-arguments
type Context struct {
	Req  *http.Request
	Res  http.ResponseWriter
	Args []string
}

// A Route is a request middleware
type Route struct {
	method    string
	host      string
	path      *regexp.Regexp
	exclusive bool
	callback  func(*Context)
}

// Construct a new Route instance
func NewRoute() *Route {
	r := &Route{}
	r.SetPath("*").SetMethod("*").SetHost("*").SetFunc(func(c *Context) {
		http.NotFound(c.Res, c.Req)
	})
	return r
}

// set the "method" of the route
// `*` means any request method
func (this *Route) SetMethod(m string) *Route {
	this.method = strings.TrimSpace(strings.ToUpper(m))
	return this
}

// set the host of the router
// `*` means any host
func (this *Route) SetHost(v string) *Route {
	this.host = strings.TrimSpace(v)
	return this
}

// set the "path" of the route
// `*` means any path
func (this *Route) SetPath(p string) *Route {
	if p != "*" {
		p = regexp.MustCompilePOSIX(`/+`).ReplaceAllString((`/` + strings.TrimSpace(p) + `/`), `/`)
		this.path = regexp.MustCompilePOSIX(`^` + (p) + `$`)
	} else {
		this.path = regexp.MustCompilePOSIX(`(.*?)?`)
	}
	return this
}

// Just stop when the current request matches this route and don't run other similar routes ?
func (this *Route) SetExclusive(s bool) *Route {
	this.exclusive = s
	return this
}

// set the "func" of the route
func (this *Route) SetFunc(fn func(*Context)) *Route {
	this.callback = fn
	return this
}

// set the "http.Handler" of the route
func (this *Route) SetHandler(h http.Handler) *Route {
	return this.SetFunc(func(c *Context) { h.ServeHTTP(c.Res, c.Req) })
}

// Main middlewares container
type App struct {
	routes []*Route
}

// Create a new instance
func NewApp() *App {
	this := new(App)
	this.routes = []*Route{}
	return this
}

// Returns a new Route
func (this *App) Factory() *Route {
	r := NewRoute()
	this.routes = append(this.routes, r)
	return r
}

// Add previously created route(s)
func (this *App) Add(routes ...*Route) *App {
	this.routes = append(this.routes, routes...)
	return this
}

// Dispatch all registered routes
func (this *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method := strings.ToUpper(r.Method)
	host := strings.SplitN(r.Host, `:`, 2)[0]
	path := regexp.MustCompilePOSIX(`/+`).ReplaceAllString(`/`+r.URL.Path+`/`, `/`)
	for _, route := range this.routes {
		if route.path.MatchString(path) {
			if (route.method == `*`) || (regexp.MustCompilePOSIX(`^` + (route.method) + `$`).MatchString(method)) {
				if (route.host == `*`) || (regexp.MustCompilePOSIX(`^` + (route.host) + `$`).MatchString(host)) {
					args := []string{}
					if route.host != `*` {
						args = append(args, regexp.MustCompilePOSIX(`^`+(route.host)+`$`).FindAllStringSubmatch(host, -1)[0][1:]...)
					}
					args = append(args, route.path.FindAllStringSubmatch(path, -1)[0][1:]...)
					route.callback(&Context{r, w, args})
					if route.exclusive {
						break
					}
				}
			}
		}
	}
}

// serving HTTP traffic
func (this *App) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, this)
}

// serving HTTPS traffic
func (this *App) ListenAndServeTLS(addr string, certFile string, keyFile string) error {
	return http.ListenAndServeTLS(addr, certFile, keyFile, this)
}

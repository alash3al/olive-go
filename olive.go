// A tiny http framework perfect for building web services .
/**
	package main

	import "net/http"
	import "github.com/alash3al/olive-go"

	func main() {
		olive.New().GET("/", func(ctx *olive.Context) bool {
			ctx.Res.Write([]byte("index"))
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

*/
package olive

// we only need these package
import("log"; "regexp"; "strings"; "net/http")

// a route callback
type callback func(*Context)bool

// a group function
type gfn func(*App)

// a route property
type route struct {
	method, path string
	cb callback
}

// A request context and properties
type Context struct {
	Req		*http.Request
	Res		http.ResponseWriter
	Params 	[]string
	Vars	map[string]interface{}
}

// Our main structure
type App struct {
	routes 	[]route
	parent	string
}

// Create a new instance 
func New() *App {
	return &App{
		routes: []route{},
		parent: "/",
	}
}

// Group some of routes under the specified path/pattern  
// the group function will pass the current instance of App
// to the grouper .
func (self *App) Group(path string, fn gfn) *App {
	old := self.parent
	self.parent = regexp.MustCompile(`/+`).ReplaceAllString("/" + self.parent + "/" + strings.TrimSpace(path) + "/", "/")
	fn(self)
	self.parent = old
	return self
}

// Handle the specified custom method for the specified path  
// NOTE: method could be a "regexp string"  
// NOTE: path could be a "regexp string"  
func (self *App) METHOD(method, path string, cb callback) *App {
	self.routes = append(self.routes, route{
		method: strings.ToUpper(strings.TrimSpace(method)),
		path: regexp.MustCompile(`/+`).ReplaceAllString("/" + self.parent + "/" + strings.TrimSpace(path) + "/", "/"),
		cb: cb,
	})
	return self
}

// Handle ANY request method for the specified path with the specified callback;  
// NOTE: this is based on "func(self *App) METHOD()"
func (self *App) ANY(path string, cb callback) *App {
	return self.METHOD("(.*?)", path, cb)
}

// Handle GET request for the specified path with the specified callback;  
// NOTE: this is based on "func(self *App) METHOD()"  
func (self *App) GET(path string, cb callback) *App {
	return self.METHOD("GET", path, cb)
}

// Handle POST request for the specified path with the specified callback;  
// NOTE: this is based on "func(self *App) METHOD()"  
func (self *App) POST(path string, cb callback) *App {
	return self.METHOD("POST", path, cb)
}

// Handle PUT request for the specified path with the specified callback;  
// NOTE: this is based on "func(self *App) METHOD()"  
func (self *App) PUT(path string, cb callback) *App {
	return self.METHOD("PUT", path, cb)
}

// Handle PATCH request for the specified path with the specified callback;  
// NOTE: this is based on "func(self *App) METHOD()"  
func (self *App) PATCH(path string, cb callback) *App {
	return self.METHOD("PATCH", path, cb)
}

// Handle HEAD request for the specified path with the specified callback;  
// NOTE: this is based on "func(self *App) METHOD()"  
func (self *App) HEAD(path string, cb callback) *App {
	return self.METHOD("HEAD", path, cb)
}

// Handle DELETE request for the specified path with the specified callback;  
// NOTE: this is based on "func(self *App) METHOD()"  
func (self *App) DELETE(path string, cb callback) *App {
	return self.METHOD("DELETE", path, cb)
}

// Handle OPTIONS request for the specified path with the specified callback;  
// NOTE: this is based on "func(self *App) METHOD()"  
func (self *App) OPTIONS(path string, cb callback) *App {
	return self.METHOD("OPTIONS", path, cb)
}

// Handle TRACE request for the specified path with the specified callback;  
// NOTE: this is based on "func(self *App) METHOD()"  
func (self *App) TRACE(path string, cb callback) *App {
	return self.METHOD("TRACE", path, cb)
}

// Handle CONNECT request for the specified path with the specified callback;  
// NOTE: this is based on "func(self *App) METHOD()"  
func (self *App) CONNECT(path string, cb callback) *App {
	return self.METHOD("CONNECT", path, cb)
}

// Dispatch all registered routes
func (self App) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	defer (func(){
		if e := recover(); e != nil {
			log.Println(e)
		}
	})()
	ctx := &Context{Req: req, Res: res}
	current_path := regexp.MustCompile(`/+`).ReplaceAllString("/" + strings.TrimSpace(req.URL.Path) + "/", "/")
	current_method := req.Method
	rlen := len(self.routes)
	var _route route
	var re *regexp.Regexp
	var i int
	for i = 0; i < rlen; i ++ {
		_route = self.routes[i]
		if ! regexp.MustCompile("^" + _route.method + "$").MatchString(current_method) {
			continue
		}
		re = regexp.MustCompile("^" + _route.path + "$")
		if ! re.MatchString(current_path) {
			continue
		}
		ctx.Params = re.FindAllStringSubmatch(current_path, -1)[0][1:]
		if ! _route.cb(ctx) {
			break
		}
	}
}

// An Alias to http.ListenAndServe
func (self App) Listen(addr string) error {
	return http.ListenAndServe(addr, self)
}

// An Alias to http.ListenAndServeTLS
func (self App) ListenTLS(addr string, certFile string, keyFile string) error {
	return http.ListenAndServeTLS(addr, certFile, keyFile, self)
}

// Convert any http.Handler compatible handler to an internal callback
func Handler(h http.Handler, rtrn bool) callback {
	return func(ctx *Context) bool {
		h.ServeHTTP(ctx.Res, ctx.Req)
		return rtrn
	}
}

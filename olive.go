// Olive (v2.0) A tiny http framework perfect for building web services, created by Mohammed Al Ashaal (http://alash3al.xyz) under MIT License .
/**
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
			ctx.SetBody("i'm the parent \n")
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

*/
package olive

// we only need these package
import("io"; "regexp"; "strings"; "net/url"; "net/http"; "encoding/json"; "html/template")

// a route callback
type handler func(*Context)bool

// a group function
type grouper func(*App)

// a route property
type route struct {
	method, path, vhost string
	cb handler
}

// ---------------------

// A request context and properties
type Context struct {
	Req		*http.Request
	Res		http.ResponseWriter
	Params 	[]string
	Vars	map[string]interface{}
}

// Set the status code
func (self Context) SetStatus(code int) {
	self.Res.WriteHeader(code)
}

// Set a http header field
func (self Context) SetHeader(k, v string) {
	self.Res.Header().Set(k, v)
}

// Add a header field
func (self Context) AddHeader(k, v string) {
	self.Res.Header().Add(k, v)
}

// Get the specified header field "from the request"
func (self Context) GetHeader(k string) string {
	return self.Req.Header.Get(k)
}

// Delete the specified header field
func (self Context) DelHeader(k string) {
	self.Res.Header().Del(k)
}

// Return the url query vars
func (self Context) GetQuery() url.Values {
	return self.Req.URL.Query()
}

// Read the request body into the provided variable, it will read the body as json
// if the specified "v" isn't ([]byte, io.Writer) .
func (self Context) GetBody(v interface{}, max int64) (err error) {
	switch v.(type) {
		case []byte:
			_, err = http.MaxBytesReader(self.Res, self.Req.Body, max).Read(v.([]byte))
		case io.Writer:
			_, err = io.Copy(v.(io.Writer), http.MaxBytesReader(self.Res, self.Req.Body, max))
		default:
			err = json.NewDecoder(http.MaxBytesReader(self.Res, self.Req.Body, max)).Decode(v)
	}
	return err
}

// Write to the client, this function will detect the type of the data,
// it will send the data as json if the specified input isn't (string, []byte, *template.Template and io.Reader),
// it execute the input as html template if the first argument is *template.Template and second argument is template's data .
func (self Context) SetBody(d ... interface{}) {
	if len(d) == 0 {
		panic("Olive:Context Calling Send() without any arguments !")
	}
	switch d[0].(type) {
		case []byte:
			self.Res.Write(d[0].([]byte))
		case string:
			self.Res.Write([]byte(d[0].(string)))
		case io.Reader:
			io.Copy(self.Res, d[0].(io.Reader))
		case *template.Template:
			self.SetHeader("Content-Type", "text/html; charset=utf-8")
			if len(d) < 2 {
				err := (d[0].(*template.Template)).Execute(self.Res, nil)
				if err != nil {
					panic(err)
				}
			} else {
				err := (d[0].(*template.Template)).Execute(self.Res, d[1])
				if err != nil {
					panic(err)
				}
			}
		default:
			self.SetHeader("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(self.Res).Encode(d[0])
	}
}

// ---------------------

// Our main structure
type App struct {
	routes 	[]route
	parent	string
	vhost 	string
}

// Create a new instance 
func New() *App {
	return &App{
		routes: []route{},
		parent: "/",
		vhost: "(.*?)",
	}
}

// Group multiple routes by the specified section and its pattern with its grouper,
// olive currently supports by "host" or "path" only .
func (self *App) GroupBy(by, pattern string, fn grouper) *App {
	by = strings.ToLower(strings.TrimSpace(by))
	if by == "host" {
		old := self.vhost
		self.vhost = strings.TrimSpace(strings.ToLower(pattern))
		fn(self)
		self.vhost = old
	} else if by == "path" {
		old := self.parent
		self.parent = regexp.MustCompile(`/+`).ReplaceAllString("/" + self.parent + "/" + strings.TrimSpace(pattern) + "/", "/")
		fn(self)
		self.parent = old
	}
	return self
}

// Handle the specified custom method for the specified path;  
// NOTE: method could be a "regexp string";  
// NOTE: path could be a "regexp string";  
func (self *App) METHOD(method, path string, cb handler) *App {
	self.routes = append(self.routes, route{
		method: strings.ToUpper(strings.TrimSpace(method)),
		path: regexp.MustCompile(`/+`).ReplaceAllString("/" + self.parent + "/" + strings.TrimSpace(path) + "/", "/"),
		vhost: self.vhost,
		cb: cb,
	})
	return self
}

// Handle ANY request method for the specified path with the specified callback;  
// NOTE: this is based on "func(self *App) METHOD()"
func (self *App) ANY(path string, cb handler) *App {
	return self.METHOD("(.*?)", path, cb)
}

// Handle GET request for the specified path with the specified callback;  
// NOTE: this is based on "func(self *App) METHOD()"  
func (self *App) GET(path string, cb handler) *App {
	return self.METHOD("GET", path, cb)
}

// Handle POST request for the specified path with the specified callback;  
// NOTE: this is based on "func(self *App) METHOD()"  
func (self *App) POST(path string, cb handler) *App {
	return self.METHOD("POST", path, cb)
}

// Handle PUT request for the specified path with the specified callback;  
// NOTE: this is based on "func(self *App) METHOD()"  
func (self *App) PUT(path string, cb handler) *App {
	return self.METHOD("PUT", path, cb)
}

// Handle PATCH request for the specified path with the specified callback;  
// NOTE: this is based on "func(self *App) METHOD()"  
func (self *App) PATCH(path string, cb handler) *App {
	return self.METHOD("PATCH", path, cb)
}

// Handle HEAD request for the specified path with the specified callback;  
// NOTE: this is based on "func(self *App) METHOD()"  
func (self *App) HEAD(path string, cb handler) *App {
	return self.METHOD("HEAD", path, cb)
}

// Handle DELETE request for the specified path with the specified callback;  
// NOTE: this is based on "func(self *App) METHOD()"  
func (self *App) DELETE(path string, cb handler) *App {
	return self.METHOD("DELETE", path, cb)
}

// Handle OPTIONS request for the specified path with the specified callback;  
// NOTE: this is based on "func(self *App) METHOD()"  
func (self *App) OPTIONS(path string, cb handler) *App {
	return self.METHOD("OPTIONS", path, cb)
}

// Handle TRACE request for the specified path with the specified callback;  
// NOTE: this is based on "func(self *App) METHOD()"  
func (self *App) TRACE(path string, cb handler) *App {
	return self.METHOD("TRACE", path, cb)
}

// Handle CONNECT request for the specified path with the specified callback;  
// NOTE: this is based on "func(self *App) METHOD()"  
func (self *App) CONNECT(path string, cb handler) *App {
	return self.METHOD("CONNECT", path, cb)
}

// Dispatch all registered routes
func (self App) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ctx := &Context{Req: req, Res: res}
	current_path := regexp.MustCompile(`/+`).ReplaceAllString("/" + strings.TrimSpace(req.URL.Path) + "/", "/")
	current_method := req.Method
	current_host := strings.TrimSpace(strings.ToLower(strings.SplitN(req.Host, ":", 2)[0]))
	rlen := len(self.routes)
	var _route route
	var re_vhost, re_path *regexp.Regexp
	var i int
	for i = 0; i < rlen; i ++ {
		_route = self.routes[i]
		if ! regexp.MustCompile("^" + _route.method + "$").MatchString(current_method) {
			continue
		}
		re_vhost = regexp.MustCompile("^" + _route.vhost + "$")
		if ! re_vhost.MatchString(current_host) {
			continue
		}
		re_path = regexp.MustCompile("^" + _route.path + "$")
		if ! re_path.MatchString(current_path) {
			continue
		}
		ctx.Params = append(re_vhost.FindAllStringSubmatch(current_host, -1)[0][1:], re_path.FindAllStringSubmatch(current_path, -1)[0][1:] ...)
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

// Convert any http.Handler compatible handler to an internal handler,
// and rtrn is what will be returned for "[don't]run the next matched route with the same method and patttern"
func Handler(h http.Handler, rtrn bool) handler {
	return func(ctx *Context) bool {
		h.ServeHTTP(ctx.Res, ctx.Req)
		return rtrn
	}
}

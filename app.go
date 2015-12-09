/*
	Olive "Go", an advanced web-server framework written in Go .  

	```
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
*/
package olive

import (
	"net/http"
	"strings"
	"regexp"
)

// A App is a request middleware
type App struct {
	method    	string
	hostname    	string
	path      	string
	useragent 	string
	exclusive 	bool
	callback  	func(*Context)
	sub 		[]*App
}

// Construct a new App instance
func NewApp() *App {
	r := &App{}
	r.SetPath("*").SetMethod("*").SetHostname("*").SetUserAgent("*").SetExclusive(false).sub = []*App{}
	return r
}

// Add a sub-router, useful for nested routes of routes
// it returns the new sub route .
func (this *App) NewSubApp() *App {
	sub := NewApp()
	sub.SetPath(this.path).SetHostname(this.hostname).SetMethod(this.method).callback = func(c *Context) { sub.ServeHTTP(c.Res, c.Req) }
	this.sub = append(this.sub, sub)
	return sub
}

// Set the "method" of the route
// `*` means any request method
func (this *App) SetMethod(m string) *App {
	this.method = strings.TrimSpace(strings.ToUpper(m))
	return this
}

// Set the hostname of the router
// `*` means any hostname
func (this *App) SetHostname(v string) *App {
	this.hostname = strings.TrimSpace(v)
	return this
}

// Set the "path" of the router,
// `*` means any path .
func (this *App) SetPath(p string) *App {
	p = strings.TrimSpace(p)
	if p != "*" {
		p = regexp.MustCompile(`/+`).ReplaceAllString((`/` + p + `/`), `/`)
		this.path = regexp.MustCompile(`/+`).ReplaceAllString((`/` + strings.TrimSpace(p) + `/`), `/`)
	} else {
		this.path = p
	}
	return this
}

// Set the user-agent of this router,
// `*` means any path .
func (this *App) SetUserAgent(ua string) *App {
	this.useragent = strings.TrimSpace(ua)
	return this
}

// Just stop when the current request matches this route and don't run other similar routes ?
func (this *App) SetExclusive(s bool) *App {
	this.exclusive = s
	return this
}

// Add the child-handler-func of the route
func (this *App) HandleFunc(path string, fn func(*Context)) *App {
	sub := NewApp()
	parentPath := this.path
	if this.path == `*` {
		parentPath = `/`
	}
	sub.SetMethod(this.method).SetHostname(this.hostname).SetPath(regexp.MustCompile(`/+`).ReplaceAllString(parentPath + `/` + path, `/`)).callback = fn
	this.sub = append(this.sub, sub)
	return sub
}

// Add a sub-handler
func (this *App) Handle(path string, h http.Handler) *App {
	return this.HandleFunc(path, func(c *Context){ h.ServeHTTP(c.Res, c.Req) })
}

// Check whether this route matches the provided context or not 
func (this App) Match(ctx *Context) bool {
	var method = ctx.Req.Method
	var hostname = strings.SplitN(ctx.Req.Host, `:`, 2)[0]
	var path = regexp.MustCompile(`/+`).ReplaceAllString(`/`+ ctx.Req.URL.Path +`/`, `/`)
	var ua = ctx.Req.Header.Get(`User-Agent`)
	if this.useragent != "*" && ! regexp.MustCompile(`^(?i)`+ this.useragent +`$`).MatchString(ua) {
		return false
	}
	if this.method != "*" && ! regexp.MustCompile(`^(?i)`+ this.method +`$`).MatchString(method) {
		return false
	}
	if this.hostname != "*" && ! regexp.MustCompile(`^(?i)`+ this.hostname +`$`).MatchString(hostname) {
		return false
	}
	if this.path != "*" && ! regexp.MustCompile(`^`+ this.path +`$`).MatchString(path) {
		return false
	}
	return true
}

// Call this route callback .
func (this App) Call(ctx *Context) bool {
	var hostname = strings.SplitN(ctx.Req.Host, `:`, 2)[0]
	var path = regexp.MustCompile(`/+`).ReplaceAllString(`/`+ ctx.Req.URL.Path +`/`, `/`)
	ctx.Args.Hostname = regexp.MustCompile(`^(?i)`+ this.hostname +`$`).FindAllStringSubmatch(hostname, -1)[0][1:]
	ctx.Args.Path = regexp.MustCompile(`^`+ this.path +`$`).FindAllStringSubmatch(path, -1)[0][1:]
	this.callback(ctx)
	return this.exclusive
}

// Apply all childes .
func (this App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := new(Context)
	ctx.Res, ctx.Req = w, r
	for _, sub := range this.sub {
		if sub.Match(ctx) && sub.Call(ctx) {
			break
		}
	}
}

// serving HTTP traffic
func (this App) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, this)
}

// serving HTTPS traffic
func (this App) ListenAndServeTLS(addr string, certFile string, keyFile string) error {
	return http.ListenAndServeTLS(addr, certFile, keyFile, this)
}

// serving both http & https traffic
func (this App) ListenAndServeBoth(httpAddr, httpsAddr, certFile string, keyFile string) error {
	err := make(chan error, 1)
	go func(){
		err <- this.ListenAndServe(httpAddr)
	}()
	go func(){
		err <- this.ListenAndServeTLS(httpsAddr, certFile, keyFile)
	}()
	select {
		case e := <- err: {
			return e
		}
	}
}

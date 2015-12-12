/* A lightweight simple web application micro framework */
/*
```
	app := olive.NewApp()
	app.Hostname(`api.localtest.me`)
	app.HandleFunc(`/`, func(ctx *olive.Context){
		ctx.WriteJSON(map[string]interface{}{
			"status": true,
			"message": "working",
		})
	})
	app.Path(`/sub`, false)
	app.HandleFunc(`page`, func(ctx *olive.Context) { //> api.localtest.me/sub/page
		ctx.WriteString(`sub/page`)
	})
	app.Path(`/sub2`, true)
	app.HandleFunc(`page`, func(ctx *olive.Context) { //> api.localtest.me/sub/sub2/page
		ctx.WriteString(`sub/sub/page`)
	})
	app.Hostname(`www.localtest.me`)
	app.Handle(`?.*`, http.FileServer(http.Dir(`.`))) //> www.localtest.me/*
	app.ListenAndServe(`:80`)
```
*/
package olive

import (
	"net/http"
	"regexp"
)

// A App is a request middle-ware .
type App struct {
	routes 		[]*Route
	path		string
	hostname 	string
}

// Construct a new App instance .
func NewApp() *App {
	this := &App{}
	this.routes = []*Route{}
	this.path = `/`
	this.hostname = `*`
	return this
}

// Global hostname .
func (this *App) Hostname(hostname string) *App {
	this.hostname = hostname
	return this
}

// Global path
func (this *App) Path(path string, append bool) *App {
	if append {
		this.path += `/` + path
	} else {
		this.path = path
	}
	return this
}

// Add a custom handler .
func (this *App) HandleFunc(path string, fn func(*Context)) *Route {
	r := new(Route)
	r.path = regexp.MustCompile(`/+`).ReplaceAllString(`/` + this.path + `/` + path + `/`, `/`)
	r.callback = fn
	r.Hostname(this.hostname).Method(`*`).Exclusive(true)
	this.routes = append(this.routes, r)
	return r
}

// Add a http.Handler .
func (this *App) Handle(path string, h http.Handler) *Route {
	return this.HandleFunc(path, func(c *Context){ h.ServeHTTP(c.Res, c.Req) })
}

// Find route that matches the current context .
func (this App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := new(Context)
	ctx.Res, ctx.Req = w, r
	for _, route := range this.routes {
		if route.matches(ctx) && route.apply(ctx) {
			break
		}
	}
}

// serving HTTP traffic .
func (this App) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, this)
}

// serving HTTPS traffic .
func (this App) ListenAndServeTLS(addr string, certFile string, keyFile string) error {
	return http.ListenAndServeTLS(addr, certFile, keyFile, this)
}

// serving both http & https traffic .
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

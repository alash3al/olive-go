// Olive - a lightweight simple Golang web application framework .
// By Mohammed Al Ashaal "alash3al.xyz" .
// Version 1.0.0 .
package olive

// Our Requirements
import (
	"net/http"
	"strings"
	"regexp"
)

// A Context is the a request + response + <some properties>
type Context struct {
	Req		*http.Request
	Res 	http.ResponseWriter
	Args	[]string
}

// A Route is an olive handler with some properties
type Route struct {
	method		string
	vhost		string
	path		*regexp.Regexp
	callback	func(Context)
}

// set the "method" of the route
// `*` means any request method
func (this *Route) SetMethod(m string) *Route {
	this.method = strings.TrimSpace(strings.ToUpper(m))
	return this
}

// set the vhost "subdomain" of the router
// `*` means any host/vhost
func (this *Route) SetVhost(v string) *Route {
	this.vhost = strings.TrimSpace(v)
	return this
}

// set the "path" of the route
func (this *Route) SetPath(p string) *Route {
	p = regexp.MustCompilePOSIX(`/+`).ReplaceAllString((`/` + strings.TrimSpace(p) + `/`), `/`)
	this.path = regexp.MustCompilePOSIX(`^` + (p) + `$`)
	return this
}

// set the main "handler" of the route
func (this *Route) SetHandler(fn func(Context)) *Route {
	this.callback = fn
	return this
}

// Our main middleware
type Handler struct {
	routes []*Route
}

// Constructor
func NewHandler() *Handler {
	this := new(Handler)
	this.routes = []*Route{}
	return this
}

// Register a new Route
func (this *Handler) HandleFunc(fn func(Context)) *Route {
	r := &Route{}
	r.SetMethod(`*`).SetVhost(`*`).SetPath(`/`).SetHandler(fn)
	this.routes = append(this.routes, r)
	return r
}

// Dispatch all routes
func (this *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method := strings.ToUpper(r.Method)
	host := strings.SplitN(r.Host, `:`, 2)[0]
	path := regexp.MustCompilePOSIX(`/+`).ReplaceAllString(`/` + r.URL.Path + `/`, `/`)
	for i := 0; i < len(this.routes); i ++ {
		route := this.routes[i]
		if route.path.MatchString(path) {
			if (route.method == `*`) || (regexp.MustCompilePOSIX(`^` + (route.method) + `$`).MatchString(method)) {
				if (route.vhost == `*`) || (regexp.MustCompilePOSIX(`^` + (route.vhost) + `$`).MatchString(host)) {
					args := []string{}
					if route.vhost != `*` {
						args = append(args, regexp.MustCompilePOSIX(`^` + (route.vhost) + `$`).FindAllStringSubmatch(host, -1)[0][1:]...)
					}
					args = append(args, route.path.FindAllStringSubmatch(path, -1)[0][1:]...)
					route.callback(Context{r, w, args})
					break
				}
			}
		}
	}
}

// serving HTTP traffic
func (this *Handler) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, this)
}

// serving HTTPS traffic
func (this *Handler) ListenAndServeTLS(addr string, certFile string, keyFile string) error {
	return http.ListenAndServeTLS(addr, certFile, keyFile, this)
}

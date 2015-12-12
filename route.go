package olive

import (
	"regexp"
	"strings"
)

// A single route properties
type Route struct {
	method		string
	hostname	string
	path		string
	exclusive	bool
	callback	func(*Context)
}

// Set the hostname(s) of the route .
func (this *Route) Hostname(h ... string) *Route {
	this.hostname = strings.ToLower(strings.Join(h, `|`))
	return this
}

// Set the method(s) of the route
func (this *Route) Method(m ... string) *Route {
	this.method = strings.ToUpper(strings.Join(m, `|`))
	return this
}

// Whether to continue after this route or not .
func (this *Route) Exclusive(e bool) *Route {
	this.exclusive = e
	return this
}

func (this *Route) matches(ctx *Context) bool {
	var method = ctx.Req.Method
	var hostname = strings.SplitN(ctx.Req.Host, `:`, 2)[0]
	var path = regexp.MustCompile(`/+`).ReplaceAllString(`/`+ ctx.Req.URL.Path +`/`, `/`)
	if this.method != "*" && ! regexp.MustCompile(`^(?i)`+ this.method +`$`).MatchString(method) {
		return false
	}
	if this.hostname != "*" && ! regexp.MustCompile(`^(?i)`+ this.hostname +`$`).MatchString(hostname) {
		return false
	}
	if this.path != "/*/" && ! regexp.MustCompile(`^`+ this.path +`$`).MatchString(path) {
		return false
	}
	return true
}

func (this *Route) apply(ctx *Context) bool {
	var hostname = strings.SplitN(ctx.Req.Host, `:`, 2)[0]
	var path = regexp.MustCompile(`/+`).ReplaceAllString(`/`+ ctx.Req.URL.Path +`/`, `/`)
	ctx.Args.Hostname = regexp.MustCompile(`^(?i)`+ this.hostname +`$`).FindAllStringSubmatch(hostname, -1)[0][1:]
	ctx.Args.Path = regexp.MustCompile(`^`+ this.path +`$`).FindAllStringSubmatch(path, -1)[0][1:]
	this.callback(ctx)
	return this.exclusive
}

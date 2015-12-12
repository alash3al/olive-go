package olive

import (
	"io"
	"net/url"
	"net/http"
	"encoding/json"
)

// A Context is the a request + response + <some properties>
type Context struct {
	Req  *http.Request
	Res  http.ResponseWriter
	Args struct {
		Hostname	[]string
		Path		[]string
	}
}

// Construct a new context .
func NewContext(req *http.Request, res http.ResponseWriter) *Context {
	this := new(Context)
	this.Req = req
	this.Res = res
	this.Args = struct {
		Hostname	[]string
		Path		[]string
	}{}
	return this
}

// Set the status code .
func (this *Context) Status(code int) *Context {
	this.Res.WriteHeader(code)
	return this
}

// Set a header key-value .
func (this *Context) Set(k string, v string) *Context {
	this.Res.Header().Set(k, v)
	return this
}

// Append a header value to the specified key .
func (this *Context) Append(k string, v string) *Context {
	this.Res.Header().Add(k, v)
	return this
}

// Write a []byte array to the response body .
func (this *Context) Write(d []byte) *Context {
	this.Res.Write(d)
	return this
}

// Write stream to the client .
func (this *Context) WriteStream(stream io.Reader) *Context {
	io.Copy(this.Res, stream)
	return this
}

// Write a string to the response body .
func (this *Context) WriteString(d string) *Context {
	this.Res.Write([]byte(d))
	return this
}

// Write the specified interface{} as a json-data .
func (this *Context) WriteJSON(d interface{}) *Context {
	this.Res.Header().Set(`Content-Type`, `application/json; charset=UTF-8`)
	j, _ := json.MarshalIndent(d, ``, `   `)
	this.Res.Write(j)
	return this
}

// Read the request url query arguments .
func (this Context) ReadURLQuery() (url.Values, error) {
	return url.ParseQuery(this.Req.URL.RawQuery)
}

// Read maxlen of the request body as json into the specified interface{} .
func (this Context) ReadBodyJSON(o interface{}, maxlen int) error {
	defer this.Req.Body.Close()
	in := make([]byte, maxlen)
	n, e := this.Req.Body.Read(in)
	if n == 0 {
		return e
	}
	return json.Unmarshal(in[0:n], &o)
}

// Read maxlen of the request body as a url.Values .
func (this Context) ReadBodyQuery(maxlen int) (url.Values, error) {
	defer this.Req.Body.Close()
	in := make([]byte, maxlen)
	n, e := this.Req.Body.Read(in)
	if n == 0 {
		return nil, e
	}
	return url.ParseQuery(string(in[0:n]))
}

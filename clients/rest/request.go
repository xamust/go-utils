package rest

type Request interface {
	Method() string
	Body() any
	Address() string
}

type request struct {
	req    any
	method string
	url    string
}

func (r *request) Method() string {
	return r.method
}

func (r *request) Body() any {
	return r.req
}

func (r *request) Address() string {
	return r.url
}

func NewRequest(req any, method string, url string) Request {
	return &request{
		req:    req,
		method: method,
		url:    url,
	}
}

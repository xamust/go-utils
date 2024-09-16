package rest

import (
	"context"
	"crypto/tls"
	"net/http"
	"strconv"
	"time"

	"github.com/xamust/go-utils/encoder"
	"github.com/xamust/go-utils/encoder/json"
)

var (
	defaultTransport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	DefaultTimeout = 5 * time.Second
)

type Option func(*Options)

type Options struct {
	Transport   *http.Transport
	Timeout     time.Duration
	CallOptions CallOptions
}

func NewOptions(opts ...Option) Options {
	options := Options{
		Transport:   defaultTransport,
		Timeout:     DefaultTimeout,
		CallOptions: NewCallOption(),
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}

func WithTransport(t *http.Transport) Option {
	return func(options *Options) {
		options.Transport = t.Clone()
	}
}

func WithTimeout(t time.Duration) Option {
	return func(options *Options) {
		options.Timeout = t
	}
}

func WithCallOption(f CallOption) Option {
	return func(options *Options) {
		f(&options.CallOptions)
	}
}

type basicAuth struct {
	login string
	pass  string
}

// debugProxy to use check doc - https://confluence.berekebank.kz/display/SYN/debug+Proxy
type debugProxy struct {
	enableProxy bool
	url         string
}

func (d debugProxy) AddHeader(h http.Header, addr string) {
	h.Add("debugURL", addr)
	h.Add("isProxy", strconv.FormatBool(d.enableProxy))
}

type CallOption func(*CallOptions)

type CallOptions struct {
	Context            context.Context
	AuthToken          string
	Codec              encoder.Codec
	BasicAuth          *basicAuth
	Header             http.Header
	ContentType        string
	AllowFailedCheckCT bool
	EnableDebugProxy   *debugProxy
	AddResponseHeader  func(header http.Header) error
	ChangeHTTPDoer     HTTPDoer
}

func ChangeHTTPDoer(doer HTTPDoer) CallOption {
	return func(options *CallOptions) {
		options.ChangeHTTPDoer = doer
	}
}

// AddResponseHeader - example:
//
// h := make(http.Header)
//
//	f := func(copy http.Header) error {
//		h.Add("key", copy.Get("key"))
//	}
func AddResponseHeader(f func(header http.Header) error) CallOption {
	return func(options *CallOptions) {
		options.AddResponseHeader = f
	}
}

// EnableDebugProxy - url `http://debugproxy-synapse.apps.ocp-t.sberbank.kz/debug`
func EnableDebugProxy(url string, enableProxy bool) CallOption {
	return func(options *CallOptions) {
		options.EnableDebugProxy = &debugProxy{enableProxy, url}
	}
}

func BasicAuth(user, pass string) CallOption {
	return func(options *CallOptions) {
		options.BasicAuth = &basicAuth{user, pass}
	}
}

func WithAuth(auth string) CallOption {
	return func(options *CallOptions) {
		options.AuthToken = auth
	}
}

func CallCodec(c encoder.Codec) CallOption {
	return func(options *CallOptions) {
		options.Codec = c
	}
}

func WithContext(ctx context.Context) CallOption {
	return func(options *CallOptions) {
		options.Context = ctx
	}
}

func SetHeader(h http.Header) CallOption {
	return func(options *CallOptions) {
		options.Header = h
	}
}

func ContentType(ct string) CallOption {
	return func(options *CallOptions) {
		options.ContentType = ct
	}
}

// AllowFailContentType - if true, then don't check content type
func AllowFailContentType(f bool) CallOption {
	return func(option *CallOptions) {
		option.AllowFailedCheckCT = f
	}
}

func NewCallOption(opts ...CallOption) CallOptions {
	opt := CallOptions{
		Context:            context.Background(),
		Header:             make(http.Header),
		Codec:              json.NewCodec(),
		AllowFailedCheckCT: true,
	}

	for _, o := range opts {
		o(&opt)
	}
	return opt
}

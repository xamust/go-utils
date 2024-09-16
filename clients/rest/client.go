package rest

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/xamust/go-utils/errors"
	"github.com/xamust/go-utils/logger"
	"github.com/xamust/go-utils/metadata"
	"github.com/xamust/go-utils/metadata/request_id"

	"github.com/labstack/echo/v4"
)

const (
	MsgSentTo          = `The request was sent to: "%s"`
	MsgErrorSent       = `Error when sending a request to: "%s"`
	MsgErrorFrom       = `Error when receiving a response from: "%s"`
	MsgSuccessResponse = `The response was successfully received from "%s"`
)

type Client interface {
	Init(opts ...Option) error

	Call(ctx context.Context, req Request, rsp any, opts ...CallOption) error

	Timeout() time.Duration
}

type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type client struct {
	cli HTTPDoer

	opts Options
}

func (c *client) Init(opts ...Option) error {
	for _, o := range opts {
		o(&c.opts)
	}

	c.cli = &http.Client{
		Transport: c.opts.Transport,
		Timeout:   c.opts.Timeout,
	}

	return nil
}

func (c *client) Call(ctx context.Context, req Request, desc any, opts ...CallOption) error {
	log := logger.FromContextLogger(ctx)
	md, _ := metadata.FromContextHeader(ctx)

	callOpts := c.opts.CallOptions
	for _, opt := range opts {
		opt(&callOpts)
	}

	reqHttp, err := newHTTPRequest(ctx, req.Method(), req.Address(), req.Body(), callOpts)
	if err != nil {
		return err
	}

	doer := c.cli
	if callOpts.ChangeHTTPDoer != nil {
		doer = callOpts.ChangeHTTPDoer
	}

	rspHttp, err := doer.Do(reqHttp)
	metadata.SetTimeContextHeader(ctx)
	if err != nil {
		log.Error(ctx, "", logger.Operation(fmt.Sprintf(MsgErrorSent, md.Header.ReceiverSystem)))
		return err
	}
	defer func(c io.Closer) {
		_ = rspHttp.Body.Close()
	}(rspHttp.Body)

	if st, ok := desc.(StatusCode); ok {
		st.SetStatus(rspHttp.StatusCode)
	}

	ctx = logger.NewContextEvent(ctx, md.Header.ReceiverSystem, logger.AppSystem)
	if rspHttp.StatusCode >= http.StatusMultipleChoices {
		log.Error(ctx, "", logger.Operation(fmt.Sprintf(MsgErrorSent, md.Header.ReceiverSystem)))
		return errors.NewResponseClient(rspHttp.StatusCode, rspHttp.Body)
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return parseResponse(ctx, rspHttp, desc, callOpts)
	}
}

func newHTTPRequest(ctx context.Context, method, addr string, src any, opt CallOptions) (*http.Request, error) {
	log := logger.FromContextLogger(ctx)
	md, _ := metadata.FromContextHeader(ctx)

	bodyReq, err := opt.Codec.Marshal(src)
	if err != nil {
		return nil, err
	}

	ctx = logger.NewContextEvent(ctx, logger.AppSystem, md.Header.ReceiverSystem)

	header := make(http.Header)
	if opt.EnableDebugProxy != nil {
		opt.EnableDebugProxy.AddHeader(header, addr)
		addr = opt.EnableDebugProxy.url
	}
	if len(opt.Header) > 0 {
		for k, v := range opt.Header {
			for index := range v {
				header.Add(k, v[index])
			}
		}
	}

	reqUID, _ := request_id.FromContext(ctx)

	header.Set(echo.HeaderXRequestID, reqUID)

	ct := opt.ContentType
	if len(ct) == 0 {
		ct = opt.Codec.ContentType()
	}
	header.Set(echo.HeaderContentType, ct)

	if len(opt.AuthToken) > 0 {
		header.Set(echo.HeaderAuthorization, opt.AuthToken)
	}

	var reqHttp *http.Request
	if len(bodyReq) > 0 {
		log.Info(ctx, string(bodyReq), logger.Operation(fmt.Sprintf(MsgSentTo, md.Header.ReceiverSystem)), logger.HttpHeaders(header))
		reqHttp, err = http.NewRequestWithContext(ctx, method, addr, bytes.NewReader(bodyReq))
	} else {
		log.Info(ctx, addr, logger.Operation(fmt.Sprintf(MsgSentTo, md.Header.ReceiverSystem)), logger.HttpHeaders(header))
		reqHttp, err = http.NewRequestWithContext(ctx, method, addr, nil)
	}
	if err != nil {
		return nil, err
	}

	reqHttp.Header = header

	if ba := opt.BasicAuth; ba != nil {
		reqHttp.SetBasicAuth(ba.login, ba.pass)
	}

	return reqHttp, nil
}

func parseResponse(ctx context.Context, rspHttp *http.Response, desc any, callOpt CallOptions) (err error) {
	log := logger.FromContextLogger(ctx)
	md, _ := metadata.FromContextHeader(ctx)
	all, err := io.ReadAll(rspHttp.Body)
	if err != nil {
		return err
	}

	if len(all) == 0 {
		log.Info(ctx, "empty response", logger.Operation(fmt.Sprintf(MsgSuccessResponse, md.Header.ReceiverSystem)), logger.HttpHeaders(rspHttp.Header))
		return
	} else {
		if !callOpt.AllowFailedCheckCT {
			if !strings.Contains(rspHttp.Header.Get(echo.HeaderContentType), callOpt.Codec.String()) {
				log.Error(ctx, string(all), logger.Operation(fmt.Sprintf(MsgErrorSent, md.Header.ReceiverSystem)))
				return fmt.Errorf("error parse body: expected type: %s, have: %s", callOpt.Codec.String(), rspHttp.Header.Get(echo.HeaderContentType))
			}
		}

		log.Info(ctx, string(all), logger.Operation(fmt.Sprintf(MsgSuccessResponse, md.Header.ReceiverSystem)), logger.HttpHeaders(rspHttp.Header))

		if err = callOpt.Codec.Unmarshal(all, desc); err != nil {
			log.Error(ctx, "", logger.Operation(fmt.Sprintf(MsgErrorFrom, md.Header.ReceiverSystem)))
		}
	}

	if rspHttp.Header != nil && len(rspHttp.Header) > 0 {
		if callOpt.AddResponseHeader != nil {
			if err = callOpt.AddResponseHeader(rspHttp.Header); err != nil {
				log.Error(ctx, "callOpt.AddResponseHeader(rspHttp.Header)", logger.Operation(fmt.Sprintf(MsgErrorFrom, md.Header.ReceiverSystem)))
			}
		}
	}

	return
}

func (c *client) Timeout() time.Duration {
	return c.opts.Timeout
}

func NewClient(o ...Option) (Client, error) {
	c := &client{
		cli:  &http.Client{},
		opts: NewOptions(o...),
	}

	if err := c.Init(); err != nil {
		return nil, err
	}

	return c, nil
}

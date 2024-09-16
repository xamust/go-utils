package rest

import (
	"crypto/tls"
	"net/http"
	"net/url"

	"github.com/xamust/go-utils/server"
)

func NewWrapper(params server.EndpointParams, opts ...Option) (Client, error) {
	if params.Timeout.Duration <= 0 {
		params.Timeout.Duration = DefaultTimeout
	}
	optsWrap := append([]Option{}, WithTimeout(params.Timeout.Duration))

	if len(params.Proxy) != 0 {
		proxyUrl, err := url.Parse(params.Proxy)
		if err != nil {
			return nil, err
		}

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: params.Insecure},
			Proxy:           http.ProxyURL(proxyUrl),
		}
		optsWrap = append(optsWrap, WithTransport(tr))
	}

	if len(params.Login) != 0 {
		optsWrap = append(optsWrap, WithCallOption(BasicAuth(params.Login, params.Password)))
	}

	optsWrap = append(optsWrap, opts...)

	return NewClient(optsWrap...)
}

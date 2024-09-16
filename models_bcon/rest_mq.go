package models_bcon

import (
	"context"
	"encoding/xml"

	"github.com/xamust/go-utils/errors"
	"github.com/xamust/go-utils/metadata"
)

type ResponseMQ struct {
	ResponseCode int    `json:"responseCode"`
	Message      string `json:"message"`
}

func (r *ResponseMQ) Parsing(ctx context.Context, dst any) error {
	var rspMQ DataMQ
	if err := xml.Unmarshal([]byte(r.Message), &rspMQ); err != nil {
		return err
	}

	metadata.SetRsUidContextHeader(ctx, rspMQ.GetUID())

	if errRsp := rspMQ.GetErrorRsp(); errRsp != nil {
		return errRsp
	}

	if err := xml.Unmarshal(rspMQ.GetData(), dst); err != nil {
		return err
	}

	return nil
}

type DataMQ struct {
	MessageUID       string           `xml:"MessageUID"`
	SourceMessageUID string           `xml:"SourceMessageUID"`
	SystemCode       string           `xml:"SystemCode"`
	ServiceCode      string           `xml:"ServiceCode"`
	MessageDateTime  string           `xml:"MessageDateTime"`
	FilialCode       string           `xml:"FilialCode"`
	ErrorData        *errors.ErrorRsp `xml:"ErrorData,omitempty"`
	ResponseData     InnerData        `xml:"ResponseData,omitempty"`
}

type InnerData struct {
	Data []byte `xml:",innerxml"`
}

func (g *DataMQ) GetUID() string {
	return g.MessageUID
}

func (g *DataMQ) GetErrorRsp() *errors.ErrorRsp {
	return g.ErrorData
}

func (g *DataMQ) GetData() []byte {
	return g.ResponseData.Data
}

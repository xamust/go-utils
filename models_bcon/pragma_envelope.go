package models_bcon

import (
	"github.com/xamust/go-utils/metadata"
)

type PragmaEnvelope struct {
	MessageUID      string
	SystemCode      string
	ServiceCode     string
	MessageDateTime string
	FilialCode      string
	RequestData     any `xml:"RequestData"`
}

func (p *PragmaEnvelope) SetMeta(h metadata.Header) *PragmaEnvelope {
	p.MessageUID = h.RqUid
	p.SystemCode = h.SourceSystem
	p.ServiceCode = h.Service
	p.MessageDateTime = h.RqTm
	p.FilialCode = "99" //maybe other

	return p
}

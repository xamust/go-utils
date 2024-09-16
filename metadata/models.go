package metadata

type Header struct {
	Service        string `json:"service" example:"Название эндпоинта" validate:"required"`
	SourceSystem   string `json:"sourceSystem" example:"Название системы - отправителя" validate:"required"`
	RqUid          string `json:"rqUid" example:"UUID(уникальный идентификатор операции)" validate:"required"`
	OperUID        string `json:"operUid" example:"UUID(в основном такой же как RqUid)"`
	RqTm           string `json:"rqTm" example:"2023-10-31T13:24:52.192Z" validate:"required"`
	ReceiverSystem string `json:"receiverSystem" example:"Название системы - получателя"`

	//required: false
	Platform string `json:"platform,omitempty"`
	//required: false
	RsUid string `json:"rsUid,omitempty"`
	//required: false
	RsTm string `json:"rsTm,omitempty"`
}

type HeaderSource struct {
	Header Header `json:"header,dive"`
}

type HeaderReq struct {
	Service        string `json:"service" example:"Название эндпоинта" validate:"required"`
	SourceSystem   string `json:"sourceSystem" example:"Название системы - отправителя" validate:"required"`
	RqUid          string `json:"rqUid" example:"UUID(уникальный идентификатор операции)" validate:"required"`
	OperUID        string `json:"operUid" example:"UUID(в основном такой же как RqUid)" validate:"required"`
	RqTm           string `json:"rqTm" example:"2023-10-31T13:24:52.192Z" validate:"required"`
	ReceiverSystem string `json:"receiverSystem" example:"Название системы - получателя"`
}

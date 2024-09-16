package mq

/*
import (
	"context"
	"testing"

	"bitbucket.sberbank.kz/bcon/ibmmqgo"

	"github.com/stretchr/testify/assert"
)

func Test_ClientCall(t *testing.T) {
	t.Skip("")

	cfg := &Config{
		Addr: "rest2mq-synapse.apps.ocp-t.sberbank.kz:9090",
		Params: &ibmmqgo.Config{
			QueueManagerChannel: "ESB.GW.SVRCONN",
			QueueManagerName:    "MKZ.ESB.DI1",
			QueueManagerHost:    "172.16.99.66",
			QueueManagerPort:    1415,
			Timeout:             120,
		},
	}

	cli := NewClient(cfg)

	err := cli.Init(context.Background())
	assert.Nil(t, err)
	defer cli.Close()

	vec := ibmmqgo.Vector{
		RequestQueue:  "BCON.PRAGMANEW.RO.REQUEST",
		ResponseQueue: "BCON.PRAGMANEW.RO.RESPONSE",
		ServiceName:   "FILE DBZ",
	}

	msg := `
<PragmaEnvelope>
<MessageUID>065b8633-a015-488b-9fe2-e914b1e9e9e3</MessageUID>
<SystemCode>FICO</SystemCode>
<ServiceCode>FICOAcctPosition</ServiceCode>
<MessageDateTime>2024-01-17T11:27:37.099Z</MessageDateTime>
<FilialCode>99</FilialCode>
<RequestData>
	<KZAcctPositionRq>
		<Acct>
			<FilialCode>03</FilialCode>
			<AcctNum>KZ26914032204KZ02TK0</AcctNum>
			<Currency>KZT</Currency>
		</Acct>
		<Period>
			<BegDate>2023-06-15T00:00:00+06:00</BegDate>
			<EndDate>2023-08-18T00:00:00+06:00</EndDate>
		</Period>
	</KZAcctPositionRq>
</RequestData>
</PragmaEnvelope>
`

	rsp, err := cli.Call(context.TODO(), msg, vec)
	assert.Nil(t, err)
	assert.NotNil(t, rsp)
}

*/

package mq

/*
import (
	pb "github.com/xamust/go-utils/clients/grpc/mq/generate/proto"
	"github.com/xamust/go-utils/logger"
	"github.com/xamust/go-utils/metadata"
	"bitbucket.sberbank.kz/bcon/ibmmqgo"
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type Client interface {
	Init(ctx context.Context) error

	Call(ctx context.Context, msg string, vector ibmmqgo.Vector) (*pb.BasicResponse, error)

	Close() error
}

type client struct {
	params *ibmmqgo.Config
	addr   string

	conn *grpc.ClientConn
}

func NewClient(cfg *Config) Client {
	if cfg == nil {
		logger.DefaultLogger.Fatal(context.Background(), "empty config")
	}

	return &client{
		addr:   cfg.Addr,
		params: cfg.Params,
	}
}

func (c *client) Init(ctx context.Context) (err error) {
	//ctx, _ = context.WithTimeout(ctx, 120*time.Second)
	if c.conn, err = grpc.Dial(c.addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithTimeout(120*time.Second)); err != nil {
		return err
	}

	return err
}

func (c *client) Call(ctx context.Context, msg string, vector ibmmqgo.Vector) (*pb.BasicResponse, error) {
	client := pb.NewMqGrpcServiceClient(c.conn)

	reqUID := uuid.New().String()
	if meta, ok := metadata.FromContextHeader(ctx); ok && len(meta.Header.RqUid) > 0 {
		reqUID = meta.Header.RqUid
	}

	req := &pb.BasicRequest{
		QueueManagerHost:    c.params.QueueManagerHost,
		QueueManagerPort:    uint32(c.params.QueueManagerPort),
		QueueManagerChannel: c.params.QueueManagerChannel,
		QueueManagerName:    c.params.QueueManagerName,
		Timeout:             uint32(c.params.Timeout),
		RequestQueue:        vector.RequestQueue,
		ResponseQueue:       vector.ResponseQueue,
		ServiceName:         vector.ServiceName,
		RqUid:               reqUID,
		Message:             msg,
	}

	return client.SendAndReceive(ctx, req)
}

func (c *client) Close() error {
	return c.conn.Close()
}

*/

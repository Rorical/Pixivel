package levelgo

import (
	"context"
	"errors"
	"log"

	levelrpc "Pixivel/internal/database/levelgo/levelrpc"

	grpc "google.golang.org/grpc"
)

type LevelRpcClient struct {
	Host   string
	conn   *grpc.ClientConn
	Client *levelrpc.LevelRpcServiceClient
}

func RpcClient(host string) *LevelRpcClient {
	return &LevelRpcClient{
		Host: host,
	}
}

//Connect the server
func (self *LevelRpcClient) Connect() {
	conn, err := grpc.Dial(self.Host, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Connect Fail: %v", err)
	}
	self.conn = conn
	client := levelrpc.NewLevelRpcServiceClient(conn)
	self.Client = &client
}

func (self *LevelRpcClient) Get(key []byte) ([]byte, error) {
	reply, err := (*self.Client).Get(context.Background(), &levelrpc.GetRequest{Key: key})
	if err != nil {
		log.Fatalf("%v", err)
	}
	return reply.Value, ErrCodeToError(reply.Error)

}

func (self *LevelRpcClient) Set(key []byte, value []byte) error {
	reply, err := (*self.Client).Set(context.Background(), &levelrpc.SetRequest{Key: key, Value: value})
	if err != nil {
		log.Fatalf("%v", err)
	}
	return ErrCodeToError(reply.Error)
}

func (self *LevelRpcClient) Has(key []byte) (bool, error) {
	reply, err := (*self.Client).Has(context.Background(), &levelrpc.GetRequest{Key: key})
	if err != nil {
		log.Fatalf("%v", err)
	}
	return reply.Value, ErrCodeToError(reply.Error)
}

func (self *LevelRpcClient) Del(key []byte) error {
	reply, err := (*self.Client).Del(context.Background(), &levelrpc.GetRequest{Key: key})
	if err != nil {
		log.Fatalf("%v", err)
	}
	return ErrCodeToError(reply.Error)
}

func (self *LevelRpcClient) IsErrNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

func (self *LevelRpcClient) Close() {
	self.conn.Close()
}

package databasePG

import (
	"emperror.dev/errors"
	pb "github.com/je4/mediaserver/v2/pkg/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClientPlain(addr string) (*Client, error) {
	var opts = []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot dial '%s'", addr)
	}
	grpcClient := pb.NewDatabaseClient(conn)
	dbClient := &Client{
		DatabaseClient: grpcClient,
		conn:           conn,
	}
	return dbClient, nil
}

type Client struct {
	conn *grpc.ClientConn
	pb.DatabaseClient
}

func (dbClient *Client) Close() error {
	return dbClient.conn.Close()
}

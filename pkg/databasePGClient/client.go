package databasePGClient

import (
	"emperror.dev/errors"
	pb "github.com/je4/mediaserver/v2/pkg/protos"
	"google.golang.org/grpc"
)

func NewDatabaseClient(addr string) (*databaseClient, error) {
	var opts []grpc.DialOption
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot dial '%s'", addr)
	}
	grpcClient := pb.NewDatabaseClient(conn)
	dbClient := &databaseClient{
		DatabaseClient: grpcClient,
		conn:           conn,
	}
	return dbClient, nil
}

type databaseClient struct {
	conn *grpc.ClientConn
	pb.DatabaseClient
}

func (dbClient *databaseClient) Close() error {
	return dbClient.conn.Close()
}

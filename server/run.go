package server

import (
	"context"
	"net"

	"google.golang.org/grpc"

	raftdv1 "github.com/amjadjibon/raftd/gen/raftd/v1"
)

func Run(
	ctx context.Context,
	raftDir string,
	raftBind string,
	raftNodeID string,
	grpcAddr string,
) error {
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	raftd, err := NewRaftd(raftDir, raftBind, raftNodeID)
	if err != nil {
		return err
	}

	raftdv1.RegisterRaftServiceServer(grpcServer, raftd)
	raftdv1.RegisterKVServiceServer(grpcServer, raftd)

	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
	}()

	if err = grpcServer.Serve(lis); err != nil {
		return err
	}

	return nil
}

package server

import (
	"context"
	"encoding/json"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	raftdv1 "github.com/amjadjibon/raftd/gen/raftd/v1"
	"github.com/amjadjibon/raftd/store"
)

const (
	transportMaxPool    = 5
	transportTimeout    = 10 * time.Second
	snapshotRetainCount = 3
)

type Raftd struct {
	store      *store.Store
	fsm        *FSM
	raftEngine *raft.Raft
	raftBoltDB *raftboltdb.BoltStore
}

var _ raftdv1.RaftServiceServer = (*Raftd)(nil)
var _ raftdv1.KVServiceServer = (*Raftd)(nil)

func NewRaftd(
	raftDir string,
	raftBind string,
	raftNodeID string,
) (*Raftd, error) {
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(raftNodeID)

	addr, err := net.ResolveTCPAddr("tcp", raftBind)
	if err != nil {
		return nil, err
	}

	transport, err := raft.NewTCPTransport(
		raftBind,
		addr,
		transportMaxPool,
		transportTimeout,
		os.Stderr,
	)
	if err != nil {
		return nil, err
	}

	snapshotStore, err := raft.NewFileSnapshotStore(
		raftDir,
		snapshotRetainCount,
		os.Stderr,
	)
	if err != nil {
		return nil, err
	}

	boltStore, err := raftboltdb.NewBoltStore(filepath.Join(raftDir, "raft.db"))
	if err != nil {
		return nil, err
	}

	fsm := NewFSM(store.New())

	raftEngine, err := raft.NewRaft(
		config,
		fsm,
		boltStore,
		boltStore,
		snapshotStore,
		transport,
	)
	if err != nil {
		return nil, err
	}

	configuration := raft.Configuration{
		Servers: []raft.Server{
			{
				ID:      raft.ServerID(raftNodeID),
				Address: transport.LocalAddr(),
			},
		},
	}

	raftEngine.BootstrapCluster(configuration)

	return &Raftd{
		store:      store.New(),
		fsm:        fsm,
		raftEngine: raftEngine,
		raftBoltDB: boltStore,
	}, nil
}

// Join implements raftdv1.RaftServiceServer.
func (s *Raftd) Join(ctx context.Context, req *raftdv1.JoinRequest) (*raftdv1.JoinResponse, error) {
	config := s.raftEngine.GetConfiguration()
	if config.Error() != nil {
		return nil, status.Errorf(codes.Internal, "failed to get configuration: %v", config.Error())
	}

	for _, server := range config.Configuration().Servers {
		if server.ID == raft.ServerID(req.Id) || server.Address == raft.ServerAddress(req.Address) {
			return nil, status.Errorf(codes.AlreadyExists, "server already joined")
		}

		future := s.raftEngine.AddVoter(raft.ServerID(req.Id), raft.ServerAddress(req.Address), 0, time.Second)
		if err := future.Error(); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to add voter: %v", err)
		}
	}

	future := s.raftEngine.AddVoter(raft.ServerID(req.Id), raft.ServerAddress(req.Address), 0, time.Second)
	if err := future.Error(); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add voter: %v", err)
	}

	return &raftdv1.JoinResponse{}, nil
}

// Leave implements raftdv1.RaftServiceServer.
func (s *Raftd) Leave(context.Context, *raftdv1.LeaveRequest) (*raftdv1.LeaveResponse, error) {
	panic("unimplemented")
}

// Status implements raftdv1.RaftServiceServer.
func (s *Raftd) Status(context.Context, *raftdv1.StatusRequest) (*raftdv1.StatusResponse, error) {
	peers := make([]*raftdv1.Peer, 0, len(s.raftEngine.GetConfiguration().Configuration().Servers))
	for _, peer := range s.raftEngine.GetConfiguration().Configuration().Servers {
		peers = append(peers, &raftdv1.Peer{Id: string(peer.ID)})
	}
	return &raftdv1.StatusResponse{
		Leader: string(s.raftEngine.Leader()),
		Peers:  peers,
	}, nil
}

// Set implements raftdv1.KVServiceServer.
func (s *Raftd) Set(ctx context.Context, req *raftdv1.SetRequest) (*raftdv1.SetResponse, error) {
	if s.raftEngine.State() != raft.Leader {
		return nil, status.Errorf(codes.FailedPrecondition, "not the leader")
	}

	cmd := &raftdv1.Command{
		Op:    "set",
		Key:   req.Key,
		Value: req.Value,
	}

	data, err := json.Marshal(cmd)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to marshal command: %v", err)
	}

	if err := s.raftEngine.Apply(data, time.Second); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to apply command: %v", err)
	}

	return &raftdv1.SetResponse{}, nil
}

// Get implements raftdv1.KVServiceServer.
func (s *Raftd) Get(ctx context.Context, req *raftdv1.GetRequest) (*raftdv1.GetResponse, error) {
	value, err := s.store.Get(req.Key)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "key not found: %v", err)
	}
	return &raftdv1.GetResponse{Value: value}, nil
}

// Delete implements raftdv1.KVServiceServer.
func (s *Raftd) Delete(ctx context.Context, req *raftdv1.DeleteRequest) (*raftdv1.DeleteResponse, error) {
	if s.raftEngine.State() != raft.Leader {
		return nil, status.Errorf(codes.FailedPrecondition, "not the leader")
	}

	cmd := &raftdv1.Command{
		Op:  "del",
		Key: req.Key,
	}

	data, err := json.Marshal(cmd)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to marshal command: %v", err)
	}

	resp := s.raftEngine.Apply(data, time.Second)
	if resp.Error() != nil {
		return nil, status.Errorf(codes.Internal, "failed to apply command: %v", resp.Error())
	}

	return &raftdv1.DeleteResponse{}, nil
}

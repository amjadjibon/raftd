package server

import (
	"encoding/json"
	"io"

	"github.com/hashicorp/raft"

	raftdv1 "github.com/amjadjibon/raftd/gen/raftd/v1"
	"github.com/amjadjibon/raftd/store"
)

type FSM struct {
	store *store.Store
}

var _ raft.FSM = (*FSM)(nil)

func NewFSM(store *store.Store) *FSM {
	return &FSM{store: store}
}

// Apply implements raft.FSM.
func (f *FSM) Apply(raftLog *raft.Log) interface{} {
	switch raftLog.Type {
	case raft.LogCommand:
		var c raftdv1.Command
		if err := json.Unmarshal(raftLog.Data, &c); err != nil {
			return err
		}

		switch c.Op {
		case "set":
			return f.store.Set(c.Key, c.Value)
		case "del":
			return f.store.Delete(c.Key)
		}
	}
	return nil
}

// Restore implements raft.FSM.
func (f *FSM) Restore(snapshot io.ReadCloser) error {
	defer func() {
		_ = snapshot.Close()
	}()

	dec := json.NewDecoder(snapshot)
	for dec.More() {
		var c raftdv1.Command
		if err := dec.Decode(&c); err != nil {
			return err
		}

		switch c.Op {
		case "set":
			if err := f.store.Set(c.Key, c.Value); err != nil {
				return err
			}
		case "del":
			if err := f.store.Delete(c.Key); err != nil {
				return err
			}
		}
	}

	return nil
}

// Snapshot implements raft.FSM.
func (f *FSM) Snapshot() (raft.FSMSnapshot, error) {
	return newSnapshot(f.store), nil
}

type snapshot struct {
	store *store.Store
}

func (s snapshot) Persist(sink raft.SnapshotSink) error {
	err := func() error {
		b, err := json.Marshal(s.store)
		if err != nil {
			return err
		}

		if _, err := sink.Write(b); err != nil {
			return err
		}

		return sink.Close()
	}()

	if err != nil {
		sink.Cancel()
	}

	return err
}

func (s snapshot) Release() {}

func newSnapshot(store *store.Store) raft.FSMSnapshot {
	return &snapshot{store: store}
}

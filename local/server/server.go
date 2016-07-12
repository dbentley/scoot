package server

import (
	"github.com/scootdev/scoot/local/protocol"
	"github.com/scootdev/scoot/runner"
	"github.com/scootdev/scoot/snapshots/filer"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

// Create a protocol.LocalScootServer
func NewServer(runner runner.Runner, snapshotter filer.Snapshotter) (protocol.LocalScootServer, error) {
	return &Server{runner, snapshotter}, nil
}

type Server struct {
	runner      runner.Runner
	snapshotter filer.Snapshotter
}

// TODO(dbentley): how to cancel

// Serve  serves the Scoot instance in scootdir with logic handler s.
func Serve(s protocol.LocalScootServer, scootdir string) error {
	socketPath := protocol.SocketForDir(scootdir)
	l, err := net.Listen("unix", socketPath)
	if err != nil {
		return err
	}
	defer l.Close()
	server := grpc.NewServer()
	protocol.RegisterLocalScootServer(server, s)
	server.Serve(l)
	return nil
}

func (s *Server) Echo(ctx context.Context, req *protocol.EchoRequest) (*protocol.EchoReply, error) {
	return &protocol.EchoReply{Pong: "Pong: " + req.Ping}, nil
}

func (s *Server) Run(ctx context.Context, req *protocol.Command) (*protocol.ProcessStatus, error) {
	cmd := runner.NewCommand(req.Argv, req.Env, time.Duration(req.Timeout))
	status, err := s.runner.Run(cmd)
	if err != nil {
		return nil, err
	}
	return protocol.FromRunnerStatus(status), nil
}

func (s *Server) Status(ctx context.Context, req *protocol.StatusQuery) (*protocol.ProcessStatus, error) {
	status, err := s.runner.Status(runner.RunId(req.RunId))
	if err != nil {
		return nil, err
	}

	return protocol.FromRunnerStatus(status), nil
}

func (s *Server) SnapshotCreate(ctx context.Context, req *protocol.SnapshotCreateReq) (*protocol.SnapshotCreateResp, error) {
	log.Println("server.SnapshotCreate")
	id, err := s.snapshotter.Snapshot(req.FromDir)
	log.Printf("server.SnapshotCreate done %q %v", id, err)
	if err != nil {
		return nil, err
	}

	r := &protocol.SnapshotCreateResp{}
	r.Id = id
	return r, nil
}

package server

import (
	"github.com/scootdev/scoot/daemon/protocol"
	"github.com/scootdev/scoot/runner"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"time"
)

// Create a protocol.ScootDaemonServer
func NewServer(runner runner.Runner) (protocol.ScootDaemonServer, error) {
	return &Server{runner}, nil
}

type Server struct {
	runner runner.Runner
}

// TODO(dbentley): how to cancel
// Serve  serves the Scoot instance in scootdir with logic handler s.
func Serve(s protocol.ScootDaemonServer, scootdir string) error {
	socketPath := protocol.SocketForDir(scootdir)
	l, err := listen(socketPath)
	if err != nil {
		return err
	}
	defer l.Close()
	server := grpc.NewServer()
	protocol.RegisterScootDaemonServer(server, s)
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

package main

import (
	"flag"
	"github.com/scootdev/scoot/local/protocol"
	"github.com/scootdev/scoot/local/server"
	"github.com/scootdev/scoot/runner/execer"
	"github.com/scootdev/scoot/runner/execer/fake"
	"github.com/scootdev/scoot/runner/execer/os"
	"github.com/scootdev/scoot/runner/local"
	"github.com/scootdev/scoot/snapshots/git"
	"log"
)

var execerType = flag.String("execer_type", "sim", "execer type; os or sim")

// A Local Scoot server.
func main() {
	flag.Parse()
	scootdir, err := protocol.LocateScootDir()
	if err != nil {
		log.Fatal("Error locating Scoot instance: ", err)
	}
	var ex execer.Execer
	switch *execerType {
	case "sim":
		ex = fake.NewSimExecer(nil)
	case "os":
		ex = os.NewExecer()
	default:
		log.Fatalf("Unknown execer type %v", *execerType)
	}
	runner := local.NewSimpleRunner(ex)

	gitRunner := git.NewExecRunner()
	snapshots, err := git.NewSnapshots(scootDir, gitRunner)
	if err != nil {
		log.Fatal("Cannot create Git Snapshots: ", err)
	}

	snapshotter, err := git.NewSnapshotter(snapshots, scootDir, gitRunner)
	if err != nil {
		log.Fatal("Cannot create Git Snapshotter: ", err)
	}

	server, err := server.NewServer(runner, snapshotter)
	if err != nil {
		log.Fatal("Cannot create Scoot server: ", err)
	}
	err = server.Serve(s, scootDir)
	if err != nil {
		log.Fatal("Error serving Local Scoot: ", err)
	}
}

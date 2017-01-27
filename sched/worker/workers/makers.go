package workers

import (
	"time"

	"github.com/scootdev/scoot/cloud/cluster"
	"github.com/scootdev/scoot/os/temp"
	"github.com/scootdev/scoot/runner/execer/execers"
	"github.com/scootdev/scoot/runner/runners"
	"github.com/scootdev/scoot/sched/worker"
	"github.com/scootdev/scoot/snapshot"
)

// Makes a worker suitable for using as an in-memory worker.
func MakeInmemoryWorker(node cluster.Node, tmp *temp.TempDir, filer snapshot.Filer) worker.Worker {
	ex := execers.NewDoneExecer()
	r := runners.NewSingleRunner(ex, filer, runners.NewNullOutputCreator(), tmp)
	chaos := runners.NewChaosRunner(r)
	chaos.SetDelay(time.Duration(500) * time.Millisecond)
	return NewServiceWorker(chaos, 0, time.Second)
}

// Makes a worker that uses a SimExecer. This is suitable for testing.
func MakeSimWorker(tmp *temp.TempDir, filer snapshot.Filer) worker.Worker {
	ex := execers.NewSimExecer()
	r := runners.NewSingleRunner(ex, filer, runners.NewNullOutputCreator(), tmp)
	return NewServiceWorker(r, 0, time.Second)
}

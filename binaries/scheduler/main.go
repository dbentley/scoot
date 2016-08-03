package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/apache/thrift/lib/go/thrift"
	clusterdef "github.com/scootdev/scoot/cloud/cluster"
	clusterimpl "github.com/scootdev/scoot/cloud/cluster/memory"
	"github.com/scootdev/scoot/saga"
	queueimpl "github.com/scootdev/scoot/sched/queue/memory"
	"github.com/scootdev/scoot/sched/scheduler"
	"github.com/scootdev/scoot/sched/worker"
	//"github.com/scootdev/scoot/sched/worker"

	"github.com/scootdev/scoot/sched/worker/fake"
	"github.com/scootdev/scoot/sched/worker/rpc"
	// "github.com/scootdev/scoot/sched/worker/rpc"
	"github.com/scootdev/scoot/scootapi/server"
)

//TODO: we'll want more flexibility with startup configuration, maybe something like:
//  {"cluster": {"type": "static"},
//   "initial_nodes": ["localhost:2345", "localhost:2346"],
//   "workerFactory": "thrift"}
//
var addr = flag.String("addr", "localhost:9090", "Bind address for api server.")
var workers = flag.String("workers", "", "Comma separated list of workers (host:port,...)|NUM:mem.")

func main() {
	log.Println("Starting Cloud Scoot API Server & Scheduler")

	flag.Parse()
	workersList := strings.Split(*workers, ",")
	workerNodes := []clusterdef.Node{}
	if len(workersList) == 1 && (workersList[0] == "" || strings.Contains(workersList[0], ":mem")) {
		//Keep the original behavior for now if no workers specified on cmdline.
		numNodes := 10
		if workersList[0] != "" {
			numNodes, _ = strconv.Atoi(strings.Split(workersList[0], ":")[0])
		}
		workersList = []string{}
		for idx := 0; idx < numNodes; idx++ {
			workersList = append(workersList, fmt.Sprintf("inmemory%d", idx))
		}
	}
	for _, worker := range workersList {
		//TODO: methods to set/get an actual addr from cluster.Node, using Node.Id for now.
		workerNodes = append(workerNodes, clusterimpl.NewIdNode(worker))
	}
	inmemory := strings.Contains(string(workerNodes[0].Id()), "inmem")

	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	transportFactory := thrift.NewTTransportFactory()

	// Create Cluster
	// TODO: replace with actual cluster implementation, currently dummy in memory cluster
	cluster := clusterimpl.NewCluster(workerNodes, nil)

	// Create Saga Log
	// TODO: Replace with Durable SagaLog, currently In Memory Only
	sagaLog := saga.MakeInMemorySagaLog()

	workerFactory := fake.MakeWaitingNoopWorker
	if !inmemory {
		workerFactory = func(node clusterdef.Node) worker.Worker {
			return rpc.NewThriftWorker(transportFactory, protocolFactory, string(node.Id()))
		}
	}
	// TODO: Replace with Durable WorkQueue, currently in Memory Only
	workQueue := queueimpl.NewSimpleQueue(1000)
	s := scheduler.NewSchedulerFromCluster(cluster, workQueue.Chan(), sagaLog, workerFactory)
	handler := server.NewHandler(workQueue, s.Sagas())

	log.Println("Starting API Server")
	err := server.Serve(handler, *addr, transportFactory, protocolFactory)
	if err != nil {
		log.Fatal("Error serving Scoot API: ", err)
	}
}

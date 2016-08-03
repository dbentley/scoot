package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/scootdev/scoot/common/endpoints"
	"github.com/scootdev/scoot/common/stats"
	execer "github.com/scootdev/scoot/runner/execer/os"
	localrunner "github.com/scootdev/scoot/runner/local"
	runnerworker "github.com/scootdev/scoot/workerapi/runner"
	"github.com/scootdev/scoot/workerapi/server"
)

var thriftPort = flag.Int("thrift_port", 9090, "port to serve thrift on")
var httpPort = flag.Int("http_port", 9091, "port to serve http on")

func main() {
	flag.Parse()
	stat, _ := stats.NewCustomStatsReceiver(stats.NewFinagleStatsRegistry, 15*time.Second)
	stat = stat.Precision(time.Millisecond)
	endpoints.RegisterStats("/admin/metrics.json", stat)
	endpoints.RegisterHealthCheck("/")
	go endpoints.Serve(fmt.Sprintf(":%d", *httpPort))

	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	transportFactory := thrift.NewTTransportFactory()

	stats := stat.Scope("workerserver")
	r := localrunner.NewSimpleRunner(execer.NewExecer())
	version := func() string { return "" }
	w := runnerworker.MakeWorker(r)
	handler := server.NewHandler(stats, w, version)
	err := server.Serve(handler, fmt.Sprintf(":%d", *thriftPort), transportFactory, protocolFactory)
	if err != nil {
		log.Fatal("Error serving Worker Server: ", err)
	}
}

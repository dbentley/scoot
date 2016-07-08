package scootapi

// We should use go generate to run:
//   thrift --gen go:package_prefix=github.com/scootdev/scoot/cloud/api/gen-go/,thrift_import=github.com/apache/thrift/lib/go/thrift scoot.thrift
// Right now we don't because thrift is hard to install programmatically.

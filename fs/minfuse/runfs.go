package minfuse

import (
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	_ "net/http/pprof"

	"github.com/scootdev/scoot/fs/min"
	"github.com/scootdev/scoot/fuse"
	"github.com/scootdev/scoot/snapshot"
)

type Options struct {
	Src          string
	Mountpoint   string
	StrategyList []string
	Async        bool
	Trace        bool
	ThreadUnsafe bool
	MaxReadahead uint32
}

func SetupLog() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}

func InitFlags() (*Options, error) {
	// Note on serial vs threadpool option:
	// Serial means that no locking will be done at all and we will serve from a single goroutine.
	// Threadpool means that we will lock where necessary and serve from multiple goroutines.
	src := flag.String("src_root", "", "source directory to mirror")
	mountpoint := flag.String("mountpoint", "", "directory to mount at")
	trace := flag.Bool("trace", false, "whether to trace execution")
	serveStrategy := flag.String("serve_strategy", "",
		"Options are any of async|sync;serial|threadpool;readahead_mb=N. Default is 'async;serial;readahead_mb=4'")

	flag.Parse()
	if *src == "" || *mountpoint == "" {
		return nil, errors.New("Both src_root and mountpoint must be set")
	}

	opts := Options{
		Src:          *src,
		Mountpoint:   *mountpoint,
		Trace:        *trace,
		StrategyList: strings.Split(*serveStrategy, ";"),
		Async:        true,
		ThreadUnsafe: true,
		MaxReadahead: uint32(4 * 1024 * 1024),
	}
	for _, strategy := range opts.StrategyList {
		if strategy == "" {
			continue
		}
		switch {
		case strategy == "async":
			opts.Async = true
		case strategy == "sync":
			opts.Async = false
		case strategy == "serial":
			opts.ThreadUnsafe = true
		case strategy == "threadpool":
			opts.ThreadUnsafe = false
		case strings.HasPrefix(strategy, "readahead_mb="):
			readaheadStr := strings.Split(strategy, "=")[1]
			if readahead, err := strconv.ParseFloat(readaheadStr, 32); err == nil {
				opts.MaxReadahead = uint32(readahead * 1024 * 1024)
				break
			}
			fallthrough
		default:
			log.Fatal("Unrecognized strategy", strategy)
		}
	}
	log.Print(opts)
	return &opts, nil
}

func Runfs(opts *Options) {
	snap := snapshot.NewFileBackedSnapshot(opts.Src, "only")
	minfs := NewSlimMinFs(snap)

	if opts.Trace {
		fuse.Trace = true
		snapshot.Trace = true
		min.Trace = true
	}

	options := []fuse.MountOption{
		fuse.DefaultPermissions(),
		fuse.MaxReadahead(opts.MaxReadahead),
		fuse.FSName("slimfs"),
		fuse.Subtype("fs"),
		fuse.VolumeName("slimfs"),
	}
	if opts.Async {
		options = append(options, fuse.AsyncRead())
	}

	log.Print("About to Mount")
	fuse.Unmount(opts.Mountpoint)
	conn, err := fuse.Mount(opts.Mountpoint, fuse.MakeAlloc(), options...)
	if err != nil {
		log.Fatal("Couldn't mount", err)
	}

	var done chan error
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigchan
		log.Printf("Canceling")
		if done != nil {
			done <- errors.New("Caller canceled")
		}
	}()

	go func() {
		log.Println("pprof exit: ", http.ListenAndServe("localhost:6060", nil))
	}()

	defer func() {
		if err := fuse.Unmount(opts.Mountpoint); err != nil {
			log.Printf("error in call to Unmount(%s): %s", opts.Mountpoint, err)
			return
		}
		log.Printf("called Umount on %s", opts.Mountpoint)
	}()

	// Serve returns immediately and we wait for the first entry from the done channel before exiting main.
	// We only care about the first error from either the signal handler or from the first serve thread to return.
	// Exiting main will cause the remaining read threads to exit.
	log.Print("About to Serve")
	done = min.Serve(conn, minfs, opts.ThreadUnsafe)
	err = <-done
	log.Printf("Returning (might take a few seconds), err=%v", err)
}

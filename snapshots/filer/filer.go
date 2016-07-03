package filer

// Filer lets clients move snapshots between a snapshot store and a local filesystem.

type Snapshotter interface {
	// Snapshot takes a snapshot of path, returning the id of the created snapshot or an error if it could not be created.
	Snapshot(path string) (string, error)
}

type Checkouter interface {
	// Checkout checks the snapshot identified by id into the local filesystem. It returns the path of the checkout or an error.
	Checkout(id string) (string, error)

	// Release releases a checkout directory that was previously checked out. This allows an implementation to reuse a directory.
	// An implementation may allow checkouts to expire even if Release is not called. (e.g., by writing into a temp directory that is occasionally cleaned)
	Release(path string)
}

package gitdb

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"

	snap "github.com/scootdev/scoot/snapshot"
	"github.com/scootdev/scoot/snapshot/bundlestore"
)

// BundlestoreConfig defines how to talk to Bundlestore
type BundlestoreConfig struct {
	Store bundlestore.Store
}

type bundlestoreBackend struct {
	cfg *BundlestoreConfig
}

const bundlestoreIDText = "bs"

// "bs-gc-<bundle>-<stream>-<sha>"

func (b *bundlestoreBackend) parseID(id snap.ID, kind snapshotKind, extraParts []string) (snapshot, error) {
	if len(extraParts) != 3 {
		return nil, fmt.Errorf("cannot parse snapshot ID: expected 5 parts in bundlestore id: %s", id)
	}
	bundleName, streamName, sha := extraParts[0], extraParts[1], extraParts[2]

	if err := validSha(sha); err != nil {
		return nil, err
	}

	return &bundlestoreSnapshot{kind: kind, sha: sha, bundleName: bundleName, streamName: streamName}, nil
}

func (b *bundlestoreBackend) upload(s snapshot, db *DB) (snapshot, error) {
	// We only have to upload a localSnapshot
	switch s := s.(type) {
	case *tagsSnapshot:
		return s, nil
	case *streamSnapshot:
		// TODO(dbentley): we should upload to bundlestore if this commit is so recent
		// it might not be on every worker already.
		return s, nil
	case *bundlestoreSnapshot:
		return s, nil
	case *localSnapshot:
		return b.uploadLocalSnapshot(s, db)
	default:
		return nil, fmt.Errorf("cannot upload %v: unknown type %T", s, s)
	}
}

// git bundle create takes a rev list; it requires that it include a ref
// so we can't just do:
// git bundle create 545c88d71d40a49ebdfb1d268c724110793330d2..3060a3a519888957e13df75ffd09ea50f97dd03b
// instead, we have to write a temporary ref
// (nb: the subtracted revisions can be a commit, not a ref, so you can do:
// git bundle create 545c88d71d40a49ebdfb1d268c724110793330d2..master
// )
const bundlestoreTempRef = "reserved_scoot/bundlestore/__temp_for_writing"

func (b *bundlestoreBackend) uploadLocalSnapshot(s *localSnapshot, db *DB) (sn snapshot, err error) {
	// the sha of the commit we're going to use as the ref
	commitSha := s.sha

	// the revList to create the bundle
	// unless we can find a merge base, we'll just include the commit
	revList := bundlestoreTempRef

	// the name of the stream that this bundle requires
	streamName := ""

	switch s.kind {
	case kindGitCommitSnapshot:
		// For a git commit, we want a bundle that has just the diff compared to the stream
		// so find the merge base with our stream

		// The generated bundle will require either no prereqs or a commit that is in the stream
		if db.stream.cfg != nil && db.stream.cfg.RefSpec != "" {
			streamHead, err := db.dataRepo.Run("rev-parse", db.stream.cfg.RefSpec)
			if err != nil {
				return nil, err
			}

			mergeBase, err := db.dataRepo.RunSha("merge-base", streamHead, commitSha)

			// if err != nil, it just means we don't have a merge-base
			if err == nil {
				revList = fmt.Sprintf("%s..%s", mergeBase, bundlestoreTempRef)
				streamName = db.stream.cfg.Name
			}

		}
	case kindFSSnapshot:
		// For an FSSnapshot (which is stored as a git tree), create a git commit
		// with no parent.
		// (Eventually we could get smarter, e.g., if it's storing the output of running
		// cmd foo, we could try to find another run of foo and use that as a parent
		// to reduce the delta)

		// The generated bundle will require no prereqs.
		commitSha, err = db.dataRepo.RunSha("commit-tree", commitSha, "-m",
			"commit to distribute GitDB FSSnapshot via bundlestore")
	default:
		return nil, fmt.Errorf("unknown Snapshot kind: %v", s.kind)
	}

	// update the ref
	if _, err := db.dataRepo.Run("update-ref", bundlestoreTempRef, commitSha); err != nil {
		return nil, err
	}

	d, err := db.tmp.TempDir("bundle-")
	if err != nil {
		return nil, err
	}

	// we can't use tmpDir.TempFile() because we need the file to not exist
	bundleFilename := path.Join(d.Dir, fmt.Sprintf("commit-%s.bundle", s.SHA()))

	// create the bundle
	if _, err := db.dataRepo.Run("bundle", "create", bundleFilename, revList); err != nil {
		return nil, err
	}

	f, err := os.Open(bundleFilename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if err := b.cfg.Store.Write(s.sha, f); err != nil {
		return nil, err
	}

	return &bundlestoreSnapshot{sha: s.SHA(), kind: s.Kind(), bundleName: s.sha, streamName: streamName}, nil
}

type bundlestoreSnapshot struct {
	sha        string
	kind       snapshotKind
	bundleName string
	streamName string
}

func (s *bundlestoreSnapshot) ID() snap.ID {
	return snap.ID(strings.Join([]string{bundlestoreIDText, string(s.kind), s.bundleName, s.streamName, s.sha}, "-"))
}
func (s *bundlestoreSnapshot) Kind() snapshotKind { return s.kind }
func (s *bundlestoreSnapshot) SHA() string        { return s.sha }

func (s *bundlestoreSnapshot) Download(db *DB) error {
	if err := db.shaPresent(s.SHA()); err == nil {
		return nil
	}

	// TODO(dbentley): keep stats about bundlestore downloading
	// TODO(dbentley): keep stats about how long it takes to unbundle
	filename, err := s.downloadBundle(db)
	if err != nil {
		return err
	}

	// unbundle optimistically
	// this will succeed if we have all of the prerequisite objects
	if _, err := db.dataRepo.Run("bundle", "unbundle", filename); err == nil {
		return db.shaPresent(s.sha)
	}

	// we couldn't unbundle
	// see if it's because we're missing prereqs
	exitError := err.(*exec.ExitError)
	if exitError == nil || !strings.Contains(string(exitError.Stderr), "error: Repository lacks these prerequisite commits:") {
		return err
	}

	// we are missing prereqs, so let's try updating the stream that's the basis of the bundle
	// this likely happened because:
	// we're in a worker that started at time T1, when master pointed at commit C1
	// at time T2, a commit C2 was created in our stream
	// at time T3, a user ingested a git commit C3 whose ancestor is C2
	// GitDB in their scoot-snapshot-db picked a merge-base of C2, because T3-T2 was sufficiently
	// large (say, a half hour) that it's reasonable to assume its easy to get.
	// Now we've got the bundle for C3, which depends on C2, but we only have C1, so we have to
	// update our stream.
	if err := db.stream.updateStream(s.streamName, db); err != nil {
		return err
	}

	if _, err := db.dataRepo.Run("bundle", "unbundle", filename); err != nil {
		return err
	}

	return db.shaPresent(s.sha)
}

func (s *bundlestoreSnapshot) downloadBundle(db *DB) (filename string, err error) {
	d, err := db.tmp.TempDir("bundle-")
	if err != nil {
		return "", err
	}
	bundleFilename := path.Join(d.Dir, fmt.Sprintf("%s.bundle", s.bundleName))
	f, err := os.Create(bundleFilename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	r, err := db.bundles.cfg.Store.OpenForRead(s.bundleName)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(f, r); err != nil {
		return "", err
	}

	return f.Name(), nil
}

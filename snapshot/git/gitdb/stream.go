package gitdb

import (
	"fmt"

	snap "github.com/scootdev/scoot/snapshot"
)

// A Stream is a sequence of SnapshotWithHistory's that updates via a
// Git ref that can be fetched from a Git remote
type StreamConfig struct {
	// Name (used in IDs (so it should be short)
	// e.g. sm for a Stream following Source (repo)'s Master (branch)
	Name string

	// Remote to fetch from (e.g. https://github.com/scootdev/scoot)
	Remote string

	// Name of ref to follow in data repo (e.g. refs/remotes/upstream/master)
	RefSpec string
}

const streamIDText = "stream"
const streamIDFmt = "%s-%s-%s-%s"

type streamBackend struct {
	cfg *StreamConfig
}

func (b *streamBackend) parseID(id snap.ID, kind snapshotKind, extraParts []string) (*streamSnapshot, error) {
	if len(extraParts) != 2 {
		return nil, fmt.Errorf("cannot parse snapshot ID: expected 4 extraParts in stream id: %s", id)
	}
	streamName, sha := extraParts[0], extraParts[1]

	if err := validSha(sha); err != nil {
		return nil, err
	}

	return &streamSnapshot{streamName: streamName, kind: kind, sha: sha}, nil
}

// streamSnapshot represents a Snapshot that lives in a Stream
type streamSnapshot struct {
	sha        string
	kind       snapshotKind
	streamName string
}

func (s *streamSnapshot) ID() snap.ID {
	return snap.ID(fmt.Sprintf(streamIDFmt, streamIDText, s.kind, s.streamName, s.sha))
}
func (s *streamSnapshot) Kind() snapshotKind { return s.kind }
func (s *streamSnapshot) SHA() string        { return s.sha }

func (s *streamSnapshot) Download(db *DB) error {
	if err := db.shaPresent(s.SHA()); err == nil {
		// Already present!
		return nil
	}

	// TODO(dbentley): what if we've already fetched recently? We should figure out some way to
	// prevent that

	if db.stream == nil {
		return fmt.Errorf("cannot download snapshot %s: no streams configured", s.ID())
	}
	if s.streamName != db.stream.cfg.Name {
		return fmt.Errorf("cannot download snapshot %s: does not match stream %s", s.ID(), db.stream.cfg.Name)
	}

	// TODO(dbentley): keep stats about fetching (when we do it, last time we did it, etc.)
	if _, err := db.dataRepo.Run("fetch", db.stream.cfg.Remote); err != nil {
		return err
	}

	return db.shaPresent(s.SHA())
}

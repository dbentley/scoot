package gitdb

import (
	"fmt"

	snap "github.com/scootdev/scoot/snapshot"
)

const localIDText = "local"
const localIDFmt = "%s-%s-%s"

type localBackend struct {
	db *DB
}

// localSnap holds a reference to a value that is in the local DB
type localSnap struct {
	sha  string
	kind snapshotKind
}

// parse id as a local ID, with kind and remaining parts (after scheme and kind were parsed)
func (b *localBackend) parseID(id snap.ID, kind snapshotKind, parts []string) (*localSnap, error) {
	if len(parts) != 1 {
		return nil, fmt.Errorf("cannot parse snapshot ID: expected 3 parts in local id %s", id)
	}
	sha := parts[0]
	if err := validSha(sha); err != nil {
		return nil, err
	}

	return &localSnap{kind: kind, sha: sha}, nil
}

func (s *localSnap) ID() snap.ID {
	return snap.ID(fmt.Sprintf(localIDFmt, localIDText, s.kind, s.sha))
}
func (s *localSnap) Kind() snapshotKind { return s.kind }
func (s *localSnap) SHA() string        { return s.sha }

func (s *localSnap) Download(db *DB) error {
	// a localSnap is either present already or we have no way to download it
	return db.shaPresent(s.SHA())
}

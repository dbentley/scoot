package rsm

type SeqN int64

type Metadata struct {
	Seq SeqN
}

type Message struct {
	Data string
	Meta Metadata
}

type GetOpts struct {
}

type GetResponse struct {
	LastEntry SeqN
	Msgs      []Message
}

type Log interface {
	Get(lastSeen SeqN) ([]Message, error)

	Write(data string) (Metadata, error)
}

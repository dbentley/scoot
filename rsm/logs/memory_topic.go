package logs

import (
	"github.com/scootdev/scoot/rsm"
)

func NewMemoryTopicLog() *MemoryTopicLog {
	r := &MemoryTopicLog{
		reqCh: make(chan req),
	}
	go r.loop()
	return r
}

type MemoryTopicLog struct {
	reqCh     chan req
	msgs      []rsm.Message
	listeners []chan []rsm.Message
}

type req interface{}
type writeReq struct {
	data   string
	respCh chan rsm.Metadata
}
type getReq struct {
	lastSeen rsm.SeqN
	respCh   chan []rsm.Message
}

func (l *MemoryTopicLog) Close() error {
	close(l.reqCh)
	return nil
}

func (l *MemoryTopicLog) Write(data string) (rsm.Metadata, error) {
	respCh := make(chan rsm.Metadata)
	l.reqCh <- writeReq{data: data, respCh: respCh}
	return <-respCh, nil
}

func (l *MemoryTopicLog) Get(lastSeen rsm.SeqN) ([]rsm.Message, error) {
	respCh := make(chan []rsm.Message)
	l.reqCh <- getReq{lastSeen: lastSeen, respCh: respCh}
	return <-respCh, nil
}

func (l *MemoryTopicLog) loop() {
	for req := range l.reqCh {
		switch req := req.(type) {
		case writeReq:
			msg := rsm.Message{
				Data: req.data,
				Meta: rsm.Metadata{
					Seq: rsm.SeqN(len(l.msgs)),
				},
			}
			l.msgs = append(l.msgs, msg)
			req.respCh <- msg.Meta
			for _, listener := range l.listeners {
				listener <- []rsm.Message{msg}
			}
			l.listeners = nil
		case getReq:
			lastWritten := 0
			if len(l.msgs) > 0 {
				lastWritten = len(l.msgs)
			}
			if req.lastSeen < rsm.SeqN(int64(lastWritten)) {
				resp := append([]rsm.Message(nil),
					l.msgs[req.lastSeen+1:]...)
				req.respCh <- resp
			} else {
				l.listeners = append(l.listeners, req.respCh)
			}
		}
	}
}

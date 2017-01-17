package rsm

type ModelAndError struct {
	mdl Model
	err error
}

type MsgsAndError struct {
	msgs []Message
	err  error
}

type Manager interface {
	Apply(t Transition) (Model, Result, error)

	Get() (Model, error)

	Listen(lastSeen SeqN) chan ModelAndError
}

type mReq interface{}

func NewManagerImpl(log Log, machine Machine) *ManagerImpl {

	m := &ManagerImpl{
		log:       log,
		machine:   machine,
		lastSeen:  SeqN(-1),
		current:   machine.Empty(),
		listeners: make(map[SeqN]chan applyResult),

		reqCh:    make(chan mReq),
		updateCh: make(chan MsgsAndError),
	}
	go m.loop()
	go m.get(SeqN(-1))

	return m
}

type ManagerImpl struct {
	log     Log
	machine Machine

	lastSeen SeqN
	current  Model
	err      error

	reqCh    chan mReq
	updateCh chan MsgsAndError

	listeners map[SeqN]chan applyResult
}

func (m *ManagerImpl) Apply(t Transition) (Model, Result, error) {
	respCh := make(chan applyResult)
	m.reqCh <- applyReq{t: t, respCh: respCh}
	resp := <-respCh
	return resp.m, resp.r, resp.err
}

func (m *ManagerImpl) Get() (Model, error) {
	respCh := make(chan ModelAndError)
	m.reqCh <- respCh
	resp := <-respCh
	return resp.mdl, resp.err
}

func (m *ManagerImpl) Listen(lastSeen SeqN) chan ModelAndError {
	return nil
}

type applyReq struct {
	t      Transition
	respCh chan applyResult
}

type applyResult struct {
	m   Model
	r   Result
	err error
}

func (m *ManagerImpl) loop() {
	for m.reqCh != nil {
		select {
		case req, ok := <-m.reqCh:
			if !ok {
				m.reqCh = nil
				continue
			}
			switch req := req.(type) {
			case chan ModelAndError:
				if m.err != nil {
					req <- ModelAndError{err: m.err}
				} else {
					req <- ModelAndError{mdl: m.current.Copy()}
				}
			case applyReq:
				m.apply(req.t, req.respCh)

			}
		case msgsAndErr := <-m.updateCh:
			if m.err != nil {
				continue
			}
			if msgsAndErr.err != nil {
				m.err = msgsAndErr.err
			} else {
				m.process(msgsAndErr.msgs)
				go m.get(m.lastSeen)
			}
		}
	}

	if m.err != nil {
		// drain the in-flight get
		<-m.updateCh
	}
}

func (m *ManagerImpl) get(lastSeen SeqN) {
	msgs, err := m.log.Get(lastSeen)
	m.updateCh <- MsgsAndError{msgs, err}
}

func (m *ManagerImpl) apply(t Transition, respCh chan applyResult) {
	data, err := m.machine.Encode(t)
	if err != nil {
		respCh <- applyResult{err: err}
		return
	}
	meta, err := m.log.Write(data)
	if err != nil {
		respCh <- applyResult{err: err}
		return
	}

	m.listeners[meta.Seq] = respCh
}

func (m *ManagerImpl) process(msgs []Message) (err error) {
	for _, msg := range msgs {
		listener := m.listeners[msg.Meta.Seq]
		if listener != nil {
			delete(m.listeners, m.lastSeen)
		}

		if err := m.updateModel(msg, listener); err != nil {
			return err
		}
	}
	return nil
}

func (m *ManagerImpl) updateModel(msg Message, l chan applyResult) error {
	t, err := m.machine.Decode(msg.Data)
	if err != nil {
		if l != nil {
			l <- applyResult{err: err}
		}
		return err
	}

	r := m.current.Apply(t)
	m.lastSeen = msg.Meta.Seq
	if l != nil {
		l <- applyResult{r: r, m: m.current.Copy()}
	}
	return nil
}

type ID string

type MultiManager interface {
	Lookup(ID, Model) (Manager, error)
}

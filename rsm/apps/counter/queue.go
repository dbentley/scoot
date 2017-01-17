package main

import (
	"fmt"

	"github.com/scootdev/scoot/rsm"
)

type TaskID int64

type AgentID string

type TaskIDAndError struct {
	id  TaskID
	err error
}

const INVALID TaskID = -1

type Queue struct {
	mgr rsm.Manager
}

func NewQueue(mgr rsm.Manager) (*Queue, error) {
	_, result, err := mgr.Apply(&setCapacityTransition{37})
	if err != nil {
		return nil, err
	}

	if result := result.(*setCapacityResult); result.err != nil {
		return nil, result.err
	}
	return &Queue{mgr: mgr}, nil
}

func (q *Queue) AddTask(priority int, weight int) (TaskID, error) {
	_, res, err := q.mgr.Apply(&addTaskTransition{
		Priority: priority,
		Weight:   weight,
	})
	if err != nil {
		return INVALID, err
	}

	result := res.(*addTaskResult)
	return result.id, result.err
}

type QueueBusyError struct {
	capacity  int
	load      int
	newWeight int
}

func (e QueueBusyError) Error() string {
	return fmt.Sprintf("cannot add task: load %v + weight %v > capacity %v", e.load, e.newWeight, e.capacity)
}

func (e QueueBusyError) String() string {
	return e.Error()
}

func (q *Queue) Process(agent string) (TaskID, error) {
	_, res, err := q.mgr.Apply(&processTransition{Agent: AgentID(agent)})
	if err != nil {
		return INVALID, err
	}
	result := res.(*processResult)
	return result.id, result.err
}

type Status struct {
	WeightAhead int
	Processor   AgentID
}

func (q *Queue) Status(id TaskID) (st Status, err error) {
	m, err := q.mgr.Get()
	if err != nil {
		return st, err
	}
	qm := m.(*QueueModel)

	agent, ok := qm.done[id]
	if ok {
		return Status{WeightAhead: -1, Processor: agent}, nil
	}

	var idx int
	var def taskModel

	for i, t := range qm.pending {
		if t.id == id {
			idx = i
			def = t
			break
		}
	}

	ahead := 0
	for i, t := range qm.pending {
		if i == idx {
			continue
		}
		if t.priority > def.priority ||
			(t.priority == def.priority && i < idx) {
			ahead += t.weight
		}
	}

	return Status{WeightAhead: ahead}, nil
}

func (q *Queue) ListenForDone() chan TaskIDAndError {
	return nil
}

type taskModel struct {
	id       TaskID
	priority int
	weight   int
}

type QueueMachine struct {
}

type QueueModel struct {
	pending  []taskModel
	done     map[TaskID]AgentID
	capacity int
	nextID   TaskID
}

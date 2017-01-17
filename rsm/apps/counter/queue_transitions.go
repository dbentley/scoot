package main

import (
	"encoding/json"
	"fmt"

	"github.com/scootdev/scoot/rsm"
)

func (q *QueueMachine) Empty() rsm.Model {
	return &QueueModel{
		done:     make(map[TaskID]AgentID),
		capacity: -1,
		nextID:   TaskID(1),
	}
}

type jsonRep struct {
	SetCap  *setCapacityTransition
	AddTask *addTaskTransition
	Process *processTransition
}

func (q *QueueMachine) Encode(t rsm.Transition) (string, error) {
	var r jsonRep
	switch t := t.(type) {
	case *setCapacityTransition:
		r.SetCap = t
	case *addTaskTransition:
		r.AddTask = t
	case *processTransition:
		r.Process = t
	}
	d, err := json.Marshal(r)

	return string(d), err
}

func (q *QueueMachine) Decode(data string) (rsm.Transition, error) {
	r := jsonRep{}
	err := json.Unmarshal([]byte(data), &r)
	if err != nil {
		return nil, err
	}
	if r.SetCap != nil {
		return r.SetCap, nil
	}
	if r.AddTask != nil {
		return r.AddTask, nil
	}
	if r.Process != nil {
		return r.Process, nil
	}
	return nil, fmt.Errorf("none set in %q", data)
}

func (q *QueueModel) Copy() rsm.Model {
	done := make(map[TaskID]AgentID)
	for k, v := range q.done {
		done[k] = v
	}
	return &QueueModel{
		pending:  append([]taskModel(nil), q.pending...),
		done:     done,
		capacity: q.capacity,
		nextID:   q.nextID,
	}
}

func (q *QueueModel) Apply(t rsm.Transition) rsm.Result {
	switch t := t.(type) {
	case *setCapacityTransition:
		if q.capacity != -1 {
			return &setCapacityResult{fmt.Errorf("capacity already set: %v", q.capacity)}
		}
		q.capacity = t.Capacity
		return &setCapacityResult{}
	case *addTaskTransition:
		load := 0
		for _, t := range q.pending {
			load += t.weight
		}
		if t.Priority < 0 {
			return &addTaskResult{
				id:  INVALID,
				err: fmt.Errorf("priority must be non-negative; was %v", t.Priority),
			}
		}
		if load+t.Weight > q.capacity {
			return &addTaskResult{
				id:  INVALID,
				err: QueueBusyError{q.capacity, load, t.Weight},
			}
		}
		q.pending = append(q.pending, taskModel{q.nextID, t.Priority, t.Weight})
		id := q.nextID
		q.nextID = q.nextID + 1
		return &addTaskResult{id: id}
	case *processTransition:
		if len(q.pending) == 0 {
			return &processResult{
				id: INVALID,
			}
		}
		pri := -1
		idx := -1
		var task taskModel
		for i, t := range q.pending {
			if t.priority > pri {
				idx = i
				task = t
			}
		}
		q.pending = append(q.pending[:idx], q.pending[idx+1:]...)
		q.done[task.id] = t.Agent
		return &processResult{
			id:  task.id,
			err: nil,
		}
	default:
		panic(fmt.Errorf("Unknown type %T %+v", t, t))

	}
}

type setCapacityTransition struct {
	Capacity int
}

func (t *setCapacityTransition) Transition() {}

type setCapacityResult struct {
	err error
}

func (r *setCapacityResult) Result() {}

type addTaskTransition struct {
	Priority int
	Weight   int
}

func (t *addTaskTransition) Transition() {}

type addTaskResult struct {
	id  TaskID
	err error
}

func (r *addTaskResult) Result() {}

type processTransition struct {
	Agent AgentID
}

func (t *processTransition) Transition() {}

type processResult struct {
	id  TaskID
	err error
}

func (r *processResult) Result() {}

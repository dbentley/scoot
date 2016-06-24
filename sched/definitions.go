package sched

/*
 * Message representing a Job, Scoot can Schedule
 */
type Job struct {
	Id      string
	JobType string
	Tasks   []Task
}

/*
 * Message representing a Task that is part of a Scoot Job.
 */
type Task struct {
	Id         string
	Command    []string // Argv to run
	SnapshotId string
}